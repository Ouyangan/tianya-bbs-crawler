package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/sirupsen/logrus"
)

var log = logrus.StandardLogger()

func init() {
	initLog()
}
func main() {
	run()
}

func run() {
	var fileName string
	var baseUrl string
	var duration int32
	fmt.Println("作者:ananan")
	fmt.Println("代码地址:github.com/Ouyangan/tianya-bbs-crawler")
	fmt.Println("联系方式:981017952@qq.com")
	fmt.Println("===========使用介绍===========")
	fmt.Println("请依次输入四个参数:")
	fmt.Println("1.生成文件名称")
	fmt.Println("2.帖子地址")
	fmt.Println("3.每隔爬取间隔时间(毫秒) 推荐1000毫秒以上,防止被封")
	fmt.Println("输入完成后将会在同级目录下生成txt文件,请等待任务完成")
	fmt.Println()
	fmt.Println("参考示例:")
	fmt.Println()
	fmt.Println("天涯")
	fmt.Println("http://bbs.tianya.cn/post-house-252774-1.shtml")
	fmt.Println("1000")
	fmt.Println("=============================")
	fmt.Println()
	fmt.Println("请输入生成文件名称并回车:")
	fmt.Scanln(&fileName)
	fmt.Println("请输入帖子地址并回车:")
	fmt.Scanln(&baseUrl)
	fmt.Println("每页爬取间隔时间并回车")
	fmt.Scanln(&duration)
	start(fileName, baseUrl, duration)
}

func start(fileName string, url string, duration int32) {
	fi := createFile(fileName)
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	baseUrl, totalPage := parsingTopicInfo(url)
	log.Infoln("总计:", totalPage, "页")
	for i := 1; i <= totalPage; i++ {
		log.Println("开始解析:第", i, "页")
		parsing(baseUrl, i, fi)
		time.Sleep(time.Duration(duration) * time.Millisecond)
		log.Println("完成解析:第", i, "页")
	}
}

//解析网页
func parsing(baseUrl string, pageNumber int, file *os.File) {
	var text strings.Builder
	var url = baseUrl + strconv.Itoa(pageNumber) + ".shtml"
	createPageHeader(pageNumber, url, &text)
	doc := httpGetUtil(url)
	doc.Find(".atl-item").Each(func(i int, selection *goquery.Selection) {
		parsingFloor(selection, &text)
		parsingAuthorAndReplayTime(selection, &text)
		parsingContent(selection, &text)
	})
	file.WriteString(text.String())
}

//解析url前缀,总页码
func parsingTopicInfo(url string) (string, int) {
	doc := httpGetUtil(url)
	var totalPage string
	var baseUrl string
	doc.Find(".atl-pages").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			selection.Find("a")
			size := selection.Find("a").Size()
			selection.Find("a").Each(func(i int, selection *goquery.Selection) {
				if i == size-2 {
					totalPage = selection.Text()
					baseUrl, _ = selection.Attr("href")
				}
			})
		}
	})
	baseUrl = strings.Replace(baseUrl, ".shtml", "", -1)
	baseUrl = strings.Replace(baseUrl, totalPage, "", -1)
	baseUrl = "http://bbs.tianya.cn" + baseUrl
	count, _ := strconv.Atoi(totalPage)
	return baseUrl, count
}

func createPageHeader(pageNumber int, url string, text *strings.Builder) {
	text.WriteString("\r\n\r\n")
	text.WriteString("第")
	text.WriteString(strconv.Itoa(pageNumber))
	text.WriteString("页 链接:")
	text.WriteString(url)
	text.WriteString("\r\n\r\n")
}

//解析回帖作者,时间
func parsingAuthorAndReplayTime(selection *goquery.Selection, text *strings.Builder) {
	selection.Find(".atl-head").Find(".atl-info").Find("span").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			str, _ := selection.Find("a").Attr("uname")
			str = "作者:" + str
			text.WriteString(str)
			text.WriteString(" ")
		} else if i == 1 {
			str, _ := selection.Html()
			text.WriteString(str)
			text.WriteString("\r\n")
		}
	})
}

//解析楼层
func parsingFloor(selection *goquery.Selection, text *strings.Builder) {
	selection.Find(".atl-reply").Find("span").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			text.WriteString(selection.Text())
			text.WriteString(" ")
		}
	})
}

//解析回帖内容
func parsingContent(selection *goquery.Selection, text *strings.Builder) {
	replacer := strings.NewReplacer("\n", "", "　", "", "	", "")
	content := selection.Find(".atl-con-bd").Find(".bbs-content").Text()
	contentStr := replacer.Replace(content)
	text.WriteString(contentStr)
	text.WriteString("\r\n\r\n")
}

//http请求
func httpGetUtil(url string) (*goquery.Document) {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal("请求地址发生错误:", err, ",请检查地址是否正确")
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("请求地址发生错误:%d %s", res.StatusCode, res.Status, "请检查地址是否正确")
	}
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return doc
}

func createFile(fileName string) *os.File {
	fi, err := os.Create(fileName + ".txt")
	if err != nil {
		log.Fatal(err)
	}
	fi.WriteString("============================")
	line(fi, 1)
	fi.WriteString("创建时间:" + time.Now().Format("2006-01-02 15:04:05"))
	line(fi, 1)
	fi.WriteString("作者:ananan")
	line(fi, 1)
	fi.WriteString("代码地址:github.com/Ouyangan/tianya-bbs-crawler")
	line(fi, 1)
	fi.WriteString("联系方式:981017952@qq.com")
	line(fi, 1)
	fi.WriteString("============================")
	line(fi, 1)
	return fi
}

//换行
func line(fi *os.File, count int) {
	for i := 0; i < count; i ++ {
		fi.WriteString("\r\n")
	}
}

func initLog() {
	log.Level = logrus.DebugLevel
	log.Formatter = &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     false,
		FullTimestamp:   true,
	}
	file, err := os.Create("tianya-bbs-crawler.log")
	if err != nil {
		log.Warn("创建日志异常,", err.Error())
		log.Out = os.Stdout
	} else {
		mw := io.MultiWriter(file, os.Stdout)
		log.Out = mw
	}
}
