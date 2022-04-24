package lru

// LRU缓存的基本数据结构
type LRUCache struct {
	capacity int                     //缓存的容量
	cache    map[string]*DLinkedNode //真正缓存数据的地方
	list     *DLinkedList            //这里其实是有序链表，头部放热点数据，需要淘汰就先淘汰尾部的
	//某条记录被移除时的回调函数，可以为 nil，这里应用场景可以是记录那些被LRU淘汰机制淘汰的数据。
	OnEvicted func(key string, value Value)
}

//初始化LRU缓存
func Constructor(capacity int, onEvicted func(string, Value)) *LRUCache {
	return &LRUCache{
		cache:     map[string]*DLinkedNode{},
		capacity:  capacity,
		list:      initDLinkedList(),
		OnEvicted: onEvicted,
	}
}

//根据key获取数据
func (lru *LRUCache) Get(key string) (value Value, ok bool) {
	if node, ok := lru.cache[key]; ok {
		lru.list.moveToHead(node)
		return node.value, true
	}
	return
}

//key存在就更新，不存在就添加，添加的时候如果容量满了就淘汰最近最久没用的
func (lru *LRUCache) Put(key string, value Value) {
	if node, ok := lru.cache[key]; ok {
		node.value = value
		lru.list.moveToHead(node)
	} else {
		if len(lru.cache) == lru.capacity {
			removed := lru.list.removeTail()
			delete(lru.cache, removed.key)
			//删除时触发自己设置的回调函数
			if lru.OnEvicted != nil {
				lru.OnEvicted(removed.key, removed.value)
			}
		}
		node := initDLinkedNode(key, value)
		lru.list.addToHead(node)
		lru.cache[key] = node
	}
}

//获取目前缓存中有多少数据
func (lru *LRUCache) Len() int {
	return len(lru.cache)
}

//为了缓存中存储数据的通用性，允许缓存值是实现了 Value 接口的任意类型，
//该接口只包含了一个方法 Len() int，返回缓存值所占用内存字节大小。
type Value interface {
	Len() int
}

//使用双向链表是为了删除操作，删除节点需要前驱，双链表这个操作是O（1）
//同时存了kv是为了容量不够时删除哈希表中的最后一个节点，这时候需要用k。
type DLinkedNode struct {
	key        string
	value      Value
	prev, next *DLinkedNode
}
type DLinkedList struct {
	head, tail *DLinkedNode
}

func initDLinkedNode(key string, value Value) *DLinkedNode {
	return &DLinkedNode{
		key:   key,
		value: value,
	}
}
func initDLinkedList() *DLinkedList {
	head, tail := initDLinkedNode("", nil), initDLinkedNode("", nil)
	head.next = tail
	tail.prev = head
	return &DLinkedList{
		head: head,
		tail: tail,
	}
}
func (lru *DLinkedList) addToHead(node *DLinkedNode) {
	node.prev = lru.head
	node.next = lru.head.next
	lru.head.next.prev = node
	lru.head.next = node
}
func (lru *DLinkedList) removeNode(node *DLinkedNode) {
	node.prev.next = node.next
	node.next.prev = node.prev
}
func (lru *DLinkedList) moveToHead(node *DLinkedNode) {
	lru.removeNode(node)
	lru.addToHead(node)
}
func (lru *DLinkedList) removeTail() *DLinkedNode {
	node := lru.tail.prev
	lru.removeNode(node)
	return node
}
