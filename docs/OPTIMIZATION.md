# xprinter 优化方案

| # | 优先级 | 问题 | 位置 | 建议 | 工作量 |
|---|---|---|---|---|---|
| 1 | 高 | `Print` 没有返回值，所有错误只用 `fmt.Println` 打印；也不解析响应体（芯烨即使返回 HTTP 200，body 里的 `code` 也可能表示失败），调用方无法知道打印是否成功 | `print.go` Print | 改为 `Print(ctx, content) (*PrintResp, error)`，检查 `code==0` | S |
| 2 | 高 | 用 `&http.Client{}`，没有超时 | `print.go:128,216` | 共享一个带 `Timeout: 10s` 的 client，请求带 ctx | S |
| 3 | 中 | 小票模板用的是 `html/template`，会把用户备注里的 `&`、`<` 转成 HTML 实体（`&amp;`），打出来是乱码；而用户输入里如果有 `<BR>`、`<QRCODE>` 这类打印指令，又能注入排版指令 | `tp/order.go:6` | 改用 `text/template`，对用户字段过滤掉 `<`、`>` 等打印机控制标签 | S |
| 4 | 中 | 模板里写死了门店名 “黄李记”，二维码地址写死为 `http://www.xpyun.net` | `tp/order.go:11`、Load | 改从 `ReceiptData` 或 `MetaInfo` 传入 | S |
| 5 | 中 | `Printer.UserKey` 带 `json:"user_key"`，打印机配置一旦被接口返回就会泄露密钥 | `print.go` Printer | 改成 `json:"-"`，另外定义入参 DTO | S |
| 6 | 中 | 测试会真实调用外部接口，没有 `.env` 时 `log.Fatal`，CI 无法运行 | `print_test.go` | 用 `httptest.Server` 模拟；真实调用改为 `-tags=integration` | S |
| 7 | 低 | 金额 `Price float64`，`PrintReq.Money float64` | `tp/order.go`、`print.go` | 内部用分（int64），渲染时再格式化 | S |
| 8 | 低 | `ioutil.ReadAll` 已废弃，用了 logrus 而不是 open4go/log；Go 1.21.5 太旧 | `print.go`、`go.mod` | 改用 `io.ReadAll`，统一日志库和 Go 版本；补 CI | S |
