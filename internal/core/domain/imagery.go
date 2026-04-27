package core_domain

import "time"

type NewImagery struct {
	ImageryRequest
	Id         int
	Created_at time.Time
	Updated_at time.Time
}

type ImageryResponce struct {
	NewImagery
	Image []byte
}

type ImageryRequest struct {
	User_id string
	City    string
	From    string
	To      string
}
