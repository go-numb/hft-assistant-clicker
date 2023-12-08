package keys

import (
	"fmt"

	"github.com/rs/zerolog/log"

	hook "github.com/robotn/gohook"
)

type Client struct{}

func New() *Client {
	return &Client{}
}

func (p *Client) Monitor() {
	// registor 監視対象のキーとショートカットを登録する
	registor()

	// monitor 常駐して登録したキーを監視する
	monitor()
}

func registor() {
	s := hook.Start()

	fmt.Println("[ctrl + shift + Q] を押すことでアプリが停止します")
	hook.Register(hook.KeyDown, []string{"q", "ctrl", "shift"}, func(e hook.Event) {
		log.Info().Msg("[ctrl-shift-q]が押されました、アプリを停止します")
		hook.End()
	})

	// fmt.Println("TEST: wを登録")
	// hook.Register(hook.KeyDown, []string{"w"}, func(e hook.Event) {
	// 	fmt.Println("wが押されました、登録された処理は有りません")
	// })

	<-hook.Process(s)
}

func monitor() {
	evChan := hook.Start()
	defer hook.End()

	for ev := range evChan {
		fmt.Println("hook: ", ev)
		if ev.Kind == hook.HookEnabled {
			break
		}
	}

	log.Fatal().Msg("停止コマンド")
}
