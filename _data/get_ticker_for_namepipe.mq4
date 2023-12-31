//+------------------------------------------------------------------+
//|                                                   get_ticker.mq4 |
//|                                         Copyright 2023, hft-assistant-clicker. |
//|                                             https://www.mql5.com |
//+------------------------------------------------------------------+
#property copyright "Copyright 2023, hft-assistant-clicker."
#property link      "https://www.mql5.com"
#property version   "1.00"
#property strict

// Document: https://lawn-tech.jp/mql4pipe_en.html

// 変数長指定
string pipe_name;
int pipe = INVALID_HANDLE;


string _sym = "";
string _temp = "";

MqlTick tick;

//+------------------------------------------------------------------+
//| Expert initialization function                                   |
//+------------------------------------------------------------------+
int OnInit()
  {
//---

  // connected file
  // read my symbol name
  _sym = StringSubstr(Symbol(),0,6);
  pipe_name = _sym;
  Print("connected file to ", pipe_name);
   
  pipe = FileOpen("\\\\.\\pipe\\" + _sym, FILE_WRITE | FILE_BIN | FILE_ANSI);

  // use millisec
  EventSetMillisecondTimer(1);

//---
   return(INIT_SUCCEEDED);
  }
//+------------------------------------------------------------------+
//| Expert deinitialization function                                 |
//+------------------------------------------------------------------+
void OnDeinit(const int reason)
  {
//---
  FileClose(pipe);
  EventKillTimer();
   
  }
//+------------------------------------------------------------------+
//| Expert tick function                                             |
//+------------------------------------------------------------------+
void OnTick()
  {
    //---

  // get ticker for ltp & volume & timestamp
  if (!SymbolInfoTick(_sym, tick)) {
    Print("Error SymbolInfoTick: ", GetLastError());
    return;
  }
    
    // create string for write file
  _temp = DoubleToString(tick.last) + "," +
          DoubleToString(Ask) + "," +
          DoubleToString(Bid) + "," +
          DoubleToString(tick.volume) + "," +
          IntegerToString(tick.time_msc);

  // check open file
  if ( pipe == INVALID_HANDLE ) {
    printf( "[INFO] line:%d, check open file: %s" , __LINE__ , pipe_name);
  } else {
    FileWriteString(pipe, _temp+"\r\n");
  }

  }
//+------------------------------------------------------------------+


