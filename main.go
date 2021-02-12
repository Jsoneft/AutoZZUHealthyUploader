package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/robfig/cron"
	"gopkg.in/gomail.v2"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
	"unsafe"
)

type SMTPInfo struct {
	Host     string
	Port     int
	IsSSL    bool
	UserName string
	Password string
	From     string
}

type Email struct {
	*SMTPInfo
}

var (
	Tos               []string
	FirstPostMatchReg = `"(h.+)"`
)

// NewEmail 初始化 STMP 服务器配置 自己配了之后写入
func NewEmail() *Email {
	return &Email{&SMTPInfo{
		Host:     "mail.*****.**",
		Port:     465,
		IsSSL:    true,
		UserName: "*******",
		Password: "********",
		From:     "*****",
	}}
}

func (e Email) SendMail(subject, body string) error {
	// Tos  收件人

	Tos = []string{
		"********@qq.com",
	}
	for _, to := range Tos {
		err := e.SendMailOneByOne(to, subject, body)
		if err != nil {
			return err
		}
		log.Printf("发给 %v", to)
		time.Sleep(time.Duration(5) * time.Second)
	}
	return nil
}

func (e Email) SendMailOneByOne(to, subject, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", e.From)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)
	dialer := gomail.NewDialer(e.Host, e.Port, e.UserName, e.Password)
	dialer.TLSConfig = &tls.Config{
		InsecureSkipVerify: e.IsSSL,
	}
	return dialer.DialAndSend(m)
}

func Bytes2String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func FirstPOST() (string, error) {
	resp, err := http.Post("https://jksb.v.zzu.edu.cn/vls6sss/zzujksb.dll/login", "application/x-www-form-urlencoded", strings.NewReader("uid=**********&upw=*********&smbtn=%E8%BF%9B%E5%85%A5%E5%81%A5%E5%BA%B7%E7%8A%B6%E5%86%B5%E4%B8%8A%E6%8A%A5%E5%B9%B3%E5%8F%B0&hh28=861"))
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	RawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	body := Bytes2String(RawBody)
	compile := regexp.MustCompile(FirstPostMatchReg)
	matches := compile.FindStringSubmatch(body)
	log.Printf("body : %v", body)
	log.Printf("match :%v", matches[0])
	if matches != nil {
		return matches[1], nil
	}
	return "", errors.New("body parse err")
}

func SecondPost(lastURL string) error {
	u, err := url.Parse(lastURL)
	if err != nil {
		return err
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return err
	}
	//log.Printf("%s %s", m["ptopid"][0],m["sid"][0])
	resp, err := http.Post("https://jksb.v.zzu.edu.cn/vls6sss/zzujksb.dll/jksb", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("day6=b&did=1&door=&men6=a&ptopid=%s&sid=%s", m["ptopid"][0], m["sid"][0])))
	defer resp.Body.Close()
	if err != nil {
		return err
	}
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	body := Bytes2String(rawBody)
	log.Printf("第二次打卡body = %v", body)
	return nil

}

func ThirdPost(lastURL string) (string, error) {
	u, err := url.Parse(lastURL)
	if err != nil {
		return "", err
	}
	m, err := url.ParseQuery(u.RawQuery)
	if err != nil {
		return "", err
	}
	resp, err := http.Post("https://jksb.v.zzu.edu.cn/vls6sss/zzujksb.dll/jksb", "application/x-www-form-urlencoded", strings.NewReader(fmt.Sprintf("myvs_1=否&myvs_2=否&myvs_3=否&myvs_4=否&myvs_5=否&myvs_6=否&myvs_7=否&myvs_8=否&myvs_9=否&myvs_10=否&myvs_11=否&myvs_12=否&myvs_13a=41&Btn3=获取地市&myvs_13b=4101&myvs_13c=河南省.郑州市.金水区&myvs_14=否&memo22=请求超时&did=2&day6=b&men6=a&sheng6=41&jingdu=113.64&weidu=34.71&ptopid=%s&sid=%s", m["ptopid"][0], m["sid"][0])))
	defer resp.Body.Close()
	if err != nil {
		return "", err
	}
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	body := Bytes2String(rawBody)
	log.Printf("第三次打卡body = %v", body)
	return body, err
}

func Process() (string, error) {

	FirstGetURL, err := FirstPOST()

	if err != nil {
		return "", err
	}

	err = SecondPost(FirstGetURL)
	if err != nil {
		return "", err
	}
	bodyMsg, err := ThirdPost(FirstGetURL)
	if err != nil {
		return "", err
	}
	return bodyMsg, err
}
func work() {
	email := NewEmail()
	Body, err := Process()
	cnt := 1
	if err != nil {
		errMsg := fmt.Sprintf("来自左佳逊的打卡自动提醒服务， 于%v 尝试打卡，不幸， 打卡失败 err = %v, 请及时报备zjx QQ:1018437256", time.Now(), err)
		log.Fatal(errMsg)
		MailErr := email.SendMail("打卡失败提醒", errMsg)
		for MailErr != nil {
			log.Printf("第%v次重发邮件 (当前打卡失败)", cnt)
			cnt += 1
			time.Sleep(5 * time.Second)
			MailErr = email.SendMail("打卡失败提醒", fmt.Sprintf("来自左佳逊的打卡自动提醒服务， 于%v 尝试打卡，不幸， 打卡失败 err = %v, 请及时报备zjx QQ:1018437256", time.Now(), err))
		}
	} else {
		sucMsg := fmt.Sprintf("来自左佳逊的打卡自动提醒服务， 于%v 尝试打卡，打卡成功, 不信你看 %v", time.Now(), Body)
		log.Printf(sucMsg)
		MailErr := email.SendMail("打卡成功", sucMsg)
		for MailErr != nil {
			log.Printf("第%v次重发邮件 (打卡成功)", cnt)
			cnt += 1
			time.Sleep(5 * time.Second)
			MailErr = email.SendMail("打卡成功", fmt.Sprintf("来自左佳逊的打卡自动提醒服务， 于%v 尝试打卡，打卡成功, 不信你看 %v", time.Now(), Body))
		}
		if MailErr == nil {
			log.Printf("邮件发送成功 body: %s \n 应该含有 ***同学，感谢你今日上报健康状况！安立民书记、陈思坤院长将对上报状况进行审核。记着明天继续来报。 字样", Body)
		} else {
			log.Printf("邮件发送失败 %v", MailErr)
		}
	}
}

func main() {
	log.Printf("started...")

	c := cron.New()
	c.AddFunc("0 10 0,1 * * ?", work)
	c.Start()
	select {}
}
