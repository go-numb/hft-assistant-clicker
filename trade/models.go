package trade

import (
	"context"
	"hft-assistant-clicker/setting"
	"time"

	"github.com/go-numb/go-pipe-for-mt4"
	"github.com/go-numb/go-trade-utilities/controllers"
	"github.com/rs/zerolog/log"
)

type Client struct {
	ctx     context.Context
	Setting *setting.Settings

	isOneOrderOneKillPermission bool
	Tick                        *pipe.Ticker
	Controller                  *controllers.CentralController
}

func New(ctx context.Context, setdata *setting.Settings) *Client {
	return &Client{
		ctx:     ctx,
		Setting: setdata,
		Tick:    &pipe.Ticker{},

		isOneOrderOneKillPermission: true,
		Controller:                  controllers.New(setdata.Logic.TermMillisec, 10000, setdata.Logic.Pips),
	}
}

func (c *Client) Trading() {
	ch := make(chan []byte, 1024)

	// 配信受け取り
	go c.Read(ch)

	// order section
	go c.Entry(ch)
	go c.Exit()

	// 値が更新されていればSetting変数を変更しBreak
	go c.CheckRead()

	// マウスオプションはc.Setting.IsMoveObject=True時のみ動作
	// マウスオプション MoveSmooth
	go c.MouseOptions()
	// マウスオプション MouseControl Toggle [KEY ESC]
	go c.MouseControl()

	<-c.ctx.Done()
	log.Info().Msg("trade ended")
}

// Read EA更新値を読み込む
func (c *Client) Read(ch chan []byte) {
	// Read OneLine
	if err := pipe.Pipe(
		c.ctx,
		c.Setting.Filename,
		c.Setting.TermMillisec,
		ch,
	); err != nil {
		log.Err(err).Msg("")
	}
}

// CheckRead 値の取得が行えたことをGUIと共有
func (c *Client) CheckRead() {
	for {
		if c.Tick.Bid != c.Tick.Ask {
			log.Info().Msg("update value")
			c.Setting.IsGetValue = true
			break
		}
		time.Sleep(time.Second)
	}
}

// One order
func (c *Client) One() {
	c.isOneOrderOneKillPermission = false
}

// Kill no position
func (c *Client) Kill() {
	c.isOneOrderOneKillPermission = true
}

// NoPosition
func (c *Client) NoPosition() bool {
	return c.isOneOrderOneKillPermission
}
