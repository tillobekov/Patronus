package util

type OrderStatus string

const (
	OrderStatusACTIVE   OrderStatus = "Active"
	OrderStatusCANCELED OrderStatus = "Canceled"
	OrderStatusFILLED   OrderStatus = "Filled"
)

type OrderType string

const (
	OrderTypeLIMIT  OrderType = "LIMIT"
	OrderTypeMARKET OrderType = "MARKET"
)
