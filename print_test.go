package xprinter

import (
	"github.com/joho/godotenv"
	"github.com/open4go/xprinter/tp"
	"os"
	"testing"
)

func TestPrinter_Print(t *testing.T) {
	// Load .env file
	// 该测试会调用真实打印机云端接口，需要 .env 中的 SN / User / UserKey；没有时跳过（CI 中没有凭据）
	if err := godotenv.Load(".env"); err != nil {
		t.Skip("skip: .env with real printer credentials not found")
	}

	type fields struct {
		Sn      string
		User    string
		UserKey string
		Debug   string
	}
	type args struct {
		content string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			"test",
			fields{
				os.Getenv("SN"),
				os.Getenv("User"),
				os.Getenv("UserKey"),
				"0",
			},
			args{
				tp.Render(tp.LoadExample()),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &Printer{
				Sn:      tt.fields.Sn,
				User:    tt.fields.User,
				UserKey: tt.fields.UserKey,
				Debug:   tt.fields.Debug,
			}
			p.Print(tt.args.content)
		})
	}
}
