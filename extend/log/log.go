package log

import (
	"gateway/extend/conf"
	"log"
	"os"
	"time"
)

var DebugLog *log.Logger

func Init() {
	ym := time.Now().Format("2006-01")
	day := time.Now().Format("02")
	fileName := conf.LogConf.LogPath + "storage/" + ym + "/" + day + "/" + level + ".log"

	logFile, logFileErr := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0544)
	if logFileErr != nil {
		log.Println(logFileErr)
	}
	defer logFile.Close()
	DebugLog = log.New(logFile, "【INOF】", log.Llongfile|log.Ltime|log.Ldate)
}
//func LogRecord(msg string, level string) error {
//	ym := time.Now().Format("2006-01")
//	day := time.Now().Format("02")
//	fileName := conf.LogConf.LogPath + "storage/" + ym + "/" + day + "/" + level + ".log"
//	fmt.Println(fileName)
//	logFile, logFileErr := os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0544)
//	if logFileErr != nil {
//		return logFileErr
//	}
//	defer logFile.Close()
//	log.SetOutput(file)
//	return nil
//}
