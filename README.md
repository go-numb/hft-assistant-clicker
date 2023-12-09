# HFT Assistant Click Order開発者ドキュメント

## 概要
HFT Assistant Click Orderは、Oanda MT4/5、EA（Expert Advisor）と連携し、ボラティリティの監視や注文の執行(マウスコントロールによるクリック)を行うことができます。

当ドキュメントは開発者用メモを含むドキュメントです。  
プログラムの概要、インストールと操作のステップバイステップの詳細説明、トラブルシューティングのためのFAQ（よくある質問）セクションを提供します。

## システム要件
- WindowsOS
- Google Chrome（PORT参照時）またはMicrosoft Edge（WebView2使用時）
- Oanda MT4/5プラットフォームと内部EA（Expert Advisor）。
- CPU 3.0Hz、MEMORY 4GB（主にブラウザの要件として使用します。）
- インターネット接続
- PORT: 8080の使用法 ※マルチポート実装で自動PORT選択になりました。

## 一時ファイル配布のインストール手順
1. 付属のリンクからソースファイルをダウンロードします： [ダウンロード](https://github.com/go-numb/hft-assistant-clicker)
2. 任意の場所に配置します。
3. フォルダ内に `main.go` があることを確認する（_data は注文履歴保存用）。
4. ビルドする `$ go build` 
5. 5. `logic.exe` をダブルクリックし、アプリケーションウィンドウを開きます。開かない場合や設定画面が表示されない場合は、Chromeなどのブラウザを起動し、http://localhost:8080 にアクセスしてください。※マルチポート実装で自動PORT選択になりました。


## フォルダ内説明
HTF-ASSISTANT-CLICKER
├ _data/: デバッグ用のログファイルを保存するフォルダ
├ api/: GUIとの通信部分
├ key/: プログラム強制終了ショートカットキーの監視
├ publisher/: [重要] 配布時に触るファイル群を集約したフォルダ
│ ├ install-script.sh: 実行するとビルド環境を整備するシェルスクリプト
│ ├ PUBLISHER.md: 配布用説明書。ビルド環境や方法を記載。
│
├ setting/: 設定ファイルを取り扱うファイル群を集約したフォルダ
├ trade/: 発注部分を取り扱うファイル群を集約したフォルダ。詳細はソースコードに記載しています。
├ view/: GUI部分を取り扱うファイル群を集約したフォルダであり、静的なファイル群は_buildに整理済み。
├ main.go: 上記各部を組み合わせて司る起動またはビルド用ファイル
├ logic.exe: ビルドすると出力されるアプリファイル。_buildに移動させてください。
├ README.md: 開発者用メモ(非整形)
└ go.mod,go.sum: パッケージ管理情報


## 使用者設定
1. Oanda MT4/5を開き、_dataフォルダ内にあるEA `get_ticker_for_namepipe` を起動します（必要に応じてフォルダを移動してください）。
   1. logic.exeのコンソールログで現在の価格表示を確認します
2. 必要に応じて設定やオブジェクト登録をカスタマイズし、[実行]をクリックしてクリックアクセスを許可する。

## FAQ
Q: アプリケーションを実行する際に、ファイアウォールの設定は必要ですか？

A：WebView2アプリケーションはファイアウォールの設定を必要としませんが、WebView2を使用したアプリケーションがファイアウォールによってブロックされる場合があります。このような場合は、ファイアウォールの設定を確認し、アプリケーションが許可されていることを確認するか、ファイアウォールを一時的に無効にしてアプリケーションが問題なく実行されるかどうかをテストしてください。

Q: Windows11で "net user "でSIDが表示されない場合、どのようにしたら自分のSIDを見つけることができますか？

A: コマンドプロンプトまたはPowerShellで次のコマンドを使用して、SIDを見つけます：`wmic useraccount where "name='%username%'" get sid`.

## 参考

![デモ取引損益履歴プロット](https://github.com/go-numb/hft-assistant-clicker/blob/images/demo.jpg)

## Author

[@_numbP](https://twitter.com/_numbP)

## License

[MIT](https://github.com/go-numb/hft-assistant-clicker/blob/master/LICENSE)