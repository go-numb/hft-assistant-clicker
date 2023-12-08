package trade

import (
	"hft-assistant-clicker/setting"
	"time"

	hook "github.com/robotn/gohook"
	"github.com/rs/zerolog/log"
)

const (
	MOUSEUNDOTERM int = 60
)

// MouseOptions マウス座標を中心付近にスムーズかつランダムに、ディレイありで戻す
// 設定の途中変更無効
func (c *Client) MouseOptions() {
	if !c.Setting.IsMoveObject {
		return
	}

	t := time.NewTicker(time.Duration(MOUSEUNDOTERM) * time.Second)
	defer t.Stop()

	for {
		<-t.C
		time.Sleep(time.Duration(setting.Intn(MOUSEUNDOTERM)) * time.Second)

		log.Info().Msg("move mouse to (center + randomXY)")
		if c.Setting.IsMouseCenter {
			c.Setting.MouseMovingSmooth(true, 20)
		}
	}

}

// MouseControl マウスコントロール（MouseOptions, _waitMouseSwicher）をスイッチする
func (c *Client) MouseControl() {
	if !c.Setting.IsMoveObject {
		return
	}

	for {
		ok := hook.AddEvents("esc")
		if ok {
			c.Setting.IsMouseCenter = !c.Setting.IsMouseCenter
			log.Info().Msgf("push ESC, toggle c.Setting.IsMouseCenter  %v\n", c.Setting.IsMouseCenter)
		}
	}
}
