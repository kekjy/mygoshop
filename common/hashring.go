package common

import (
	"hash/crc32"
	"sort"
	"sync"
)

type HashRing struct {
	nodes       []uint32                 // 哈希环上的所有节点
	nodemap     map[uint32]string        // 哈希映射
	HashFn      func(data []byte) uint32 // 哈希函数
	virtualNode int
	sync.RWMutex
}

func NewHashRing() *HashRing {
	return &HashRing{
		nodes:       make([]uint32, 0),
		nodemap:     make(map[uint32]string),
		HashFn:      crc32.ChecksumIEEE,
		virtualNode: 20,
	}
}

func (hr *HashRing) Add(id string) {
	hr.Lock()
	defer hr.Unlock()
	for range hr.virtualNode {
		hash := hr.HashFn([]byte(id))
		hr.nodes = append(hr.nodes, hash)
		hr.nodemap[hash] = id
	}
	sort.Slice(hr.nodes, func(i, j int) bool {
		return hr.nodes[i] < hr.nodes[j]
	})
}

func (hr *HashRing) Get(key string) string {
	hr.RLock()
	defer hr.RUnlock()
	if len(hr.nodes) == 0 {
		return "EMPTY"
	}
	hash := hr.HashFn([]byte(key))
	idx := sort.Search(len(hr.nodes), func(i int) bool {
		return hr.nodes[i] >= hash
	})

	if idx == len(hr.nodes) {
		return hr.nodemap[hr.nodes[0]]
	}
	return hr.nodemap[hr.nodes[idx]]
}
