package records

import "time"

type Record struct {
	Name        string    `json:"name"`
	Url       string    `json:"url"`
	Timestamp time.Time `json:"timestamp"`
	Size      int64     `json:"size"`
}


type RecordsResponse struct {
	Status string `json:"status"`
	Recordings []Record `json:"records"`
	Count      int         `json:"count"`
}
