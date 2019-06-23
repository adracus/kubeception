package pointers

func Int64(i int64) *int64 {
	return &i
}

func Int32(i int32) *int32 {
	return &i
}

func Bool(b bool) *bool {
	return &b
}

func DerefBoolOrDefault(b *bool, defaultValue bool) bool {
	if b == nil {
		return defaultValue
	}
	return *b
}

func DerefInt32OrDefault(i *int32, defaultValue int32) int32 {
	if i == nil {
		return defaultValue
	}
	return *i
}

func DerefInt64OrDefault(i *int64, defaultValue int64) int64 {
	if i == nil {
		return defaultValue
	}
	return *i
}
