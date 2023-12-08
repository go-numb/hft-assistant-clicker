package setting

import (
	"math/rand"
	"time"

	"github.com/go-numb/go-mouse-click"
	"github.com/go-vgo/robotgo"
	"gonum.org/v1/gonum/stat"
)

var (
	s rand.Source
	r *rand.Rand
)

func init() {
	s = rand.NewSource(time.Now().UnixNano())
	r = rand.New(s)
}

func Intn(n int) int {
	if n <= 0 {
		n = 1
	}
	return r.Intn(n) - 1
}

// Target ターゲットオブジェクト
type Target struct {
	Name             string `toml:"name,omitempty" json:"name,omitempty"`
	RandomPixel      int    `toml:"random_pixel,omitempty" json:"random_pixel,omitempty"`
	LoopOrderN       int    `toml:"loop_order_n,omitempty" json:"loop_order_n,omitempty"`
	LoopWaitMillisec int    `toml:"loop_wait_millisec,omitempty" json:"loop_wait_millisec,omitempty"`
	// Restruct Button
	Buttons []*Button `toml:"button,omitempty" json:"button,omitempty"`
	// add Area
	Area *mouse.Corners `toml:"area,omitempty" json:"area,omitempty"`
}

type Button struct {
	Y int `toml:"y,omitempty" json:"y,omitempty"`
	X int `toml:"x,omitempty" json:"x,omitempty"`
}

func SetTarget() *mouse.Corners {
	return mouse.GetFourCorners()
}

// RestructOrder ターゲット位置にマウスクリックをシミュレート
func (t *Target) RestructOrder() {
	for i := 0; i < t.LoopOrderN; i++ {
		// mouse.Click(t.X, t.Y, t.RandomPixel)
		time.Sleep(time.Duration(t.LoopWaitMillisec) * time.Millisecond)
	}
}

// Order ターゲット位置にマウスクリックをシミュレート
func (t *Target) Order() {
	var (
		diff_x = t.Area.MaxX() - t.Area.MinX()
		diff_y = t.Area.MaxY() - t.Area.MinY()
		rx, ry int
	)

	// 決済注文
	for i := 0; i < t.LoopOrderN; i++ {
		rx = Intn(diff_x)
		ry = Intn(diff_y)

		robotgo.Move(t.Area.MinX()+rx, t.Area.MinY()+ry)
		robotgo.Click()
		time.Sleep(time.Duration(t.LoopWaitMillisec) * time.Millisecond)
	}
}

// _objectsXYMean 平均のX座標とY座標を計算し、すべてのターゲットオブジェクトの平均位置を求める
func (s *Settings) _objectsXYMean() (ymean, xmean float64) {
	var (
		xs, ys []float64
	)

	s.m.RLock()

	xs = append(xs, 0.5*float64(s.Objects["entry_buy"].Area.MinX()+s.Objects["entry_buy"].Area.MaxX()))
	ys = append(ys, 0.5*float64(s.Objects["entry_buy"].Area.MinY()+s.Objects["entry_buy"].Area.MaxY()))

	xs = append(xs, 0.5*float64(s.Objects["entry_sell"].Area.MinX()+s.Objects["entry_sell"].Area.MaxX()))
	ys = append(ys, 0.5*float64(s.Objects["entry_sell"].Area.MinY()+s.Objects["entry_sell"].Area.MaxY()))

	xs = append(xs, 0.5*float64(s.Objects["exit"].Area.MinX()+s.Objects["exit"].Area.MaxX()))
	ys = append(ys, 0.5*float64(s.Objects["exit"].Area.MinY()+s.Objects["exit"].Area.MaxY()))

	s.m.RUnlock()

	return stat.Mean(ys, nil), stat.Mean(xs, nil)
}

// MouseUndo すべてのターゲットオブジェクトの平均位置にマウスを移動
func (s *Settings) MouseUndo() {
	ymean, xmean := s._objectsXYMean()

	robotgo.Move(int(xmean), int(ymean))
}

// MouseMovingSmooth すべてのターゲットオブジェクトの平均位置の近くにランダムにマウスを移動する機能
func (s *Settings) MouseMovingSmooth(isRandom bool, randomRange int) {
	ymean, xmean := s._objectsXYMean()

	n := time.Now().UnixNano()
	x := Intn(randomRange)
	y := Intn(randomRange)

	high := r.Float64()
	delay := r.Float64() * float64(r.Intn(3))

	f := func(xy int) int {
		if n%2 == 0 {
			return -xy
		}
		return xy
	}

	robotgo.MoveSmooth(int(xmean)+f(x), int(ymean)+f(y), high, delay)
}

// _delay 指定されたミリ秒だけ実行を一時停止するためのヘルパー関数
func _delay(millisec int) {
	time.Sleep(time.Duration(millisec) * time.Millisecond)
}
