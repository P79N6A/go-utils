package sinawatch

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const alertUri = "/v1/alert/send"

type SinaWatch struct {
	IsIConnect    bool
	needAuthorize bool
	path          string
	Kid           string
	Password      string
	Host          string
	Port          int
	Timeout       int
	ApiExpired    int
	Ip            string
}

type Content struct {
	Subject string
	Content string
	Html    string
}

type Operation struct {
	Sv      string
	Service string
	Object  string
}

type Receiver struct {
	Mail        string
	Weibo       string
	Wechat      string
	Sms         string
	Ivr         string
	Push        string
	MailGroup   string
	WeiboGroup  string
	WechatGroup string
	SmsGroup    string
	IvrGroup    string
	PushGroup   string
}

func (receiver Receiver) IsEmpty() (empty bool) {

	if receiver.Mail == "" && receiver.MailGroup == "" &&
		receiver.Weibo == "" && receiver.WeiboGroup == "" &&
		receiver.Wechat == "" && receiver.WechatGroup == "" &&
		receiver.Sms == "" && receiver.SmsGroup == "" &&
		receiver.Ivr == "" && receiver.IvrGroup == "" &&
		receiver.Push == "" && receiver.PushGroup == "" {
		return true
	}

	return false
}

func (receiver Receiver) IsMailEmpty() (empty bool) {

	if receiver.Mail == "" && receiver.MailGroup == "" {
		return true
	}

	return false
}

func (receiver Receiver) IsSmsEmpty() (empty bool) {

	if receiver.Sms == "" && receiver.SmsGroup == "" {
		return true
	}

	return false
}

func (receiver Receiver) IsWeiboEmpty() (empty bool) {

	if receiver.Weibo == "" && receiver.WeiboGroup == "" {
		return true
	}

	return false
}

func (receiver Receiver) IsWechatEmpty() (empty bool) {

	if receiver.Wechat == "" && receiver.WechatGroup == "" {
		return true
	}

	return false
}

func (receiver Receiver) IsIvrEmpty() (empty bool) {

	if receiver.Ivr == "" && receiver.IvrGroup == "" {
		return true
	}

	return false
}

func (receiver Receiver) IsPushEmpty() (empty bool) {

	if receiver.Push == "" && receiver.PushGroup == "" {
		return true
	}

	return false
}

//http://wiki.intra.sina.com.cn/pages/viewpage.action?pageId=7162793
func (sinaWatch *SinaWatch) SendAlert(operation Operation, content Content, receiver Receiver, autoMerge int) (err error) {

	data := map[string]string{
		"url":        sinaWatch.Host,
		"auto_merge": string(autoMerge),
		"sv":         operation.Sv,
		"service":    operation.Service,
		"object":     operation.Object,
		"subject":    content.Subject,
		"content":    content.Content,
		"html":       content.Html,
		"mailto":     receiver.Mail,
		"msgto":      receiver.Sms,
		"ivrto":      receiver.Ivr,
		"weiboto":    receiver.Weibo,
		"wechatto":   receiver.Wechat,
		"pushto":     receiver.Push,
		"gmailto":    receiver.MailGroup,
		"gmsgto":     receiver.SmsGroup,
		"givrto":     receiver.IvrGroup,
		"gweiboto":   receiver.WeiboGroup,
		"gwechatto":  receiver.WechatGroup,
		"gpushto":    receiver.PushGroup,
	}

	sinaWatch.path = alertUri;
	res, err := sinaWatch.post(data)

	if err != nil {
		return err
	}

	code := res["code"]

	code, err = code.(json.Number).Int64()

	if err != nil {
		return err
	}

	if code != 0 && res["message"] != nil {
		return errors.New(fmt.Sprintf("sinawatch request error : %s", res["message"]))
	}

	return nil
}

func NewSinaWatch(kid string, password string, isIConnect bool) SinaWatch {

	host := "http://connect.monitor.sina.com.cn"

	if isIConnect {
		host = "http://iconnect.monitor.sina.com.cn"
	}

	sinaWatch := SinaWatch{
		Host:          host,
		Port:          80,
		Timeout:       1,
		needAuthorize: true,
		Kid:           kid,
		Password:      password,
		ApiExpired:    60,
	}
	return sinaWatch;
}

func (sinaWatch *SinaWatch) bindAuth(kid string, passwd string, api_expired int) {
	sinaWatch.needAuthorize = true
	sinaWatch.Kid = kid
	sinaWatch.Password = passwd
	sinaWatch.ApiExpired = api_expired
}

