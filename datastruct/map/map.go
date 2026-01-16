package database

import "sync"

// MyConcurrentMap 是一个基于 RWMutex 的并发安全 map
// 读操作并发，写操作互斥
type MyConcurrentMap struct {
	sync.RWMutex
	mp map[int]int
}

// MakeMyConcurrentMap 创建一个并发安全的 map
func MakeMyConcurrentMap() *MyConcurrentMap {
	return &MyConcurrentMap{
		mp: make(map[int]int),
	}
}

// Get 读取 key，对应读锁
func (m *MyConcurrentMap) Get(key int) (int, bool) {
	m.RLock()
	defer m.RUnlock()

	v, ok := m.mp[key]
	return v, ok
}

// Put 写入 key-value，对应写锁
func (m *MyConcurrentMap) Put(key int, val int) {
	m.Lock()
	defer m.Unlock()

	m.mp[key] = val
}

// Delete 删除 key，对应写锁
func (m *MyConcurrentMap) Delete(key int) {
	m.Lock()
	defer m.Unlock()

	delete(m.mp, key)
}
