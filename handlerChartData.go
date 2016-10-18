package main

import (
    "github.com/gin-gonic/gin"
    "encoding/json"
    "net/http"
    "math/rand"
    "gopkg.in/mgo.v2"
    "fmt"
    "./storage"
    "time"
    "log"
)

type ChartDataQuery struct {
    Type  int64    `json:"type"`
    Area  int      `json:"area"`
    Point int      `json:"point"`
    Begin string   `json:"begin"`
    End   string   `json:"end"`
}

type ChartData struct {
    X     int64    `json:"x" bson:"_id,omitempty"`
    Y     int64    `json:"y" bson:"username"`
    Area  int64    `json:"area" bson:"area, omitempty"`
    Point int64   `json:"point" bson:"pt, omitempty"`
}

type DataItem struct {
    //ID        bson.ObjectId `bson:"-,omitempty" json:"-"`
    Timestamp int64  `bson:"ts" json:"x"`
    Area      int  `bson:"area" json:"area"`
    Point     int  `bson:"pt" json:"pt"`
    Value     float64 `bson:"val" json:"y""`
}

func random(min, max int64) int64 {
    //rand.Seed(time.Now().Unix())
    return rand.Int63n(max - min) + min
}

func getChartData(c *gin.Context) {
    response := newResponse(c)

    body := &ChartDataQuery{}
    err := json.NewDecoder(c.Request.Body).Decode(&body)
    if err != nil {
        response.Errors = append(response.Errors, err.Error())
        response.Fail()
        return
    }
    log.Printf("request body: %+v\n", body)
    //log.Printf("query body: %v\n", body)
    //if err != nil {
    //    response.Errors = append(response.Errors, err.Error())
    //    response.Fail()
    //    return
    //}

    session, err := mgo.Dial("localhost:27017")
    defer session.Close()

    if err != nil {
        fmt.Println(err)
        return
    }

    //db := session.DB("data_viz")
    //collection := db.C("data");

    //log.Println(collection.FullName)

    //var result []DataItem
    //err = collection.Find(bson.M{"pt": 1, "area": 1}).All(&result)
    //if err != nil {
    //    response.Errors = append(response.Errors, err.Error())
    //    response.Fail()
    //    return
    //}

    //t := time.Unix(0, 0).Format("2006-1-2")
    body.Begin = time.Unix(0, 0).Format("2006-1-2")
    begin, err := time.Parse("2006-1-2", body.Begin)
    end, err := time.Parse("2006-1-2", body.End)

    log.Println("query date: ", begin, " to ", end)

    data, err := storage.GetChartData(body.Area, body.Point, begin, end)

    //if len(data) == 0 {
    //    response.Errors = append(response.Errors, "No data")
    //    response.Fail()
    //    return
    //}

    for i := range data {
        data[i].Timestamp *= 1000
    }

    c.JSON(http.StatusOK, data)
}

func renderChartData(c *gin.Context) {
    c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

