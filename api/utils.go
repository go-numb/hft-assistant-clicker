package api

import (
	"fmt"
	"strings"
)

// _toMap マウス座標の更新パラメータをGUIが使いやすい型にする
func (cli *Client) _toMap(key string) map[string]any {
	return map[string]any{
		"msg": fmt.Sprintf("[SUCCESS] %s 登録完了", key),
		"area": map[string]any{
			"min_x": cli.Setting.Objects[key].Area.MinX(),
			"min_y": cli.Setting.Objects[key].Area.MinY(),
			"max_x": cli.Setting.Objects[key].Area.MaxX(),
			"max_y": cli.Setting.Objects[key].Area.MaxY(),
		},
		"random_pixel":       cli.Setting.Objects[key].RandomPixel,
		"loop_order_n":       cli.Setting.Objects[key].LoopOrderN,
		"loop_wait_millisec": cli.Setting.Objects[key].LoopWaitMillisec,
	}
}

// PipsSwitcher GUIで定義するpips（0.01=1pips）をUSDシンボルの場合、*0.01する
// プログラム定義値をUSDシンボルの場合、*100する
// GUIとAPIの通信（GetとPost）での使用
func (cli *Client) PipsSwitcher(fromGUI bool) {
	if !IsUSDoller(cli.Setting.Symbol) {
		return
	}

	// 対USDシンボル
	// GUIからのUSDシンボル値は *0.01
	if fromGUI {
		cli.Setting.Logic.Pips = cli.Setting.Logic.Pips * 0.01
		return
	}

	cli.Setting.Logic.Pips = cli.Setting.Logic.Pips * 100
}

func IsUSDoller(s string) bool {
	return strings.HasSuffix(strings.ToUpper(s), "USD")
}
