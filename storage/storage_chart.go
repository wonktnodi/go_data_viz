package storage

import (
    "../model"
    "errors"
    "time"
    "log"
)

//import "time"

func GetChartData(area, point int, begin, end time.Time) ([]model.ChartDataItem, error) {
    if StorageDB == nil {
        return nil, errors.New("invalidate parameter")
    }
    if begin.Equal(end) {
        offset, _ :=time.ParseDuration("24h")
        end = end.Add(offset)
    }
    log.Println("query date: ", begin, " to ", end)
    var ret []model.ChartDataItem
    StorageDB.Model(&model.ChartDataItem{}).Where("timestamp >= ? and timestamp < ? and point = ? and area = ?",
        begin.Unix(), end.Unix(), point, area).Find(&ret)
    return ret, StorageDB.Error
}