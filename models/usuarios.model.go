package models

import "time"

type Usuario struct {
	Id        int       `json:"id" gorm:"primaryKey"`
	UserId    string    `json:"user_id"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (u *Usuario) TableName() string {
	// custom table name, this is default
	return "megavirada24.usuarios"
}
