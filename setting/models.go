package setting

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/go-numb/go-mouse-click"
	"github.com/pelletier/go-toml/v2"
	"github.com/rs/zerolog/log"
)

const (
	TEMPFILE = "./logic.exe/d969831eb8a99cff8c02e681f43289e5d3d69664.toml"
)

type Settings struct {
	m          sync.RWMutex `toml:"-" json:"-"`
	isTrade    bool         `toml:"-" json:"-"`
	IsGetValue bool         `toml:"-" json:"-"`
	// EA Pipefile
	Filename string `toml:"filename" json:"filename"`
	// 値の取得（更新）間隔
	TermMillisec int `toml:"term_millisec" json:"term_millisec"`

	// TARGET Object to file
	IsMouseCenter    bool               `toml:"is_mouse_center" json:"is_mouse_center"`
	IsMoveObject     bool               `toml:"is_move_object" json:"is_move_object"`
	Objects          map[string]*Target `toml:"objects" json:"objects"`
	ClickRandomPixel int                `toml:"click_random_pixel" json:"click_random_pixel"`

	// Customer
	Symbol string `toml:"symbol" json:"symbol"`

	// Customer logic
	Logic *Logic `toml:"logic" json:"logic"`

	// Program Timeout Hour
	Timeout int `toml:"timeout" json:"timeout"`
}

type types int

const (
	Normal types = iota
	OnlyEntryBuy
	OnlyEntrySell
	OnlyExit
	OnlyExitBuy
	OnlyExitSell
)

type Logic struct {
	Type types `toml:"type" json:"type"`
	// Pips ボラティリティ検知用
	// 間隔
	TermMillisec int     `toml:"term_millisec" json:"term_millisec"`
	Pips         float64 `toml:"pips" json:"pips"`

	// Interval to ExitOrder
	// ランダムに変えたい、変えれるように
	IntervalIsRandom bool `toml:"interval_is_random" json:"interval_is_random"`
	// Interval Unit in Sec
	Interval int `toml:"interval" json:"interval"`

	// order
	OneShotSize float64 `toml:"one_shot_size" json:"one_shot_size"`
}

func New() *Settings {
	s := Read()
	if s != nil {
		return s
	}

	// 初期設定
	s = &Settings{
		isTrade:          false,
		Filename:         `\\.\pipe\USDJPY`,
		TermMillisec:     10,
		Objects:          make(map[string]*Target),
		ClickRandomPixel: 5,
		Symbol:           "USDJPY",
		Logic:            &Logic{Type: Normal, TermMillisec: 10, Pips: 0.01, IntervalIsRandom: false, Interval: 10, OneShotSize: 1},
		Timeout:          24,
	}

	str := []string{"entry_buy", "entry_sell", "exit"}
	var temp *Target
	for i := 0; i < len(str); i++ {
		if !strings.Contains(str[i], "exit") {
			temp = &Target{
				Name:    str[i],
				Buttons: make([]*Button, 0),
				Area: &mouse.Corners{
					X1: 0,
					Y1: 0,
					X2: 0,
					Y2: 0,
				},
				RandomPixel:      3,
				LoopOrderN:       1,
				LoopWaitMillisec: 0,
			}
		} else {
			temp = &Target{
				Name:    str[i],
				Buttons: make([]*Button, 0),
				Area: &mouse.Corners{
					X1: 0,
					Y1: 0,
					X2: 0,
					Y2: 0,
				},
				RandomPixel:      3,
				LoopOrderN:       10,
				LoopWaitMillisec: 1000,
			}
		}

		s.Objects[str[i]] = temp
	}

	return s
}

func (s *Settings) IsTrade() bool {
	s.m.RLock()
	defer s.m.RUnlock()

	return s.isTrade
}

func (s *Settings) UpdateIsTrade(b bool) {
	s.m.Lock()
	defer s.m.Unlock()

	s.isTrade = b
}

func (s *Settings) ToggleIsTrade() {
	s.m.Lock()
	defer s.m.Unlock()

	s.isTrade = !s.isTrade
}

func (p *Settings) Close() error {
	return p.Write()
}

// Write 設定ファイルの保存
func (p *Settings) Write() error {
	// Windows app temporaly dir
	filename := TempDir(TEMPFILE)

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	log.Info().Msg("save setting to file")
	return toml.NewEncoder(f).Encode(p)
}

// Read 設定ファイルの読み込み
func Read() *Settings {
	// Windows app temporaly dir
	filename := TempDir(TEMPFILE)

	f, err := os.Open(filename)
	if err != nil {
		log.Err(err).Msg("file open not find")
		return nil
	}
	defer f.Close()

	var s = new(Settings)
	if err := toml.NewDecoder(f).Decode(s); err != nil {
		log.Err(err).Msg("read setting file")
		return nil
	}

	for key, val := range s.Objects {
		for _, v := range val.Buttons {
			fmt.Printf("key:%s,  %d, %d\n", key, v.X, v.Y)
		}
	}

	return s
}

// TempDir Windows app temporaly dir
func TempDir(filename string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return ""
	}

	fullpath := filepath.Join(dir, filename)

	return fullpath
}
