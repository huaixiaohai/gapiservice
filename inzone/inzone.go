package inzone

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/huaixiaohai/lib/log"

	"github.com/PuerkitoBio/goquery"
)

//func GetLuckUsers(cookie string) ([][]string, error) {
//	dailyLuck, err := getDailyLuckUsers(cookie)
//	if err != nil {
//		return nil, err
//	}
//	seriesLuck, err := getSeriesLuckUsers(cookie)
//	if err != nil {
//		return nil, err
//	}
//	res := make([][]string, 0)
//	res = append(res, dailyLuck)
//	res = append(res, seriesLuck...)
//	return res, nil
//}

var InvalidInzoneCookieErr = errors.New("无效的inzone cookie")

type Luck struct {
	Label string
	UUIDs []string
}

func IsValid(cookie string) bool {
	buf, err := GetIndex(cookie)
	if err != nil {
		log.Error("cookie失效", err.Error())
	}
	fmt.Println("IsValid buf", string(buf))
	return err == nil
}

func GetUserName(cookie string) (string, error) {
	buf, err := get("https://wx0.yinzuo.cn/index.php/FinishInfo/finish_info.html", cookie)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	context := doc.Find(".core").Find(".fz_16").Text()
	context = strings.ReplaceAll(context, "\t", "")
	context = strings.ReplaceAll(context, "\n", "")
	context = context[0:9]
	ss := strings.Split(context, " ")
	return ss[0], nil
}

func GetPhone(cookie string) (string, error) {
	buf, err := get("https://wx0.yinzuo.cn/index.php/Index/qdindex.html", cookie)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	context := doc.Find(".gr_hy").Find(".gr_ri_txt").Text()
	context = strings.ReplaceAll(context, "\t", "")
	context = strings.ReplaceAll(context, "\n", "")
	return context, nil
}

func GetCID(cookie string) (string, error) {
	buf, err := get("https://wx0.yinzuo.cn/index.php/MaoTCT/indexnew.html", cookie)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	context := doc.Find(".uf.uf-ac.bottomData").Find(".uf.uf-f1.uf-col.uf-as.leftCon").Find(".uf.uf-pc.uf-f1.bottomDataLine").Text()
	fmt.Println("buf...", string(buf))
	fmt.Println("context", context)
	if !strings.Contains(context, "CID: ") {
		return "", errors.New("CID 查找不到")
	}
	cid := strings.ReplaceAll(context, "CID:", "")
	cid = strings.ReplaceAll(cid, " ", "")

	return cid, nil
}

func IsLuck(cookie string) (bool, error) {
	buf, err := get("https://wx0.yinzuo.cn/index.php/MaoTCT/mycounpon.html", cookie)
	if err != nil {
		fmt.Println(err.Error())
		return false, err
	}
	//var doc *goquery.Document
	//doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
	//if err != nil {
	//	fmt.Println(err.Error())
	//	return false, err
	//}
	return !strings.Contains(string(buf), "暂无资格"), nil
}

func GetIndex(cookie string) ([]byte, error) {
	buf, err := get("https://wx0.yinzuo.cn/index.php/MaoTCT/Index.html", cookie)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return buf, err
}

// GetDailyLuckUsers 获取每日获奖名单那
func GetDailyLuckUsers(cookie string) ([]*Luck, error) {
	buf, err := get("https://wx0.yinzuo.cn/index.php/MaoTCT/luckshow.html", cookie)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	uuids := make([]string, 0)
	doc.Find(".luckyLists").Find(".theLists").Find(".uf").Each(func(i int, s *goquery.Selection) {
		uuid := nameFormat(s.Children().First().Text()) + strings.ReplaceAll(s.Children().Last().Text(), " ", "")
		uuids = append(uuids, uuid)
	})
	if len(uuids) <= 0 {
		return nil, errors.New("名单还未公布")
	}

	res := []*Luck{{
		Label: "每日",
		UUIDs: uuids,
	}}
	return res, nil
}

// GetSeriesLuckUsers 获取每日获奖名单那
func GetSeriesLuckUsers(cookie string) ([]*Luck, error) {
	buf, err := get("https://wx0.yinzuo.cn/index.php/MaoTCTx/luckshow.html", cookie)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	var doc *goquery.Document
	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	res := make([]*Luck, 0)
	doc.Find(".hasLuckys").Children().Each(func(i int, s *goquery.Selection) {
		luck := &Luck{
			Label: fmt.Sprintf("系列%d", i),
			UUIDs: make([]string, 0),
		}
		s.Find(".luckyLists").Find(".theLists").Find(".uf").Each(func(i int, s *goquery.Selection) {
			uuid := nameFormat(s.Children().First().Text()) + strings.ReplaceAll(s.Children().Last().Text(), " ", "")
			luck.UUIDs = append(luck.UUIDs, uuid)
		})
		res = append(res, luck)
	})

	if len(res) <= 0 {
		return nil, errors.New("系列名单还未公布")
	}
	return res, nil
}

// GetAppY 抽签报名
func GetAppY(cookie string) error {
	_, err := post("https://wx0.yinzuo.cn/index.php/Home/MaoTCT/getappy", cookie)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	return err
}

func query(method, url string, cookie string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	req.Header.Set("host", "wx0.yinzuo.cn")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x6309001c) XWEB/6500")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//req.Header.Set("Referer", "http://wx0.yinzuo.cn/index.php/MaoTCT/indexnew.html")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", cookie)
	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if strings.Contains(resp.Request.URL.String(), "open.weixin.qq.com") {
		return nil, InvalidInzoneCookieErr
	}

	var buf []byte
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func get(url string, cookie string) ([]byte, error) {
	return query("GET", url, cookie)
}
func post(url string, cookie string) ([]byte, error) {
	return query("POST", url, cookie)
}

func nameFormat(str string) string {
	str = strings.ReplaceAll(str, " ", "")
	str = strings.ReplaceAll(str, " ", "")
	return str
}
