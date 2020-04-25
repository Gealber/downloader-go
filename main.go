// +build linux

package main

import (
	"github.com/Gealber/downloader-go/downloadtools"
	"flag"
	"log"
	"os"
	"path"
)

var (
	url     string
	name    string
	threads int
)

func main() {
	flag.StringVar(&url, "u", "", "url to be fetched")
	flag.StringVar(&name, "n", "", "name of the file to be fetched")
	flag.IntVar(&threads, "t", 4, "threads to be used in the connection")

	flag.Parse()
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatalln(err)
	}
	currentDir = path.Join(currentDir, "Downloads")
	downloadtools.HandleDownload(name, url, currentDir, threads)

	log.Println("You must provide an url and a name at least")
}
