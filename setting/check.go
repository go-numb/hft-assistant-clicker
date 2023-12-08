package setting

// StartChekcer GUIから開始リクエスト時、稼働似必要な項目をチェックする関数
func (s *Settings) StartChecker() (descript string, isOK bool) {
	if s.Filename == "" {
		return "[ERROR] ファイルパスの設定がありません", false
	} else if s.TermMillisec < 1 {
		return "[ERROR] 設定値が小さすぎます", false
	} else if !s.IsGetValue {
		return "[ERROR] 価格値の得られません。Oanda及びEAの起動や設定をご確認ください", false
	}

	if s.Logic.TermMillisec < 1 {
		return "[ERROR] 設定値が小さすぎます", false
	}

	// check Mouse objects
	if s.Logic.Type != OnlyEntrySell {
		if s.Objects["entry_buy"].Area.X1 < 1 ||
			s.Objects["entry_buy"].Area.X2 < 1 ||
			s.Objects["entry_buy"].Area.Y1 < 1 ||
			s.Objects["entry_buy"].Area.Y2 < 1 {
			return "[ERROR] 設定値が初期状態または小さすぎます", false
		}
	}

	if s.Logic.Type != OnlyEntryBuy {
		if s.Objects["entry_sell"].Area.X1 < 1 ||
			s.Objects["entry_sell"].Area.X2 < 1 ||
			s.Objects["entry_sell"].Area.Y1 < 1 ||
			s.Objects["entry_sell"].Area.Y2 < 1 {
			return "[ERROR] 設定値が初期状態または小さすぎます", false
		}
	}

	if s.Objects["exit"].Area.X1 < 1 ||
		s.Objects["exit"].Area.X2 < 1 ||
		s.Objects["exit"].Area.Y1 < 1 ||
		s.Objects["exit"].Area.Y2 < 1 {
		return "[ERROR] 設定値が初期状態または小さすぎます", false
	}

	return "", true
}
