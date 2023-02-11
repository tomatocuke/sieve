package sieve

import (
	"strings"
	"sync"
)

// ==== 检测关键词 =====

type Sieve struct {
	mu sync.RWMutex
	// DFA算法
	trie *root
}

func New() *Sieve {
	s := &Sieve{
		trie: newNode(),
	}
	return s
}

// 批量添加关键词，选择性打标签
func (s *Sieve) Add(words []string, category uint8) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		s.trie.AddWord(w, category)
	}
}

// 移除关键词
func (s *Sieve) Remove(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		s.trie.RemoveWord(w)
	}
}

// 搜索关键词，返回第一个匹配到的关键词和其类型
func (s *Sieve) Search(text string) (string, uint8) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ws := []rune(text)
	start, end, category := s.Index(ws)
	return string(ws[start:end]), category
}

// 替换全部匹配到的关键词
func (s *Sieve) Replace(text string, symbol rune) string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ws := []rune(text)

	var start, end, offset, counter int
	for counter < 5 {
		counter++

		offset = end
		start, end, _ = s.Index(ws[offset:])
		if end == 0 {
			break
		}

		start += offset
		end += offset
		for i := start; i < end; i++ {
			ws[i] = symbol
		}
	}

	if counter >= 5 {
		return strings.Repeat(string(symbol), len(ws))
	}

	return string(ws)
}

func (s *Sieve) Index(ws []rune) (start int, end int, category uint8) {
	node := s.trie

	var k rune
	for i, w := range ws {
		k = filter(w)
		if k == 0 {
			continue
		}

		// 查询是否存在该字符
		node = node.GetChild(k)

		if node != nil {
			if start == 0 {
				start = i
			}
			if node.IsEnd {
				end = i
				category = node.Category
			}
		} else {
			if end == 0 {
				start = 0
				node = s.trie
			} else {
				break
			}
		}
	}

	// 结尾匹配一半
	if end == 0 {
		start = 0
	} else {
		end += 1 // 为了切片右边闭合
	}

	return
}

func filter(w rune) rune {
	// 非ascii或小写字母
	if w > 255 || (w >= 'a' && w <= 'z') {
		return w
	}
	// 大写转小写
	if w <= 'Z' && w >= 'A' {
		return w + 32
	}

	return 0
}
