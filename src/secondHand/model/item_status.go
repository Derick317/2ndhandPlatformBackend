package model

type ItemStatusType uint

const (
	Available TagType = iota
	SellerModifying
	OnOrder // some buyer is about to buy but has not pay for it yet
	Deleted // deleted by seller
)
