package ting13

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"github.com/huaixiaohai/lib/log"

	"github.com/huaixiaohai/gapiservice/pb"
)

func GetNovelResources(ctx context.Context, id string) ([]*pb.NovelResource, error) {

	var pageIndex int
	for {
		pageIndex++
		urlStr := fmt.Sprintf("https://www.ting13.com/youshengxiaoshuo/15687/%d.html", pageIndex)
		resp, err := http.Get(urlStr)
		if err != nil {
			log.Error(err)
			return nil, err
		}
		buf, err := ioutil.ReadAll(resp.Body)
		var doc *goquery.Document
		doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
		if err != nil {
			log.Error(err)
			return nil, err
		}
		context := doc.Find("#playlist").Each().Text()
	}
}
