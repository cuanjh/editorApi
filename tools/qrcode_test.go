package tools

import (
	"fmt"
	"testing"

	"github.com/boombuler/barcode/qr"
)

func TestQrCodeEncode(t *testing.T) {

	qr := NewQrCode(
		"http://www.talkmate.com",
		150,
		150,
		qr.M,
		qr.Auto,
	)
	qr.Encode("qrcode/")
}

func TestMerge(t *testing.T) {
	m := &Merge{
		BgFilePath:     "qrcode/",
		BgFileName:     "bg.jpg",
		QrCodeFilePath: "qrcode/",
		MergeFilePath:  "qrcode/",
		QrCodeFileName: MD5V("http://www.talkmate.com") + ".jpg",
		Pt: &Pt{
			X: 520,
			Y: 976,
		},
	}
	p, s, e := m.Generate("bgQr.jpg")
	fmt.Println(p, s, e)
}
