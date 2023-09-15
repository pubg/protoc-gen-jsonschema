package utils

func UInt32(i uint32) *int {
	ii := int(i)
	return &ii
}

func Int32(i int32) *int {
	ii := int(i)
	return &ii
}

func Bool(b bool) *bool {
	return &b
}
