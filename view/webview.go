package view

import (
	_ "embed"
	"fmt"
	"strings"

	webview "github.com/webview/webview_go"
)

var (
	//go:embed _static/index.html
	basic string
)

// Webview embedでGUIファイルのInclude
func WebView(isProduction bool, port int) {
	w := webview.New(!isProduction)
	defer w.Destroy()

	w.SetTitle("HFT Assistant")
	w.SetSize(int(1920*0.5), int(1080*0.5), webview.HintNone)

	// 指定ポートが8080出なかった場合HTML及びJavascriptを書き換える
	if port != 8080 {
		basic = strings.ReplaceAll(basic, ":8080", fmt.Sprintf(":%d", port))
	}

	w.SetHtml(basic)
	// w.Navigate("http://localhost:8080")
	w.Run()
}
