package main

import (
	"bytes"
	"encoding/base64"
	"github.com/skip2/go-qrcode"
	"strings"
	"syscall/js"
)

// function sanitize must remove from string characters < and >
func sanitize(s string) string {
	if len(s) > 150 {
		s = s[:150]
	}

	s = strings.ReplaceAll(s, "<", "")
	s = strings.ReplaceAll(s, ">", "")
	return s
}

func generateQRCode(this js.Value, p []js.Value) interface{} {
	text := p[0].String()
	var png []byte
	png, err := qrcode.Encode(sanitize(text), qrcode.Medium, 256)
	if err != nil {
		return err.Error()
	}

	buf := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buf)
	encoder.Write(png)
	encoder.Close()

	js.Global().Get("document").Call("getElementById", "qrcode").Set("src", "data:image/png;base64,"+buf.String())
	return nil
}

func main() {
	c := make(chan struct{}, 0)

	js.Global().Set("generateQRCode", js.FuncOf(generateQRCode))

	<-c
}
