package book

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"
	"ting55-downloader/pkg/console"
	"ting55-downloader/pkg/ua"
)

var ting55Uri = "https://ting55.com/book/"

type Book struct {
	Id         int
	Title      string
	Number     int
	Cover      string
	Category   string
	Author     string
	Announcer  string
	Status     string
	CreateTime string
	IsMobile   bool
}

func Parse(bookId int) *Book {
	url := fmt.Sprintf("%s%d", ting55Uri, bookId)
	log.Println(console.Green(fmt.Sprintf("Start Parsing Book From %s", url)))
	uaNew := ua.New()
	UA, err := uaNew.Random()
	if err != nil {
		log.Fatal(console.Red(fmt.Sprintf("Random Gen UA Fail,%s", err.Error())))
	}
	fmt.Println(UA)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(console.Red(fmt.Sprintf("http.NewRequest Fail,%s", err.Error())))
	}
	req.Header.Set("User-Agent", UA)
	res, err := (&http.Client{}).Do(req)
	if err != nil {
		log.Fatal(console.Red(fmt.Sprintf("Parsing %s Fail,%s", url, err.Error())))
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(console.Red(fmt.Sprintf("Read Body Fail,%s", err.Error())))
	}
	html := string(body)
	isMobile := !strings.Contains(html, "手机恋听网")
	book := &Book{
		Id:       bookId,
		IsMobile: isMobile,
	}
	book.Init(html)
	return book
}

func (b *Book) Init(html string) {
	var re *regexp.Regexp
	if b.IsMobile {
		re = regexp.MustCompile(`class="bookinfo".*?class="bimg".*?src="(.*?)".*?alt.*?class="binfo".*?<h1>(.*?)</h1>.*?<p>类型：(.*?)</p>.*?<p>作者：(.*?)</p>.*?<p>播音：.*?<a.*?>(.*?)</a>.*?</p>.*?<p>时间：(.*?)</p>.*?<p>状态：(.*?)</p>.*?class="intro".*?class="playlist".*?class="plist">(.*?)</div>`)
	} else {
		re = regexp.MustCompile(`class="bookinfo".*?class="bimg".*?src="(.*?)".*?alt.*?class="binfo".*?<h1>(.*?)</h1>.*?<p>类别：(.*?)</p>.*?<p>作者：(.*?)</p>.*?<p>播音：.*?<a.*?>(.*?)</a>.*?</p>.*?<p>状态：(.*?)</p>.*?<p>时间：(.*?)</p>.*?class="intro".*?class="playlist".*?<ul>(.*?)</ul>`)
	}
	match := re.FindAllStringSubmatch(html, -1)
	if match == nil {
		log.Fatal(console.Red("No book information was matched"))
	}
	b.Cover = match[0][1]
	b.Title = match[0][2]
	b.Category = match[0][3]
	b.Author = match[0][4]
	b.Announcer = match[0][5]
	if b.IsMobile {
		b.CreateTime = match[0][6]
		b.Status = match[0][7]
	} else {
		b.Status = match[0][6]
		b.Category = match[0][7]
	}
	fmt.Println(match[0][8])
}
