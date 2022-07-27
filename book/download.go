package book

import (
	"bufio"
	"fmt"
	"github.com/cheggaaa/pb/v3"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"sync"
	"time"
	"ting55-downloader/pkg/console"
)

type downloader struct {
	book          *Book
	threadNum     int
	downloadMode  int
	downloadPath  string
	jobFinishChan chan bool
	mutex         sync.Mutex
}

var ModeMap = map[int]string{
	1: "Overwrite",
	2: "Skip",
	3: "Breakpoint Recovery",
	4: "Keep Original",
}

func GetProgressBarTemplate(name string) pb.ProgressBarTemplate {
	s := `{{ "%s:" }} {{ bar . "[" "#" (cycle . "=>" "=>" "=>" "=>" ) "." "]" | green}} {{percent . | green}} {{speed . | green }}`
	tmpl := fmt.Sprintf(s, name)
	return pb.ProgressBarTemplate(tmpl)
}

func GetProgressBar(name string, size int64) *pb.ProgressBar {
	pbt := GetProgressBarTemplate(name)
	bar := pb.New64(size).SetTemplate(pbt)
	bar.SetRefreshRate(time.Millisecond * 100)
	return bar
}

// FileIsExist 文件或文件夹是否存在
func FileIsExist(file string) bool {
	_, err := os.Stat(file)
	return err == nil || os.IsExist(err)
}

func NewDownloader(book *Book, num int, mode int, dpath string) *downloader {
	dpath = path.Join(dpath, book.Title)
	if exist := FileIsExist(dpath); !exist {
		if err := os.MkdirAll(dpath, 0775); err != nil {
			log.Fatal(console.Red("Config file create Error"))
		}
	}
	return &downloader{
		book:          book,
		threadNum:     num,
		downloadMode:  mode,
		downloadPath:  dpath,
		jobFinishChan: make(chan bool, 1),
	}
}

func (d *downloader) Download() {
	book := d.book
	log.Println(console.Green(fmt.Sprintf("Start Download Book,BookId:%d BookName:%s", book.Id, book.Title)))
	// 处理封面
	if book.Cover != "" {
		go d.downloadCover()
	}
	// 保存基本信息到文件
	go d.saveBookInfoToFile()

	//pool := pb.NewPool()
	//pool.Start()
	//defer func() {
	//	pool.Stop()
	//}()

	// 投递任务
	n := d.threadNum
	if d.threadNum > book.Number {
		n = book.Number
	}
	jobCh := make(chan int, n)
	statusCh := make(chan int, n)
	go d.deliver(jobCh)
	go func() {
		for i := 0; i < n; i++ {
			no, ok := <-jobCh
			if ok {
				go d.downloadJob(no, statusCh)
			}
		}
		for i := 0; i < book.Number; i++ {
			<-statusCh
			// fmt.Printf("statusCh no:%d\n", s)
			b, ok := <-jobCh
			if ok {
				go d.downloadJob(b, statusCh)
			} else {
				// fmt.Println("jobCh Closed")
			}
		}
		d.jobFinishChan <- true
	}()
	<-d.jobFinishChan
}

// 下载封面
func (d *downloader) downloadCover() error {
	book := d.book
	res, err := http.Get(book.Cover)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return err
	}
	body := res.Body
	coverPath := path.Join(d.downloadPath, "封面图.png")
	file, err := os.Create(coverPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}
	return nil
}

//func (d *downloader) downloadCover() error {
//	book := d.book
//	res, err := http.Get(book.Cover)
//	if err != nil {
//		return err
//	}
//	if res.StatusCode != http.StatusOK {
//		return err
//	}
//	length := res.Header.Get("Content-Length")
//	size, _ := strconv.ParseInt(length, 10, 64)
//	d.mutex.Lock()
//	bar := GetProgressBar("封面图", size)
//	pool.Add(bar)
//	d.mutex.Unlock()
//	defer bar.Finish()
//	body := res.Body
//	reader := bar.NewProxyReader(body)
//	coverPath := path.Join(d.downloadPath, "封面图.png")
//	file, err := os.Create(coverPath)
//	if err != nil {
//		return err
//	}
//	writer := io.Writer(file)
//	_, err = io.Copy(writer, reader)
//	if err != nil {
//		return err
//	}
//	return nil
//}

// 保存书的基本信息
func (d *downloader) saveBookInfoToFile() error {
	book := d.book
	file, err := os.OpenFile(path.Join(d.downloadPath, "book.txt"), os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0775)
	if err != nil {
		return err
	}
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(fmt.Sprintf("%s：%s\n", "名称", book.Title))
	write.WriteString(fmt.Sprintf("%s：%d\n", "节数", book.Number))
	write.WriteString(fmt.Sprintf("%s：%s\n", "类别", book.Category))
	write.WriteString(fmt.Sprintf("%s：%s\n", "作者", book.Author))
	write.WriteString(fmt.Sprintf("%s：%s\n", "播音", book.Announcer))
	write.WriteString(fmt.Sprintf("%s：%s\n", "状态", book.Status))
	write.WriteString(fmt.Sprintf("%s：%s\n", "时间", book.CreateTime))
	write.Flush()
	return nil
}

func (d *downloader) deliver(jobCh chan int) {
	book := d.book
	// 投递任务
	go func() {
		for i := 1; i <= book.Number; i++ {
			fmt.Sprintf("%d\n", i)
			jobCh <- i
		}
		close(jobCh)
	}()
}

func (d *downloader) downloadJob(no int, statusCh chan int) error {
	fmt.Printf("download chapter:%d\n", no)
	time.Sleep(time.Millisecond * 100)
	err := d.DownloadAudio(no)
	fmt.Println(err)
	statusCh <- no
	return nil
}

func (d *downloader) DownloadAudio(no int) error {
	book := d.book
	chapter, err := book.GetChapter(no)
	if err != nil {
		return err
	}
	url, err := chapter.GetChapterAudioUrl()
	if err != nil {
		return err
	}
	ext := path.Ext(url) //获取文件后缀

	res, err := http.Get(url)
	defer res.Body.Close()
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return err
	}
	body := res.Body
	audioPath := path.Join(d.downloadPath, fmt.Sprintf("%s-第%d章%s", book.Title, no, ext))
	file, err := os.Create(audioPath)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, body)
	if err != nil {
		return err
	}
	return nil
}
