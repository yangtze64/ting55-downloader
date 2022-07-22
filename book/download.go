package book

import (
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"log"
	"time"
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
	tmpl := `{{ green "With funcs:" }} {{ bar . "[" "#" (cycle . "=>" "=>" "=>" "=>" ) "." "]" | green}} {{speed . | green }} {{percent .}} {{string . "my_green_string" | green}} {{string . "my_blue_string" | blue}}`
	bar := pb.ProgressBarTemplate(tmpl).Start(10000)
	for i := 0; i < 10000; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	bar.Set(pb.SIBytesPrefix, true)
	bar.Finish()
}
