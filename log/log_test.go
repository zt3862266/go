package log

import (
	"log"
	"testing"
)

func TestWrite(t *testing.T) {

	for i := 0; i < 2; i++ {
		log.Printf("hahaasdfasdfasdfasdfasdfasdfasdfhahaasdfasdfasdfasdfasdfasdfasdfhahaasdfasdfasdfasdfasdfasdfasdf")

		Info("this is a info log,this is a info log,this is a info log,this is a info log,this is a info log,this is a info log")

		Trace("this is a trace log this is a trace log this is a trace log this is a trace log this is a trace log this is a trace log this is a trace log")

		Warn("this is a warn log,this is a warn log,this is a warn log,this is a warn log,this is a warn log,this is a warn log")

		Error("this is a error log,this is a error log,this is a error log,this is a error log,this is a error log,this is a error log,")
	}

}
