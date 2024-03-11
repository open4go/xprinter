package xprinter

import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"time"
)

const (
	printAPI = "https://open.xpyun.net/api/openapi/xprinter/print"
)

type Printer struct {
	Sn      string `json:"sn"`
	User    string `json:"user"`
	UserKey string `json:"user_key"`
	Debug   string `json:"debug"`
}

type PrintReq struct {
	Sn        string  `json:"sn"`
	Content   string  `json:"content"`
	Copies    int     `json:"copies"`
	Voice     int     `json:"voice"`
	User      string  `json:"user"`
	Timestamp string  `json:"timestamp"`
	Sign      string  `json:"sign"`
	Debug     string  `json:"debug"`
	Money     float64 `json:"money"`
}

// SignNow 参数说明：例如：user=acc、UserKEY=abc、timestamp=acbc，那么先拼成字符串accabcacbc，再将此字符串进行SHA1加密，得到sign。
func (p *Printer) SignNow() (string, string) {
	// Calculate SHA1 hash of JSON payload
	hash := sha1.New()
	t := fmt.Sprintf("%d", time.Now().Unix())
	hash.Write([]byte(fmt.Sprintf("%s%s%s", p.User, p.UserKey, t)))
	hashSum := hash.Sum(nil)
	signature := fmt.Sprintf("%x", hashSum)
	return signature, t
}

// Print
//
//	{
//	   "sn": "XPY123456789A",
//	   "content": "----货到付款----",
//	   "copies": 1,
//	   "voice": 2,
//	   "user": "testuser",
//	   "timestamp": "1565417654",
//	   "sign": "82bdcbe2cf6ac4923339b13c2aad1f95ddf0b0a8",
//	   "debug": "0"
//	}
func (p *Printer) Print(content string) {

	sign, t := p.SignNow()
	printRequest := PrintReq{
		Sn:        p.Sn,
		Content:   content,
		Copies:    1,
		Voice:     2,
		User:      p.User,
		Timestamp: t, // or any other timestamp format you prefer
		Sign:      sign,
		Debug:     "0",
		//Money:     22.12,
	}

	// Convert struct to JSON
	jsonPayload, err := json.Marshal(printRequest)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", printAPI, bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client
	client := &http.Client{}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(resp.Body)

	// Print response status
	fmt.Println("Response Status:", resp.Status)
}
