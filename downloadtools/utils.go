package downloadtools

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func checkError(err error, option string) {
	if err != nil {
		if option == "fatal" {
			log.Fatalln(err.Error())
		} else {
			log.Panicln(err.Error())
		}
	}
}

//RangeAndSize is used to know if the server
//accept to request an specific part of the file.
//And also to retrieve the length of the file
func RangeAndSize(url string) (bool, string) {
	client := &http.Client{}
	response, err := client.Head(url)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer response.Body.Close()

	size := response.Header.Get("Content-Length")
	if response.Header.Get("Accept-Ranges") == "" {
		return false, size
	}
	return true, size
}


//JoinFiles join the files in the current directory
//to a final file with the gicen name
func JoinFiles(name string) error {
	finalFile, err := os.OpenFile(name, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Panicln(err.Error())
	}
	defer finalFile.Close()

	files, err := ioutil.ReadDir(".")
	buf := make([]byte, 4096)
	for _, f := range files {
		tempFile, err := os.Open(f.Name())
		if err != nil {
			log.Panicln(err.Error())
		}

		if f.Name() != finalFile.Name() && f.Name() != "meta" {
			_,err := io.CopyBuffer(finalFile, tempFile,buf)
			if err != nil {
				return err
			}
			tempFile.Close()
			os.Remove(f.Name())
		}
		tempFile.Close()
	}
	return nil
}

//SetBenchmarks set the marks for parts
func SetBenchmarks(size int, threads int) (int, int, int, int) {
	step := size / threads
	start := 0
	end := step - 1
	rest := size % threads
	return start, step, end, rest
}