package core_domain

import "time"

type NewImagery struct {
	ImageryRequest
	Id         int       `json:"id"`
	Created_at time.Time `json:"createdAt"`
	Updated_at time.Time `json:"updatedAt"`
}

type ImageryResponce struct {
	NewImagery
	Image []byte
}

type ImageryRequest struct {
	User_id string `json:"userId"`
	City    string `json:"city" validate:"required"`
	From    string `json:"from" validate:"required"`
	To      string `json:"to" validate:"required"`
}
