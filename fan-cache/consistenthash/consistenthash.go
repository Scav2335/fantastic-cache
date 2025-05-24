package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
type Hash func(data []byte) uint32

type Map struct {
	hash     Hash
	replicas int
	keys     []int
	hashMap  map[int]string
}

func New(replicas int, hash Hash) *Map {
	m := &Map{
		hash:     hash,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE
	}

	return m
}

// Add adds some keys to the hash.
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			vh := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, vh)
			m.hashMap[vh] = key
		}
	}
	sort.Ints(m.keys)
}

// Get gets the closest item in the hash to the provided key.
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}

	h := int(m.hash([]byte(key)))

	// If there is no such index, Search returns n.
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= h
	})

	return m.hashMap[m.keys[idx%len(m.keys)]]
}
