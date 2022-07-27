package book

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"ting55-downloader/pkg/console"
	"ting55-downloader/pkg/request"
	"ting55-downloader/pkg/ua"
)

var (
	Host              = "ting55.com"
	MobileHost        = "m.ting55.com"
	Protocol          = "https://"
	BookUri           = "/book/%d"
	ChapterUri        = "/book/%d-%d"
	AudioReqUri       = "/nlinka"
	MobileAudioReqUri = "/glink"
)

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

type Chapter struct {
	BookId   int
	Page     int
	IsPay    int
	UA       string
	L        string
	XT       string
	Host     string
	Origin   string
	Referer  string
	IP       string
	IsMobile bool
}

func Parse(bookId int) *Book {
	url := fmt.Sprintf("%s%s%s", Protocol, Host, fmt.Sprintf(BookUri, bookId))
	log.Println(console.Green(fmt.Sprintf("Start Parsing Book From %s", url)))
	uaNew := ua.New()
	UA, err := uaNew.Random()
	if err != nil {
		log.Fatal(console.Red(fmt.Sprintf("Random Gen UA Fail,%s", err.Error())))
	}
	fmt.Printf("UA:%s\n", UA)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(console.Red(fmt.Sprintf("http.NewRequest Fail,%s", err.Error())))
	}
	req.Header.Set("User-Agent", UA)
	res, err := (&http.Client{}).Do(req)
	defer res.Body.Close()
	if err != nil {
		log.Fatal(console.Red(fmt.Sprintf("Parsing %s Fail,%s", url, err.Error())))
	}
	if res.StatusCode != http.StatusOK {
		log.Fatal(console.Red(fmt.Sprintf("res.StatusCode Is %d Not Is %d", res.StatusCode, http.StatusOK)))
	}
	body, err := ioutil.ReadAll(res.Body)
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
	cover := match[0][1]
	if !strings.HasPrefix(cover, "http:") && !strings.HasPrefix(cover, "https:") {
		cover = "http:" + cover
	}
	b.Cover = cover
	title := match[0][2]
	if strings.HasSuffix(title, "有声小说") {
		title = strings.TrimSuffix(title, "有声小说")
	}
	b.Title = title
	b.Category = match[0][3]
	b.Author = match[0][4]
	b.Announcer = match[0][5]
	if b.IsMobile {
		b.CreateTime = match[0][6]
		b.Status = match[0][7]
	} else {
		b.Status = match[0][6]
		b.CreateTime = match[0][7]
	}
	re = regexp.MustCompile(`<a.*?>(.*?)</a>`)
	l := re.FindAllStringSubmatch(match[0][8], -1)
	if l == nil {
		log.Fatal(console.Red("Set number not matched, Please try again"))
	}
	b.Number = len(l)
}

func (b *Book) GetChapter(no int) (*Chapter, error) {
	url := fmt.Sprintf("%s%s%s", Protocol, Host, fmt.Sprintf(ChapterUri, b.Id, no))
	uaNew := ua.New()
	// uaNew.Use(ua.Chrome, ua.Firefox, ua.Safari)
	UA, _ := uaNew.Random()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", UA)
	ip := request.GenIpaddr()
	req.Header.Set("X-Forwarded-For", ip)
	res, err := (&http.Client{}).Do(req)
	defer res.Body.Close()
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	html := string(body)
	if html == "" {
		return nil, errors.New("book chapter body is empty")
	}
	isMobile := !strings.Contains(html, "手机恋听网")
	var re *regexp.Regexp
	if isMobile {
		re = regexp.MustCompile(`<meta name="_c" content="(.*?)"`)
	} else {
		re = regexp.MustCompile(`<meta name="_c" content="(.*?)".*?<meta name="_l" content="(.*?)"`)
	}
	match := re.FindAllStringSubmatch(html, -1)
	if match == nil {
		return nil, errors.New("no book chapter information was matched")
	}
	xt := match[0][1]
	l := "1"
	if !isMobile {
		l = match[0][2]
	}
	host := Host
	if isMobile {
		host = MobileHost
	}
	origin := fmt.Sprintf("%s%s", Protocol, host)
	referer := fmt.Sprintf("%s%s%s", Protocol, host, fmt.Sprintf(ChapterUri, b.Id, no))

	chapter := &Chapter{
		BookId:   b.Id,
		Page:     no,
		IsPay:    0,
		UA:       UA,
		IsMobile: isMobile,
		L:        l,
		XT:       xt,
		Host:     host,
		Origin:   origin,
		Referer:  referer,
		IP:       ip,
	}
	return chapter, nil
}

func (c *Chapter) GetChapterAudioUrl() (string, error) {
	var url string
	if c.IsMobile {
		url = fmt.Sprintf("%s%s%s", Protocol, MobileHost, MobileAudioReqUri)
	} else {
		url = fmt.Sprintf("%s%s%s", Protocol, Host, AudioReqUri)
	}
	data := map[string]int{
		"bookId": c.BookId,
		"isPay":  c.IsPay,
		"page":   c.Page,
	}
	//bytesData, _ := json.Marshal(data)
	var str string
	for k, v := range data {
		str += fmt.Sprintf("%s=%s&", k, strconv.Itoa(v))
	}
	str = strings.TrimSuffix(str, "&")
	contentLength := len(str)
	req, err := http.NewRequest("POST", url, strings.NewReader(str))
	if err != nil {
		return "", err
	}

	req.Header.Set("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(contentLength))
	req.Header.Set("User-Agent", c.UA)
	req.Header.Set("Host", c.Host)
	req.Header.Set("Origin", c.Origin)
	req.Header.Set("Referer", c.Referer)
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("X-Forwarded-For", c.IP)
	req.Header.Set("xt", c.XT)
	if !c.IsMobile {
		req.Header.Set("l", c.L)
	}
	res, err := (&http.Client{}).Do(req)
	defer res.Body.Close()
	if err != nil {
		return "", err
	}
	if res.StatusCode != http.StatusOK {
		return "", err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	resultStr := string(body)
	if resultStr == "" {
		return "", errors.New("request chapter audio result is empty")
	}
	result := make(map[string]interface{})
	err = json.Unmarshal([]byte(resultStr), &result)
	if err != nil {
		return "", err
	}
	status, ok := result["status"]
	if !ok {
		return "", errors.New("request chapter audio result no `status` is short")
	}
	resUrl, ok := result["url"]
	if !ok {
		return "", errors.New("request chapter audio result no `url` is short")
	}
	if status.(float64) != 1 || resUrl.(string) == "" {
		return "", errors.New("request chapter audio result `status` is not 1 or `url` is empty")
	}
	return result["url"].(string), nil
}
