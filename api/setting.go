package api

import (
	"fmt"
	"hft-assistant-clicker/setting"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/labstack/echo/v4"
)

// Toggle APIリクエストを受け、開始及び一時停止をトグルする
func (cli *Client) Toggle(c echo.Context) error {
	param := c.QueryParam("q")
	fmt.Println(param)

	changed := "[INFO] 変更なし"

	if param == "start" && !cli.Setting.IsTrade() {
		// 各項目が過不足なく設定されているかをチェックして、不十分ならば不足項目をGUIに返す
		if descript, isOK := cli.Setting.StartChecker(); isOK {
			changed = "[SUCCESS] 取引開始"
			cli.Setting.UpdateIsTrade(true)
		} else {
			changed = descript
		}
	} else if param == "stop" && cli.Setting.IsTrade() {
		changed = "[SUCCESS] 取引一時停止"
		cli.Setting.UpdateIsTrade(false)
	} else {
		if cli.Setting.IsTrade() {
			changed += "、現在取引中"
		} else {
			changed += "、現在停止中"
		}
	}

	// 変更が行われた設定の情報をレスポンスとして返します。
	return c.JSON(http.StatusOK, changed)
}

// SetMouse マウスの設定を変更するためのエンドポイント
func (cli *Client) SetMouse(c echo.Context) error {
	q := c.QueryParam("q")

	// ClientからのPostRequestを構造体にBind
	var params = make(map[string]*setting.Target)
	if err := c.Bind(&params); err != nil {
		return err
	}

	// 初期値設定
	changed := map[string]any{
		"msg": "[INFO] 変更なし",
		"y":   1,
		"x":   1,
	}

	// クライアント側で変更された設定値を現在の設定に反映
	if strings.ToLower(q) == "entry_buy" {
		cli.Setting.Objects["entry_buy"].Area = setting.SetTarget()
		cli.Setting.Objects["entry_buy"].RandomPixel = params["entry_buy"].RandomPixel
		cli.Setting.Objects["entry_buy"].LoopOrderN = params["entry_buy"].LoopOrderN
		cli.Setting.Objects["entry_buy"].LoopWaitMillisec = params["entry_buy"].LoopWaitMillisec

		changed = cli._toMap("entry_buy")
	} else if strings.ToLower(q) == "entry_sell" {
		cli.Setting.Objects["entry_sell"].Area = setting.SetTarget()
		cli.Setting.Objects["entry_sell"].RandomPixel = params["entry_sell"].RandomPixel
		cli.Setting.Objects["entry_sell"].LoopOrderN = params["entry_sell"].LoopOrderN
		cli.Setting.Objects["entry_sell"].LoopWaitMillisec = params["entry_sell"].LoopWaitMillisec

		changed = cli._toMap("entry_sell")
	} else if strings.ToLower(q) == "exit" {
		cli.Setting.Objects["exit"].Area = setting.SetTarget()
		cli.Setting.Objects["exit"].RandomPixel = params["exit"].RandomPixel
		cli.Setting.Objects["exit"].LoopOrderN = params["exit"].LoopOrderN
		cli.Setting.Objects["exit"].LoopWaitMillisec = params["exit"].LoopWaitMillisec

		changed = cli._toMap("exit")
	}

	log.Info().Msgf("%v", changed)

	// 変更が行われた設定の情報をレスポンスとして返します。
	return c.JSON(http.StatusOK, changed)
}

// Read 現在の設定値を返す
func (cli *Client) Read(c echo.Context) error {
	log.Info().Msgf("comming, %v", cli.Setting.Logic)

	// USDシンボルのpipsをGUI用に100倍する
	cli.PipsSwitcher(false)

	return c.JSON(http.StatusOK, cli.Setting)
}

