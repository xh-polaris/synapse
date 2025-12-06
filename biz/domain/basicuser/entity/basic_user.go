package entity

import "time"

// gender
const (
	UnKnow uint8 = 0
	Male   uint8 = 1
	Female uint8 = 2
)

type BasicUser struct {
	ID        string
	SchoolID  *string
	Code      *string
	Phone     *string
	Password  *string
	Name      string
	Gender    uint8
	Extra     *map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Auth struct {
	AuthType string
	AuthId   string
}

type School struct {
	ID        string
	Name      string
	Extra     *map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}
