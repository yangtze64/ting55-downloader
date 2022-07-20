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
	mode         = flag.Int("m", 1, "download mod optional 1:Overwrite,2:Skip,3:Keep Original")
	maxThreadNum = runtime.NumCPU() * 2
	modeMap      = map[int]string{
		1: "Overwrite",
		2: "Skip",
		3: "Keep Original",
	}
)

func main() {
	flag.Parse()
	bookId := *id
	downloadPath := *dp
	threadNum := *num
	downloadMode := *mode
	if bookId == 0 {
		log.Fatal(console.Red("请设置bookId"))
	}
	if threadNum > maxThreadNum {
		log.Fatal(console.Red("thread num 不可超过CUP核心数"))
	}
	downloadModeVal, ok := modeMap[downloadMode]
	if !ok {
		log.Fatal(console.Red("download mod 1:Overwrite,2:Skip,3:Keep Original"))
	}
	log.Println(console.Yellow(fmt.Sprintf("bookId:%d", bookId)))
	log.Println(console.Yellow(fmt.Sprintf("threadNum:%d", threadNum)))
	log.Println(console.Yellow(fmt.Sprintf("downloadMode:%s", downloadModeVal)))
	log.Println(console.Yellow(fmt.Sprintf("downloadPath:%s", downloadPath)))
	bookInfo := book.Parse(bookId)
	fmt.Println(bookInfo)
}