func (sinaWatch *SinaWatch) setTimeout(timeout int) {
	sinaWatch.Timeout = timeout
}

func (sinaWatch *SinaWatch) hash(r *http.Request, body string) {
	contentMd5 := ""

	if method := r.Method; method == "POST" {

		md5Ctx := md5.New()
		md5Ctx.Write([]byte(body))
		cipherStr := md5Ctx.Sum(nil)
		contentMd5 = hex.EncodeToString(cipherStr)
		r.Header.Add("Content-MD5", contentMd5)

	}

	contentType := r.Header.Get("Content-type")
	expired := fmt.Sprintf("%v%v", time.Now().Unix(), sinaWatch.Timeout)
	r.Header.Add("Expires", expired)

	if len(sinaWatch.Ip) != 0 {
		r.Header.Add("x-sinawatch-ip", sinaWatch.Ip)
	}

	headers := []string{}
	for key, value := range r.Header {

		ret1, _ := regexp.Match("x-sinawatch-", []byte(strings.ToLower(key)))
		ret2, _ := regexp.Match("x-sina-", []byte(strings.ToLower(key)))

		if ret1 || ret2 {
			item := fmt.Sprintf("%v:%v\n", strings.ToLower(key), strings.TrimSpace(value[0]))
			headers = append(headers, item)
		}
	}

	canonicalizedamzheaders := strings.Join(headers, "")
	canonicalizedresource := sinaWatch.path

	stringtosignlist := []string{}
	stringtosignlist = append(stringtosignlist, r.Method)
	stringtosignlist = append(stringtosignlist, contentMd5)
	stringtosignlist = append(stringtosignlist, contentType)
	stringtosignlist = append(stringtosignlist, expired)
	stringtosignlist = append(stringtosignlist, canonicalizedamzheaders+canonicalizedresource)
	stringtosign := strings.Join(stringtosignlist, "\n")

	key := []byte(sinaWatch.Password)
	mac := hmac.New(sha1.New, key)
	mac.Write([]byte(stringtosign))
	sign := base64.StdEncoding.EncodeToString(mac.Sum(nil))[5:15]
	authorization := fmt.Sprintf("sinawatch %s:%s", sinaWatch.Kid, sign)
	r.Header.Add("Authorization", authorization)

}

func (sinaWatch *SinaWatch) encodeMultipartFormdata(fileds map[string]string, files []string) (string, string) {
	/*
	 fields is a sequence of (name, value) elements for regular form fields.
	 files is a sequence of (name, filename, value) elements for data to be uploaded as files
	 Return (content_type, body) ready for httplib.HTTP instance
	*/
	boundary := fmt.Sprintf("----------0x%0x", time.Now().UnixNano()/1000)
	crlf := "\r\n"
	L := []string{}
	for key, value := range fileds {
		L = append(L, fmt.Sprintf("--%v", boundary))
		L = append(L, fmt.Sprintf("Content-Disposition: form-data; name=\"%s\"", key))
		L = append(L, "")
		L = append(L, value)
	}
	L = append(L, fmt.Sprintf("--%v--", boundary))
	L = append(L, "")
	body := strings.Join(L, crlf)
	contentType := fmt.Sprintf("multipart/form-data; boundary=%v", boundary)
	return contentType, body
}

func (sinaWatch *SinaWatch) post(data map[string]string) (map[string]interface{}, error) {
	watchurl := fmt.Sprintf("%s:%d", sinaWatch.Host, sinaWatch.Port)
	u, _ := url.ParseRequestURI(watchurl)
	u.Path = sinaWatch.path
	urlStr := u.String()
	files := []string{}
	contentType, body := sinaWatch.encodeMultipartFormdata(data, files)

	client := &http.Client{
		Timeout: time.Duration(sinaWatch.Timeout) * time.Second,
	}
	r, err := http.NewRequest("POST", urlStr, bytes.NewBufferString(body))
	if err != nil {
		return nil, errors.New(fmt.Sprintf("request %s err: %v", urlStr, err))
	}

	//add header
	r.Header.Add("Content-Type", contentType)
	r.Header.Add("Content-Length", strconv.Itoa(len(body)))
	if sinaWatch.needAuthorize {
		sinaWatch.hash(r, body)
	}

	resp, err := client.Do(r)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("request %s err: %v", urlStr, err))
	}

	//body
	ret, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("request %s err: %v", urlStr, err))
	}

	defer resp.Body.Close()
	//json
	retMap := make(map[string]interface{})

	d := json.NewDecoder(bytes.NewReader(ret))
	d.UseNumber()
	_ = d.Decode(&retMap)
	return retMap, nil
}
