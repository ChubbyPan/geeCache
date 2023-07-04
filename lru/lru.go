package lru

import "container/list"

// lru的 cache 由一个双向链表和存储kv的map组成
type Cache struct {
	maxBytes  int64                         // 容量
	nbytes    int64                         // 使用容量
	ll        *list.List                    // 双向链表
	cache     map[string]*list.Element      // 单个缓存单位
	onEvicted func(key string, value Value) // 移除数据后的回调函数
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		onEvicted: onEvicted,
	}
}

// lru 基本CRUD
// 缓存命中
func (c *Cache) Get(key string) (value Value, ok bool) {
	if elm, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elm)
		v := elm.Value.(*entry)
		return v.value, true
	}
	return
}

//缓存淘汰
func (c *Cache) RemoveOldest() {
	elm := c.ll.Back()
	if elm != nil {
		c.ll.Remove(elm)
		v := elm.Value.(*entry)
		delete(c.cache, v.key)
		c.nbytes -= int64(len(v.key)) + int64(v.value.Len())
		if c.onEvicted != nil {
			c.onEvicted(v.key, v.value)
		}
	}
}

// 新增数据
func (c *Cache) Add(key string, value Value) {
	// 已经存在, 更新value
	if elm, ok := c.cache[key]; ok {
		c.ll.MoveToFront(elm)
		v := elm.Value.(*entry)
		c.nbytes += int64(len(v.key)) + int64(v.value.Len())
		v.value = value
	} else {
		elm := c.ll.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = elm
		c.nbytes += int64(len(key)) + int64(value.Len())
	}

	// 超过大小
	for c.maxBytes != 0 && c.maxBytes < c.nbytes {
		c.RemoveOldest()
	}
}

// 获取当前cache数据单位量
func (c *Cache) Len() int {
	return c.ll.Len()
}
