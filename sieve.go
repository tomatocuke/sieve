package sieve

import (
	"bufio"
	"io"
	"os"
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
		trie: &node{},
	}
	return s
}

// 简单添加关键词
func (s *Sieve) Add(words []string) int {
	return s.AddWithTag(words, 0, true)
}

// 从文本添加关键词，打标签并设定是否强制替换
func (s *Sieve) AddByFile(filename string, tag uint8, canReplace bool) (int, error) {
	const delim = '\n'
	words := make([]string, 0, 2048)

	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	br := bufio.NewReader(f)
	for {
		b, err := br.ReadBytes(delim)
		words = append(words, string(b))
		if err == io.EOF {
			break
		}
	}

	i := s.AddWithTag(words, tag, canReplace)
	return i, nil
}

// 添加关键词，打标签并设定是否强制替换
func (s *Sieve) AddWithTag(words []string, tag uint8, canReplace bool) (i int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		if s.trie.AddWord(w, tag, canReplace) {
			i++
		}
	}

	s.len += i

	return
}

// 移除关键词
func (s *Sieve) Remove(words []string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	var i int
	for _, w := range words {
		if s.trie.RemoveWord(w) {
			i++
		}
	}
	s.len -= i
	return i
}

func (s *Sieve) Len() int {
	return s.len
}

// 是否含有关键词
func (s *Sieve) Has(text string) bool {
	word, _ := s.Search(text)
	return word != ""
}

// 替换所有关键词
func (s *Sieve) Replace(text string) string {
	result, _ := s.ReplaceAndCheckTags(text, nil)
	return result
}

// 返回文本中第一个关键词及其类型
func (s *Sieve) Search(text string) (string, uint8) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ws := []rune(text)
	start, end, tag, _ := s.trie.Search(ws)
	return string(ws[start:end]), tag
}

// 替换文本的关键词，检查是否含有特定标签
func (s *Sieve) ReplaceAndCheckTags(text string, tags []uint8) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var (
		start, end, offset int

		ws         = []rune(text)
		canReplace bool
		hasTag     bool
		tag        uint8
	)

	for {

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

	return string(ws), hasTag
}
