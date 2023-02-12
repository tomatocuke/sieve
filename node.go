package sieve

const (
	// 通配符
	symbolStar rune = '*'
)

type Tag uint8

// 节点
type node struct {
	// 是否结束
	IsEnd bool
	// 标签
	Tag Tag
	// 替换
	CanReplace bool
	// 通配符长度
	SymbolStarLen uint8
	// 联想字符
	Children map[rune]*node
}

// 根节点
type root = node

func newNode() *node {
	return &node{}
}

// 添加关键词
func (r *root) AddWord(word string, tag Tag, canReplace bool) {
	n := r
	var x uint8
	for i, w := range word {
		if w == symbolStar {
			// 不接受首字符是通配符
			if i == 0 {
				break
			}
			x++
		} else {
			w = trans(w)
			if w > 0 {
				// 解决添加相同关键词，通配符数量叠加问题
				if n.SymbolStarLen < x {
					n.SymbolStarLen = x
				}
				n = n.addChild(w)
				x = 0
			}
		}
	}

	// 非根节点才修改，防止无效关键词修改根节点
	if n != r {
		n.IsEnd = true
		n.Tag = tag
		n.CanReplace = canReplace
		if x > 0 && n.SymbolStarLen < x {
			n.SymbolStarLen = x
		}
	}
}

// 删除关键词
func (r *root) RemoveWord(word string) {
	path := []rune(word)
	ptrs := make([]*node, len(path))
	n := r

	ok := false
	// 正向检验关键词是否存在
	for i, w := range path {
		ptrs[i] = n
		n, ok = n.Children[w]
		if !ok {
			return
		}
	}

	n.IsEnd = false
	for i := len(path) - 1; i >= 0; i-- {
		if i > 0 && !ptrs[i].IsEnd && len(ptrs[i].Children) == 0 {
			delete(ptrs[i-1].Children, path[i-1])
		}
	}

}

// 获取子字符节点
func (n *node) GetChild(w rune) *node {
	child, ok := n.Children[w]
	if ok {
		return child
	}
	return nil
}

// 添加单个字符
func (n *node) addChild(w rune) *node {
	if n.Children == nil {
		n.Children = make(map[rune]*node)
	} else {
		child, ok := n.Children[w]
		if ok {
			return child
		}
	}

	child := newNode()
	n.Children[w] = child
	return child
}

func trans(w rune) rune {
	if w > 255 || w == symbolStar || (w >= 'a' && w <= 'z') || (w >= '0' && w <= '9') {
		return w
	}

	if w >= 'A' && w <= 'Z' {
		return w + 32
	}

	return -1
}
