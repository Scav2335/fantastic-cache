package cache

type Cache interface {
	Set(key string, value Value)
	Get(key string) (Value, bool)
	Len() int
	RemoveOldest()
}
