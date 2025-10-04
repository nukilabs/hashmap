// Package hashmap provides a hash table implementation with case-insensitive
// string hashing, matching Chromium's WTF HashMap behavior.
package hashmap

import (
	"iter"

	"github.com/nukilabs/hashmap/traits"
)

const (
	initialCapacity = 8
	maximumLoad     = 2 // Expands at 50% load factor
)

// Pair represents a key-value pair stored in the hash table.
type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

// HashMap is a hash table using quadratic probing for collision resolution
// and case-insensitive hashing for string keys.
type HashMap[K comparable, V any] struct {
	table    []*Pair[K, V]
	size     int
	capacity int
}

// New creates a new HashMap with the default initial capacity.
func New[K comparable, V any]() *HashMap[K, V] {
	return &HashMap[K, V]{
		table:    make([]*Pair[K, V], initialCapacity),
		capacity: initialCapacity,
	}
}

// hash computes the hash value for a key.
// For strings, uses case-insensitive hashing.
func (h *HashMap[K, V]) hash(key K) uint32 {
	switch k := any(key).(type) {
	case string:
		return traits.CaseFoldingHash(k)
	default:
		return 0
	}
}

// index returns the bucket index for a hash value.
func (h *HashMap[K, V]) index(hash uint32) int {
	return int(hash & uint32(h.capacity-1))
}

// find locates the slot for a key using quadratic probing.
// Returns the index and whether the key was found.
func (h *HashMap[K, V]) find(key K) (int, bool) {
	hash := h.hash(key)
	idx := h.index(hash)
	count := 0

	for {
		if h.table[idx] == nil {
			return idx, false
		}

		if h.table[idx].Key == key {
			return idx, true
		}

		count++
		if count >= h.capacity {
			break
		}
		idx = (idx + count) & (h.capacity - 1)
	}

	return idx, false
}

// rehash grows the table and rehashes all existing elements.
func (h *HashMap[K, V]) rehash() {
	old := h.table
	h.capacity *= 2
	h.table = make([]*Pair[K, V], h.capacity)
	h.size = 0

	for _, pair := range old {
		if pair != nil {
			h.Set(pair.Key, pair.Value)
		}
	}
}

// Set inserts or updates a key-value pair.
// If the key exists, only the value is updated.
// If the key is new, both key and value are inserted.
func (h *HashMap[K, V]) Set(key K, value V) {
	if (h.size+1)*maximumLoad >= h.capacity {
		h.rehash()
	}

	idx, found := h.find(key)
	if found {
		h.table[idx].Value = value
		return
	}

	h.table[idx] = &Pair[K, V]{
		Key:   key,
		Value: value,
	}
	h.size++
}

// Get retrieves the value for a key.
// Returns the value and true if found, zero value and false otherwise.
func (h *HashMap[K, V]) Get(key K) (V, bool) {
	idx, found := h.find(key)
	if !found {
		var zero V
		return zero, false
	}
	return h.table[idx].Value, true
}

// Contains checks whether a key exists in the map.
func (h *HashMap[K, V]) Contains(key K) bool {
	_, found := h.find(key)
	return found
}

// Delete removes a key-value pair from the map.
// Returns true if the key was found and deleted.
func (h *HashMap[K, V]) Delete(key K) bool {
	idx, found := h.find(key)
	if !found {
		return false
	}

	h.table[idx] = nil
	h.size--
	return true
}

// Clear removes all elements from the map.
func (h *HashMap[K, V]) Clear() {
	h.table = make([]*Pair[K, V], initialCapacity)
	h.capacity = initialCapacity
	h.size = 0
}

// Size returns the number of key-value pairs in the map.
func (h *HashMap[K, V]) Size() int {
	return h.size
}

// Capacity returns the current capacity of the underlying table.
func (h *HashMap[K, V]) Capacity() int {
	return h.capacity
}

// Iter returns an iterator over key-value pairs.
func (h *HashMap[K, V]) Iter() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for _, pair := range h.table {
			if pair != nil {
				if !yield(pair.Key, pair.Value) {
					return
				}
			}
		}
	}
}
