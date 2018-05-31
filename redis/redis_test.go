package redis

import (
	"testing"
	. "rong360.com/framework/config"
	. "rong360.com/framework/log"
)

func TestSet(t *testing.T){

	LoadEnv()
	SetRongLogFile(GlobalEnv.Log.File)
	InitRedis()
	rc := new(RongCache)
	rc.Set("zhangtao","tao",600)
	val,err := rc.Get("zhangtao")
	if err != nil{
		t.Errorf("get error:%v",err)
	}
	if val != "tao" {
		t.Errorf("wrong  value:%v",val)
	}

}
