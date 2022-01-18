/*
The client package defines a client, which is a collection of torrents.
Each torrent stores the metainfo of the torrent and the peer connection
info.
*/
package client

import (
	"fmt"
	"strings"
)

type Client struct {
	Torrents []torrent
}

func (c *Client) AddTorrent(input string) error {
	torrent, err := newTorrent(input)
	if err != nil {
		return err
	}

	// Later, to make this more efficient, do the following check immediately
	// after setting the metainfo
	for _, t := range c.Torrents {
		name := t.metainfo.Info.Name
		if name == torrent.metainfo.Info.Name {
			return fmt.Errorf("torrent %s already exists", name)
		}
	}

	c.Torrents = append(c.Torrents, *torrent)
	return nil
}

func (c *Client) RemoveTorrent(prefix string) error {
	for i, t := range c.Torrents {
		name := t.metainfo.Info.Name
		if strings.HasPrefix(name, prefix) {
			fmt.Printf("Removed torrent %s\n", name)
			c.Torrents = append(c.Torrents[:i], c.Torrents[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("no torrent matches prefix %s", prefix)
}

func (c Client) ShowTorrents() {
	for i, t := range c.Torrents {
		fmt.Printf("%d. %s\n", i+1, t.metainfo.Info.Name)
	}
}
