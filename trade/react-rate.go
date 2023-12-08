package trade

import (
	"fmt"
	"os"
	"time"
)

const (
	FILENAME = "./_data/reacts01.csv"
)

var f *os.File

type Performance struct {
	start time.Time
}

func NewPerformance() *Performance {
	f, _ = os.Create(FILENAME)
	return &Performance{}
}

func ClosePerformance() {
	f.Close()
}

// Start から End()までの時間を計測します
func (p *Performance) Start() {
	p.start = time.Now()
}

func (p *Performance) End() {
	d := time.Since(p.start).Seconds()
	f.WriteString(fmt.Sprintf("%f\n", d))
}
