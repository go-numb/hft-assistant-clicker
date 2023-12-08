package trade

import (
	"fmt"
	"os"
	"time"

	"github.com/go-numb/go-trade-utilities/orders"
	"github.com/gocarina/gocsv"
	"github.com/rs/zerolog/log"
)

// Save ボラ検知履歴の保存
func (c *Client) Save() error {
	f, _ := os.Create(fmt.Sprintf("./_data/%s-ltp-%d-%f.csv", c.Setting.Symbol, c.Setting.Logic.TermMillisec, c.Setting.Logic.Pips))
	defer f.Close()

	if err := gocsv.Marshal([]orders.Order{
		{},
	}, f); err != nil {
		log.Err(err).Msg("orders write to initial file")
	}

	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		<-ticker.C

		buyl, selll := c.Controller.HistoricalController.Length()
		pos := make([]orders.Order, int(buyl)+int(selll))
		var i int
		c.Controller.HistoricalController.Orders.Range(func(key, value any) bool {
			pos[i] = value.(orders.Order)
			pos[i].OID = key.(string)
			i++

			c.Controller.HistoricalController.Delete(key)

			return true
		})

		if err := gocsv.MarshalWithoutHeaders(&pos, f); err != nil {
			log.Err(err).Msg("")
			continue
		}
		log.Info().Msg("log saved to temporary file")

	}
}
