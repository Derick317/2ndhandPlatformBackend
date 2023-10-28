package model

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"unique" form:"email"`
	Username string `form:"username"`
	Password string `form:"password"`
	Enabled  bool   `gorm:"default:true"`
}

type Item struct {
	ID          uint `gorm:"primaryKey"`
	SellerId    uint
	Price       float32
	Tag         TagType
	Description string        `gorm:"type:text"`
	ImageUrls   ImageUrlsType `gorm:"type:text"`
	// false: available for every one to use; true: lock is held
	Lock   bool           `gorm:"default:false"`
	Status ItemStatusType `gorm:"default:0"`
}
