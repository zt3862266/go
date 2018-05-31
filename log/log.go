package log

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"log"
	"fmt"
)

var(
	Log *Logger = nil
)

type Logger struct{

	Trace	*log.Logger
	Info	*log.Logger
	Warn	*log.Logger
	Error	*log.Logger
}

func SetRongLogFile(file string){
	if Log != nil{
		return
	}
	Log = &Logger{
		Trace: &log.Logger{},
		Info: &log.Logger{},
		Warn: &log.Logger{},
		Error: &log.Logger{},
	}
	Log.Trace.SetOutput(
		&lumberjack.Logger{
			Filename:file+".trace",
			MaxSize:1024,
			LocalTime:true,
			Compress:true,
		})
	Log.Info.SetOutput(
		&lumberjack.Logger{
			Filename:  file + ".info",
			MaxSize:   1024,
			LocalTime: true,
			Compress:  true,
		})

	Log.Warn.SetOutput(
		&lumberjack.Logger{
			Filename:  file + ".warn",
			MaxSize:   1024,
			LocalTime: true,
			Compress:  true,
		})
	Log.Error.SetOutput(
		&lumberjack.Logger{
			Filename:  file + ".error",
			MaxSize:   1024,
			LocalTime: true,
			Compress:  true,
		})
	Log.Trace.SetFlags(log.LstdFlags | log.Lshortfile)
	Log.Info.SetFlags(log.LstdFlags | log.Lshortfile)
	Log.Warn.SetFlags(log.LstdFlags | log.Lshortfile)
	Log.Error.SetFlags(log.LstdFlags | log.Lshortfile)

	//默认输出到 info 日志
	log.SetOutput(&lumberjack.Logger{
		Filename:  file + ".info",
		MaxSize:   1024,
		LocalTime: true,
		Compress:  true,
	})
	log.SetFlags(log.LstdFlags | log.Lshortfile)

}


func Trace(format string,v ...interface{}){
	Log.Trace.Output(2,fmt.Sprintf(format,v...))
}

func Info(format string, v ...interface{}) {
	Log.Info.Output(2,fmt.Sprintf(format,v...))
}

func Warn(format string, v ...interface{}) {
	Log.Warn.Output(2,fmt.Sprintf(format,v...))
}
func Error(format string, v ...interface{}) {
	Log.Error.Output(2,fmt.Sprintf(format,v...))
}