package model


type ChartDataItem struct {
    Timestamp int64  `json:"x"`
    Area      int  `json:"area"`
    Point     int  `json:"pt"`
    Value     int64 `json:"y""`
}
