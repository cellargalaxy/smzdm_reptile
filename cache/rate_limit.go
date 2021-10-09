package cache

import "time"

func RateLimit(key string, duration time.Duration) bool {
	_, ok := get(key)
	if ok {
		return ok
	}
	set(key, key, duration)
	return false
}
