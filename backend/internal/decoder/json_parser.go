package decoder

import (
	"github.com/Avi18971911/kafka-window/backend/internal/kafka/model"
	"github.com/valyala/fastjson"
)

func parseFastJSON(value *fastjson.Value) model.JSONValue {
	if value == nil {
		return model.JSONValue{NullVal: true}
	}
	switch value.Type() {
	case fastjson.TypeString:
		str := string(value.GetStringBytes())
		return model.JSONValue{StringVal: &str}
	case fastjson.TypeNumber:
		num := value.GetFloat64()
		return model.JSONValue{NumberVal: &num}
	case fastjson.TypeTrue:
		b := true
		return model.JSONValue{BoolVal: &b}
	case fastjson.TypeFalse:
		b := false
		return model.JSONValue{BoolVal: &b}
	case fastjson.TypeObject:
		obj := make(map[string]model.JSONValue)
		value.GetObject().Visit(func(key []byte, v *fastjson.Value) {
			obj[string(key)] = parseFastJSON(v)
		})
		return model.JSONValue{ObjectVal: obj}
	case fastjson.TypeArray:
		var arr []model.JSONValue
		for _, val := range value.GetArray() {
			arr = append(arr, parseFastJSON(val))
		}
		return model.JSONValue{ArrayVal: arr}
	case fastjson.TypeNull:
		return model.JSONValue{NullVal: true}
	default:
		return model.JSONValue{}
	}
}

func ParseString(jsonString string) (model.JSONValue, error) {
	p := fastjson.Parser{}
	value, err := p.Parse(jsonString)
	if err != nil {
		return model.JSONValue{}, err
	}
	return parseFastJSON(value), nil
}
