package torrent

import (
	"crypto/sha1"
	"encoding/hex"
	"io"
	"strings"
	"time"

	"github.com/marksamman/bencode"
)

type BaseInfoDict struct {
	Name        string
	PieceLength int64
	Pieces      [][]byte
}

type File struct {
	Length int64
	Path   string
}

type InfoDictSingleFile struct {
	Base   BaseInfoDict
	Length int64
}

type InfoDictMultipleFiles struct {
	Base  BaseInfoDict
	Files []File
}

type Metainfo struct {
	Announce         string
	AnnounceList     []string
	Comment          string
	CreatedBy        string
	CreationDate     string
	HasMultipleFiles bool
	Info             interface{}
	InfoHash         string
}

func unixTimeInterfaceToUtcString(i interface{}) string {
	utcTime := time.Unix(i.(int64), 0).String()

	return utcTime
}

func hasMultipleFiles(torrentDict map[string]interface{}) bool {
	_, ok := torrentDict["info"].(map[string]interface{})["files"]

	return ok
}

func getAnnounceList(torrentDict map[string]interface{}) []string {
	if torrentDict["interface"] != nil {
		announceListInterface := torrentDict["announce-list"].([]interface{})
		announceList := make([]string, 0)

		for _, urls := range announceListInterface {
			url := urls.([]interface{})[0]
			announceList = append(announceList, url.(string))
		}

		return announceList
	}

	return nil
}

func getComment(torrentDict map[string]interface{}) string {
	commentInterface := torrentDict["comment"]

	if commentInterface != nil {
		return commentInterface.(string)
	}

	return ""
}

func getCreatedBy(torrentDict map[string]interface{}) string {
	createdByInterface := torrentDict["created by"]

	if createdByInterface != nil {
		return createdByInterface.(string)
	}

	return ""
}

func getPieces(piecesInterface interface{}) [][]byte {
	pieces := make([][]byte, 0)
	piecesBytes := []byte(piecesInterface.(string))
	hashLength := 20

	for i := 0; i < len(piecesBytes); i += hashLength {
		pieces = append(pieces, piecesBytes[i:i+hashLength])
	}

	return pieces
}

func getFiles(filesInterface []interface{}) []File {
	files := make([]File, 0)

	for _, fileInterface := range filesInterface {
		fileDict := fileInterface.(map[string]interface{})
		pathInterface := fileDict["path"].([]interface{})
		pathParts := make([]string, 0)

		for _, pathPart := range pathInterface {
			pathParts = append(pathParts, pathPart.(string))
		}

		file := File{
			Length: fileDict["length"].(int64),
			Path:   strings.Join(pathParts, "/"),
		}

		files = append(files, file)
	}

	return files
}

func getInfoHash(infoDictInterface map[string]interface{}) string {
	encodedInfoDict := bencode.Encode(infoDictInterface)
	encryptedInfoDict := sha1.Sum(encodedInfoDict)
	infohash := hex.EncodeToString(encryptedInfoDict[:])

	return infohash
}

func getInfoDict(infoDictInterface map[string]interface{}, hasMultipleFiles bool) interface{} {
	baseInfoDict := BaseInfoDict{
		Name:        infoDictInterface["name"].(string),
		Pieces:      getPieces(infoDictInterface["pieces"]),
		PieceLength: infoDictInterface["piece length"].(int64),
	}

	if !hasMultipleFiles {
		return InfoDictSingleFile{
			Base:   baseInfoDict,
			Length: infoDictInterface["length"].(int64),
		}
	}

	return InfoDictMultipleFiles{
		Base:  baseInfoDict,
		Files: getFiles(infoDictInterface["files"].([]interface{})),
	}
}

func dictToMetainfo(torrentDict map[string]interface{}) Metainfo {
	metainfo := Metainfo{
		Announce:         torrentDict["announce"].(string),
		AnnounceList:     getAnnounceList(torrentDict),
		Comment:          getComment(torrentDict),
		CreatedBy:        getCreatedBy(torrentDict),
		CreationDate:     unixTimeInterfaceToUtcString(torrentDict["creation date"]),
		HasMultipleFiles: hasMultipleFiles(torrentDict),
		InfoHash:         getInfoHash(torrentDict["info"].(map[string]interface{})),
	}

	metainfo.Info = getInfoDict(torrentDict["info"].(map[string]interface{}), hasMultipleFiles(torrentDict))

	return metainfo
}

func Decode(reader io.Reader) (Metainfo, error) {
	var metainfo Metainfo
	metainfoDict, err := bencode.Decode(reader)

	if err != nil {
		return metainfo, err
	}

	metainfo = dictToMetainfo(metainfoDict)

	return metainfo, nil
}
