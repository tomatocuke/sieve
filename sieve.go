package sieve

import (
	"strings"
	"sync"
)

// ==== 检测关键词 =====
type Sieve struct {
	mu sync.RWMutex
	// DFA算法
	trie *node
	// 关键词数量
	len int
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
		if s.trie.AddWord(w, 0, true) {
			s.len++
		}
	}
}

// 添加，打标签并设定是否强制替换
func (s *Sieve) AddWithTag(words []string, tag uint8, canReplace bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		if s.trie.AddWord(w, tag, canReplace) {
			s.len++
		}
	}
}

// 移除关键词
func (s *Sieve) Remove(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		if s.trie.RemoveWord(w) {
			s.len--
		}
	}
}

func (s *Sieve) Len() int {
	return s.len
}

// 搜索关键词，返回第一个匹配到的关键词和其类型
func (s *Sieve) Search(text string) (string, uint8) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ws := []rune(text)
	start, end, tag, _ := s.trie.Search(ws)
	return string(ws[start:end]), tag
}

// 替换匹配到的关键词
func (s *Sieve) Replace(text string) string {
	result, _ := s.ReplaceAndCheckTags(text, nil)
	return result
}

// 替换文本的关键词，检查是否含有特定标签
func (s *Sieve) ReplaceAndCheckTags(text string, tags []uint8) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var (
		start, end, offset, counter int

		ws         = []rune(text)
		canReplace bool
		hasTag     bool
		tag        uint8
	)

	for counter < 5 {
		counter++

		offset = end
		start, end, tag, canReplace = s.trie.Search(ws[offset:])
		if end == 0 {
			break
		}

		start += offset
		end += offset

		if canReplace {
			// fmt.Println("替换:", string(ws), "=>", string(ws[start:end]))
			for i := start; i < end; i++ {
				ws[i] = symbolStar
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
		return strings.Repeat(string(symbolStar), len(ws)), hasTag
	}

	return string(ws), hasTag
}
