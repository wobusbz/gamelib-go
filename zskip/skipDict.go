package zskip

type ZskDict struct {
	dict  map[string]uint64
	sl    *SkipList
	level int
}

const _defaultLevel = 10

var G_ZskDict = NewZskDict(_defaultLevel)

func NewZskDict(level int) *ZskDict {
	return &ZskDict{dict: map[string]uint64{}, sl: NewSkipList(level), level: level}
}

func (z *ZskDict) ZslSet(k string, score uint64) {
	if _, ok := z.dict[k]; ok {
		z.sl.ZslDelete(score, k)
	}
	z.sl.ZslInsert(score, k)
	z.dict[k] = score
}

func (z *ZskDict) ZslSetEvictFront(k string, score uint64, l int) {
	z.ZslSet(k, score)
	for z.ZslLen() > l {
		node := z.sl.DeleteFirst()
		if node == nil {
			break
		}
		delete(z.dict, node.value)
	}
}

func (z *ZskDict) ZslSetEvictBack(k string, score uint64, l int) {
	z.ZslSet(k, score)
	for z.ZslLen() > l {
		node := z.sl.DeleteLast()
		if node == nil {
			break
		}
		delete(z.dict, node.value)
	}
}

func (z *ZskDict) ZslDelete(k string) {
	score, ok := z.dict[k]
	if !ok {
		return
	}
	z.sl.ZslDelete(score, k)
	delete(z.dict, k)
}

func (z *ZskDict) ZslRange(start, end int, reverse bool) []string {
	var ks []string
	for _, node := range z.sl.ZslRange(start, end, reverse) {
		ks = append(ks, node.value)
	}
	return ks
}

func (z *ZskDict) ZslLen() int { return z.sl.ZslLen() }

func (z *ZskDict) TestPrint() {
	z.sl.TestPrint()
}

func (z *ZskDict) TestPrint2() {
	z.sl.TestPrint2()
}
