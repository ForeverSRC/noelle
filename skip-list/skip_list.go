package skip_list

import (
	"github.com/ForeverSRC/noelle/skip-list/level"
	"golang.org/x/exp/constraints"
)

type Element[K constraints.Ordered, V any] struct {
	Key   K
	Value V
}

type node[K constraints.Ordered, V any] struct {
	forward []*node[K, V]
	elem    Element[K, V]
}

func (n *node[K, V]) helpGC() {
	n.forward = nil
}

type SkipList[K constraints.Ordered, V any] struct {
	head           *node[K, V]
	level          int
	maxLevel       int
	levelGenerator level.Generator
}

func NewSkipList[K constraints.Ordered, V any](lg level.Generator) *SkipList[K, V] {
	sl := &SkipList[K, V]{
		level:          1,
		levelGenerator: lg,
		maxLevel:       lg.MaxLevel(),
	}

	sl.head = &node[K, V]{
		forward: make([]*node[K, V], sl.maxLevel+1),
	}

	return sl
}

func (sl *SkipList[K, V]) Search(searchKey K) (res V, found bool) {
	x := sl.head
	for i := sl.level; i >= 1; i-- {
		for x != nil && x.forward[i] != nil && x.forward[i].elem.Key < searchKey {
			x = x.forward[i]
		}
	}

	x = x.forward[1]
	if x != nil && x.elem.Key == searchKey {
		return x.elem.Value, true
	} else {
		return res, false
	}

}

func (sl *SkipList[K, V]) Insert(key K, value V) {
	searchKey := key
	update := make([]*node[K, V], sl.maxLevel+1)

	x := sl.head
	for i := sl.level; i >= 1; i-- {
		for x != nil && x.forward[i] != nil && x.forward[i].elem.Key < searchKey {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[1]
	elem := Element[K, V]{Key: key, Value: value}
	if x != nil && x.elem.Key == searchKey {
		x.elem = elem
	} else {
		sl.addNode(elem, update)
	}

}

func (sl *SkipList[K, V]) Delete(key K) {
	searchKey := key
	update := make([]*node[K, V], sl.maxLevel+1)

	x := sl.head
	for i := sl.level; i >= 1; i-- {
		for x != nil && x.forward[i] != nil && x.forward[i].elem.Key < searchKey {
			x = x.forward[i]
		}
		update[i] = x
	}

	x = x.forward[1]

	if x != nil && x.elem.Key == searchKey {
		for i := 1; i <= sl.level; i++ {
			if update[i].forward[i] != x {
				break
			}

			update[i].forward[i] = x.forward[i]
		}

		x.helpGC()

		for sl.level > 1 && sl.head.forward[sl.level] == nil {
			sl.level--
		}
	}

}

func (sl *SkipList[K, V]) addNode(elem Element[K, V], update []*node[K, V]) {
	lvl := sl.levelGenerator.GetLevel()
	if lvl > sl.level {
		for i := sl.level + 1; i <= lvl; i++ {
			update[i] = sl.head
		}

		sl.level = lvl
	}

	n := &node[K, V]{elem: elem}
	n.forward = make([]*node[K, V], lvl+1)

	for i := 1; i <= lvl; i++ {
		n.forward[i] = update[i].forward[i]
		update[i].forward[i] = n
	}

}

func (sl *SkipList[K, V]) LevelIteration(lvl int) []Element[K, V] {
	if lvl > sl.level {
		return []Element[K, V]{}
	}

	res := make([]Element[K, V], 0, sl.level*2)

	iter := sl.head.forward[lvl]

	for iter != nil {
		res = append(res, iter.elem)
		iter = iter.forward[lvl]
	}

	return res

}

func (sl *SkipList[K, V]) GetLevel() int {
	return sl.level
}
