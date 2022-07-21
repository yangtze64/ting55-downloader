package book

import (
	"fmt"
	"log"
	"ting55-downloader/pkg/console"
)

type downloader struct {
	threadNum    int
	downloadMode int
	downloadPath string
}

var ModeMap = map[int]string{
	1: "Overwrite",
	2: "Skip",
	3: "Breakpoint Recovery",
	4: "Keep Original",
}

func NewDownloader(num int, mode int, dpath string) *downloader {
	return &downloader{
		threadNum:    num,
		downloadMode: mode,
		downloadPath: dpath,
	}
}

func (d *downloader) Download(book *Book) {
	log.Println(console.Green(fmt.Sprintf("Start Download Book,BookId:%d BookName:%s", book.Id, book.Title)))

}
