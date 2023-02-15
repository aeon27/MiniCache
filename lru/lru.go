package lru

import "container/list"

//
type Value interface {
	Len() int //用于返回值所占用的内存大小
}

// 键值对entry是双向链表节点的数据类型，在链表中仍保存key的好处在于，淘汰链表的节点时，可以以key去删除字典里对应的映射。
type entry struct {
	key   string
	value Value
}

// Cache是基于LRU实现的缓存。并发访问不安全。
type Cache struct {
	maxBytes  int64                         // 允许使用的最大内存
	uesdBytes int64                         // 当前已使用的内存
	ll        *list.List                    // 采用Go标准库的双向链表
	cache     map[string]*list.Element      // 缓存
	OnEvicted func(key string, value Value) // 是某条记录被移除时的回调函数，可以为nil
}

// Cache的实例化函数
func NewCache(maxBytes int64, onEvicted func(string, Value)) *Cache {
	return &Cache{
		maxBytes:  maxBytes,
		ll:        list.New(),
		cache:     make(map[string]*list.Element),
		OnEvicted: onEvicted,
	}
}

// 查找
func (c *Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	}
	return
}

// 删除
func (c *Cache) RemoveOldest() {
	ele := c.ll.Back()
	if ele != nil {
		c.ll.Remove(ele)
		kv := ele.Value.(*entry)
		delete(c.cache, kv.key)
		c.uesdBytes -= int64(len(kv.key)) + int64(kv.value.Len())

		// 如果回调函数存在，则回调
		if c.OnEvicted != nil {
			c.OnEvicted(kv.key, kv.value)
		}
	}
}

// 新增
func (c *Cache) Add(key string, value Value) bool {
	if ele, ok := c.cache[key]; ok {
		return false
	} else {
		ele = c.ll.PushFront(&entry{key, value})
		c.cache[key] = ele
		c.uesdBytes += int64(len(key)) + int64(value.Len())
	}

	// 检查缓存是否溢出，若溢出则进行淘汰
	c.CheckCache()

	return true
}

// 修改
func (c *Cache) Update(key string, value Value) bool {
	if ele, ok := c.cache[key]; !ok {
		return false
	} else {
		c.ll.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.uesdBytes += int64(value.Len()) - int64(kv.value.Len())
		kv.value = value
	}

	// 检查缓存是否溢出，若溢出则进行淘汰
	c.CheckCache()

	return true
}

// 用于检查缓存是否溢出
func (c *Cache) CheckCache() {
	for c.maxBytes != 0 && c.maxBytes < c.uesdBytes {
		c.RemoveOldest()
	}
}

func (c *Cache) Len() int {
	return c.ll.Len()
}
