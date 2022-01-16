package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nigeltroy/bittorrent-client-go/client"
)

const logFileName string = "logs.txt"

func main() {
	logFile, err := os.OpenFile(logFileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)

	if err != nil {
		log.Fatal(err)
	}

	defer logFile.Close()
	log.SetOutput(logFile)
	log.Printf("\n\nNew program run\n")

	var exited bool
	var input, command string
	var inputArgs []string
	var torrentClient client.Client

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("BitTorrent Client CLI")
	fmt.Println("-----------------------------")
	fmt.Println()
	fmt.Println("Type 'help' for valid commands")

	for !exited {
		fmt.Println()
		fmt.Print("Enter command: ")

		scanner.Scan()
		input = scanner.Text()

		if input == "exit" {
			exited = true
			continue
		}

		inputArgs = strings.Fields(input)
		command = inputArgs[0]

		switch command {
		case "help":
			fmt.Println("help: prints valid commands")
			fmt.Println("exit: exits program")
			fmt.Println("print: prints all torrents")
			fmt.Println("add <file path>: adds torrent from <file path>")
			fmt.Println("remove <id>: remove torrent with id <id>")
		case "exit":
			exited = true
		case "print":
			torrentClient.ShowTorrents()
		case "add":
			if len(inputArgs) != 2 {
				log.Println(errors.New("not enough input arguments to add torrent"))
				fmt.Println("Not enough input arguments supplied to add torrent")
				continue
			}

			err = torrentClient.AddTorrent(inputArgs[1])
			if err != nil {
				log.Println(err)
				fmt.Println(err)
			}
		case "remove":
			if len(inputArgs) != 2 {
				log.Println(errors.New("not enough input arguments to remove torrent"))
				fmt.Println("Id not supplied to remove torrent")
				continue
			}

			id, err := strconv.Atoi(inputArgs[1])
			if err != nil {
				log.Println(err)
				fmt.Println("Id supplied is not a valid integer")
				continue
			}

			err = torrentClient.RemoveTorrent(id)
			if err != nil {
				log.Println(err)
				fmt.Println(err)
			}
		}
	}
}
