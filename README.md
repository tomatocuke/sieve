

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

var (
	filter = sieve.New()
)

func main() {

	// 空字符串、符号 会被过滤掉。
	filter.Add([]string{"", "&3,)", "apple", "苹果", "苹果手机"}, 1)
	// 可重复添加类型被覆盖
	filter.Add([]string{"苹果手机"}, 2)
	fmt.Printf("\n添加: apple、苹果、苹果手机")

	search([]string{
		"&3,)",          // 不会搜索到
		"翻译:App|le pie", // 忽略大写和夹杂的符号
		"什么app",         // 匹配到一半不算
		"1起去买个苹果手机",     // 包含中文和字符。做了优化，可以匹配到最长的「苹果手机」而不是「苹果」
	})

	// 移除
	fmt.Printf("\n删除: 苹果")
	filter.Remove([]string{"", "苹果"})

	search([]string{
		"买个苹果手机", // 可以
		"买个苹果派",  // 不可以
	})

	text := "苹果和苹果手机我都要！"
	str := filter.Replace(text, '*')
	fmt.Printf("\n替换「%s」 => 「%s」\n\n", text, str)
}

func search(texts []string) {
	for _, t := range texts {
		r, c := filter.Search(t)
		if r == "" {
			fmt.Printf("\n搜索: %s ，未匹配到关键词", t)
		} else {
			fmt.Printf("\n搜索: %s ，关键词: %s ， 类型:%d", t, r, c)
		}
	}
}
```
打印结果
```sh
添加: apple、苹果、苹果手机
搜索: &3,) ，未匹配到关键词
搜索: 翻译:App|le pie ，关键词: App|le ， 类型:1
搜索: 什么app ，未匹配到关键词
搜索: 1起去买个苹果手机 ，关键词: 苹果手机 ， 类型:2
删除: 苹果
搜索: 买个苹果手机 ，关键词: 苹果手机 ， 类型:2
搜索: 买个苹果派 ，未匹配到关键词
替换「苹果和苹果手机我都要！」 => 「苹果和****我都要！」
```


### 提供方法
- `New()`: 生成检测对象`*Sieve`
- `(*Sieve).Add([]string, uint8)`: 批量添加关键词，选择性打标签
- `(*Sieve).Remove([]string)`: 批量移除关键词
- `(*Sieve).Search()`: 搜索文本中的关键词，返回匹配到的第一个关键词及其标签
- `(*Sieve).Replace()`: 替换文本中的所有关键词