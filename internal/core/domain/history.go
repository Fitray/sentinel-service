package core_domain

import "time"

type FilterRequest struct {
	City    string
	From    time.Time
	To      time.Time
	User_id string
}
