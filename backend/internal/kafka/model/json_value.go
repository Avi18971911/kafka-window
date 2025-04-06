package model

type JSONValue struct {
	StringVal *string
	NumberVal *float64
	BoolVal   *bool
	NullVal   bool
	ObjectVal map[string]JSONValue
	ArrayVal  []JSONValue
	Type      JSONType
}

type JSONType string

const (
	STRING JSONType = "string"
	NUMBER JSONType = "number"
	BOOL   JSONType = "bool"
	OBJECT JSONType = "object"
	ARRAY  JSONType = "array"
	NULL   JSONType = "null"
)
