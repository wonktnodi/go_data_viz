package main

import (
    "github.com/gin-gonic/gin"
    "gopkg.in/mgo.v2/bson"

    "log"
    "encoding/json"
    "net/http"
    "time"
    "math/rand"
)


type ChartDataQuery struct {
	ID                   bson.ObjectId `json:"-" bson:"_id,omitempty"`
	Type                 int64        `json:"type" bson:"username"`
}

type ChartData struct {
    X    int64    `json:"x" bson:"_id,omitempty"`
    Y    int64    `json:"y" bson:"username"`
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

    //err := c.Bind(body)
    log.Printf("query body: %v\n", body)
    //if err != nil {
    //    response.Errors = append(response.Errors, err.Error())
    //    response.Fail()
    //    return
    //}
    resp := make([]ChartData, 1000)

    t := time.Now()
    tval := t.Unix()

    for i, _ := range resp {
        resp[i].X = tval + (int64)(i * 1000)
        resp[i].Y = random(-1000, 1000)
    }

    c.JSON(http.StatusOK, resp)

}

func renderChartData(c *gin.Context) {
	c.HTML(http.StatusOK, c.Request.URL.Path, c.Keys)
}

