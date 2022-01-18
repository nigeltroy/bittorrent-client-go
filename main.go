package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/nigeltroy/bittorrent-client-go/client"
)

func initializeLogger() *os.File {
	const logFileName string = "logs.txt"

	f, err := os.OpenFile(logFileName, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(f)
	return f
}

func printHelp() {
	fmt.Println("exit: exits program")
	fmt.Println("print: prints all torrents")
	fmt.Println("add <file path>: adds torrent from <file path>")
	fmt.Println("remove <prefix>: removes first torrent that starts with <prefix>")
}

func runClient() {
	log.Println("Program bittorrent-client-go has started...")
	fmt.Println("Type 'help' for valid commands")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)
	torrentClient := client.Client{}
	exited := false
	for !exited {
		scanner.Scan()
		input := strings.Fields(scanner.Text())
		if len(input) == 0 {
			continue
		}

		cmd := input[0]
		switch cmd {
		case "exit":
			exited = true
		case "help":
			printHelp()
		case "print":
			torrentClient.ShowTorrents()
		case "add":
			if len(input) != 2 {
				fmt.Print("not enough input arguments")
				continue
			}

			err := torrentClient.AddTorrent(input[1])
			if err != nil {
				fmt.Println(err)
			}
		case "remove":
			if len(input) != 2 {
				fmt.Print("not enough input arguments")
				continue
			}

			err := torrentClient.RemoveTorrent(input[1])
			if err != nil {
				fmt.Println(err)
			}
		}

		fmt.Println()
	}
}

func main() {
	f := initializeLogger()
	defer f.Close()

	runClient()
}
