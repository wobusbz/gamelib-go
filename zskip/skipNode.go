package zskip

type SkipLevel struct {
	forward *SkipNode
	span    int
}

type SkipNode struct {
	score    uint64
	value    string
	backward *SkipNode
	levels   []*SkipLevel
}

func newSkipNode(level int, score uint64, value string) *SkipNode {
	node := &SkipNode{score: score, value: value, levels: make([]*SkipLevel, level)}
	for i := range level {
		node.levels[i] = &SkipLevel{}
	}
	return node
}
