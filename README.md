

### 介绍
- 使用DFA算法实现关键词/敏感词检测
- 快速，对检测文本最多只完整遍历一次；动态增删查询并发安全
- 支持多种场景。支持中英文，忽略英文大小写，忽略夹杂符号，支持关键词分类



### 使用

```sh
go get github.com/tomatocuke/sieve
```

```go
package main 

import (
	"fmt"
	"github.com/tomatocuke/sieve"
)

// 自定义关键词分类
var categories = [...]string{"", "分类A", "分类B"}

const (
	categoryYellow uint8 = iota + 1
	categoryForce
)

func main() {
	filter := sieve.New()

	in := []string{"&3,)", "苹果", "苹果手机"}
	fmt.Printf("\n\n==== 添加 「&3,)」、「苹果」、「苹果手机」、[apple] =====\n")
	filter.Add(in, categoryYellow)
  // 重复添加apple，类型变了
	filter.Add([]string{"apple"}, categoryForce)

	test1 := map[string]bool{
		"&3,)":          false, // 字符不录入系统，也不匹配
		"翻译:App|le pie": true,  // 包含英文和字符，忽略大小写
		"什么app":         false, // 不完全匹配
		"1起去买个苹果手机":     true,  // 包含中文和字符，可以匹配到最长的苹果手机而不是苹果
	}
	check(filter, test1)

	fmt.Printf("\n\n==== 移除「苹果」=====\n")
	// 移除
	filter.Remove([]string{"", "苹果"})

	test2 := map[string]bool{
		"买个苹果手机": true, 
		"买个苹果派":  false, 
	}
	check(filter, test2)

	text := "苹果和苹果手机我都要！"
	fmt.Printf("\n\n==== 替换「%s」=====\n", text)
	str := filter.Replace(text, '*')
	fmt.Printf("「%s」 => 「%s」\n\n", text, str)
}

func check(filter *sieve.Sieve, test map[string]bool) {
	for k, v := range test {
		r, c := filter.Search(k)
		b := r != ""
		if b != v {
			panic("\n\n校验字符串「" + k + "」不符合预期\n\n")
		}
		if b {
			fmt.Printf("原字符:「%s」, 检测到:「%s」 分类:「%s」\n", k, r, categories[c])
		}
	}
}

```
打印结果
```sh
==== 添加 「&3,)」、「苹果」、「苹果手机」、[apple] =====
原字符:「翻译:App|le pie」, 检测到:「App|le」 分类:「分类B」
原字符:「1起去买个苹果手机」, 检测到:「苹果手机」 分类:「分类A」


==== 移除「苹果」=====
原字符:「买个苹果手机」, 检测到:「苹果手机」 分类:「分类A」


==== 替换「苹果和苹果手机我都要！」=====
「苹果和苹果手机我都要！」 => 「苹果和****我都要！」
```


### 提供方法
- `New()`: 生成检测对象`*Sieve`
- `(*Sieve).Add([]string, uint8)`: 批量添加关键词，选择性打标签
- `(*Sieve).Remove([]string)`: 批量移除关键词
- `(*Sieve).Search()`: 搜索文本中的关键词，返回匹配到的第一个关键词及其标签
- `(*Sieve).Replace()`: 替换文本中的所有关键词