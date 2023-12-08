package trade

import (
	"fmt"
	"os"
	"time"
)

// Historical 秒足を保存
func (c *Client) Historical(isHistorical bool) {
	if !isHistorical {
		return
	}

	f, _ := os.OpenFile(fmt.Sprintf("./_data/ticker-%s.csv", c.Setting.Symbol), os.O_WRONLY|os.O_APPEND, 0644)
	defer f.Close()

	t := time.NewTicker(1 * time.Second)
	defer t.Stop()

L:
	for {
		select {
		case <-c.ctx.Done():
			break L

		case <-t.C:
			f.WriteString(fmt.Sprintf("%f,%f,%s\n", c.Tick.Ask, c.Tick.Bid, c.Tick.Timestamp.String()))

		}
	}
}
