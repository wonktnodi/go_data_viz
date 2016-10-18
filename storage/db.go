package storage

import (
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
    _ "github.com/jinzhu/gorm/dialects/mysql"
    "../model"
    "fmt"
    "../config"
)

var StorageDB *gorm.DB

// DSN returns the Data Source Name
func DSN(ci config.MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Name + ci.Parameter
}

func InitDB(cfg config.Configuration) (error) {
    sqlConnectionStr := DSN(cfg.Database.MySQL)
    db, err := gorm.Open("mysql", sqlConnectionStr)
    if err != nil {
        return err
    }
    db.LogMode(true)

    db.AutoMigrate(&model.User{}, &model.LoginAttempt{}, &model.ChartDataItem{})
    StorageDB = db

    NewUserStorage(StorageDB)
    return nil
}

func CloseDB() {
    if StorageDB != nil {
        StorageDB.Close()
        StorageDB = nil
    }
}
