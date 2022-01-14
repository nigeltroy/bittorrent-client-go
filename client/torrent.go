package client

import (
	"errors"
	"io"
	"log"

	"github.com/marksamman/bencode"
)

type state int

const (
	started state = iota
	stopped
)

type torrent struct {
	id       int
	state    state
	metainfo metainfo
	tracker  tracker
	peers    []peer
}

func (t *torrent) setMetainfo(r io.Reader) error {
	decodedStream, err := bencode.Decode(r)
	if err != nil {
		return err
	}

	metainfo, err := extractMetainfoFromDecodedStream(decodedStream)
	if err != nil {
		return err
	}

	t.metainfo = *metainfo

	return nil
}

func (t *torrent) setTracker() error {
	return nil
}

func (t *torrent) setPeers() error {
	return nil
}

func createTorrent(id int, r io.Reader, torrents []torrent) (*torrent, error) {
	log.Println("Creating torrent...")

	torrent := torrent{
		id:    id,
		state: stopped,
	}

	err := torrent.setMetainfo(r)
	if err != nil {
		return nil, err
	}
	log.Println("Successfully set metainfo")

	for _, t := range torrents {
		if t.metainfo.info.name == torrent.metainfo.info.name {
			return nil, errors.New("torrent already exists in client")
		}
	}

	err = torrent.setTracker()
	if err != nil {
		return nil, err
	}
	log.Println("Successfully set tracker")

	err = torrent.setPeers()
	if err != nil {
		return nil, err
	}
	log.Println("Successfully set peers")

	return &torrent, nil
}
