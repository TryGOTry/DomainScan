/*
* @Author: Try
* @Date:   2021/5/5 18:22
 */
package scan

import (
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/miekg/dns"
	"regexp"
	"strings"
	"time"
)

var (
	title = `<title>([\s\S]+?)</title>`
)

type Webinfo struct {
	StatusCode int
	Title      string
	Server     string
	Powered    string
	Body       string
	Res        string //成功的结果
	Bodylen    int    //返回包长度
	Ipadd      string
	Okdomain   string  //成功的域名
}

func Goscan(url string, flServerAddr string, timeout int64) (Webinfo, error) {
	//time.Sleep(time.Duration(timesleep) * time.Second) //设置延时时间
	var t string
	var Web Webinfo
	t = url
	a2 := strings.Replace(url, ":80/", "", -1)
	t = "http://" + a2
	if strings.Contains(a2, ":443/") {
		b := strings.Replace(url, ":443/", "", -1)
		var msg dns.Msg
		fqdn := dns.Fqdn(b)
		msg.SetQuestion(fqdn, dns.TypeA)
		in, err := dns.Exchange(&msg, flServerAddr+":53")
		if err != nil {
			//panic(err)
			return Web, err
		}
		if len(in.Answer) < 1 {
			//fmt.Println("No records")
			return Web, err
		}
		for _, answer := range in.Answer {
			if a, ok := answer.(*dns.A); ok {
				//fmt.Println(a.A)
				Web.Okdomain = b
				Web.Ipadd = fmt.Sprintf("%s ", a.A)
			}
		}
	} else {
		var msg dns.Msg
		fqdn := dns.Fqdn(a2)
		msg.SetQuestion(fqdn, dns.TypeA)
		in, err := dns.Exchange(&msg, flServerAddr+":53")
		if err != nil {
			//panic(err)
			return Web, err
		}
		if len(in.Answer) < 1 {
			//fmt.Println("No records")
			//return
			return Web, err
		}
		for _, answer := range in.Answer {
			if a, ok := answer.(*dns.A); ok {
				//fmt.Println(a.A)
				Web.Okdomain = a2
				Web.Ipadd = fmt.Sprintf("%s ", a.A)
			}
		}
	}
	//fmt.Println(Web.Ipadd)
	if Web.Ipadd != "127.0.0.1" {
		//fmt.Println(Web.Ipadd)
		if strings.Contains(url, ":443/") {
			b := strings.Replace(url, ":443/", "", -1)
			t = "https://" + b
		} else {
			a := strings.Replace(url, ":80/", "", -1)
			t = "http://" + a
		}
		//fmt.Println(t)
		client := resty.New().SetTimeout(time.Duration(timeout) * time.Second).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true}) //忽略https证书错误，设置超时时间
		client.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36")
		resp, err := client.R().EnableTrace().Get(t) //开始请求扫描
		if err != nil {
			//log.Println(err)
			return Web, err
		}
		str := resp.Body()
		body := string(str)
		//fmt.Println(body)
		if strings.Contains(body, "<title>") {
			re1 := regexp.MustCompile(title) //正则取标题
			titlename := re1.FindAllStringSubmatch(body, 1)
			if len(titlename) > 0 {
				Web.Title = titlename[0][1]
			}
		}

		//fmt.Println(b)
		//fmt.Println("Resolved address is ", addr.String())
		Web.StatusCode = resp.StatusCode()
		Web.Powered = resp.Header().Get("X-Powered-By")
		//Web.Title = titlename[0][1]
		Web.Server = resp.Header().Get("server")
		Web.Body = body
		Web.Res = t
		Web.Bodylen = len(body)
		return Web, nil
	}
	return Web, errors.New("err")
}
