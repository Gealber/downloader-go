package humanize

import (
	"fmt"
	"strings"
	"sync"	
)

//WriteCounter counts the number of bytes written to it.
//It implements to the io.Writer interface and we can pass
//this into io.TeeReader() which will report progress on each write cycle.
type WriteCounter struct {
	mu    sync.Mutex
	FileSize int
	Total uint64
}

func (wc *WriteCounter) Write(p []byte) (int, error) {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	n := len(p)
	wc.Total += uint64(n)	
	wc.PrintProgress()	
	return n, nil
}

//PrintProgress print the current progress of the download
func (wc *WriteCounter) PrintProgress() {
	// Clear the line by using a character return to go back to the start and remove
	// the remaining characters by filling it with spaces
	fmt.Printf("\r%s", strings.Repeat(" ", 30))

	// Return again and print current status of download
	// We use the humanize package to print the bytes in a meaningful way (e.g. 10 MB)
	fmt.Printf("\r%sDownloading... %s%s%s/%s%s", "\033[0;32m","\033[0;36m",Bytes(wc.Total),"\033[0;32m",Bytes(uint64(wc.FileSize)),"\033[0;37m")
}
