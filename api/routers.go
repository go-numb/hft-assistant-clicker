package api

import (
	"fmt"
	"hft-assistant-clicker/setting"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Client struct {
	Setting *setting.Settings
}

// Routers プログラム実行部分とクライアント（GUI）の通信を行うローカルAPIサーバー
func (cli *Client) Routers(port int) {
	e := echo.New()
	e.Use(middleware.CORS())

	// 静的ファイルを置く場所
	e.Static("/", "./view/_static")

	// 以下、API通信を行う各パス
	// 各設定値の読み込み
	e.GET("/api/setting", cli.Read)
	// 各設定値の更新
	e.POST("/api/setting", cli.Update)
	// マウス座標確認
	e.GET("/api/mouse", cli.Confirm)
	// マウス座標取得
	e.POST("/api/mouse", cli.SetMouse)
	// 開始及び一時停止
	e.GET("/api/program", cli.Toggle)

	// ローカルAPIサーバーの開始
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
