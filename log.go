package errors

import "fmt"

// GetLogKV ...
func GetLogKV(err error) []interface{} {
	if err == nil {
		return nil
	}

	switch errTyped := err.(type) {
	case *Error:
		if errTyped.IsNil() {
			return nil
		}

		return append(
			append(
				[]interface{}{"err", errTyped.message},
				errTyped.GetLogKV()...,
			),
			GetLogKV(errTyped.Cause())...,
		)

	case Causer:
		return append(
			[]interface{}{"err", err.Error()},
			GetLogKV(errTyped.Cause())...,
		)

	default:
		return []interface{}{
			"err_type", fmt.Sprintf("%T", err),
			"err", err.Error(),
		}
	}
}

func KVtoString(kvs []interface{}) string {
	dbl := (DataBagList)(dataBagsFromSlice(kvs...))
	return dbl.String()
}
