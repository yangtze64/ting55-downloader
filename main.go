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
	maxThreadNum = runtime.NumCPU() * 2
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
		log.Fatal(console.Red("thread num 不可超过CUP*2核心数"))
	}
	downloadModeVal, ok := book.ModeMap[downloadMode]
	if !ok {
		log.Fatal(console.Red("download mod 1:Overwrite,2:Skip,3:Keep Original"))
	}
	log.Println(console.Yellow(fmt.Sprintf("bookId:%d, threadNum:%d, downloadPath:%s, downloadMode:%s", bookId, threadNum, downloadPath, downloadModeVal)))
	bookInfo := book.Parse(bookId)
	fmt.Printf("%#+v\n", bookInfo)

	downloader := book.NewDownloader(threadNum, downloadMode, downloadPath)
	downloader.Download(bookInfo)
	//count := 100
	//count2 := 100

	// go func() {
	// 	bar1 := pb.StartNew(count2)
	// 	for i := 0; i < count2; i++ {
	// 		bar1.Increment()
	// 		time.Sleep(time.Millisecond)
	// 	}
	// 	bar1.Finish()
	// }()
	// go func() {
	// 	bar := pb.StartNew(count)
	// 	for i := 0; i < count; i++ {
	// 		bar.Increment()
	// 		time.Sleep(time.Millisecond)
	// 	}
	// 	bar.Finish()
	// }()
	//bar1 := pb.New(count)
	//bar2 := pb.New(count2)
	//p := pb.NewPool(bar1, bar2)
	//
	//p.Start()
	//ch := make(chan int, 2)
	//go func() {
	//	for i := 0; i < count; i++ {
	//		bar1.Increment()
	//		time.Sleep(time.Millisecond * 30)
	//	}
	//	ch <- 1
	//}()
	//go func() {
	//	for i := 0; i < count2; i++ {
	//		bar2.Increment()
	//		time.Sleep(time.Millisecond * 50)
	//	}
	//	ch <- 1
	//}()
	//<-ch
	//<-ch
	//p.Stop()
	// create and start new bar
	// bar := pb.StartNew(count)

	// start bar from 'default' template
	// bar := pb.Default.Start(count)

	// start bar from 'simple' template
	// bar := pb.Simple.Start(count)

	// start bar from 'full' template
	// bar := pb.Full.Start(count)

	// for i := 0; i < count; i++ {
	// 	bar.Increment()
	// 	time.Sleep(time.Millisecond)
	// }

	// finish bar
	// bar.Finish()
}
