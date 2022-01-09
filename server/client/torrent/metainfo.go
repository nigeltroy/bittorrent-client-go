package torrent

import (
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/marksamman/bencode"
)

type metadata struct {
	comment      string
	createdBy    string
	creationDate string
}

type file struct {
	length int64
	path   string
}

type info struct {
	name             string
	hasMultipleFiles bool
	pieceLength      int64
	pieces           [][]byte
	length           int64  // single file mode
	files            []file // multiple file mode
}

type infoHash struct {
	hexString        string
	urlEncodedString string
}

type metainfo struct {
	announce     string
	announceList []string
	metadata     metadata
	info         info
	infoHash     infoHash
}

func (info info) GetLength() (int64, error) {
	hasMultipleFiles := info.hasMultipleFiles
	length := info.length

	if hasMultipleFiles {
		return 0, errors.New("length requested from torrent with multiple files")
	} else if !hasMultipleFiles && length == 0 {
		return 0, errors.New("length has not been set or is equal to zero")
	}

	return info.length, nil
}

func extractAnnounceFromDecodedStream(decodedStream map[string]interface{}) (string, error) {
	announce, announceExists := decodedStream["announce"]

	if !announceExists {
		return "", errors.New("announce not found in decoded file contents")
	}

	return announce.(string), nil
}

func extractAnnounceListFromDecodedStream(decodedStream map[string]interface{}) []string {
	announceListInterface, announceListExists := decodedStream["announce-list"]

	if !announceListExists {
		return nil
	}

	announceList := make([]string, 0)

	for _, urls := range announceListInterface.([]interface{}) {
		url := urls.([]interface{})[0].(string)
		announceList = append(announceList, url)
	}

	return announceList
}

func extractCommentFromDecodedStream(decodedStream map[string]interface{}) string {
	commentInterface, commentExists := decodedStream["comment"]

	if !commentExists {
		return ""
	}

	return commentInterface.(string)
}

func extractCreatedByFromDecodedStream(decodedStream map[string]interface{}) string {
	createdByInterface, createdByExists := decodedStream["created by"]

	if !createdByExists {
		return ""
	}

	return createdByInterface.(string)
}

func extractCreationDateFromDecodedStream(decodedStream map[string]interface{}) string {
	creationDateInterface, creationDateExists := decodedStream["creation date"]

	if !creationDateExists {
		// Creation date is a mandatory key, but since it isn't critical to
		// the operation of this client, it is just logged and omitted.
		log.Println("Creation date not found in decoded file contents.")
		return ""
	}

	creationDate := time.Unix(creationDateInterface.(int64), 0).String()

	return creationDate
}

func createMetadata(decodedStream map[string]interface{}) metadata {
	return metadata{
		comment:      extractCommentFromDecodedStream(decodedStream),
		createdBy:    extractCreatedByFromDecodedStream(decodedStream),
		creationDate: extractCreationDateFromDecodedStream(decodedStream),
	}
}

func extractNameFromInfoDictInterface(infoDictInterface map[string]interface{}) (string, error) {
	nameInterface, nameExists := infoDictInterface["name"]

	if !nameExists {
		return "", errors.New("name not found in info dict of decoded file contents")
	}

	name := nameInterface.(string)

	return name, nil
}

func extractHasMultipleFilesFromInfoDictInterface(infoDictInterface map[string]interface{}) (bool, error) {
	_, filesExist := infoDictInterface["files"]
	_, lengthExists := infoDictInterface["length"]

	// Check if files !XOR length exists in info dict
	if filesExist == lengthExists {
		return false, fmt.Errorf("both files and length exist or do not exist -> value: %t", filesExist)
	}

	hasMultipleFiles := filesExist

	return hasMultipleFiles, nil
}

func extractPieceLengthFromInfoDictInterface(infoDictInterface map[string]interface{}) (int64, error) {
	pieceLengthInterface, pieceLengthExists := infoDictInterface["piece length"]

	if !pieceLengthExists {
		return 0, errors.New("piece length not found in info dict of decoded file contents")
	}

	pieceLength := pieceLengthInterface.(int64)

	return pieceLength, nil
}

