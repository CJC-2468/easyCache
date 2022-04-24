// Package consistenthash provides an implementation of a ring hash.
package consistenthash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

// Hash maps bytes to uint32
//这里就是把数据映射到2的32次方的那个虚拟环上
//定义了函数类型 Hash，采取依赖注入的方式，允许用于替换成自定义的 Hash 函数，也方便测试时替换，
type Hash func(data []byte) uint32

// Map constains all hashed keys
//一致性哈希算法的主数据结构
type Map struct {
	hash     Hash
	replicas int            //虚拟节点倍数
	keys     []int          //sorted，哈希环，存放的是所有虚拟节点
	hashMap  map[int]string //虚拟节点与真实节点的映射表,键是虚拟节点的哈希值，值是真实节点的名称。
}

// New creates a Map instance
//允许自定义虚拟节点倍数和 Hash 函数
func New(replicas int, f Hash) *Map {
	m := &Map{
		hash:     f,
		replicas: replicas,
		hashMap:  make(map[int]string),
	}
	if m.hash == nil { //默认为 crc32.ChecksumIEEE 算法。
		m.hash = crc32.ChecksumIEEE
	}
	return m
}

// Add adds some keys to the hash.
//传入 0 或 多个真实节点的名称
//对每一个真实节点 key，对应创建 m.replicas 个虚拟节点，虚拟节点的名称是strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点。
//使用 m.hash() 计算虚拟节点的哈希值，使用 append(m.keys, hash) 添加到环上。
//在 hashMap 中增加虚拟节点和真实节点的映射关系。
//最后一步，环上的哈希值排序。
func (m *Map) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < m.replicas; i++ {
			hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
			m.keys = append(m.keys, hash)
			m.hashMap[hash] = key
		}
	}
	//环上的哈希值排序。
	sort.Ints(m.keys)
}

// 根据key获取最近的下一个真实节点名称
func (m *Map) Get(key string) string {
	if len(m.keys) == 0 {
		return ""
	}
	//计算 key 的哈希值。
	hash := int(m.hash([]byte(key)))

	//顺时针找到第一个匹配的虚拟节点的下标 idx
	//该函数使用二分查找的方法，会从[0, n)中取出一个值index，
	//index为[0, n)中最小的使函数f(index)为True的值，并且f(index+1)也为True。
	//如果无法找到该index值，则该方法为返回n，etcd中的键值范围查询就用到了该方法。
	index := sort.Search(len(m.keys), func(i int) bool {
		return m.keys[i] >= hash
	})

	//如果 idx == len(m.keys)，说明应选择 m.keys[0]，因为 m.keys 是一个环状结构，所以用取余数的方式来处理这种情况。
	//通过 hashMap 映射得到真实的节点。
	return m.hashMap[m.keys[index%len(m.keys)]]
}

// Remove use to remove a key and its virtual keys on the ring and map
func (m *Map) Remove(key string) {
	for i := 0; i < m.replicas; i++ {
		//虚拟节点的名称是：strconv.Itoa(i) + key，即通过添加编号的方式区分不同虚拟节点。
		hash := int(m.hash([]byte(strconv.Itoa(i) + key)))
		index := sort.SearchInts(m.keys, hash)
		m.keys = append(m.keys[:index], m.keys[index+1:]...)
		delete(m.hashMap, hash)
	}
}
