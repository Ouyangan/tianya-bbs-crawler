package main

import (
	"net/http"
	"log"
	"github.com/PuerkitoBio/goquery"
	"strings"
	"strconv"
	"os"
	"time"
)

func main() {
	fi, err := os.Create("/Users/ouyangan/go/src/tianya_spider/content.txt")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()
	fi.WriteString(time.Now().String())
	fi.WriteString("\r\n")
	for i := 109; i <= 109; i++ {
		parsing(i, fi)
		time.Sleep(1 * time.Second)
	}
}

func parsing(pageNumber int, file *os.File) {
	var text strings.Builder
	var url = "http://bbs.tianya.cn/post-house-252774-" + strconv.Itoa(pageNumber) + ".shtml"
	text.WriteString("=============================================")
	text.WriteString("\r\n")
	text.WriteString("\r\n")
	text.WriteString("第")
	text.WriteString(strconv.Itoa(pageNumber))
	text.WriteString("页 链接:")
	text.WriteString(url)
	text.WriteString("\r\n")
	text.WriteString("\r\n")
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	log.Println(res.Body)
	doc, err := goquery.NewDocumentFromReader(res.Body)
	getTotalPage(doc)
	if err != nil {
		log.Fatalf(err.Error())
	}
	doc.Find(".atl-item").Each(func(i int, selection *goquery.Selection) {
		parsingFloor(selection, &text)
		parsingAuthorAndReplayTime(selection, &text)
		parsingContent(selection, &text)
	})
	log.Println(text.String())

	file.WriteString(text.String())
}

//获取总页数
func getTotalPage(doc *goquery.Document) int {
	selection := doc.Find(".atl-pages").Find("a")
	total := selection.Nodes
	for i, n := range total {
		log.Println(i)
		log.Println(n.Data)
	}
	return 0
}

func creatHeader() {

}

//解析回帖作者,时间
func parsingAuthorAndReplayTime(selection *goquery.Selection, text *strings.Builder) {
	selection.Find(".atl-head").Find(".atl-info").Find("span").Each(func(i int, selection *goquery.Selection) {
		if i == 0 {
			str, _ := selection.Find("a").Attr("uname")
			if str == "kkndme" {
				str = "楼主:" + str
			} else {
				str = "作者:" + str
			}
			text.WriteString(str)
			text.WriteString(" ")
		} else if i == 1 {
			str, _ := selection.Html()
			text.WriteString(str)
			text.WriteString("\n")
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
	text.WriteString("\r\n")
	text.WriteString("\r\n")
}
