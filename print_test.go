package xprinter

import (
	"github.com/joho/godotenv"
	"github.com/open4go/xprinter/tp"
	"log"
	"os"
	"testing"
)

func TestPrinter_Print(t *testing.T) {
	// Load .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
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
				tp.RenderNow(),
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
