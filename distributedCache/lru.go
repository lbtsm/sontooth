package distributedCache

import (
	"container/list"
	"testing"
)

type cache struct {
	maxByte, usedByte int64 // 最大储存字节，已经使用的字节数量
	twoList           list.List
	data              map[string]*list.Element // 只要是用到map，一定要考虑并发
}

type entry struct {
	key   string
	value Value
}

type Value interface {
	Len() int
}

func New(maxByte int64) *cache {
	c := &cache{
		maxByte: maxByte,
		twoList: list.List{},
		data:    make(map[string]*list.Element),
	}
	return c
}

func (c *cache) Get(key string) (Value, bool) {
	ele, ok := c.data[key]
	if !ok {
		return nil, false
	}
	c.twoList.MoveToFront(ele) // 将这个元素移动到队列尾部
	return ele.Value.(*entry).value, ok
}

func (c *cache) Add(key string, value Value) {
	ele, ok := c.data[key]
	if !ok {
		push := c.twoList.PushFront(&entry{key: key, value: value})
		c.data[key] = push
		c.usedByte += int64(value.Len()) + int64(len(key))
	} else {
		c.twoList.MoveToFront(ele) // 移动到最后面
		entry := ele.Value.(*entry)
		c.usedByte += int64(value.Len()) - int64(entry.value.Len()) // 这里会不会产生并发
		entry.value = value
	}
	// 判断是否需要进行删除，近期很少使用的缓存
	for c.maxByte > 0 && c.maxByte < c.usedByte {
		c.lru()
	}
}

func (c *cache) lru() {
	ele := c.twoList.Back()
	if ele == nil {
		return
	}
	c.twoList.Remove(ele)
	entry := ele.Value.(*entry)
	delete(c.data, entry.key)
	c.usedByte -= int64(len(entry.key)) + int64(entry.value.Len())
}

func (c *cache) Len() int {
	return c.twoList.Len()
}

func (c *cache) BianLi(ele *list.Element, t testing.TB) {
	if ele.Next() != nil {
		c.BianLi(ele.Next(), t)
	}
}
