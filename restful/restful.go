package restful

import (
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"strconv"
)

type Stringer interface {
	String() string
}

var NilValue = reflect.Value{}

type argumentParser func([]string) (reflect.Value, error)

func parseBool(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseBool(values[0])
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(n), nil
}

func parseInt(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseInt(values[0], 10, 0)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(int(n)), nil
}

func parseInt8(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseInt(values[0], 10, 8)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(int8(n)), nil
}

func parseInt16(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseInt(values[0], 10, 16)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(int16(n)), nil
}

func parseInt32(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseInt(values[0], 10, 32)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(int32(n)), nil
}

func parseInt64(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseInt(values[0], 10, 64)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(int64(n)), nil
}

func parseUint(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseUint(values[0], 10, 0)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(uint(n)), nil
}

func parseUint8(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseUint(values[0], 10, 8)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(uint8(n)), nil
}

func parseUint16(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseUint(values[0], 10, 16)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(uint16(n)), nil
}

func parseUint32(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseUint(values[0], 10, 32)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(uint32(n)), nil
}

func parseUint64(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseUint(values[0], 10, 64)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(uint64(n)), nil
}

func parseFloat32(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseFloat(values[0], 32)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(float32(n)), nil
}

func parseFloat64(values []string) (v reflect.Value, err error) {
	n, err := strconv.ParseFloat(values[0], 64)
	if err != nil {
		return NilValue, err
	}
	return reflect.ValueOf(float64(n)), nil
}

func parseString(values []string) (v reflect.Value, err error) {
	return reflect.ValueOf(values[0]), nil
}

type restfulHandler struct {
	Func       reflect.Value
	ArgParsers []argumentParser
	ArgNames   []string
}

func (h *restfulHandler) encodeResponse(w http.ResponseWriter, response interface{}) {
	enc := json.NewEncoder(w)
	err := enc.Encode(response)
	if err != nil {
		panic(err) // Let the recovery handler deal with it
	}
}

func (h *restfulHandler) encodeErrorResponse(w http.ResponseWriter, err error) {
	h.encodeResponse(w, map[string]interface{}{"error": err.Error()})
}

func (h *restfulHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer func() {
		v := recover()
		if v != nil {
			http.Error(w, v.(Stringer).String(), http.StatusInternalServerError)
		}
	}()

	r.ParseForm()

	arguments := make([]reflect.Value, len(h.ArgParsers))

	for i, parser := range h.ArgParsers {
		name := h.ArgNames[i]
		values, ok := r.Form[name]
		if !ok {
			h.encodeErrorResponse(w, fmt.Errorf("Required argument '%s' not found", name))
			return
		}

		value, err := parser(values)
		if err != nil {
			h.encodeErrorResponse(w, err)
			return
		}

		arguments[i] = value
	}

	results := h.Func.Call(arguments)

	var response interface{}

	switch len(results) {
	case 0:

	case 1:
		response = results[0].Interface()

		err, ok := response.(error)
		if ok {
			h.encodeErrorResponse(w, err)
			return
		}

	default:
		responseSlice := make([]interface{}, len(results))

		for i, result := range results {
			responseSlice[i] = result.Interface()
		}

		response = responseSlice
	}

	h.encodeResponse(w, map[string]interface{}{"response": response})
}

func determineParser(kind reflect.Kind) (parser argumentParser) {
	switch kind {
	case reflect.Bool:
		return parseBool
	case reflect.Int:
		return parseInt
	case reflect.Int8:
		return parseInt8
	case reflect.Int16:
		return parseInt16
	case reflect.Int32:
		return parseInt32
	case reflect.Int64:
		return parseInt64
	case reflect.Uint:
		return parseUint
	case reflect.Uint8:
		return parseUint8
	case reflect.Uint16:
		return parseUint16
	case reflect.Uint32:
		return parseUint32
	case reflect.Uint64:
		return parseUint64
	case reflect.Float32:
		return parseFloat32
	case reflect.Float64:
		return parseFloat64
	case reflect.String:
		return parseString
	}

	return nil
}

func RestfulHandler(f interface{}, argNames []string) (rh http.Handler) {
	v := reflect.ValueOf(f)
	t := v.Type()

	// NumIn, NumOut etc. will panic for us if the argument isn't a function.

	if t.NumIn() != len(argNames) {
		panic("Length of argNames does not match number of parameters to function")
	}

	argParsers := make([]argumentParser, t.NumIn())

	for i := 0; i < t.NumIn(); i++ {
		inType := t.In(i)
		kind := inType.Kind()

		var parser argumentParser

		if kind == reflect.Slice {
			subParser := determineParser(inType.Elem().Kind())
			if subParser == nil {
				panic(fmt.Sprintf("Invalid argument type: []%s", inType.Kind().String()))
			}

			parser = func(values []string) (v reflect.Value, err error) {
				results := reflect.MakeSlice(inType, len(values), len(values))

				for i, value := range values {
					result, err := subParser([]string{value})
					if err != nil {
						return NilValue, err
					}

					results.Index(i).Set(result)
				}

				return results, nil
			}

		} else {
			parser = determineParser(kind)
			if parser == nil {
				panic(fmt.Sprintf("Invalid argument type: %s", inType.Kind().String()))
			}
		}

		argParsers[i] = parser
	}

	return &restfulHandler{
		Func:       v,
		ArgParsers: argParsers,
		ArgNames:   argNames,
	}
}
