package model

type TagType uint

const (
	Book TagType = iota
	Electronics
	Stationery
	Others
	TagCounter
)
