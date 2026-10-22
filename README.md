# xprinter — 芯烨云打印机（xpyun）客户端与小票模板

> 仓库：https://github.com/open4go/xprinter ｜ 模块：`github.com/open4go/xprinter` ｜ Go 1.21.5 ｜ 最近提交：2026-05-01（最新 tag v0.0.14）

## 项目简介
调用芯烨云开放平台（`https://open.xpyun.net`，可以用环境变量 `PRINTER_HOST` 覆盖）的打印接口和打印机状态查询接口；签名算法是 `SHA1(user + userKey + timestamp)`；`tp` 包提供小票数据结构和渲染模板（含二维码、条码）。

## 目录结构
```
xprinter/
├── print.go       Printer{Sn,User,UserKey,Debug}、SignNow、Print(content)、Check() (*PrinterStatusResp, error)
├── print_test.go  读取 .env 后真实调用打印接口
└── tp/order.go    ReceiptData/Item/MetaInfo、LoadExample、Load、Render、RenderFromClient（自定义模板）
```

## 使用示例
```go
p := &xprinter.Printer{Sn: sn, User: user, UserKey: key, Debug: "0"}
content := tp.Render(tp.Load(items, addr, phone, note, total))
p.Print(content)
st, err := p.Check() // Data: 打印机状态
```

## 依赖
`sirupsen/logrus`、`joho/godotenv`（测试）

## 构建与测试
`go build`、`go vet` 通过；**`go test` 失败**：`print_test.go` 需要本地 `.env`，否则 `log.Fatal("Error loading .env file")`。govulncheck 没有发现问题。

## 与其他仓库的关系
被 `req5rsp` 依赖（v0.0.9，最新是 v0.0.14），打印相关的 DTO 在 `req5rsp/req/printer.go`。
