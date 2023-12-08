package main

import (
	"fmt"
	"hft-assistant-clicker/api"
	"hft-assistant-clicker/keys"
	"hft-assistant-clicker/setting"
	"hft-assistant-clicker/trade"
	"hft-assistant-clicker/view"
	"net"

	"github.com/go-numb/go-mouse-click"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"golang.org/x/net/context"
)

const (
	// 配布用チェックの有無
	IsProduction = false

	// [TEST&DEBUG]秒足（BidAskのみ）を制作する
	IsHistorical = false

	// PORT = 8080
	PORT = 0 // ポート自動割当

	ENTRYBUY  = "entry_buy"
	ENTRYSELL = "entry_sell"
	EXIT      = "exit"
)

var (
	tmpSet      *setting.Settings
	OBJECTLABEL = []string{ENTRYBUY, ENTRYSELL, EXIT}
)

func init() {
	// 設定を読み込む
	tmpSet = setting.New()

	// ログ出力をUnixMillisecondに
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	// ログレベルをINFOに
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	// プロダクションモードでは、AppFileがなく、設定ファイルがない状態を見込む
	if IsProduction {
		// ログレベルをERRORから
		// 不必要なログを排除
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	}

	// 買いと売りのボタンが存在しない場合、初期値を登録する
	if len(tmpSet.Objects) == 0 {
		for i := 0; i < len(OBJECTLABEL); i++ {
			tmpSet.Objects[OBJECTLABEL[i]].Buttons = []*setting.Button{}
			tmpSet.Objects[OBJECTLABEL[i]].Area = &mouse.Corners{
				X1: 0,
				Y1: 0,
				X2: 0,
				Y2: 0,
			}
			tmpSet.Objects[OBJECTLABEL[i]].RandomPixel = 1
			tmpSet.Objects[OBJECTLABEL[i]].LoopOrderN = 1
			tmpSet.Objects[OBJECTLABEL[i]].LoopWaitMillisec = 0
		}

		// ExitOrder retry
		tmpSet.Objects[EXIT].RandomPixel = 3
		tmpSet.Objects[EXIT].LoopOrderN = 10
		tmpSet.Objects[EXIT].LoopWaitMillisec = 1000
	}
}

// アプリ起動時常に稼働する関数
// 設定値を保持し、各セクションを統括する
func main() {
	// 設定ファイルを保持する
	client := &api.Client{
		Setting: tmpSet,
	}
	// アプリを閉じたとき設定ファイルを保存する
	defer client.Setting.Close()

	log.Info().Msgf("%+v", client.Setting)

	// GUIと通信するAPI
	port, err := getAvailablePort(PORT)
	if err != nil {
		log.Fatal().Msg(err.Error())
	}
	go client.Routers(port)

	// 子関数への停止伝播
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 発注及び決済スレッド
	agent := trade.New(ctx, client.Setting)
	go agent.Trading()

	// 発注情報を保存するオプション
	go agent.Save()
	// 秒足を保存する
	go agent.Historical(IsHistorical)

	// ADD
	monitor := keys.New()
	go monitor.Monitor()

	// GUIを立ち上げ、GUIが閉じられればmain関数が終了する
	// 以下GUI関数が終了すれば、Contextで子関数へ終了が伝搬される
	// 伴い、設定ファイルが保存される
	view.WebView(IsProduction, port)
}

// getAvailablePort port変数":0"はOSに空いているポートを自動的に割り当てる
func getAvailablePort(port int) (int, error) {
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return 0, err
	}
	// 取得した空きPORTを開放する
	defer ln.Close()

	// 接続可能だったPORTを取得して返す
	port = ln.Addr().(*net.TCPAddr).Port
	return port, nil
}
