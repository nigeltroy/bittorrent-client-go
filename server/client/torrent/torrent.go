package torrent

import (
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
}

func unixTimeInterfaceToUtcString(i interface{}) string {
	utcTime := time.Unix(i.(int64), 0).String()

	return utcTime
}

func hasMultipleFiles(d map[string]interface{}) bool {
	_, ok := d["info"].(map[string]interface{})["files"]

	return ok
}

func getAnnounceList(i []interface{}) []string {
	announceList := make([]string, 0)

	for _, urls := range i {
		url := urls.([]interface{})[0]
		announceList = append(announceList, url.(string))
	}

	return announceList
}

func getPieces(i interface{}) [][]byte {
	pieces := make([][]byte, 0)
	piecesBytes := []byte(i.(string))
	hashLength := 20

	for i := 0; i < len(piecesBytes); i += hashLength {
		pieces = append(pieces, piecesBytes[i:i+hashLength])
	}

	return pieces
}

func getFiles(i []interface{}) []File {
	files := make([]File, 0)

	for _, fileInterface := range i {
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

func dictToMetainfo(d map[string]interface{}) Metainfo {
	metainfo := Metainfo{
		Announce:         d["announce"].(string),
		AnnounceList:     getAnnounceList(d["announce-list"].([]interface{})),
		Comment:          d["comment"].(string),
		CreatedBy:        d["created by"].(string),
		CreationDate:     unixTimeInterfaceToUtcString(d["creation date"]),
		HasMultipleFiles: hasMultipleFiles(d),
	}

	infoDictInterface := d["info"].(map[string]interface{})

	baseInfoDict := BaseInfoDict{
		Name:        infoDictInterface["name"].(string),
		Pieces:      getPieces(infoDictInterface["pieces"]),
		PieceLength: infoDictInterface["piece length"].(int64),
	}

	if !hasMultipleFiles(d) {
		metainfo.Info = InfoDictSingleFile{
			Base:   baseInfoDict,
			Length: infoDictInterface["length"].(int64),
		}
	} else {
		metainfo.Info = InfoDictMultipleFiles{
			Base:  baseInfoDict,
			Files: getFiles(infoDictInterface["files"].([]interface{})),
		}
	}

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
