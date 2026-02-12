package model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string `gorm:"not null"`
	Email    string `validate:"required,email"`
	Name     string `validate:"required"`
	Age      int    `validate:"gte=0,lte=130"`
}
