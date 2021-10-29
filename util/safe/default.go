package safe

import "time"

func TimeOrDefault(t *time.Time, def time.Time) time.Time {
	if t == nil {
		return def
	}
	return *t
}

func StringOrDefault(s *string, def string) string {
	if s == nil {
		return def
	}
	return *s
}

func IntOrDefault(i *int, def int) int {
	if i == nil {
		return def
	}
	return *i
}

func Int8OrDefault(i *int8, def int8) int8 {
	if i == nil {
		return def
	}
	return *i
}

func Int16OrDefault(i *int16, def int16) int16 {
	if i == nil {
		return def
	}
	return *i
}

func Int32OrDefault(i *int32, def int32) int32 {
	if i == nil {
		return def
	}
	return *i
}

func Int64OrDefault(i *int64, def int64) int64 {
	if i == nil {
		return def
	}
	return *i
}

func UintOrDefault(i *uint, def uint) uint {
	if i == nil {
		return def
	}
	return *i
}

func Uint8OrDefault(i *uint8, def uint8) uint8 {
	if i == nil {
		return def
	}
	return *i
}

func Uint16OrDefault(i *uint16, def uint16) uint16 {
	if i == nil {
		return def
	}
	return *i
}

func Uint32OrDefault(i *uint32, def uint32) uint32 {
	if i == nil {
		return def
	}
	return *i
}

func Uint64OrDefault(i *uint64, def uint64) uint64 {
	if i == nil {
		return def
	}
	return *i
}

func Float32OrDefault(f *float32, def float32) float32 {
	if f == nil {
		return def
	}
	return *f
}

func Float64OrDefault(f *float64, def float64) float64 {
	if f == nil {
		return def
	}
	return *f
}