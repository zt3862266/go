package log

import (
	"testing"
	"log"
	"io/ioutil"
	"fmt"
	"encoding/json"
	. "rong360.com/framework/config"
)


func TestWrite(t *testing.T){

	for i:=0; i<2;i++ {
		log.Printf("hahaasdfasdfasdfasdfasdfasdfasdfhahaasdfasdfasdfasdfasdfasdfasdfhahaasdfasdfasdfasdfasdfasdfasdf")

		Info("this is a info log,this is a info log,this is a info log,this is a info log,this is a info log,this is a info log")

		Trace("this is a trace log this is a trace log this is a trace log this is a trace log this is a trace log this is a trace log this is a trace log")

		Warn("this is a warn log,this is a warn log,this is a warn log,this is a warn log,this is a warn log,this is a warn log")

		Error("this is a error log,this is a error log,this is a error log,this is a error log,this is a error log,this is a error log,")
	}

}

func TestJson(t *testing.T){

	bytes,err := ioutil.ReadFile("config/env.json")
	if err !=nil{
		fmt.Printf("read file error,%v",err)
	}
	var objEvn = &EnvConfig{}
	if err = json.Unmarshal(bytes,objEvn);err !=nil{
		fmt.Printf("json unmarshal error,%s",err)
	}
	fmt.Printf("%+v",objEvn.Mysql.Database.Slave)

}
