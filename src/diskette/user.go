package main

import "time"

type User struct {
	Id               string    `json:"_id"`
	Roles            []string  `json:"roles"`
	Name             string    `json:"name"`
	Email            string    `json:"email"`
	HashedPass       string    `json:"hashedPass"`
	Lang             string    `json:"lang"`
	CreatedAt        time.Time `json:"createdAt"`
	ConfirmedAt      time.Time `json:"confirmedAt"`
	ResetKey         string    `json:"resetKey"`
	RequestedResetAt string    `json:"requestedResetAt"`
	IsSuspended      bool      `json:"isSuspended"`
	SessionKeys      []string  `json:"sessionKeys"`
}
