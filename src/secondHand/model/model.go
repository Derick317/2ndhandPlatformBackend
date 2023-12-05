package model

type OrdersType = MapType[uint, struct{}]
type ImageUrlsType = MapType[string, string]

type User struct {
	ID       uint64 `gorm:"primaryKey"`
	Email    string `gorm:"unique" form:"email"`
	Username string `form:"username"`
	Password string `form:"password"`
	Enabled  bool   `gorm:"default:true"`
	// Orders   MapType[uint64, struct{}] `gorm:"type:bytes"`
}

type Item struct {
	ID          uint64 `gorm:"primaryKey"`
	SellerId    uint64
	Title       string
	Price       float32
	Tag         TagType
	Description string         `gorm:"type:text"`
	ImageUrls   ImageUrlsType  `gorm:"type:bytes"`
	Status      ItemStatusType `gorm:"default:0"`
	Version     uint           `gorm:"default:0"` // optimistic lock for status
}

type Order struct {
	ID      uint64 `gorm:"primaryKey"`
	BuyerId uint64
	ItemId  uint64
	ExpTime int64
}
