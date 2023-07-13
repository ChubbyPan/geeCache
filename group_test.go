package geecache

import (
	"fmt"
	"log"
	"reflect"
	"testing"
)

func TestGetter(t *testing.T) {
	var f Getter = GetterFunc(func(key string) ([]byte, error) {
		return []byte(key), nil
	})

	expect := []byte("key")
	if v, _ := f.Get("key"); !reflect.DeepEqual(v, expect) {
		t.Errorf("callback failed")
	}
}

var db = map[string]string{
	"Panny":  "888",
	"Rshao":  "666",
	"Others": "6868",
}

func TestGet(t *testing.T) {
	loadCounts := make(map[string]int, len(db))
	// create new namespace
	gee := NewGroup("luckyNums", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("luckyNumsDB search key:", key)
			if v, ok := db[key]; ok {
				loadCounts[key] += 1
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	for k, v := range db {
		if view, err := gee.Get(k); err != nil || view.String() != v {
			t.Fatalf("failed to get value of %s", k)
		} // 没有正确获得期待的v，此时没有缓存，但是本地有数据
		if _, err := gee.Get(k); err != nil || loadCounts[k] > 1 {
			t.Fatalf("cache %s miss", k)
		} // 命中次数>1，但是没有缓存
	}

	// 本地没有，应该err 会返回不存在
	if view, err := gee.Get("unknown"); err == nil {
		t.Fatalf("the value of unknow should be empty, but %s got", view)
	}
}
