package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type Hash func(data []byte) uint32

// map contains all hashed keys
type Map struct {
	hash     Hash           // hash algorithm
	replicas int            // virtual node nums
	keys     []int          // hash ring, which keys stored, sorted
	hashMap  map[int]string // map virtual node and real node, key: vitual node id, value: real node id
}

// creates a Map instance
func NewMap(replicas int, fn Hash) *Map {
	m := &Map{
		replicas: replicas,
		hash:     fn,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil {
		m.hash = crc32.ChecksumIEEE // or defina another hash algrithom
	}
	return m
}

// input: real nodes' name
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key))) // caculate hash of virtual node
			m.keys = append(m.keys, hash)                      // add to hash ring
			m.hashMap[hash] = key                              // map of vitural node and real node
		}
	}
	sort.Ints(m.keys)
}

// input key, caculate key's hash, return key's real node
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	hash := int(m.hash([]byte(key)))
	// binary search for appropriate replica
	idx := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})
	return m.hashMap[m.keys[idx%len(m.keys)]]
}
