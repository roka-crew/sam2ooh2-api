package query

import "strings"

type Sorts []Sort

type Sort struct {
	By    string
	Order Order
}

type Order string

const (
	OrderASC  Order = "ASC"
	OrderDESC Order = "DESC"
)

func (o Order) String() string {
	return string(o)
}

func (o Order) IsDESC() bool {
	return o == OrderDESC
}

func NormalizeOrder(s string, defaultOrder Order) Order {
	upper := strings.ToUpper(strings.TrimSpace(s))
	switch Order(upper) {
	case OrderASC:
		return OrderASC
	case OrderDESC:
		return OrderDESC
	default:
		if defaultOrder != "" {
			return defaultOrder
		}
		return OrderASC
	}
}

func (o Order) Valid() bool {
	return o == OrderASC || o == OrderDESC
}
