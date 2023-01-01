package models

import "time"

type Post struct {
	Id        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UserId    string    `json:"user_id"`
}
