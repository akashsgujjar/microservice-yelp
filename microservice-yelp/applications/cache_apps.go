package applications

import (
	"container/heap"
	"container/list"
	"errors"
	"log"
	"math/rand"
	"sync"

	"gitlab.cs.washington.edu/syslab/cse453-welp/proto/mycache"
)

var (
	ErrItemNotFound = errors.New("mycache: cache miss")
)

// type CacheItem struct {
// 	// Key is CacheItem's key
// 	Key string

// 	// Value is the CacheItem's value
// 	Value []byte
// }

// Cache is a simple Key-Value cache interface.
type Cache interface {
	// Len returns the number of elements in the cache.
	Len() int

	// Get retrieves the value for the specified key.
	Get(key string) (*mycache.CacheItem, error)

	// Set sets the value for the specified key. If the maximum capacity of the cache is exceeded,
	// an eviction policy is applied.
	Set(item *mycache.CacheItem) error

	// Delete deletes the value for the specified key.
	Delete(key string) error

	// Clear removes all items from the cache.
	Clear()
}

// FIFOCacheApp is a simple in-memory FIFO (First-In-First-Out) key-value cache.
type FIFOCacheApp struct {
	data     map[string]*mycache.CacheItem
	order    *list.List // Use a doubly-linked list to maintain FIFO order
	capacity int
	lock     sync.Mutex
}

// NewFIFOCacheApp returns a new FIFO Cache with the specified maximum capacity.
func NewFIFOCacheApp(capacity int) *FIFOCacheApp {
	log.Println("eviction policy: FIFO cache")
	return &FIFOCacheApp{
		data:     make(map[string]*mycache.CacheItem),
		order:    list.New(),
		capacity: capacity,
	}
}

// Len returns the number of elements in the cache.
func (c *FIFOCacheApp) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

// Get retrieves the value for the specified key.
func (c *FIFOCacheApp) Get(key string) (*mycache.CacheItem, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.data[key]
	if !ok {
		return nil, ErrItemNotFound
	}
	return value, nil
}

// Set sets the value for the specified key. If the maximum capacity of the cache is exceeded,
// the oldest key-value pair will be evicted.
func (c *FIFOCacheApp) Set(item *mycache.CacheItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if len(c.data) >= c.capacity {
		// If the cache is full, evict the oldest item (front of the list)
		oldestElement := c.order.Front()
		if oldestElement != nil {
			oldestKey := oldestElement.Value.(string)
			delete(c.data, oldestKey)
			c.order.Remove(oldestElement)
		}
	}

	key := item.Key
	c.data[key] = item
	c.order.PushBack(key) // Add the new key to the back of the list
	return nil
}

// Delete deletes the value for the specified key.
func (c *FIFOCacheApp) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	_, ok := c.data[key]
	if !ok {
		return ErrItemNotFound
	}

	// Remove the key from the data map and the list
	delete(c.data, key)
	for element := c.order.Front(); element != nil; element = element.Next() {
		if element.Value.(string) == key {
			c.order.Remove(element)
			break
		}
	}
	return nil
}

// Clear removes all items from the cache.
func (c *FIFOCacheApp) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = make(map[string]*mycache.CacheItem)
	c.order.Init()
}

// RandomCacheApp is a simple in-memory key-value cache.
type RandomCacheApp struct {
	lock     sync.Mutex // Mutex for protecting concurrent access
	data     map[string]*mycache.CacheItem
	capacity int
}

// NewRandomCacheApp returns a new Cache with the specified maximum capacity.
func NewRandomCacheApp(capacity int) *RandomCacheApp {
	log.Println("eviction policy: random cache")
	return &RandomCacheApp{
		lock:     sync.Mutex{},
		data:     make(map[string]*mycache.CacheItem),
		capacity: capacity,
	}
}

// Len returns the number of elements in the cache.
func (c *RandomCacheApp) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()
	return len(c.data)
}

// Get retrieves the value for the specified key.
func (c *RandomCacheApp) Get(key string) (*mycache.CacheItem, error) {
	c.lock.Lock()
	defer c.lock.Unlock()
	value, ok := c.data[key]
	if !ok {
		return nil, ErrItemNotFound
	}
	return value, nil
}

