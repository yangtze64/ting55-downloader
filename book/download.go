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
	"strconv"
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

func NewDownloader(num int, mode int, dpath string) *downloader {
	if exist := FileIsExist(dpath); !exist {
		if err := os.MkdirAll(dpath, 0775); err != nil {
			log.Fatal(console.Red("Config file create Error"))
		}
	}
	return &downloader{
		threadNum:    num,
		downloadMode: mode,
		downloadPath: dpath,
	}
}

func (d *downloader) Download(book *Book) {
	log.Println(console.Green(fmt.Sprintf("Start Download Book,BookId:%d BookName:%s", book.Id, book.Title)))
	pool := pb.NewPool()
	pool.Start()
	coverCh := make(chan bool, 1)
	if book.Cover != "" {
		go func() {
			err := d.downloadCover(book, pool)
			if err != nil {
				coverCh <- false
				fmt.Println(err)
			} else {
				coverCh <- true
			}
		}()
	}
	go d.saveBookInfoToFile(book)
	for i := 0; i < d.threadNum; i++ {

	}
	<-coverCh
	pool.Stop()
}

func (d *downloader) downloadCover(book *Book, pool *pb.Pool) error {
	res, err := http.Get(book.Cover)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return err
	}
	length := res.Header.Get("Content-Length")
	size, _ := strconv.ParseInt(length, 10, 64)
	bar := GetProgressBar("封面图", size)
	pool.Add(bar)
	defer bar.Finish()
	body := res.Body
	reader := bar.NewProxyReader(body)
	coverPath := path.Join(d.downloadPath, "封面图.png")
	file, err := os.Create(coverPath)
	if err != nil {
		return err
	}
	writer := io.Writer(file)
	_, err = io.Copy(writer, reader)
	if err != nil {
		return err
	}
	return nil
}

func (d *downloader) saveBookInfoToFile(book *Book) error {
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
