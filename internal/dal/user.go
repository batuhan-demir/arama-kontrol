package dal

type User struct {
	Id        int    `gorm:"primarykey;autoIncrement:true;unique" json:"id"`
	Name      string `gorm:"required" json:"name"`
	Surname   string `gorm:"required" json:"surname"`
	Email     string `gorm:"unique" json:"email"`
	Phone     string `gorm:"unique" json:"phone"`
	Password  string `gorm:"required" json:"-"`
	Is_Active bool   `gorm:"default:true" json:"is_active"`
}

type UserCreate struct {
	Name     string `validate:"required"`
	Surname  string `validate:"required"`
	Email    string `validate:"required"`
	Phone    string `validate:"required"`
	Password string `validate:"required"`
}

type UserLogin struct {
	Email    string `validate:"required"`
	Password string `validate:"required"`
}
