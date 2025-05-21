package domain

type User struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `validate:"required"`
	Email string `validate:"required,email"`
	Age   int    `validate:"gte=0,lte=130"`
}