func extractPiecesFromInfoDictInterface(infoDictInterface map[string]interface{}) ([][]byte, error) {
	piecesInterface, piecesExists := infoDictInterface["pieces"]

	if !piecesExists {
		return nil, errors.New("pieces not found in info dict of decoded file contents")
	}

	piecesBytes := []byte(piecesInterface.(string))
	hashLength := 20
	pieces := make([][]byte, 0)

	for i := 0; i < len(piecesBytes); i += hashLength {
		pieces = append(pieces, piecesBytes[i:i+hashLength])
	}

	return pieces, nil
}

func extractLengthFromInfoDictInterface(infoDictInterface map[string]interface{}) int64 {
	lengthInterface, lengthExists := infoDictInterface["length"]

	if !lengthExists {
		return 0
	}

	length := lengthInterface.(int64)

	return length
}

func extractFilesFromInfoDictInterface(infoDictInterface map[string]interface{}) []file {
	filesInterface, filesExist := infoDictInterface["files"]

	if !filesExist {
		return nil
	}

	files := make([]file, 0)

	for _, fileInterface := range filesInterface.([]interface{}) {
		fileDict := fileInterface.(map[string]interface{})
		pathInterface := fileDict["path"].([]interface{})
		pathParts := make([]string, 0)

		for _, pathPart := range pathInterface {
			pathParts = append(pathParts, pathPart.(string))
		}

		file := file{
			length: fileDict["length"].(int64),
			path:   strings.Join(pathParts, "/"),
		}

		files = append(files, file)
	}

	return files
}

func createInfo(decodedStream map[string]interface{}) (*info, error) {
	infoDictInterface, infoExists := decodedStream["info"]

	if !infoExists {
		return nil, errors.New("info not found in decoded file contents")
	}

	infoDict := infoDictInterface.(map[string]interface{})

	name, err := extractNameFromInfoDictInterface(infoDict)

	if err != nil {
		return nil, err
	}

	hasMultipleFiles, err := extractHasMultipleFilesFromInfoDictInterface(infoDict)

	if err != nil {
		return nil, err
	}

	pieceLength, err := extractPieceLengthFromInfoDictInterface(infoDict)

	if err != nil {
		return nil, err
	}

	pieces, err := extractPiecesFromInfoDictInterface(infoDict)

	if err != nil {
		return nil, err
	}

	length := extractLengthFromInfoDictInterface(infoDict)
	files := extractFilesFromInfoDictInterface(infoDict)

	return &info{
		name:             name,
		hasMultipleFiles: hasMultipleFiles,
		pieceLength:      pieceLength,
		pieces:           pieces,
		length:           length,
		files:            files,
	}, nil
}

func convertInfoToInfoHash(infoDictInterface map[string]interface{}) infoHash {
	encodedInfoDict := bencode.Encode(infoDictInterface)
	encryptedInfoDict := sha1.Sum(encodedInfoDict)

	infoHash := infoHash{
		hexString:        hex.EncodeToString(encryptedInfoDict[:]),
		urlEncodedString: url.QueryEscape(string(encryptedInfoDict[:])),
	}

	return infoHash
}

func extractMetainfoFromDecodedStream(decodedStream map[string]interface{}) (*metainfo, error) {
	announce, err := extractAnnounceFromDecodedStream(decodedStream)

	if err != nil {
		return nil, err
	}

	announceList := extractAnnounceListFromDecodedStream(decodedStream)
	metadata := createMetadata(decodedStream)

	info, err := createInfo(decodedStream)

	if err != nil {
		return nil, err
	}

	// If we can create an info instance, then we know that the info key
	// exists in the file contents, so we use that directly
	infoHash := convertInfoToInfoHash(decodedStream["info"].(map[string]interface{}))

	return &metainfo{
		announce:     announce,
		announceList: announceList,
		metadata:     metadata,
		info:         *info,
		infoHash:     infoHash,
	}, nil
}

func createMetainfoFromFileContents(reader io.Reader) (*metainfo, error) {
	decodedStream, err := bencode.Decode(reader)

	if err != nil {
		return nil, err
	}

	metainfo, err := extractMetainfoFromDecodedStream(decodedStream)

	if err != nil {
		return nil, err
	}

	log.Println("Successfully created metainfo from file contents.")

	return metainfo, nil
}
