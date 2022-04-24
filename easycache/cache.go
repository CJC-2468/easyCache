package easycache

import (
	"easyCache/easycache/lru"
	"sync"
)

//对lru的封装，支持并发控制.
type cache struct {
	mu       sync.Mutex
	lru      *lru.LRUCache
	capacity int
}

//支持并发的情况下，做了初始化lru
func (c *cache) add(key string, value ByteView) {
	c.mu.Lock()
	defer c.mu.Unlock()
	//延迟初始化(Lazy Initialization),一个对象的延迟初始化意味着该对象的创建将会延迟至第一次使用该对象时
	//主要用于提高性能，并减少程序内存要求。
	if c.lru == nil {
		c.lru = lru.Constructor(c.capacity, nil)
	}
	c.lru.Put(key, value)
}

//注意读缓存的时候其实也会对lru进行一个写操作，
//因为需要把热点数据移动到链表头部,所以也要加锁
func (c *cache) get(key string) (value ByteView, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.lru == nil {
		return
	}
	if v, ok := c.lru.Get(key); ok {
		return v.(ByteView), ok
	}
	return
}
