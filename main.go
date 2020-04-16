// +build linux

package main

import (
	"downloader-go/downloadtools"
	"flag"
	"log"	
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

	downloadtools.HandleDownload(name,url,threads)
	
	log.Println("You must provide an url and a name at least")
}