// Set sets the value for the specified key. If the maximum capacity of the cache is exceeded,
// the oldest key-value pair will be evicted.
func (c *RandomCacheApp) Set(item *mycache.CacheItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	if len(c.data) >= c.capacity {
		_ = c.evictRandomKey()
	}
	key := item.Key
	c.data[key] = item
	return nil
}

// Remove deletes the value for the specified key.
func (c *RandomCacheApp) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, ok := c.data[key]
	if !ok {
		return ErrItemNotFound
	}
	delete(c.data, key)
	return nil
}

// Clear removes all items from the cache.
func (c *RandomCacheApp) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = make(map[string]*mycache.CacheItem)
}

// evictRandomKey deletes a random key-value pair from the cache and returns the evicted key.
func (c *RandomCacheApp) evictRandomKey() string {
	randomKey := ""
	randomIndex := rand.Intn(len(c.data))

	i := 0
	for key := range c.data {
		if i == randomIndex {
			randomKey = key
			break
		}
		i++
	}
	delete(c.data, randomKey)
	return randomKey
}

// LRUCacheApp is a simple in-memory LRU (Least-Recently-Used) key-value cache.
type LRUCacheApp struct {
	capacity int
	data     map[string]*list.Element
	list     *list.List
	lock     sync.Mutex
}

// NewLRUCacheApp creates a new LRUCache with the specified capacity.
func NewLRUCacheApp(capacity int) *LRUCacheApp {
	log.Println("eviction policy: LRU cache")
	return &LRUCacheApp{
		capacity: capacity,
		data:     make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Get retrieves a value from the cache based on the key.
func (c *LRUCacheApp) Get(key string) (*mycache.CacheItem, error) {
	// log.Println("LRU Cache GET")
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, ok := c.data[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*mycache.CacheItem), nil
	}
	return nil, ErrItemNotFound
}

// Set inserts or updates a value in the cache.
func (c *LRUCacheApp) Set(item *mycache.CacheItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	key := item.Key
	// Element already in cache
	if elem, ok := c.data[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value = item
	} else {
		if c.list.Len() >= c.capacity {
			// Remove the least recently used item
			lastElem := c.list.Back()
			delete(c.data, lastElem.Value.(*mycache.CacheItem).Key)
			c.list.Remove(lastElem)
		}

		newElem := c.list.PushFront(item)
		c.data[key] = newElem
	}
	return nil
}

// Delete removes a value from the cache based on the key.
func (c *LRUCacheApp) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, ok := c.data[key]; ok {
		c.list.Remove(elem)
		delete(c.data, key)
		return nil
	}
	return ErrItemNotFound
}

// Clear removes all items from the cache.
func (c *LRUCacheApp) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = make(map[string]*list.Element)
	c.list.Init()
}

// Len returns the number of items currently in the cache.
func (c *LRUCacheApp) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return len(c.data)
}

// LFUElement represents an item stored in the LFU cache.
type LFUNode struct {
	Value     *mycache.CacheItem
	Frequency int // Frequency of item access
	index     int // Index in the frequency heap
}

type frequencyHeap []*LFUNode

func (h frequencyHeap) Len() int {
	return len(h)
}

// Use < so that our heap behaves as min-heap
func (h frequencyHeap) Less(i, j int) bool {
	return h[i].Frequency < h[j].Frequency
}

func (h frequencyHeap) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
	h[i].index = i
	h[j].index = j
}

func (h *frequencyHeap) Push(x interface{}) {
	item := x.(*LFUNode)
	item.index = len(*h)
	*h = append(*h, item)
}

func (h *frequencyHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	item.index = -1 // for safety
	*h = old[0 : n-1]

	// TODO: should i do this?
	old[n-1] = nil
	return item
}

// LFUCacheApp is a concurrency-safe LFU cache.
type LFUCacheApp struct {
	capacity int
	cache    map[string]*LFUNode
	freqHeap *frequencyHeap
	mu       sync.Mutex
}

