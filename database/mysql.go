package database

import (
	"database/sql"
	"errors"
	"github.com/didi/gendry/manager"
	"github.com/didi/gendry/scanner"
	_ "github.com/go-sql-driver/mysql"
	. "github.com/zt3862266/go/config"
	. "github.com/zt3862266/go/log"
	"math/rand"
	"sync/atomic"
	"time"
)

const (
	mysqlPingIntervalInSecond = 5
	mysqlPingFailedRetry      = 3
	mysqlStatusOk             = iota
	mysqlStatusError
)

var (
	db            *rdb
	badSlaveCount int64
)

//slave 结构体
type sdb struct {
	Ndb    *sql.DB //原始的数据库连接池
	Status uint    //该连接池状态
}

type rdb struct {
	Master *sql.DB
	Slave  []sdb
}

func healthCheck() {
	for {
		for _, slave := range db.Slave {
			err := slave.Ndb.Ping()
			if err != nil {
				Error("ping failed,db:%#v", slave)
			}
			if err == nil {
				//若成功,则恢复
				if slave.Status == mysqlStatusError {
					slave.Status = mysqlStatusOk
					Info("slave ok,recover: %#v", slave)
					atomic.AddInt64(&badSlaveCount, -1)
				}
			} else {
				//若失败超过一台,不处理
				if badSlaveCount > 0 {
					continue
				}
				isRecover := false
				for i := 0; i < mysqlPingFailedRetry-1; i++ {
					err = slave.Ndb.Ping()
					if err == nil {
						isRecover = true
						break
					}
					time.Sleep(time.Second * mysqlPingIntervalInSecond)
				}
				if isRecover == false && slave.Status == mysqlStatusOk {
					slave.Status = mysqlStatusError
					Error("slave ping error,remove it. %#v", slave)
					atomic.AddInt64(&badSlaveCount, 1)
				}

			}
		}
		time.Sleep(time.Second * mysqlPingIntervalInSecond)
	}
}

//初始化主从库的连接池
func InitDbPool() {
	if db == nil {
		db = &rdb{}
	} else {
		return
	}
	scanner.SetTagName("json")
	commonConf := GlobalEnv.Mysql
	dbMasterConf := commonConf.Database.Master
	dbSlaveConf := GlobalEnv.Mysql.Database.Slave

	//init master
	dbCommonConf := GlobalEnv.Mysql
	opendb, err := manager.New(dbMasterConf.Dbname, dbMasterConf.Dbuser, dbMasterConf.Dbpass, dbMasterConf.Host).Set(
		manager.SetCharset("utf8"),
		manager.SetAllowCleartextPasswords(true),
		manager.SetInterpolateParams(true),
		manager.SetParseTime(true),
		manager.SetTimeout(time.Duration(dbCommonConf.ConnnectTimeout)*time.Second),
		manager.SetReadTimeout(time.Duration(dbCommonConf.ReadTimeout)*time.Second),
		manager.SetWriteTimeout(time.Duration(dbCommonConf.WriteTimeout)*time.Second),
	).Port(dbMasterConf.Port).Open(true)
	if err != nil {
		Error("open  mysql master error: %#v", err)
		return
	}
	opendb.SetMaxIdleConns(commonConf.MaxIdleConnections)
	opendb.SetMaxOpenConns(commonConf.MaxOpenConections)
	opendb.SetConnMaxLifetime(time.Hour)
	db.Master = opendb

	for _, slave := range dbSlaveConf {
		slavedb, err := manager.New(slave.Dbname, slave.Dbuser, slave.Dbpass, slave.Host).Set(
			manager.SetCharset("utf8"),
			manager.SetAllowCleartextPasswords(true),
			manager.SetInterpolateParams(true),
			manager.SetParseTime(true),
			manager.SetTimeout(time.Duration(dbCommonConf.ConnnectTimeout)*time.Second),
			manager.SetReadTimeout(time.Duration(dbCommonConf.ReadTimeout)*time.Second),
			manager.SetWriteTimeout(time.Duration(dbCommonConf.WriteTimeout)*time.Second),
		).Port(slave.Port).Open(true)

		slavedb.SetMaxIdleConns(commonConf.MaxIdleConnections)
		slavedb.SetMaxOpenConns(commonConf.MaxOpenConections)
		slavedb.SetConnMaxLifetime(time.Hour)

		if err != nil {
			Error("open mysql slave error: %$v,config:%#v", err, slave)
			db.Slave = append(db.Slave, sdb{
				Ndb:    slavedb,
				Status: mysqlStatusError,
			})
			atomic.AddInt64(&badSlaveCount, 1)
		} else {
			db.Slave = append(db.Slave, sdb{
				Ndb:    slavedb,
				Status: mysqlStatusOk,
			})
		}
	}
	Info("master:%#v,slave:%#v", db.Master, db.Slave)
	go healthCheck()

}

//获取从数据库连接池
func SDB() (pool *sql.DB, err error) {

	slaveLen := len(db.Slave)
	if slaveLen == 0 {
		return nil, errors.New("no valid slave")
	}
	chooseIdx := rand.Intn(slaveLen)
	for {
		if db.Slave[chooseIdx].Status == mysqlStatusOk {
			return db.Slave[chooseIdx].Ndb, nil
		} else {
			chooseIdx = (chooseIdx - 1) % slaveLen
		}
	}
}

func MDB() (pool *sql.DB, err error) {
	return db.Master, nil
}
