package errors

import (
	"bytes"
	"fmt"
)

const lostKey = "lost_key"

// DataBag ...
type DataBag struct {
	Key   string
	Value interface{}
}

// DataBagList ...
type DataBagList []DataBag

type stringable interface {
	String() string
}

// String ...
func (dbl *DataBagList) String() string {
	buf := bytes.NewBufferString("{")
	first := true
	for _, bag := range *dbl {
		if !first {
			buf.WriteString(", ")
		}
		buf.WriteString(fmt.Sprintf("%s: %v", bag.Key, bag.Value))
		first = false
	}
	buf.WriteString("}")

	return buf.String()
}

func dataBagsFromSlice(items ...interface{}) []DataBag {
	if len(items) == 0 {
		return nil
	}

	if len(items)%2 != 0 {
		n := len(items) - 1
		items = append(items[:n], lostKey, items[n])
	}

	n := len(items)
	result := make([]DataBag, 0, n/2)
	for i := 1; i < n; i += 2 {
		result = append(result, DataBag{
			Key:   formatDataKey(items[i-1]),
			Value: items[i],
		})
	}

	return result
}

func formatDataKey(key interface{}) string {
	switch keyTyped := key.(type) {
	case stringable:
		return keyTyped.String()
	case string:
		return keyTyped
	default:
		return fmt.Sprint(key)
	}
}
