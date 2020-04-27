package downloadtools

import (
	"downloader-go/downloadtools"
	"os"
	"testing"
)

var URL = "http://127.0.0.1:5000/static/videos/porn.mp4"


func BenchmarkPartDownload(b *testing.B) {
	for i := 0; i < b.N; i++ {
		downloadtools.HandleDownload("porn.mp4",URL,"download_test",1)
		os.Remove("porn.mp4")
	}
}


func BenchmarkSingleDownload(b *testing.B) {
	for i := 0; i < b.N; i++ {
		downloadtools.Download("porn.mp4", URL)
		os.Remove("porn.mp4")
	}
}