// Update HTTPリクエストのjsonデータから設定情報を更新するために使用されます。まず、jsonデータをmap[string]anyの型にバインドして、設定情報を更新するために使用されます。その後、更新された設定情報の詳細がロギングされます。最後に、HTTPレスポンスで更新されたデータが返されます。
func (cli *Client) Update(c echo.Context) error {
	// 型へのBindではなく、jsonで扱う
	// お客様は数値を抽象化したカテゴリから選択し、そのパターンを数値化してプログラムに設定する
	v := make(map[string]any)
	if err := c.Bind(&v); err != nil {
		return err
	}

	var isPipsChanged bool

	// Update settings

	// ファイルパス設定は、設定のシンボルにて生成し代用したフためリストラ
	// val, isThere := v["filename"]
	// if isThere && cli.Setting.Filename != val.(string) {
	// 	cli.Setting.Filename = val.(string)
	// }

	val, isThere := v["symbol"]
	if isThere && cli.Setting.Symbol != val.(string) {
		// EAとの通信を行うPIPEファイルパスを生成する
		cli.Setting.Filename = fmt.Sprintf(`\\.\pipe\%s`, val.(string))
		cli.Setting.Symbol = val.(string)
	}

	// // TODO: Mouse options
	// 約定率は上がるがマウスコントロールを保持する時間が長く、回避する事項が多く実装するにはコスト高いためリストラ
	// 実行コストに見合う約定率向上が見込めるため、中心極限定理と平均回帰の観点から実装が好ましいと考えられる
	// val, isThere = v["is_mouse_center"]
	// if isThere && cli.Setting.IsMouseCenter != val.(bool) {
	// 	cli.Setting.IsMouseCenter = val.(bool)
	// }
	// val, isThere = v["is_move_object"]
	// if isThere && cli.Setting.IsMoveObject != val.(bool) {
	// 	cli.Setting.IsMoveObject = val.(bool)
	// }

	val, isThere = v["logic"]
	if isThere {
		logic := val.(map[string]any)

		val, isThere := logic["type"]
		if isThere {
			switch value := val.(type) {
			case float64: // APIからのデータをもったクライアントがそのままお繰り返している場合
				// Through

			case int:
				switch value {
				case 0:
					cli.Setting.Logic.Type = setting.Normal
				case 1:
					cli.Setting.Logic.Type = setting.OnlyEntryBuy
				case 2:
					cli.Setting.Logic.Type = setting.OnlyEntrySell
				case 3:
					cli.Setting.Logic.Type = setting.OnlyExit
				case 4:
					cli.Setting.Logic.Type = setting.OnlyExitBuy
				case 5:
					cli.Setting.Logic.Type = setting.OnlyExitSell
				}

			case string:
				switch value {
				case "0":
					cli.Setting.Logic.Type = setting.Normal
				case "1":
					cli.Setting.Logic.Type = setting.OnlyEntryBuy
				case "2":
					cli.Setting.Logic.Type = setting.OnlyEntrySell
				case "3":
					cli.Setting.Logic.Type = setting.OnlyExit
				case "4":
					cli.Setting.Logic.Type = setting.OnlyExitBuy
				case "5":
					cli.Setting.Logic.Type = setting.OnlyExitSell
				}
			}
		}

		val, isThere = logic["term_millisec"]
		if isThere {
			switch value := val.(type) {
			case int:
				if cli.Setting.Logic.TermMillisec != value {
					cli.Setting.Logic.TermMillisec = value
				}

			case float64:
				if cli.Setting.Logic.TermMillisec != int(value) {
					cli.Setting.Logic.TermMillisec = int(value)
				}

			case string:
				switch value {
				case "0":
					cli.Setting.Logic.TermMillisec = 1
				case "1":
					cli.Setting.Logic.TermMillisec = 10
				case "2":
					cli.Setting.Logic.TermMillisec = 100
				case "3":
					cli.Setting.Logic.TermMillisec = 1000
				}
			}
		}

		val, isThere = logic["pips"]
		if isThere {
			// TODO: Pipsの単位 例）USDJPYとEURUSDなどの違い
			switch value := val.(type) {
			case float64: // APIからのデータをもったクライアントがそのまま送り返している場合
				// Through
				if cli.Setting.Logic.Pips != value {
					cli.Setting.Logic.Pips = value
					isPipsChanged = true
				}

			case string:
				switch value {
				case "0":
					cli.Setting.Logic.Pips = 0.01
					isPipsChanged = true
				case "1":
					cli.Setting.Logic.Pips = 0.05
					isPipsChanged = true
				case "2":
					cli.Setting.Logic.Pips = 0.1
					isPipsChanged = true
				}

			}
		}

		val, isThere = logic["interval"]
		if isThere {
			switch value := val.(type) {
			case float64: // APIからのデータをもったクライアントがそのままお繰り返している場合
				if cli.Setting.Logic.Interval != int(value) {
					cli.Setting.Logic.Interval = int(value)
				}

			case int:
				cli.Setting.Logic.Interval = value
			}
		}

		val, isThere = logic["interval_is_random"]
		if isThere && cli.Setting.Logic.IntervalIsRandom != val.(bool) {
			cli.Setting.Logic.IntervalIsRandom = val.(bool)
		}
	}

	if isPipsChanged {
		// USDシンボルのpipsをGUI用に100倍する
		cli.PipsSwitcher(true)
	}

	log.Info().Msgf("setting: %+v", cli.Setting)
	log.Info().Msgf("logic: %+v", cli.Setting.Logic)

	return c.JSON(http.StatusOK, v)
}