// NewLFUCacheApp creates a new LFUCache with the specified capacity.
func NewLFUCacheApp(capacity int) *LFUCacheApp {
	log.Println("eviction policy: LFU cache")
	freqHeap := make(frequencyHeap, 0)
	heap.Init(&freqHeap)
	return &LFUCacheApp{
		capacity: capacity,
		cache:    make(map[string]*LFUNode),
		freqHeap: &freqHeap,
	}
}

// updateFrequency updates the frequency of an LFU node in the frequency heap.
func (c *LFUCacheApp) updateFrequency(node *LFUNode) {
	node.Frequency++
	heap.Fix(c.freqHeap, node.index)
}

// Get retrieves a value from the cache based on the key.
func (c *LFUCacheApp) Get(key string) (*mycache.CacheItem, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if node, ok := c.cache[key]; ok {
		c.updateFrequency(node)
		return node.Value, nil
	}
	return nil, ErrItemNotFound
}

// Set inserts or updates a value in the cache.
func (c *LFUCacheApp) Set(item *mycache.CacheItem) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	key := item.Key
	if existingItem, ok := c.cache[key]; ok {
		existingItem.Value = item
		c.updateFrequency(existingItem)
	} else {
		if len(c.cache) >= c.capacity {
			leastFreqItem := heap.Pop(c.freqHeap).(*LFUNode)
			delete(c.cache, leastFreqItem.Value.Key)
		}

		newNode := &LFUNode{Value: item, Frequency: 1}
		c.cache[key] = newNode
		heap.Push(c.freqHeap, newNode)
	}
	return nil
}

// Delete removes a value from the cache based on the key.
func (c *LFUCacheApp) Delete(key string) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if item, ok := c.cache[key]; ok {
		heap.Remove(c.freqHeap, item.index)
		delete(c.cache, key)
		return nil
	}
	return ErrItemNotFound
}

// Clear removes all items from the cache.
func (c *LFUCacheApp) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO: add clear heap?
	c.cache = make(map[string]*LFUNode)
	heap.Init(c.freqHeap)
}

// Len returns the number of items currently in the cache.
func (c *LFUCacheApp) Len() int {
	c.mu.Lock()
	defer c.mu.Unlock()

	return len(c.cache)
}

// MRUCacheApp is a simple in-memory MRU (Most-Recently-Used) key-value cache.
type MRUCacheApp struct {
	capacity int
	data     map[string]*list.Element
	list     *list.List
	lock     sync.Mutex
}

// NewCacheApp creates a new MRUCache with the specified capacity.
func NewCacheApp(capacity int) *MRUCacheApp {
	log.Println("eviction policy: MRU cache")
	return &MRUCacheApp{
		capacity: capacity,
		data:     make(map[string]*list.Element),
		list:     list.New(),
	}
}

// Get retrieves a value from the cache based on the key.
func (c *MRUCacheApp) Get(key string) (*mycache.CacheItem, error) {
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, ok := c.data[key]; ok {
		c.list.MoveToFront(elem)
		return elem.Value.(*mycache.CacheItem), nil
	}
	return nil, ErrItemNotFound
}

// Set inserts or updates a value in the cache.
func (c *MRUCacheApp) Set(item *mycache.CacheItem) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	key := item.Key
	// Element already in cache
	if elem, ok := c.data[key]; ok {
		c.list.MoveToFront(elem)
		elem.Value = item
	} else {
		if c.list.Len() >= c.capacity {
			// Remove the most recently used item
			firstElem := c.list.Front()
			delete(c.data, firstElem.Value.(*mycache.CacheItem).Key)
			c.list.Remove(firstElem)
		}

		newElem := c.list.PushFront(item)
		c.data[key] = newElem
	}
	return nil
}

// Delete removes a value from the cache based on the key.
func (c *MRUCacheApp) Delete(key string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	if elem, ok := c.data[key]; ok {
		c.list.Remove(elem)
		delete(c.data, key)
		return nil
	}
	return ErrItemNotFound
}

// Clear removes all items from the cache.
func (c *MRUCacheApp) Clear() {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.data = make(map[string]*list.Element)
	c.list.Init()
}

// Len returns the number of items currently in the cache.
func (c *MRUCacheApp) Len() int {
	c.lock.Lock()
	defer c.lock.Unlock()

	return len(c.data)
}
