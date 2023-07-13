package geecache

import (
	"fmt"
	"log"
	"sync"
)

// 与外界交互，有三种情况
// 1. 接受一个key， 检查是否已经被缓存，如果被缓存了，直接返回
// 2. 接受一个key，没有被缓存，但是不需要从远程节点获取，调用回调函数函数，获取值并添加到缓存中，返回缓存值
// 3. 接受一个key，没有被缓存，需要用远程节点获取，更新缓存值后，与远程节点交互后返回缓存

// a getter loads data for key
type Getter interface {
	Get(key string) ([]byte, error)
}

type GetterFunc func(key string) ([]byte, error)

func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

// 封装一个cache (namespace) 提供未击中缓存时源数据的回调
type Group struct {
	name      string
	getter    Getter
	mainCache cache
}

var (
	mu     sync.Mutex
	groups = make(map[string]*Group)
)

func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:   name,
		getter: getter,
		mainCache: cache{
			cacheBytes: cacheBytes,
		},
	}
	// 保存group信息，避免重复定义
	groups[name] = g
	return g
}

func GetGroup(name string) *Group {
	mu.Lock()
	defer mu.Unlock()
	return groups[name]
}

// group‘s get method
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("Key is required")
	}
	// first situation
	if v, ok := g.mainCache.get(key); ok {
		log.Printf("[%v Cache] hit, key: %v, value: %v", g.name, key, v)
		return v, nil
	}

	return g.load(key)

}

//	second situation
func (g *Group) load(key string) (ByteView, error) {
	return g.getLocally(key)
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	// 本地无
	if err != nil {
		return ByteView{}, err
	}
	// 本地有，存入缓存
	v := ByteView{b: cloneBytes(bytes)}
	g.populateGache(key, v)
	return v, nil
}

func (g *Group) populateGache(key string, value ByteView) {
	g.mainCache.add(key, value)
}
