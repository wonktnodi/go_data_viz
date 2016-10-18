package main

import (
    "gopkg.in/mgo.v2"
    "strings"
    "github.com/jinzhu/gorm"
    "fmt"
)

var (
    db_mysql      *gorm.DB
    sqlConnection string
)

func getMongoDBInstance() *mgo.Database {
    session, err := mgo.Dial(config.MongoDB)
    if err != nil {
        EXCEPTION(err)
    }
    // if MongoDBName == "" it will check the connection url MongoDB for a dbname
    // that logic inside mgo
    return session.DB(config.dbName)
}

// attempt to get dbName from URL
// it will work on MongoLab where dbName is part of url
func getDBName(url *string) string {
    arr := strings.Split(*url, ":")
    arr = strings.Split(arr[len(arr) - 1], "/")
    return arr[len(arr) - 1]
}

// count of documents in collection
func getCount(collection *mgo.Collection, c chan int, query interface{}) {
    count, _ := collection.Find(query).Count()
    c <- count
}


// DSN returns the Data Source Name
func DSN(ci MySQLInfo) string {
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

// get mysql database handle
func getMysqlDBInstance() *gorm.DB {
    if db_mysql != nil {
        return db_mysql
    }
    sqlConnection = DSN(new_config.Database.MySQL)
    var err error
    db_mysql, err = gorm.Open("mysql", sqlConnection)
    if err == nil {
        EXCEPTION(err)
    }
    db_mysql.LogMode(true)
    return db_mysql
}

