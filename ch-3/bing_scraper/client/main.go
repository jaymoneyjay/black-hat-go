package main

import (
	"archive/zip"
	"bytes"
	"github.com/jaymoneyjay/black-hat-go/ch-3/bing_scraper/scraper"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
)

func main() {
	if len(os.Args) != 1 {
		log.Fatalln("Missing required argument. Usage: main.go domain ext")
	}
	domain := "learningculture.ch"//os.Args[1]
	filetype := "pdf"//os.Args[2]

	q := fmt.Sprintf("site:%s && filetype:%s && instreamset:(url title):%s",
		domain,
		filetype,
		filetype,
		)
	search := fmt.Sprintf("http://www.bing.com/search?q=%s", url.QueryEscape(q))
	res, err := http.Get(search)
	if err != nil {
		log.Fatalln(err)
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	defer res.Body.Close()

	//s := "html body.b_respl.b_sbText div#b_content main ol#b_results li.b_algo h2"
	debugS := "html body div#b_content main ol#b_results"
	doc.Find(debugS).Each(debugHandler)

}

func debugHandler(i int, q *goquery.Selection) {
	html, err := q.Html()
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Printf("%d: %s", i, html)
}

func handler(i int, q *goquery.Selection) {
	href, ok := q.Find("a").Attr("href")
	if !ok {
		log.Panicln()
	}
	res, err := http.Get(href)
	if err != nil {
		log.Println(err)
		return
	}

	buf, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
		return
	}

	defer res.Body.Close()

	zr, err := zip.NewReader(bytes.NewReader(buf), int64(len(buf)))
	if err != nil {
		log.Println(err)
		return
	}
	appProps, coreProps, err := scraper.NewProperties(zr)
	if err != nil {
		log.Println(err)
		return
	}

	fmt.Printf("%d: %s", i, href)
	fmt.Printf("%25s %25s - %s %s",
		coreProps.Creator,
		coreProps.LastModifiedBy,
		appProps.Application,
		appProps.GetMajorVersion(),
	)

	log.Println("Success")
	return
}