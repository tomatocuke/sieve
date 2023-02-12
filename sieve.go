package sieve

import (
	"strings"
	"sync"
)

// ==== 检测关键词 =====

const (
	replaceSymbol = '*'
)

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
func (s *Sieve) Add(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		s.trie.AddWord(w, 0, true)
	}
}

// 添加，打标签并设定是否强制替换
func (s *Sieve) AddWithTag(words []string, tag Tag, canReplace bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		s.trie.AddWord(w, tag, canReplace)
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
func (s *Sieve) Search(text string) (string, Tag) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ws := []rune(text)
	start, end, tag, _ := s.index(ws)
	return string(ws[start:end]), tag
}

// 替换匹配到的关键词
func (s *Sieve) Replace(text string) string {
	result, _ := s.ReplaceAndCheckTags(text, nil)
	return result
}

// 替换文本的关键词，检查是否含有特定标签
func (s *Sieve) ReplaceAndCheckTags(text string, tags []Tag) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var (
		start, end, offset, counter int

		ws         = []rune(text)
		canReplace bool
		hasTag     bool
		tag        Tag
	)

	for counter < 5 {
		counter++

		offset = end
		start, end, tag, canReplace = s.index(ws[offset:])
		if end == 0 {
			break
		}

		start += offset
		end += offset

		if canReplace {
			for i := start; i < end; i++ {
				ws[i] = replaceSymbol
			}
		}

		if !hasTag && len(tags) > 0 {
			for _, t := range tags {
				if t == tag {
					hasTag = true
				}
			}
		}
	}

	// 太多了直接全屏蔽
	if counter >= 5 {
		return strings.Repeat(string(replaceSymbol), len(ws)), hasTag
	}

	return string(ws), hasTag
}

func (s *Sieve) index(ws []rune) (start int, end int, tag Tag, canReplace bool) {

	node := s.trie
	jumping := false
	start = -1
	end = -1

	length := len(ws)
	for i := 0; i < length; i++ {
		w := trans(ws[i])
		if w <= 0 {
			continue
		}

		// 查询是否存在该字符
		node = node.GetChild(w)
		// 举例 「苹果」和「苹果**本」是关键词
		if node == nil {
			// 苹果笔记
			if end > -1 {
				break
			}
			// 苹方
			if start > -1 {
				start = -1
				jumping = false
			}
			node = s.trie
		} else {
			// 苹
			if start == -1 {
				start = i
			}
			// 苹果
			if node.IsEnd {
				end = i
				tag = node.Tag
				canReplace = node.CanReplace
			}
			// 当前字符「果」，向后偏移2位
			if node.SymbolStarLen > 0 && !jumping {
				jumping = true
				i += int(node.SymbolStarLen)
				end += int(node.SymbolStarLen)
				if end >= length {
					end = length - 1
				}
			}
		}
	}

	// 匹配失败，防止匹配一半的情况。
	// 匹配成功，适配数组左开右闭把end+1
	if end == -1 {
		end = 0
	} else {
		end += 1
	}

	return
}
