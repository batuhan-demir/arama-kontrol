package dal

type Number struct {
	Number string `gorm:"primarykey;unique" validate:"required" json:"number"`
	Name   string `validate:"required" json:"name"`
}
