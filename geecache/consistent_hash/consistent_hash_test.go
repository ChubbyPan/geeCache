package consistenthash

import (
	"strconv"
	"testing"
)

func TestHashing(t *testing.T) {
	hash := NewMap(3, func(data []byte) uint32 {
		i, _ := strconv.Atoi(string(data))
		return uint32(i)
	})
	//real node 6: 6, 16, 26
	//real node 4: 4, 14, 24
	//real node 2: 2, 12, 22
	hash.Add("6", "4", "2")

	testCases := map[string]string{
		"2":  "2",
		"11": "2",
		"23": "4",
		"27": "2",
	}

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s, but yielded: %s", k, v, hash.Get(k))
		}
	}
	// real node 8: 8. 18. 28
	hash.Add("8")
	testCases["27"] = "8"

	for k, v := range testCases {
		if hash.Get(k) != v {
			t.Errorf("Asking for %s, should have yielded %s, but yielded: %s", k, v, hash.Get(k))
		}
	}
}
