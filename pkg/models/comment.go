package models

import "time"

type Comment struct {
	ID        int
	NewsID    int
	ParentID  *int
	Content   string
	Author    string
	CreatedAt time.Time
}
