package inzone

//
//import (
//	"bytes"
//	"errors"
//	"fmt"
//	"io/ioutil"
//	"net/http"
//
//	"github.com/huaixiaohai/gapiservice/pb"
//
//	"github.com/PuerkitoBio/goquery"
//)
//
//// GetDailyLuckUsers 获取每日获奖名单那
//func GetDailyLuckUsers(cookie string) ([]*pb.LuckUser, error) {
//	buf, err := query("http://wx0.yinzuo.cn/index.php/MaoTCT/luckshow.html", cookie)
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//	}
//	var doc *goquery.Document
//	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//	}
//	luckUsers := make([]*User, 0)
//	doc.Find(".luckyLists").Find(".theLists").Find(".uf").Each(func(i int, s *goquery.Selection) {
//		//luckUsers = append(luckUsers, &User{
//		//	Name:  s.Children().First().Text(),
//		//	Phone: s.Children().Last().Text(),
//		//})
//		luckUsers = append(luckUsers, NewUser(s.Children().First().Text(), s.Children().Last().Text(), ""))
//
//	})
//	//for _, user := range luckUsers {
//	//	println(user.ID, user.Name, user.Phone)
//	//}
//	if len(luckUsers) <= 0 {
//		return "", nil, errors.New("名单还未公布")
//	}
//
//	buf, err = a.query("http://wx0.yinzuo.cn/index.php/MaoTCT/indexnew.html")
//	if err != nil {
//		fmt.Println(err.Error())
//		return "", luckUsers, err
//	}
//	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
//	if err != nil {
//		fmt.Println(err.Error())
//		return "", luckUsers, err
//	}
//
//	prizeName := doc.Find(".myorder.quanCenter").Find(".goodsCon").Find(".gdName").Text()
//
//	return prizeName, luckUsers, nil
//}
//
//// GetSeriesLuckUsers 获取每日获奖名单那
//func GetSeriesLuckUsers(cookie string) ([]*pb.InzoneUser, error) {
//	buf, err := query("http://wx0.yinzuo.cn/index.php/MaoTCTx/luckshow.html", cookie)
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//	}
//	var doc *goquery.Document
//	doc, err = goquery.NewDocumentFromReader(bytes.NewBuffer(buf))
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//	}
//	records := make([]*Record, 0)
//	println(len(doc.Find(".hasLuckys").Children().Nodes))
//	doc.Find(".hasLuckys").Children().Each(func(i int, s *goquery.Selection) {
//		//goods := s.Find(".goodsNameRight").Text()
//		users := make([]*User, 0)
//		s.Find(".luckyLists").Find(".theLists").Find(".uf").Each(func(i int, s *goquery.Selection) {
//			users = append(users, NewUser(s.Children().First().Text(), s.Children().Last().Text(), ""))
//		})
//		goods := fmt.Sprintf("系列%d", i+1)
//		records = append(records, &Record{
//			Goods: goods,
//			Users: users,
//		})
//	})
//
//	if len(records) <= 0 {
//		return nil, errors.New("系列名单还未公布")
//	}
//	for _, v := range records {
//		if len(v.Users) <= 0 {
//			return nil, errors.New("系列名单还未公布")
//		}
//	}
//
//	return records, nil
//}
//
//func query(url string, cookie string) ([]byte, error) {
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//	}
//	req.Header.Set("host", "wx0.yinzuo.cn")
//	req.Header.Set("Upgrade-Insecure-Requests", "1")
//	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x6309001c) XWEB/6500")
//	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
//	req.Header.Set("Referer", "http://wx0.yinzuo.cn/index.php/MaoTCT/indexnew.html")
//	//req.Header.Set("Accept-Encoding", "gzip, deflate")
//	req.Header.Set("Accept-Language", "zh-CN,zh")
//	req.Header.Set("Connection", "keep-alive")
//	req.Header.Set("Cookie", cookie)
//	client := &http.Client{}
//	var resp *http.Response
//	resp, err = client.Do(req)
//	if err != nil {
//		fmt.Println(err.Error())
//		return nil, err
//	}
//	defer func() {
//		_ = resp.Body.Close()
//	}()
//	var buf []byte
//	buf, err = ioutil.ReadAll(resp.Body)
//	if err != nil {
//		return nil, err
//	}
//	return buf, nil
//}
