package loglib

import (
	"github.com/yangyouwei/ws-tcp-monitor/configlib"
	"log"
	"os"
)

var logFile *os.File
var Logtofile *log.Logger

func InitLog(conf_vale *configlib.MainConf)  {
	if conf_vale.Log {
		//输出日志到文件
		logFile, err := os.OpenFile(conf_vale.LogFile, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
		if err != nil {
			log.Fatal(err)
		}
		Logtofile = log.New(logFile, "[ws-tcp monitor] ", log.Ldate|log.Ltime|log.LstdFlags)
	}else {
		//输出日志到标准输出
		Logtofile = log.New(os.Stdout, "[ws-tcp monitor] ", log.Ldate|log.Ltime|log.LstdFlags)
	}
}
