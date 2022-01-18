package client

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"strconv"
)

type inputType int

const (
	// Just add support for file paths for now
	path inputType = iota
	// url
	// info hash
	// magnet link
	invalid
)

type torrent struct {
	// torrents are considered clients here, so this is the peer id
	id       []byte
	metainfo metainfo
	trackers []tracker
	peers    []peer
}

func interpretInput(input string) inputType {
	_, err := os.Open(input)
	if err == nil {
		return path
	}

	return invalid
}

func createTorrentFromFileContents(path string) (*torrent, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	r := io.Reader(f)
	m, err := newMetainfo(r)
	if err != nil {
		return nil, err
	}

	return &torrent{metainfo: *m}, err
}

func createId() []byte {
	// Azureus style with arbitrary client id and version number
	base := []byte("-GG0001-")
	randSuffix := make([]byte, 0)
	for i := 0; i < 12; i++ {
		bytes := []byte(strconv.Itoa(rand.Intn(10)))
		randSuffix = append(randSuffix, bytes[0])
	}

	id := append(base, randSuffix...)
	return id
}

func newTorrent(input string) (*torrent, error) {
	var t *torrent
	var err error

	inputType := interpretInput(input)
	switch inputType {
	case path:
		t, err = createTorrentFromFileContents(input)
		if err != nil {
			return nil, err
		}
	case invalid:
		return nil, fmt.Errorf("input %s is invalid", input)
	}

	t.id = createId()
	err = t.requestPeers()
	if err != nil {
		return nil, err
	}

	return t, nil
}
