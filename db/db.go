package db

import(
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"sync"
	"xg/conf"
	"xg/log"
)

var (
	globalDB *gorm.DB
	_dbOnce sync.Once
)

func Get()*gorm.DB{
	_dbOnce.Do(func() {
		db, err := gorm.Open("mysql", conf.Get().DBConnectionString)
		if err != nil {
			panic(err)
		}
		db.SetLogger(log.Trace)
		db.LogMode(true)
		globalDB = db
	})

	return globalDB
}