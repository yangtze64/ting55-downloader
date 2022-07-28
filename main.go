package main

import (
	"flag"
	"fmt"
	"log"
	"runtime"
	"ting55-downloader/book"
	"ting55-downloader/pkg/console"
)

var (
	id           = flag.Int("i", 0, "bookId")
	dp           = flag.String("p", "download/", "download path")
	num          = flag.Int("n", 1, "thread num default 1")
	mode         = flag.Int("m", 1, "download mod optional 1:Overwrite,2:Skip,3:Breakpoint Recovery,4:Keep Original")
	chapter      = flag.Int("c", 0, "chapter n")
	maxThreadNum = runtime.NumCPU() * 2
)

func main() {
	flag.Parse()
	bookId := *id
	downloadPath := *dp
	threadNum := *num
	downloadMode := *mode
	no := *chapter
	if bookId == 0 {
		log.Fatal(console.Red("请设置bookId"))
	}
	if threadNum > maxThreadNum {
		log.Fatal(console.Red("thread num 不可超过CUP*2核心数"))
	}
	downloadModeVal, ok := book.ModeMap[downloadMode]
	if !ok {
		log.Fatal(console.Red("download mod 1:Overwrite,2:Skip,3:Keep Original"))
	}
	log.Println(console.Yellow(fmt.Sprintf("bookId:%d, threadNum:%d, downloadPath:%s, downloadMode:%s", bookId, threadNum, downloadPath, downloadModeVal)))
	bookInfo := book.Parse(bookId)
	fmt.Printf("%#+v\n", bookInfo)
	downloader := book.NewDownloader(bookInfo, threadNum, downloadMode, downloadPath)
	if no > 0 {
		if no <= bookInfo.Number {
			err := downloader.DownloadAudio(no)
			if err != nil {
				fmt.Printf("download chapter %d fail,err %s\n", no, err.Error())
			} else {
				fmt.Printf("download chapter %d success\n", no)
			}
		} else {
			log.Fatal(console.Red("book without this chapter"))
		}
	} else {
		downloader.Download()
	}
}
