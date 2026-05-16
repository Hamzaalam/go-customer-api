package models

type Customer struct {
	ID        int    `json:"id" gorm:"primaryKey;autoIncrement"`
	FirstName string `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName  string `json:"last_name" gorm:"type:varchar(100);not null"`
	Email     string `json:"email" gorm:"type:varchar(255);uniqueIndex;not null"`
}
