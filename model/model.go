package model

import "time"

type Goods struct {
	Title    string    `json:"title"`
	Price    float64   `json:"price"`
	Zhi      int       `json:"zhi"`
	Buzhi    int       `json:"buzhi"`
	Merchant string    `json:"merchant"`
	Url      string    `json:"url"`
	Date     time.Time `json:"date"`
}
