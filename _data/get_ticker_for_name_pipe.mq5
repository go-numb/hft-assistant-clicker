
// Tickerの取得: MT5では、"Symbol()"関数を使用してアクティブなシンボル名を取得できます。価格は"CopyRates()"関数で取得します。

// 関数の名前変更: MT5では、初期化関数を"OnInit()"、開放関数を"OnDeinit()"、タイマー関数を"OnTimer()"、新しいバー関数を"OnTick()"として使用する必要があります。


//+------------------------------------------------------------------+
//| get_ticker.mq5 |
//| Copyright 2023, hft-assistant-clicker. |
//| https://www.mql5.com |
//+------------------------------------------------------------------+
#property copyright "Copyright 2023, hft-assistant-clicker."
#property link "https://www.mql5.com"
#property version "1.00"
#property strict

// Variable length specification
string pipe_name;
int pipe = INVALID_HANDLE;


string _sym = "";
string _temp = "";

MqlTick tick;

//+------------------------------------------------------------------+
//| Expert initialization function |
//+------------------------------------------------------------------+
int OnInit()
{
    //---
    // Connect to the file
    // Read my symbol name
    _sym = StringSubstr(Symbol(),0,6);
    pipe_name = _sym;
    Print("Connected to file ", pipe_name);

    pipe = FileOpen("\\\\.\\pipe\\" + _sym, FILE_WRITE | FILE_BIN | FILE_ANSI);

    // Use milliseconds
    EventSetMillisecondTimer(1);

    //---
    return(INIT_SUCCEEDED);
}
//+------------------------------------------------------------------+
//| Expert deinitialization function |
//+------------------------------------------------------------------+
void OnDeinit(const int reason)
{
    //---
    FileClose(pipe);
    EventKillTimer();

}
//+------------------------------------------------------------------+
//| Expert tick function |
//+------------------------------------------------------------------+
void OnTick()
{
    //-- Get ticker for the last price, volume and timestamp
    if (!SymbolInfoTick(_sym, tick)) {
        Print("Error SymbolInfoTick: ", GetLastError());
        return;
    }

    // Create a string for writing to the file
    _temp = DoubleToString(tick.last) + "," +
    DoubleToString(SymbolInfoDouble(_sym, SYMBOL_ASK)) + "," +
    DoubleToString(SymbolInfoDouble(_sym, SYMBOL_BID)) + "," +
    DoubleToString(tick.volume) + "," +
    IntegerToString(tick.time_msc);

    // Check if the file is open
    if ( pipe == INVALID_HANDLE ) {
        printf( "[INFO] line:%d, check open file: %s" , __LINE__ , pipe_name);
    } else {
        FileWriteString(pipe, _temp+"\r\n");
    }

}
//+------------------------------------------------------------------+

// このコードは、MT5で指定されたシンボルのBid,Ask値を取得し、そのデータをパイプに書き込むように設定されています