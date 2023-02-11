package sieve

// 节点
type node struct {
	// 是否结束
	IsEnd bool
	// 分类 (非必需)
	Category uint8
	// 联想字符
	Children map[rune]*node
}

// 根节点
type root = node

func newNode() *node {
	return &node{}
}

// 添加关键词
func (r *root) AddWord(word string, category uint8) {
	n := r
	for _, w := range word {
		n = n.addChild(w)
	}
	// 非根节点才修改，防止无效关键词修改根节点
	if n != r {
		n.IsEnd = true
		n.Category = category
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
