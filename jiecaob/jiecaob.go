package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	fmtt "fmt"
	"io"
	"io/ioutil"
	logutil "log"
	"net/http"
	"os"
	"time"

	"github.com/getsentry/raven-go"
)

const logFileName = "jiecaob.log"
const userID = "2017319"
const to = "243786753@qq.com"

const DATE_TIME_FORMAT = "2006-01-02 15:04:05"
const DATE_FORMAT = "2006-01-02"

var cst *time.Location
var log *logutil.Logger

func init() {
	cst = time.FixedZone("CST", 28800)
	raven.SetDSN("https://93c90027363c4dd898db3a3c860bee77:8eb0e1445cde43cf86e8b8bd4c585eb1@sentry.szzbmy.com/23")
}

func isWorkDay(t *time.Time) bool {
	d := t.Weekday()
	if d != time.Sunday && d != time.Saturday {
		return true
	}
	return false
}

func handErr(err error) {
	if err != nil {
		panic(err)
	}
}

func parseTime(s string) (*time.Time, error) {
	t, err := time.ParseInLocation(DATE_TIME_FORMAT, s, cst)
	return &t, err
}

func isStopWork(t *time.Time) bool {
	if t.Hour() == 18 && t.Minute() >= 30 {
		return true
	}
	if t.Hour() > 18 {
		return true
	}
	return false
}

func isStartWork(t *time.Time) bool {
	if t.Hour() < 9 {
		return true
	}
	if t.Hour() == 9 && t.Minute() <= 30 {
		return true
	}
	return false
}

func getData(userID string) (body []byte, err error) {
	resp, err := http.Get("http://jiecaob.szzbmy.com/sign/getRecord?userId=" + userID)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func query(now *time.Time, userID string) {
	str := "工作日"
	if !isWorkDay(now) {
		str = "周末"
	}
	log.Print("现在时间: ", now.Format(DATE_TIME_FORMAT), " ", str)

	yesterday := time.Date(now.Year(), now.Month(), now.Day()-1, 0, 0, 0, 0, cst)
	if !isWorkDay(&yesterday) {
		log.Println("昨天不是工作日，不查询")
		return
	}

	body, err := getData(userID)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return
	}
	data := make(map[string]interface{})
	err = json.Unmarshal(body, &data)
	if err != nil {
		raven.CaptureErrorAndWait(err, nil)
		return
	}

	v, ok := data["status"].(float64)
	if !ok || int(v) != 200 {
		raven.CaptureMessageAndWait("接口调用失败", nil)
		log.Println("接口调用失败:", string(body))
		return
	}

	arr, ok := data["data"].([]interface{})
	if !ok {
		raven.CaptureMessageAndWait("没有data", nil)
		log.Println("没有data:", string(body))
		return
	}

	foundStartWork := false
	foundStopWork := false
	times := make([]*time.Time, 0, 2)

	for i, v := range arr {
		str, ok := v.(string)
		if !ok {
			log.Printf("data[%d] is not string\n", i)
			continue
		}
		t, err := parseTime(str)
		if err != nil {
			raven.CaptureErrorAndWait(err, nil)
			log.Println("解析时间失败:", str)
			continue
		}
		if t.After(yesterday) {
			times = append(times, t)
			if isStartWork(t) {
				foundStartWork = true
				log.Println("发现上班打卡: ", t.Format(DATE_TIME_FORMAT))
			}
			if isStopWork(t) {
				foundStopWork = true
				log.Println("发现下班打卡: ", t.Format(DATE_TIME_FORMAT))
			}
		}
	}

	if foundStartWork && foundStopWork {
		log.Println("打卡正常")
		return
	}

	datestr := yesterday.Format(DATE_FORMAT)
	subject := fmt.Sprintf("打卡异常 %s 工号: %s", datestr, userID)

	var b bytes.Buffer
	b.WriteString("<h1>" + subject + "</h1>")
	if len(times) > 0 {
		b.WriteString(datestr + " 的打卡记录: ")
		b.WriteString("<ul>")
		for _, t := range times {
			b.WriteString("<li>" + t.Format(DATE_TIME_FORMAT) + "</li>")
		}
		b.WriteString("</ul>")
	}
	if !foundStartWork {
		b.WriteString("<p>上班没打卡</p>")
	}
	if !foundStopWork {
		b.WriteString("<p>下班没打卡</p>")
	}
	b.WriteString(`<div style="display: inline-block; padding:30px; color:red; ">
		<a style="text-decoration: none; color: #ffffff" href="http://jiecaob.szzbmy.com/" target="_blank">
			<div style="font-family:'微软雅黑';font-size: 18px; text-decoration: none; white-space: nowrap; color: #ffffff; padding: 10px 25px; text-align: center;  margin: 0px;  background-color: #cc0001; border-radius: 3px">
				马上查询
			</div>
		</a>
	</div>`)

	log.Println(b.String())

	// var b bytes.Buffer
	// writer := bufio.NewWriter(&b)
	// tpl, err := template.ParseFiles("m.tpl")
	// handErr(err)
	// err = tpl.Execute(writer, nil)
	// handErr(err)
	// err = writer.Flush()
	// handErr(err)

	// log.Println(b.String())

	const senderName = "漏打卡提醒"
	const sender = "robot@zengming.top"

	err = SendMail(senderName, sender, to, subject, b.String())
	if err != nil {
		log.Println("发送邮件失败")
		raven.CaptureErrorAndWait(err, nil)
		return
	}
	log.Println("发送成功")
}

func main() {
	p := raven.NewPacket("hello go")
	p.Level = raven.FATAL
	p.Extra["foo"] = "foodata"
	p.Extra["objdata"] = map[string]interface{}{
		"haha": "hehe",
		"a":    1234,
		"c":    []int{1, 5},
	}

	eventID, ch := raven.Capture(p, nil)
	if eventID != "" {
		<-ch
	}
	fmt.Println("eventID", eventID)

	// raven.CaptureErrorAndWait()
	// raven.CaptureMessageAndWait()
	// raven.CapturePanicAndWait()

	f, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmtt.Println("无法创建日志文件")
		panic(err)
	}
	defer f.Close()
	log = logutil.New(io.MultiWriter(f, os.Stdout), "", logutil.LstdFlags)

	for {
		now := time.Now()
		// now = time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, cst)
		query(&now, userID)

		tomorrow := time.Date(now.Year(), now.Month(), now.Day()+1, 10, 30, 0, 0, cst)
		log.Println("下次查询时间:", tomorrow)
		<-time.After(tomorrow.Sub(now))
	}
}
