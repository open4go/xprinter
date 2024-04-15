package tp

import (
	"bytes"
	"fmt"
	"html/template"
	"time"
)

// 准备模板字符串
const receiptTemplate = `<CB>黄李记<BR><BR><BR></CB>
<L><N>-----------------------------------------------
<HB>#{{.QueueNumber}}
<L><N>下单时间:{{.OrderTime}}
订单编号: {{.OrderID}}
********************商品清单********************
<L>
<LINE p="20,36" />菜名<HT>数量<HT>单价<BR>
{{range .Items}}--------------------------------------<BR>
{{.Name}}<HT>{{.Amount}}<HT>{{.Price}}<BR>
{{end}}--------------------------------------<BR>
</L>
**********************************************
<R><B>合计：{{.Total}}元</B><BR></R><BR>
<L>门店地址：{{.StoreAddress}}<BR>
门店电话：{{.StorePhone}}<BR>
客户地址：{{.Address}}<BR>
客户电话：{{.Phone}}<BR>
支付时间：{{.PayTime}}<BR>
备注：{{.Note}}<BR>
</L>
<C><QRCODE s=8 e=L l=center>{{.Invoice}}</QRCODE><BR>
<C><BARCODE t=CODE128 w=2 h=100 p=2>{{.TradeID}}</BARCODE>
`

// ReceiptData 定义小票数据结构
type ReceiptData struct {
	Items        []Item
	Total        string
	Address      string
	Phone        string
	OrderTime    string
	Note         string
	Invoice      string
	StoreAddress string
	StorePhone   string
	OrderID      string
	QueueNumber  string
	PayTime      string
	TradeID      string
}

// Item 定义菜品条目结构
type Item struct {
	Name   string
	Amount int
	Price  float64
}

func LoadExample() ReceiptData {
	// 准备数据
	items := []Item{
		{Name: "可乐鸡翅", Amount: 2, Price: 9.99},
		{Name: "水煮鱼特辣", Amount: 1, Price: 108.00},
		{Name: "豪华版超级无敌龙虾炒饭", Amount: 1, Price: 99.90},
		{Name: "炭烤鳕鱼", Amount: 5, Price: 19.99},
	}
	total := "327.83"
	address := "珠海市香洲区xx路xx号"
	phone := "1363*****88"
	orderTime := time.Now().Format("2006-01-02 15:04:05")
	note := "少放辣 不吃香菜"
	qrCodeURL := "http://www.xpyun.net"

	// 渲染模板
	receiptData := ReceiptData{
		Items:     items,
		Total:     total,
		Address:   address,
		Phone:     phone,
		OrderTime: orderTime,
		Note:      note,
		Invoice:   qrCodeURL,
	}

	return receiptData
}

func Load(items []Item, address string, phone string, note string, total string) ReceiptData {
	// 准备数据
	orderTime := time.Now().Format("2006-01-02 15:04:05")
	qrCodeURL := "http://www.xpyun.net"

	// 渲染模板
	receiptData := ReceiptData{
		Items:     items,
		Total:     total,
		Address:   address,
		Phone:     phone,
		OrderTime: orderTime,
		Note:      note,
		Invoice:   qrCodeURL,
	}

	return receiptData
}

func Render(data ReceiptData) string {

	// 准备模板对象
	tmpl, err := template.New("receipt").Parse(receiptTemplate)
	if err != nil {
		fmt.Println("Error parsing template:", err)
		return ""
	}

	var output bytes.Buffer

	err = tmpl.Execute(&output, data)
	if err != nil {
		fmt.Println("Error executing template:", err)
		return ""
	}

	return output.String()
}
