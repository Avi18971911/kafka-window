package model

type JSONValue struct {
	StringVal *string              `json:"stringVal,omitempty"`
	NumberVal *float64             `json:"numberVal,omitempty"`
	BoolVal   *bool                `json:"boolVal,omitempty"`
	NullVal   bool                 `json:"nullVal"`
	ObjectVal map[string]JSONValue `json:"objectVal,omitempty"`
	ArrayVal  []JSONValue          `json:"arrayVal,omitempty"`
}
