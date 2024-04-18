package xprinter

// 开发文档 https://www.xpyun.net/open/index.html
import (
	"bytes"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

const (
	defaultApiHost = "https://open.xpyun.net"
	printApi       = "/api/openapi/xprinter/print"
	checkStatus    = "/api/openapi/xprinter/queryPrinterStatus"
)

// GetXPrinterRequestPath 获取打印机请求地址
func GetXPrinterRequestPath(path string) string {
	printerHost := os.Getenv("PRINTER_HOST")
	if printerHost == "" {
		return defaultApiHost + path
	}
	return printerHost + path
}

type Printer struct {
	Sn      string `json:"sn"  bson:"sn"`
	User    string `json:"user"  bson:"user"`
	UserKey string `json:"user_key"  bson:"user_key"`
	Debug   string `json:"debug"  bson:"debug"`
}

// PrintReq 打印请求参数
type PrintReq struct {
	// 打印机编号
	Sn string `json:"sn"`
	// 打印内容,打印内容使用 GBK 编码判断,不能超过12K
	Content string `json:"content"`
	// 打印份数，默认为1，取值[1-65535]，超出范围将视为无效参数
	Copies int `json:"copies"`
	// 声音播放模式，0 为取消订单模式，1 为静音模式，2 为来单播放模式，3 为有用户申请退单了，默认为 2 来单播放模式
	Voice     int    `json:"voice"`
	User      string `json:"user"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
	Debug     string `json:"debug"`
	// Money 支付金额：
	// 最多允许保留2位小数。
	// 仅用于支持金额播报的芯烨云打印机。
	Money float64 `json:"money"`
	//支付与否：
	//取值范围59~61：
	//退款 59 到账 60 消费 61。
	//仅用于支持金额播报的芯烨云打印机。
	PayMode int `json:"payMode"`
	// 支付方式：
	//取值范围41~55：
	//支付宝 41、微信 42、云支付 43、银联刷卡 44、银联支付 45、会员卡消费 46、会员卡充值
	// 47、翼支付 48、成功收款 49、嘉联支付 50、壹钱包 51、京东支付 52、快钱支付 53、威支付 54、享钱支付 55
	//仅用于支持金额播报的芯烨云打印机。
	PayType int `json:"payType"`
	// 打印模式：
	// 值为 0 或不指定则会检查打印机是否在线，如果不在线 则不生成打印订单，直接返回设备不在线状态码；如果在线则生成打印订单，并返回打印订单号。
	// 值为 1不检查打印机是否在线，直接生成打印订单，并返回打印订单号。如果打印机不在线，订单将缓存在打印队列中，打印机正常在线时会自动打印。
	Mode int `json:"mode"`
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
	req, err := http.NewRequest("POST", GetXPrinterRequestPath(printApi), bytes.NewBuffer(jsonPayload))
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

// CheckReq 打印请求查询机器状态
type CheckReq struct {
	// 打印机编号
	Sn        string `json:"sn"`
	User      string `json:"user"`
	Timestamp string `json:"timestamp"`
	Sign      string `json:"sign"`
}

type PrinterStatusResp struct {
	Msg                string `json:"msg"`
	Code               int    `json:"code"`
	Data               int    `json:"data"`
	ServerExecutedTime int    `json:"serverExecutedTime"`
}

// Check
//
//	{
//	   "sn": "XPY123456789A",
//	   "user": "testuser",
//	   "timestamp": "1565417654",
//	   "sign": "82bdcbe2cf6ac4923339b13c2aad1f95ddf0b0a8",
//	   "debug": "0"
//	}
//
// resp success
//
//	{
//	   "msg":"ok",
//	   "code":0,
//	   "data":1,
//	   "serverExecutedTime":1
//	}
//
// resp failed
//
//	{
//	   "msg":"REQUEST_PARAM_INVALID",
//	   "code":-2,
//	   "data":0,
//	   "serverExecutedTime":1
//	}
func (p *Printer) Check() (*PrinterStatusResp, error) {

	sign, t := p.SignNow()
	printRequest := CheckReq{
		Sn:        p.Sn,
		User:      p.User,
		Timestamp: t, // or any other timestamp format you prefer
		Sign:      sign,
	}

	// Convert struct to JSON
	jsonPayload, err := json.Marshal(printRequest)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return nil, err
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", GetXPrinterRequestPath(checkStatus), bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	// Create HTTP client
	client := &http.Client{}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			logrus.Error(err)
		}
	}(resp.Body)

	// Read the response body as a byte slice
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return nil, err
	}

	response := &PrinterStatusResp{}
	err = json.Unmarshal(body, response)
	if err != nil {
		fmt.Println("Error Unmarshal response:", err)
		return nil, err
	}

	// Print response status
	fmt.Println("Response Status:", resp.Status)
	return response, nil
}
