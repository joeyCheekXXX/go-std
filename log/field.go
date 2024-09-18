package log

import (
	"go.uber.org/zap"
	"time"
)

func String(key string, val string) zap.Field {
	return zap.String(key, val)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Int32(key string, val int32) zap.Field {

	return zap.Int32(key, val)
}

func Int64(key string, val int64) zap.Field {
	return zap.Int64(key, val)
}

func Uint(key string, val uint) zap.Field {
	return zap.Uint(key, val)
}

func Uint32(key string, val uint32) zap.Field {
	return zap.Uint32(key, val)
}

func Uint64(key string, val uint64) zap.Field {
	return zap.Uint64(key, val)
}

func Float32(key string, val float32) zap.Field {
	return zap.Float32(key, val)
}

func Float64(key string, val float64) zap.Field {
	return zap.Float64(key, val)
}

func Time(key string, val time.Time) zap.Field {
	return zap.Time(key, val)
}

func Duration(key string, val time.Duration) zap.Field {
	return zap.Duration(key, val)
}

func Any(key string, val any) zap.Field {
	return zap.Any(key, val)
}

func Strings(key string, val []string) zap.Field {
	return zap.Strings(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}
