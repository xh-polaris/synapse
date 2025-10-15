package entity

import "time"

type ThirdPartyUser struct {
	ID        string
	Code      string
	App       string
	CreatedAt time.Time
	UpdatedAt time.Time
}
