package query

import "strings"

type Sorts []Sort

type Sort struct {
	By    string
	Order Order
}

type Order string

const (
	ASC  Order = "ASC"
	DESC Order = "DESC"
)

func (o Order) String() string {
	return string(o)
}

func (o Order) IsDESC() bool {
	return o == DESC
}

func NormalizeOrder(s string, defaultOrder Order) Order {
	upper := strings.ToUpper(strings.TrimSpace(s))
	switch Order(upper) {
	case ASC:
		return ASC
	case DESC:
		return DESC
	default:
		if defaultOrder != "" {
			return defaultOrder
		}
		return ASC
	}
}

func (o Order) Valid() bool {
	return o == ASC || o == DESC
}
