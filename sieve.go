package sieve

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
)

// ==== 检测关键词 =====
type Sieve struct {
	mu   sync.RWMutex
	trie *node
}

func New() *Sieve {
	s := &Sieve{
		trie: &node{},
	}
	return s
}

// 简单添加关键词
func (s *Sieve) Add(words []string) (fail []string) {
	return s.add(words, 0, true)
}

// 从文本添加关键词，打标签并设定是否自动替换为*
func (s *Sieve) AddByFile(filename string, tag uint8, autoReplace bool) (fails []string, err error) {
	const delim = '\n'
	words := make([]string, 0, 2048)

	var reader io.Reader
	// 远程文件
	if strings.HasPrefix(filename, "http") {
		resp, err := http.Get(filename)
		if err != nil {
			return nil, err
		}
		reader = resp.Body
		defer resp.Body.Close()
	} else {
		f, err := os.Open(filename)
		if err != nil {
			return nil, err
		}
		reader = f
		defer f.Close()
	}

	br := bufio.NewReader(reader)
	for {
		b, err := br.ReadBytes(delim)
		words = append(words, string(b))
		if err == io.EOF {
			break
		}
	}

	fails = s.add(words, tag, autoReplace)

	return
}

// 移除关键词
func (s *Sieve) Remove(words []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		s.trie.RemoveWord(w)
	}
}

// 返回文本中第一个关键词及其标签
func (s *Sieve) Search(text string) (string, uint8) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	ws := []rune(text)
	start, end, tag, _ := s.trie.Search(ws)
	return string(ws[start:end]), tag
}

// 替换文本的关键词
func (s *Sieve) Replace(text string) (string, map[uint8][]string) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var (
		start, end, offset int

		ws          = []rune(text)
		keywords    = make(map[uint8][]string)
		tag         uint8
		autoReplace bool
	)

	for {
		offset = end
		start, end, tag, autoReplace = s.trie.Search(ws[offset:])
		if end == 0 {
			break
		}

		start += offset
		end += offset

		keywords[tag] = append(keywords[tag], string(ws[start:end]))

		if autoReplace {
			// fmt.Println("替换:", string(ws), "=>", string(ws[start:end]))
			for i := start; i < end; i++ {
				ws[i] = symbolStar
			}
		}

	}

	return string(ws), keywords
}

// 添加关键词，打标签并设定是否强制替换
func (s *Sieve) add(words []string, tag uint8, autoReplace bool) (fail []string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, w := range words {
		if !s.trie.AddWord(w, tag, autoReplace) {
			fail = append(fail, w)
		}
	}

	return
}
