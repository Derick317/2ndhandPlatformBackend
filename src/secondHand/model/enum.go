package model

type TagType = uint

const (
	Book TagType = iota
	Electronics
	Stationery
	Others
	TagCounter
)

type ItemStatusType = uint

const (
	Available TagType = iota
	SellerModifying
	OnOrder // some buyer is about to buy but has not pay for it yet
	Sold
	Deleted // deleted by seller
)
