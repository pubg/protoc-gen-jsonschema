package jsonschema

type RefId string

func (i RefId) String() string {
	return string(i)
}

func (i RefId) IsEmpty() bool {
	return i.String() == ""
}
