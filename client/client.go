package client

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type Client struct {
	Torrents []torrent
}

func interpretInput(input string) (string, error) {
	_, err := os.Open(input)
	if err == nil {
		return "path", nil
	}

	_, err = url.Parse(input)
	if err == nil {
		if strings.HasPrefix(input, "magnet") {
			return "magnetLink", nil
		}
		return "url", nil
	}

	_, err = hex.DecodeString(input)
	if err == nil {
		return "infoHash", nil
	}

	return "", errors.New("input is invalid")
}

func (c *Client) AddTorrent(input string) error {
	inputType, err := interpretInput(input)
	if err != nil {
		return err
	}

	switch inputType {
	case "url": // not supported yet
		return errors.New("urls are not supported yet")
	case "path":
		f, err := os.Open(input)
		if err != nil {
			return err
		}

		r := io.Reader(f)
		id := len(c.Torrents) + 1

		torrent, err := createTorrent(id, r, c.Torrents)
		if err != nil {
			return err
		}

		c.Torrents = append(c.Torrents, *torrent)
		log.Println(fmt.Sprintf("Added torrent %s at id %d", torrent.metainfo.info.name, id))
	case "magnetLink": // not supported yet
		return errors.New("magnet links are not supported yet")
	case "infoHash": // not supported yet
		return errors.New("info hashes are not supported yet")
	}

	return nil
}

func (c *Client) RemoveTorrent(id int) error {
	if len(c.Torrents) == 0 {
		return errors.New("client has no torrents to remove")
	} else if id > len(c.Torrents) {
		return errors.New("id to remove is out of bounds")
	}
	c.Torrents = append(c.Torrents[:id-1], c.Torrents[id:]...)

	// Reset the ids of the torrents
	for _, t := range c.Torrents[id-1:] {
		t.id -= 1
	}

	log.Println(fmt.Sprintf("Removed torrent at id %d", id))
	return nil
}

func (c *Client) ShowTorrents() {
	fmt.Println("----- All Torrents -----")
	fmt.Println()
	fmt.Println("ID     Name")
	for _, t := range c.Torrents {
		fmt.Printf("%s      %s\n", strconv.Itoa(t.id), t.metainfo.info.name)
	}
	fmt.Println()
	fmt.Println("------------------------")
}
