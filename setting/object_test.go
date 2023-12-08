package setting

import (
	"fmt"
	"testing"
	"time"

	"github.com/go-vgo/robotgo"
)

func TestMove(t *testing.T) {
	s := New()
	s.MouseMovingSmooth(true, 30)

	for i := 0; i < 3; i++ {
		s.MouseMovingSmooth(true, 100)
		time.Sleep(2 * time.Second)
	}
}

// TestSpeed 乱数生成の速度とマウスコントロールの速度を計測する
func TestSpeed(t *testing.T) {
	s := New()
	s.MouseMovingSmooth(true, 30)

	var (
		start time.Time
		temp  float64
	)

	for i := 0; i < 300; i++ {
		start = time.Now()
		_ = Intn(1000)
		temp = time.Since(start).Seconds()

		start = time.Now()
		robotgo.MoveClick(1000, 500)
		fmt.Printf("%f,%f\n", temp, time.Since(start).Seconds())

	}
}
