package trade

/*
このパッケージの「Exit」関数は、定期的にエントリー注文が約定したかどうかをチェックし、指定された時間が経過していたら清算注文を出すために使用されます。

この関数では、GUIで指定された出口までの時間間隔を使用して、時間間隔を定期的に更新することができます。また、時間間隔をランダムに設定するオプションも提供されています。

注文が約定した場合、清算注文を出し、ログに注文情報を追加し、ポジションを削除して、トレードを終了します。また、マウスポインタをオブジェクトの中心に移動するための補助的な関数も含まれています。
*/

import (
	"fmt"
	"hft-assistant-clicker/setting"
	"math/rand"
	"time"

	"github.com/go-numb/go-trade-utilities/orders"
	"github.com/go-numb/go-trade-utilities/types"
	"github.com/rs/zerolog/log"
)

var interval time.Duration

// IfRandomInterval GUIでの変更でEntryToExit Intervalを変更する
func (c *Client) IfRandomInterval() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	s := rand.NewSource(time.Now().UnixNano())
	r := rand.New(s)

	for {
		<-ticker.C
		interval = time.Duration(c.Setting.Logic.Interval) * time.Second

		// RandomeInterval設定がなければReturn
		if !c.Setting.Logic.IntervalIsRandom {
			continue
		}
		interval = time.Duration(r.Intn(c.Setting.Logic.Interval)+1) * time.Second
		log.Info().Msgf("change exit interval to %.3fsec", interval.Seconds())
	}
}

// Exit 注文タイプ[Normal, OnlyEntryBuy,OnlyEntrySell]に適用
// [OnlyExit, OnlyExitBuy,OnlyExitSell]はentry.goで実行する
func (c *Client) Exit() {
	if c.Setting.Logic.Type >= setting.OnlyExit {
		return
	}

	interval = time.Duration(c.Setting.Logic.Interval) * time.Second

	go c.IfRandomInterval()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		<-ticker.C

		// One entry order, one key
		buypos := c.Controller.PositionController.Load(types.BUY)
		if buypos.Size != 0 {
			if time.Since(buypos.CreatedAt) > interval {
				// Fixed 清算
				// 必ず約定
				if c.NoPosition() {
					continue
				}

				if c.Setting.Logic.Type == setting.OnlyEntrySell {
					continue
				}

				c.Setting.Objects["exit"].Order()
				log.Info().Msg("return sell order")
				c.Controller.HistoricalController.Set(fmt.Sprintf("buyin-%d", time.Now().UnixNano()), orders.Order{
					ID:        buypos.ID,
					Symbol:    buypos.Symbol,
					Type:      "MARKET in",
					Side:      buypos.Side,
					Price:     buypos.Price,
					Size:      buypos.Size,
					CreatedAt: buypos.CreatedAt,
				})
				c.Controller.HistoricalController.Set(fmt.Sprintf("sellout-%d", time.Now().UnixNano()), orders.Order{
					ID:        0,
					Symbol:    c.Setting.Symbol,
					Type:      "MARKET out",
					Side:      types.SELL,
					Price:     c.Tick.Ltp,
					Size:      c.Setting.Logic.OneShotSize,
					CreatedAt: time.Now(),
				})
				c.Controller.PositionController.Delete(types.BUY)
				c.Kill()
				// マウス位置をオブジェクト中心に
				if c.Setting.IsMouseCenter {
					c.Setting.MouseUndo()
				}
			}
		}

		// One entry order, one key
		sellpos := c.Controller.PositionController.Load(types.SELL)
		if sellpos.Size != 0 {
			if time.Since(sellpos.CreatedAt) > interval {
				// Fixed 清算
				// 必ず約定
				if c.NoPosition() {
					continue
				}

				if c.Setting.Logic.Type == setting.OnlyEntryBuy {
					continue
				}

				c.Setting.Objects["exit"].Order()
				log.Info().Msg("return buy order")
				c.Controller.HistoricalController.Set(fmt.Sprintf("sellin-%d", time.Now().UnixNano()), orders.Order{
					ID:        sellpos.ID,
					Symbol:    sellpos.Symbol,
					Type:      "MARKET in",
					Side:      sellpos.Side,
					Price:     sellpos.Price,
					Size:      sellpos.Size,
					CreatedAt: sellpos.CreatedAt,
				})
				c.Controller.HistoricalController.Set(fmt.Sprintf("buyout-%d", time.Now().UnixNano()), orders.Order{
					ID:        0,
					Symbol:    c.Setting.Symbol,
					Type:      "MARKET out",
					Side:      types.BUY,
					Price:     c.Tick.Ltp,
					Size:      c.Setting.Logic.OneShotSize,
					CreatedAt: time.Now(),
				})
				c.Controller.PositionController.Delete(types.SELL)
				c.Kill()
				// マウス位置をオブジェクト中心に
				if c.Setting.IsMouseCenter {
					c.Setting.MouseUndo()
				}
			}
		}
	}
}
