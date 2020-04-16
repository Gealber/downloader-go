package downloadtools

import (
	"downloader-go/humanize"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
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
	checkError(err, "panic")
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
	checkError(err, "fatal")
	defer response.Body.Close()

	file, err := os.Create(fileName + ".temp")
	checkError(err, "fatal")
	defer file.Close()

	counter := &humanize.WriteCounter{}
	_, err = io.Copy(file, io.TeeReader(response.Body, counter))
	checkError(err, "panic")

	err = os.Rename(fileName+".temp", fileName)
	checkError(err, "fatal")

	log.Println("The downloaded has finished")
}

//HandleRangeDownload handle the download of the file by making a
//range request
func HandleRangeDownload(length, url, name string, threads int) {
	size, err := strconv.Atoi(length)
	if err != nil {
		log.Panicln(err.Error())
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

		go DownloadPart(&wg, counter, fmt.Sprintf("%d.temp", i), url, part)
	}

	wg.Wait()
	fmt.Printf("\nJoining files...\n")
	JoinFiles(name)
}

//HandleDownload make all the stuff related to the download
func HandleDownload(name, url string, threads int) {
	if url != "" && name != "" && threads < 8 {
		_ = os.Mkdir("Downloads", 0700)
		err := os.Chdir("Downloads")

		if err != nil {
			log.Fatal(err)
		}

		log.Println("Download started...")

		accept, size := RangeAndSize(url)
		if accept {

			HandleRangeDownload(size, url, name, threads)
			return
		}
		Download(name, url)
		return
	}
}
