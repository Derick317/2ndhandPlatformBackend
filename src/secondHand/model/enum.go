package model

type TagType = uint64

const (
	TagAll TagType = iota
	Book
	Electronics
	Stationery
	Others
	TagCounter
)

type ItemStatusType = uint

const (
	Available ItemStatusType = iota
	SellerModifying
	OnOrder // some buyer is about to buy but has not pay for it yet
	Sold
	Deleted // deleted by seller
	StatusCounter
)
