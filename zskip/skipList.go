package zskip

import (
	"fmt"
	"math/rand"
)

type SkipList struct {
	head     *SkipNode
	tail     *SkipNode
	level    int
	maxLevel int
	length   int
}

func NewSkipList(level int) *SkipList {
	return &SkipList{level: 1, maxLevel: level, length: 0, head: newSkipNode(level, 0, "")}
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
	level := s.randomLevel()
	if level > s.level {
		for i := s.level; i < level; i++ {
			ranks[i] = 0
			updates[i] = s.head
			updates[i].levels[i].span = s.length
		}
		s.level = level
	}
	newNodes := newSkipNode(level, score, value)
	for i := range level {
		newNodes.levels[i].forward = updates[i].levels[i].forward
		updates[i].levels[i].forward = newNodes
		newNodes.levels[i].span = updates[i].levels[i].span - (ranks[0] - ranks[i])
		updates[i].levels[i].span = (ranks[0] - ranks[i]) + 1
	}
	for i := level; i < s.level; i++ {
		updates[i].levels[i].span++
	}
	if updates[0] != s.head {
		newNodes.backward = updates[0]
	} else {
		newNodes.backward = nil
	}
	if newNodes.levels[0].forward != nil {
		newNodes.levels[0].forward.backward = newNodes
	} else {
		s.tail = newNodes
	}
	s.length++
	return newNodes
}

func (s *SkipList) ZslDelete(score uint64) {
	updates := make([]*SkipNode, s.maxLevel)
	x := s.head
	for i := s.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && x.levels[i].forward.score < score {
			x = x.levels[i].forward
		}
		updates[i] = x
	}
	x = x.levels[0].forward
	if x != nil && x.score == score {
		s.ZslDeleteNode(x, updates)
	}
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
	node = nil
}

func (s *SkipList) DeleteLast() *SkipNode {
	if s.length == 0 {
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

	for i := 0; i < s.level; i++ {
		if updates[i].levels[i].forward == target {
			updates[i].levels[i].span += target.levels[i].span - 1
			updates[i].levels[i].forward = target.levels[i].forward
		} else {
			updates[i].levels[i].span -= 1
		}
	}

	if target.backward != nil {
		s.tail = target.backward
		s.tail.levels[0].forward = nil
	} else {
		s.tail = nil
	}

	for s.level > 1 && s.head.levels[s.level-1].forward == nil {
		s.level--
	}

	s.length--
	return target
}

func (s *SkipList) DeleteFirst() *SkipNode {
	if s.length == 0 {
		return nil
	}

	target := s.head.levels[0].forward
	if target == nil {
		return nil
	}

	updates := make([]*SkipNode, s.maxLevel)

	for i := 0; i < s.level; i++ {
		if s.head.levels[i].forward == target {
			s.head.levels[i].forward = target.levels[i].forward
			s.head.levels[i].span -= 1
		} else {
			s.head.levels[i].span -= 1
		}
		updates[i] = s.head
	}

	if target.levels[0].forward != nil {
		target.levels[0].forward.backward = nil
	} else {
		s.tail = nil
	}

	for s.level > 1 && s.head.levels[s.level-1].forward == nil {
		s.level--
	}

	s.length--
	return target
}

func (s *SkipList) ZslGetNodeByRank(rank int) *SkipNode {
	x := s.head
	var traversed int
	for i := s.level - 1; i >= 0; i-- {
		for x.levels[i].forward != nil && (traversed+x.levels[i].span <= rank) {
			traversed += x.levels[i].span
			x = x.levels[i].forward
		}
		if traversed == rank {
			return x
		}
	}
	return nil
}

func (s *SkipList) ZslGetNodeByRankNo(score uint64, value string) int {
	x := s.head
	rank := 0
	for i := s.level - 1; i > 0; i-- {
		for x.levels[i].forward != nil && (x.levels[i].forward.score < score ||
			(x.levels[i].forward.score == score && x.levels[i].forward.value <= value)) {
			rank += x.levels[i].span
			x = x.levels[i].forward
		}
		if x != nil && x.score == score {
			return rank
		}
	}
	return 0
}

func (s *SkipList) ZslGetNodeBySocre(score uint64) *SkipNode {
	x := s.head
	rank := 0
	for i := s.level - 1; i > 0; i-- {
		for x.levels[i].forward != nil && x.levels[i].forward.score == score {
			rank += x.levels[i].span
			x = x.levels[i].forward
		}
		if x != nil && x.score == score {
			return x
		}
	}
	return nil
}

func (s *SkipList) ZslRange(start, end int, reverse bool) []*SkipNode {
	if s.length == 0 {
		return nil
	}
	if start < 0 {
		start = 0
	}
	if end >= s.length {
		end = s.length - 1
	}
	if start > end {
		return nil
	}

	var nodes []*SkipNode

	if !reverse {
		node := s.ZslGetNodeByRank(start + 1)
		for i := start; i <= end && node != nil; i++ {
			nodes = append(nodes, node)
			node = node.levels[0].forward
		}
	} else {
		node := s.tail
		for i := 0; i < start && node != nil; i++ {
			node = node.backward
		}
		for i := start; i <= end && node != nil; i++ {
			nodes = append(nodes, node)
			node = node.backward
		}
	}

	return nodes
}

func (s *SkipList) ZslLen() int { return s.length }

func (s *SkipList) randomLevel() int {
	level := 1
	for rand.Float64() < 0.25 && level < s.maxLevel {
		level++
	}
	return level
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
