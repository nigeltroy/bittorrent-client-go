/*
Package torrent defines structs with fields based on user input and torrent
file contents in bencoded and decoded forms
*/
package torrent

import (
	"io"
	"log"
)

// Torrent struct with fields representing:
// - file contents (contents)
// - the metainfo associated with the torrent (metainfo)
type Torrent struct {
	contents io.Reader
	metainfo metainfo
}

// New() accepts an io.Reader object and returns a Torrent object
func New(reader io.Reader) (*Torrent, error) {
	metainfo, err := createMetainfoFromFileContents(reader)

	if err != nil {
		return nil, err
	}

	torrent := Torrent{
		contents: reader,
		metainfo: *metainfo,
	}

	log.Println("Successfully created a new torrent instance.")

	return &torrent, nil
}
