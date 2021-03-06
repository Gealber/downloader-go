package downloadtools

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"sync"

	"github.com/Gealber/downloader-go/humanize"
)

//DownloadPart download the part of the video or file
//that is assign to download. Takes three arguments
//tempName contain the name of the temporary file
//where is going to be stored the data, url the url
//from where is going to download, position the position
// of the file to start copying.
func DownloadPart(wg *sync.WaitGroup, counter *humanize.WriteCounter, tempName, url, part string) {
	//setting up the client to make the request	
	client := http.Client{}
	request, err := http.NewRequest("GET", url, nil)

	//setting up the requests
	request.Header.Set("Range", part)
	response, err := client.Do(request)
	checkError(err, "fatal")
	defer response.Body.Close()

	//creating the temporary file and copying
	// the response to it
	file, err := os.Create(tempName)
	checkError(err, "fatal")
	defer file.Close()

	buf := make([]byte, 4096)
	_, err = io.CopyBuffer(file, io.TeeReader(response.Body, counter), buf)
	if err != nil {
		if err == io.EOF {
			wg.Done()
			return
		}
		log.Fatal(err.Error())
	}

	wg.Done()
}

//Download is called in case the server doesn't accept ranges
//requests. In this case is necesary to download the file on
//his full length.
func Download(fileName string, url string) {
	response, err := http.Get(url)
	checkError(err, "print")
	defer response.Body.Close()

	file, err := os.Create(fileName + ".temp")
	checkError(err, "print")
	defer file.Close()

	//counter := &humanize.WriteCounter{}
	_, err = io.Copy(file, response.Body)
	checkError(err, "print")

	err = os.Rename(fileName+".temp", fileName)
	checkError(err, "print")	
}

//HandleRangeDownload handle the download of the file by making a
//range request
func HandleRangeDownload(length, url, name, downloadDir string, threads int) error {
	size, err := strconv.Atoi(length)
	if err != nil {
		return err
	}

	start, step, end, rest := SetBenchmarks(size, threads)

	var wg sync.WaitGroup
	wg.Add(threads)
	//initializing the goroutines
	counter := &humanize.WriteCounter{
		FileSize: size,
	}
	for i := 0; i < threads; i++ {
		part := fmt.Sprintf("bytes=%d-%d", start, end)
		start = end + 1
		if i == threads-1 {
			end = end + step + rest
		} else {
			end = end + step
		}
		tempName := path.Join(downloadDir,fmt.Sprintf("%d.temp", i))
		go DownloadPart(&wg, counter, tempName, url, part)
	}

	wg.Wait()
	fmt.Printf("\nJoining files...\n")
	err = JoinFiles(name,downloadDir)
	if err != nil {
		return err
	}
	return nil
}

//HandleDownload make all the stuff related to the download
func HandleDownload(name, url, downloadDir string, threads int) error {
	if url != "" && name != "" && threads < 8 {
		name = path.Join(downloadDir,name)
		log.Println("Download started...")

		accept, size := RangeAndSize(url)
		if accept {

			err := HandleRangeDownload(size, url, name,downloadDir, threads)
			if err != nil {
				return err
			}
			return nil
		}
		Download(name, url)
		return nil
	}
	return nil
}
