package trade

/*
このパッケージは、簡単なトレードロジックを実行するためのものです。
この特定の関数「Entry」では、チャンネルから最新の価格情報を受信し、価格の変動を監視し、ボラティリティがある場合に注文を出します。

注文の種類は、ロジックの設定によって異なり、GUIからオン/オフされるトグルを使用して切り替えることができます。また、注文が出た際に、注文情報をログに記録するための仕組みも含まれています。さらに、マウスポインタをオブジェクトの中心に移動するための補助的な関数も含まれています。
*/

import (
	"hft-assistant-clicker/setting"
	"time"

	"github.com/go-numb/go-pipe-for-mt4"
	"github.com/go-numb/go-trade-utilities/positions"
	"github.com/go-numb/go-trade-utilities/pricing"
	"github.com/go-numb/go-trade-utilities/types"
	"github.com/go-vgo/robotgo"
	"github.com/rs/zerolog/log"
)

var (
	isVolatility bool
)

func (c *Client) Entry(ch chan []byte) {
	var (
		isPump int

		orderID int
	)

	for {
		// Oanda MT4 EAからPIPEで値を取得
		v := <-ch
		// Byteを型に変換
		c.Tick = pipe.ByteToTicker(v)
		// Update LTP
		// IsUpdateLTP(t.Ask, t.Bid)
		c.Tick.Ltp = (c.Tick.Ask + c.Tick.Bid) * 0.5
		c.Controller.PriceController.Set(pricing.Price{
			Ltp:       c.Tick.Ltp,
			Mid:       c.Tick.Ltp,
			Ask:       c.Tick.Ask,
			Bid:       c.Tick.Bid,
			Volume:    c.Tick.Volume,
			Timestamp: c.Tick.Timestamp.UnixMilli(),
		})
		log.Info().Msgf("%f, %f", c.Tick.Ask, c.Tick.Bid)

		// ボラティリティの検知
		isPump, isVolatility = c.Controller.PriceController.Volatility()
		if isVolatility {
			log.Info().Msgf("get volatility ltp: %.3f, ask: %.3f, bid: %.3f", c.Tick.Ltp, c.Tick.Ask, c.Tick.Bid)
			if !c.Setting.IsTrade() { // ONOFF Toggle from GUI as customer
				continue
			}
			if isPump == 1 {
				// Only entry sell
				if c.Setting.Logic.Type == setting.OnlyEntrySell ||
					c.Setting.Logic.Type == setting.OnlyExit ||
					c.Setting.Logic.Type == setting.OnlyExitSell ||
					c.Setting.Logic.Type == setting.OnlyExitBuy {
					// ボラ検出内での「決済のみ及び決済買いの実行」
					if c.Setting.Logic.Type == setting.OnlyExit ||
						c.Setting.Logic.Type == setting.OnlyExitBuy {
						// 注文をリトライ(リトライが終わるまで待つ)
						c.Setting.Objects["exit"].Order()
						log.Info().Msg("exit only buy order by volatility")
						// マウス位置をオブジェクト中心に
						if c.Setting.IsMouseCenter {
							c.Setting.MouseUndo()
						}
					}
					continue
				}

				// 建玉がなければ注文開始
				if !c.NoPosition() {
					continue
				}
				c.Setting.Objects["entry_buy"].Order()
				log.Info().Msg("entry buy order")
				// set log
				c.Controller.PositionController.Set(types.BUY, positions.Position{
					ID:        orderID,
					Symbol:    c.Setting.Symbol,
					Type:      "MARKET",
					Side:      types.BUY,
					Price:     c.Tick.Ltp,
					Size:      c.Setting.Logic.OneShotSize,
					CreatedAt: time.Now(),
				})
				orderID++
				c.One()
				// マウス位置をオブジェクト中心に
				if c.Setting.IsMouseCenter {
					c._waitMouseSwicher()
				}

			} else if isPump == -1 {
				// Only buy entry
				if c.Setting.Logic.Type == setting.OnlyEntryBuy ||
					c.Setting.Logic.Type == setting.OnlyExit ||
					c.Setting.Logic.Type == setting.OnlyExitBuy ||
					c.Setting.Logic.Type == setting.OnlyExitSell {
					// ボラ検出内での「決済のみ及び決済売りの実行」
					if c.Setting.Logic.Type == setting.OnlyExit ||
						c.Setting.Logic.Type == setting.OnlyExitSell {
						// 注文をリトライ(リトライが終わるまで待つ)
						c.Setting.Objects["exit"].Order()
						log.Info().Msg("exit only sell order by volatility")
						// マウス位置をオブジェクト中心に
						if c.Setting.IsMouseCenter {
							c.Setting.MouseUndo()
						}
					}
					continue
				}

				// 建玉がなければ注文開始
				if !c.NoPosition() {
					continue
				}
				c.Setting.Objects["entry_sell"].Order()
				log.Info().Msg("entry sell order")
				// set log
				c.Controller.PositionController.Set(types.SELL, positions.Position{
					ID:        orderID,
					Symbol:    c.Setting.Symbol,
					Type:      "MARKET",
					Side:      types.SELL,
					Price:     c.Tick.Ltp,
					Size:      c.Setting.Logic.OneShotSize,
					CreatedAt: time.Now(),
				})
				orderID++
				c.One()
				// マウス位置をオブジェクト中心に
				if c.Setting.IsMouseCenter {
					c._waitMouseSwicher()
				}
			}

		}
	}
}

// _waitMouseSwicher c.Setting.IsMouseCenter = trueならば、マウスを中央に
func (c *Client) _waitMouseSwicher() {
	var rn int
	if c.Setting.Logic.Type == setting.OnlyEntryBuy {
		rn = setting.Intn(len(c.Setting.Objects["entry_buy"].Buttons))
		robotgo.Move(int(c.Setting.Objects["entry_buy"].Buttons[rn].X), int(c.Setting.Objects["entry_buy"].Buttons[rn].Y))
		return
	} else if c.Setting.Logic.Type == setting.OnlyEntrySell {
		rn = setting.Intn(len(c.Setting.Objects["entry_sell"].Buttons))
		robotgo.Move(int(c.Setting.Objects["entry_sell"].Buttons[rn].X), int(c.Setting.Objects["entry_sell"].Buttons[rn].Y))
		return
	}
	c.Setting.MouseUndo()
}
