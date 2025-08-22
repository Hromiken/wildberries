package cache

import (
	"github.com/google/uuid"
	"sync"
	"time"

	"order-notification/internal/entity"
)

// Node : элемент двусвязного списка.
type Node struct {
	key       uuid.UUID
	value     *entity.Order
	expiresAt time.Time
	prev      *Node
	next      *Node
}

// OrderCache : LRU-кеш с TTL.
type OrderCache struct {
	data     map[uuid.UUID]*Node
	capacity int
	head     *Node
	tail     *Node
	mu       sync.Mutex
	ttl      time.Duration
}

// NewCache : инициализация нового кеша.
func NewCache(cap int, ttl time.Duration) *OrderCache {
	o := &OrderCache{
		capacity: cap,
		data:     make(map[uuid.UUID]*Node),
		ttl:      ttl,
	}
	o.head = &Node{}
	o.tail = &Node{}
	o.head.next = o.tail
	o.tail.prev = o.head
	return o
}

// Get : возвращает заказ по ключу, если он есть и не протух по TTL.
func (o *OrderCache) Get(key uuid.UUID) *entity.Order {
	o.mu.Lock()
	defer o.mu.Unlock()

	node, ok := o.data[key]
	if !ok {
		return nil
	}

	// проверяем TTL
	if time.Now().After(node.expiresAt) {
		o.removeNode(node)
		delete(o.data, key)
		return nil
	}

	// «освежаем» позицию
	o.moveToHead(node)
	return node.value
}

// Set : кладёт заказ в кеш (обновляет, если уже был).
func (o *OrderCache) Set(key uuid.UUID, order *entity.Order) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if node, ok := o.data[key]; ok {
		node.value = order
		node.expiresAt = time.Now().Add(o.ttl)
		o.moveToHead(node)
		return
	}

	// если переполнено — удаляем самый старый
	if len(o.data) >= o.capacity {
		old := o.removeTail()
		delete(o.data, old.key)
	}

	// создаём новый узел
	newNode := &Node{
		key:       key,
		value:     order,
		expiresAt: time.Now().Add(o.ttl),
	}
	o.data[key] = newNode
	o.addToHead(newNode)
}

// addToHead : вставить узел сразу после head.
func (o *OrderCache) addToHead(n *Node) {
	n.prev = o.head
	n.next = o.head.next
	o.head.next.prev = n
	o.head.next = n
}

// removeNode : удалить узел.
func (o *OrderCache) removeNode(n *Node) {
	n.prev.next = n.next
	n.next.prev = n.prev
}

// removeTail : удалить самый старый элемент.
func (o *OrderCache) removeTail() *Node {
	n := o.tail.prev
	o.removeNode(n)
	return n
}

// moveToHead : переместить узел в начало.
func (o *OrderCache) moveToHead(n *Node) {
	o.removeNode(n)
	o.addToHead(n)
}
