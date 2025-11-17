package zskip

import (
	"fmt"
	"math/rand"
)

const defaultProbability = 0.25

type SkipList struct {
	head     *SkipNode
	tail     *SkipNode
	level    int
	maxLevel int
	length   int
}

func NewSkipList(maxLevel int) *SkipList {
	return &SkipList{level: 1, maxLevel: maxLevel, head: newSkipNode(maxLevel, 0, "")}
}

func (s *SkipList) randomLevel() int {
	lvl := 1
	for rand.Float64() < defaultProbability && lvl < s.maxLevel {
		lvl++
	}
	return lvl
}

func (s *SkipList) ZslInsert(score uint64, value string) *SkipNode {
	updates := make([]*SkipNode, s.maxLevel)
	ranks := make([]int, s.maxLevel)
	x := s.head

	for i := s.level - 1; i >= 0; i-- {
		if i != s.level-1 {
			ranks[i] = ranks[i+1]
		}
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && x.levels[i].forward.value < value)) {
			ranks[i] += x.levels[i].span
			x = x.levels[i].forward
		}
		updates[i] = x
	}

	newLevel := s.randomLevel()
	if newLevel > s.level {
		for i := s.level; i < newLevel; i++ {
			ranks[i] = 0
			updates[i] = s.head
			updates[i].levels[i].span = s.length
		}
		s.level = newLevel
	}

	newNode := newSkipNode(newLevel, score, value)

	for i := range newLevel {
		newNode.levels[i].forward = updates[i].levels[i].forward
		updates[i].levels[i].forward = newNode

		newNode.levels[i].span = updates[i].levels[i].span - (ranks[0] - ranks[i])
		updates[i].levels[i].span = (ranks[0] - ranks[i]) + 1
	}

	for i := newLevel; i < s.level; i++ {
		updates[i].levels[i].span++
	}

	if updates[0] != s.head {
		newNode.backward = updates[0]
	} else {
		newNode.backward = nil
	}

	if newNode.levels[0].forward != nil {
		newNode.levels[0].forward.backward = newNode
	} else {
		s.tail = newNode
	}

	s.length++
	return newNode
}

func (s *SkipList) ZslDeleteNode(node *SkipNode, updates []*SkipNode) {
	for i := 0; i < s.level; i++ {
		if updates[i].levels[i].forward == node {
			updates[i].levels[i].span += node.levels[i].span - 1
			updates[i].levels[i].forward = node.levels[i].forward
		} else {
			updates[i].levels[i].span--
		}
	}

	if node.levels[0].forward != nil {
		node.levels[0].forward.backward = node.backward
	} else {
		s.tail = node.backward
	}

	for s.level > 1 && s.head.levels[s.level-1].forward == nil {
		s.level--
	}
	s.length--
}

func (s *SkipList) ZslDelete(score uint64, value string) bool {
	updates := make([]*SkipNode, s.maxLevel)
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil &&
			(x.levels[i].forward.score < score ||
				(x.levels[i].forward.score == score && x.levels[i].forward.value < value)) {
			x = x.levels[i].forward
		}
		updates[i] = x
	}

	target := x.levels[0].forward
	if target != nil && target.score == score && target.value == value {
		s.ZslDeleteNode(target, updates)
		return true
	}
	return false
}

func (s *SkipList) DeleteLast() *SkipNode {
	if s.length == 0 || s.tail == nil {
		return nil
	}

	target := s.tail
	updates := make([]*SkipNode, s.maxLevel)
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && x.levels[i].forward != target {
			x = x.levels[i].forward
		}
		updates[i] = x
	}

	s.ZslDeleteNode(target, updates)
	return target
}

func (s *SkipList) DeleteFirst() *SkipNode {
	if s.length == 0 || s.head.levels[0].forward == nil {
		return nil
	}

	target := s.head.levels[0].forward
	updates := make([]*SkipNode, s.maxLevel)
	for i := 0; i < s.level; i++ {
		updates[i] = s.head
	}
	s.ZslDeleteNode(target, updates)
	return target
}

func (s *SkipList) ZslGetNodeByRank(rank int) *SkipNode {
	if rank < 1 || rank > s.length {
		return nil
	}
	x := s.head
	traversed := 0
	for i := s.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && traversed+x.levels[i].span <= rank {
			traversed += x.levels[i].span
			x = x.levels[i].forward
		}
		if traversed == rank {
			return x
		}
	}
	return nil
}

func (s *SkipList) ZslLen() int { return s.length }

func (s *SkipList) ZslRange(start, end int, reverse bool) []*SkipNode {
	if s.length == 0 || start > end {
		return nil
	}
	if start < 0 {
		start = 0
	}
	if end >= s.length {
		end = s.length - 1
	}
	size := end - start + 1
	nodes := make([]*SkipNode, 0, size)

	if !reverse {
		node := s.ZslGetNodeByRank(start + 1)
		for i := 0; i < size && node != nil; i++ {
			nodes = append(nodes, node)
			node = node.levels[0].forward
		}
	} else {
		node := s.tail
		for i := 0; i < start && node != nil; i++ {
			node = node.backward
		}
		for i := 0; i < size && node != nil; i++ {
			nodes = append(nodes, node)
			node = node.backward
		}
	}
	return nodes
}

func (s *SkipList) TestPrint() {
	fmt.Println("长度: ", s.ZslLen())
	for i := s.level - 1; i > 0; i-- {
		for x := s.head; x != nil && x.levels[i].forward != nil; x = x.levels[i].forward {
			fmt.Printf("Value: %v\tscore: %v\tlevel: %d\n", x.levels[i].forward.value, x.levels[i].forward.score, i)
		}
		fmt.Println()
	}
}

func (s *SkipList) TestPrint2() {
	fmt.Println("长度:", s.ZslLen())

	for x := s.head.levels[0].forward; x != nil; x = x.levels[0].forward {
		fmt.Printf("score=%v value=%v\n", x.score, x.value)
	}
}
