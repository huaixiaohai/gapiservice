package api

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/huaixiaohai/lib/log"

	"github.com/gin-gonic/gin"

	"github.com/huaixiaohai/gapiservice/dao"

	"github.com/PuerkitoBio/goquery"
	"github.com/huaixiaohai/gapiservice/config"
	"github.com/huaixiaohai/gapiservice/pb"
	"github.com/xuri/excelize/v2"

	"github.com/google/wire"
	"github.com/robfig/cron/v3"
)

var InzoneApiSet = wire.NewSet(NewInzoneApi)

func NewInzoneApi(userRepo *dao.InzoneUserRepo) *InzoneApi {
	ins := &InzoneApi{
		userRepo: userRepo,
	}

	ins.c = cron.New(cron.WithSeconds())
	// 任务列表
	_, err := ins.c.AddFunc(config.C.Cron.GetLuckListJob, ins.Run)
	if err != nil {
		panic(err)
	}
	ins.c.Start()

	return ins
}

type InzoneApi struct {
	userRepo *dao.InzoneUserRepo
	c        *cron.Cron
}

func (a *InzoneApi) Upload(ctx *gin.Context, req *pb.File) (*pb.Empty, error) {
	file, err := req.FileHeader.Open()
	if err != nil {
		return nil, err
	}
	var buf []byte
	buf, err = ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	excelFile, err := excelize.OpenReader(bytes.NewReader(buf))
	if err != nil {
		return nil, err
	}
	defer excelFile.Close()
	_, err = a.readUsersByFile(excelFile)
	if err != nil {
		return nil, err
	}

	var f *os.File
	f, err = os.Create(path.Join(config.C.BasePath, "config", "user.xlsx"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	_, err = f.Write(buf)
	if err != nil {
		return nil, err
	}
	return &pb.Empty{}, nil
}

func (a *InzoneApi) Download(ctx *gin.Context) {
	buf, err := os.ReadFile(path.Join(config.C.BasePath, "config", "user.xlsx"))
	if err != nil {
		log.Error(err)
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", "attachment; filename=users.xlsx")
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(buf)))
	_, err = ctx.Writer.Write(buf)
	if err != nil {
		log.Error(err)
	}
}

// Run run
func (a *InzoneApi) Run() {
	myUsersM, err := a.readUsers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// 每日名单
	luckUserM, err := a.getDailyLuckUsers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	sheet2Hook := config.GetSheet2Hook()

	for sheet, hook := range sheet2Hook {

		fmt.Println(fmt.Sprintf("发送表:%s ======================", sheet))

		users := myUsersM[sheet]

		// 每日
		func() {
			data := make([]*pb.InzoneUser, 0)
			for _, user := range users {
				if luckUserM[user.ID] != nil {
					data = append(data, user)
				}
			}
			a.pushLuckMsg(hook, "每日名单", []*Record{{
				Goods:       "名单",
				InzoneUsers: data,
			}})
		}()
	}

	// 系列名单
	res, err := a.getSeriesLuckUsers()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	for sheet, hook := range sheet2Hook {

		fmt.Println(fmt.Sprintf("发送表:%s ======================", sheet))

		users := myUsersM[sheet]

		// 系列
		func() {
			records := make([]*Record, 0)
			for k, v := range res {
				data := make([]*pb.InzoneUser, 0)
				for _, user := range users {
					if v[user.ID] != nil {
						data = append(data, user)
					}
				}
				records = append(records, &Record{
					Goods:       fmt.Sprintf("系列%d", k+1),
					InzoneUsers: data,
				})
			}

			a.pushLuckMsg(hook, "系列名单", records)
		}()
	}

}

// 获取每日获奖名单那
func (a *InzoneApi) getDailyLuckUsers() (map[string]*pb.InzoneUser, error) {
	for {
		buf, err := a.query("http://wx0.yinzuo.cn/index.php/MaoTCT/luckshow.html")
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

		fmt.Println("每日名单")
		m := make(map[string]*pb.InzoneUser)
		doc.Find(".luckyLists").Find(".theLists").Find(".uf").Each(func(i int, s *goquery.Selection) {
			name := strings.ReplaceAll(s.Children().First().Text(), " ", "")
			name = strings.ReplaceAll(name, " ", "")
			phone := s.Children().Last().Text()
			id := name + phone
			//luckUsers = append(luckUsers, NewUser(name+phone, name, phone, ""))
			m[id] = NewInzoneUser(id, name, phone, "")

			fmt.Println(phone + "  " + name)
		})
		if len(m) <= 0 {
			time.Sleep(time.Second)
			continue
		}
		return m, nil
	}
}

// 获取每日获奖名单那
func (a *InzoneApi) getSeriesLuckUsers() ([]map[string]*pb.InzoneUser, error) {
	for {
		buf, err := a.query("http://wx0.yinzuo.cn/index.php/MaoTCT/luckshow.html")
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
		res := make([]map[string]*pb.InzoneUser, 0)
		doc.Find(".hasLuckys").Children().Each(func(i int, s *goquery.Selection) {
			//goods := s.Find(".goodsNameRight").Text()
			fmt.Println(fmt.Sprintf("系列%d名单", i+1))
			m := make(map[string]*pb.InzoneUser)
			s.Find(".luckyLists").Find(".theLists").Find(".uf").Each(func(i int, s *goquery.Selection) {
				name := strings.ReplaceAll(s.Children().First().Text(), " ", "")
				name = strings.ReplaceAll(name, " ", "")
				phone := s.Children().Last().Text()
				id := name + phone
				m[id] = NewInzoneUser(id, name, phone, "")

				fmt.Println(phone + "  " + name)
			})
			res = append(res, m)
		})
		if len(res) <= 0 {
			return nil, errors.New("系列名单还未公布")
		}
		return res, nil
	}
}

func (a *InzoneApi) pushLuckMsg(hook, head string, records []*Record) {
	text := fmt.Sprintf("今日消息:%s \n", head)
	for _, v := range records {
		text += fmt.Sprintf("%s  (%d人)\n", v.Goods, len(v.InzoneUsers))
		for _, user := range v.InzoneUsers {
			text += "   " + user.Name + "  " + user.Phone + "  " + user.Remark + "\n"
		}
	}

	fmt.Println("打印，并且发送通知")
	fmt.Println(text)

	msg := fmt.Sprintf(`{
	   "msgtype": "text",
	   "text": {
	       "content": "%s"
	   }
	}`, text)

	r := bytes.NewBuffer([]byte(msg))
	_, err := http.Post(hook, "application/json", r)
	if err != nil {
		fmt.Println("发送失败", err.Error())
	} else {
		fmt.Println("发送成功")
	}
}

func (a *InzoneApi) query(url string) ([]byte, error) {

	user, err := a.userRepo.GetValidFirst(context.Background())
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user IsEmpty")
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	req.Header.Set("host", "wx0.yinzuo.cn")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/98.0.4758.102 Safari/537.36 NetType/WIFI MicroMessenger/7.0.20.1781(0x6700143B) WindowsWechat(0x6309001c) XWEB/6500")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("Referer", "http://wx0.yinzuo.cn/index.php/MaoTCT/indexnew.html")
	//req.Header.Set("Accept-Encoding", "gzip, deflate")
	req.Header.Set("Accept-Language", "zh-CN,zh")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Cookie", user.Cookie)
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
	var buf []byte
	buf, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return buf, nil
}

func (a *InzoneApi) readUsers() (map[string][]*pb.InzoneUser, error) {
	f, err := excelize.OpenFile(path.Join(config.C.BasePath, "config", "user.xlsx"))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return a.readUsersByFile(f)
}

// 读取excel用户信息
func (a *InzoneApi) readUsersByFile(f *excelize.File) (map[string][]*pb.InzoneUser, error) {
	res := make(map[string][]*pb.InzoneUser)
	for _, sheet := range f.GetSheetList() {
		users := res[sheet]
		if users == nil {
			res[sheet] = make([]*pb.InzoneUser, 0)
		}

		rows, err := f.GetRows(sheet)
		if err != nil {
			return nil, err
		}
		for k, row := range rows {
			if len(row) < 2 {
				continue
			}
			//println(len(row[0]), row[0], len(row[1]),row[1])
			if len(row[0])%3 != 0 || len(row[0]) == 0 {
				errStr := fmt.Sprintf("%s 第 %d 行 姓名格式不对,长度应该为3的倍数，实际长度%d", sheet, k+1, len(row[0]))
				return nil, errors.New(errStr)
			}
			if len(row[1]) != 11 {
				errStr := fmt.Sprintf("%s 第 %d 行 手机号格式不对,长度应该为11，实际长度%d", sheet, k+1, len(row[1]))
				return nil, errors.New(errStr)
			}
			var remark string
			if len(row) > 2 {
				remark = row[2]
			}
			id := ""
			if len(row[0]) == 9 {
				id += row[0][0:3] + "**"
			} else {
				id += row[0][0:3] + "*"
			}

			id += row[1][0:3] + "*****" + row[1][8:11]

			user := NewInzoneUser(id, row[0], row[1], remark)
			res[sheet] = append(res[sheet], user)
		}

	}
	return res, nil
}

func NewInzoneUser(id, name, phone, remark string) *pb.InzoneUser {
	return &pb.InzoneUser{
		ID:     id,
		Name:   name,
		Phone:  phone,
		Remark: remark,
	}
}

type Record struct {
	Goods       string
	InzoneUsers []*pb.InzoneUser
}
