/**
 *
 * LRU（Least recently used，最近最少使用）算法根据数据的历史访问记录来进行淘汰数据，其核心思想是“如果数据最近被访问过，那么将来被访问的几率也更高”。
 * LRU 最常见的实现是使用一个链表保存缓存数据。详细算法如下：
 * 1. 新数据插入到链表头部；
 * 2. 每当缓存命中（即缓存数据被访问），则将数据移到链表头部；
 * 3. 当链表满的时候，将链表尾部的数据丢弃。
 * @FileName: lru.go
 * @Auther : 杨雪波
 * @Email : yangxuebo@yeah.net
 * @CreateTime: 2018-07-24 10:44:32
 * @UpdateTime: 2018-07-24 10:44:54
 * Copyright@2018 版权所有
 */

package lruCache

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

//当有元素被剔除时，执行这个回调
type EvictCallback func(key interface{}, value interface{})

type LRU struct {
	size      int
	evictList *list.List
	items     map[interface{}]*list.Element
	onEvict   EvictCallback
	timeout   time.Duration
	lock      sync.RWMutex
}

type entry struct {
	key     interface{}
	value   interface{}
	expires time.Time
}

func NewLRU(size int, timeout time.Duration, onEvict EvictCallback) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c := &LRU{
		size:      size,
		evictList: list.New(),
		items:     make(map[interface{}]*list.Element),
		onEvict:   onEvict,
		timeout:   timeout,
		lock:      sync.RWMutex{},
	}
	return c, nil
}

func (c *LRU) Add(key, value interface{}) (evicted bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	// 如果已经存在key值，则将这个key移动到前边，将过期时间设置为下一个时间段
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		ent.Value.(*entry).value = value
		ent.Value.(*entry).expires = time.Now().Add(c.timeout)
		return false
	}

	// 如果不存在，将新元素添加到列表头部
	ent := &entry{key, value, time.Now().Add(c.timeout)}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry

	evict := c.evictList.Len() > c.size
	//如果超过了限制，则移除最老的元素
	if evict {
		c.removeOldest()
	}
	return evict
}

func (c *LRU) Get(key interface{}) (value interface{}, ok bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if ent, ok := c.items[key]; ok {
		if ent.Value.(*entry).expires.After(time.Now()) {
			c.evictList.MoveToFront(ent)
			return ent.Value.(*entry).value, true
		} else {
			//删除已过期的缓存
			c.removeElement(ent)
		}
	}
	return
}

/**
 * 判断缓存中是否包含指定key的元素
 */
func (c *LRU) Contains(key interface{}) (ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	_, ok = c.items[key]
	return ok
}

/**
 * 移除缓存中指定key的元素
 */
func (c *LRU) Remove(key interface{}) (present bool) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

/**
 *清楚缓存中所以的元素
 */
func (c *LRU) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()
	for k, v := range c.items {
		if c.onEvict != nil {
			c.onEvict(k, v.Value.(*entry).value)
		}
		delete(c.items, k)
	}
	c.evictList.Init()
}

/**
 *Keys返回缓存中所有键的一个切片，顺序从最老的到最新的
 */
func (c *LRU) Keys() []interface{} {
	c.lock.RLock()
	defer c.lock.RUnlock()
	keys := make([]interface{}, len(c.items))
	i := 0
	for ent := c.evictList.Back(); ent != nil; ent = ent.Prev() {
		keys[i] = ent.Value.(*entry).key
		i++
	}
	return keys
}

/**
*从缓存中删除最老元素
 */
func (c *LRU) RemoveOldest() (key interface{}, value interface{}, ok bool) {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
		kv := ent.Value.(*entry)
		return kv.key, kv.value, true
	}
	return nil, nil, false
}

/**
 *  从缓存中移除最老的元素
 */
func (c *LRU) removeOldest() {
	ent := c.evictList.Back()
	if ent != nil {
		c.removeElement(ent)
	}
}

/**
* 从缓存中删除指定的列表元素，并通知回调函数
 */
func (c *LRU) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	if c.onEvict != nil {
		c.onEvict(kv.key, kv.value)
	}
}

func (c *LRU) Len() int {
	c.lock.RLock()
	defer c.lock.RUnlock()
	return c.evictList.Len()
}
