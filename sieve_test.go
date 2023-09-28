package sieve

import (
	"fmt"
	"strings"
	"testing"
)

var (
	filter = New()

	myDemoKeywords = "./keyword"
)

func TestMain(t *testing.T) {
	// ===== 基础用法：添加、移除、搜索、替换 =====
	demo1()

	// ===== 进阶用法：通配符、忽略大小写、符号干扰无效 =====
	demo2()

	// ===== 特殊功能：打标签 (用于区分敏感词类型) =====
	demo3()

	// ===== 特殊功能：不替换 (仅发现) =====
	demo4()
}

func demo1() {
	// 添加
	filter.Add([]string{"苹果", "西红柿", "葡萄"})
	// 移除
	filter.Remove([]string{"葡萄"})
	const text = "我想吃葡萄和西红柿，苹果也不错"
	// 搜索 (第一个关键词)
	searchKeyword, _ := filter.Search(text)
	// 替换
	replaceText, _ := filter.Replace(text)

	fmt.Println("\n===== 基础用法：添加、移除、搜索、替换")
	fmt.Println("添加:", "苹果", "西红柿", "葡萄")
	fmt.Println("移除:", "葡萄")
	fmt.Println("测试:", text)
	fmt.Println("搜索关键词:", searchKeyword)
	fmt.Println("替换后:", replaceText)
}

func demo2() {
	const text = "FUCK!我操你x、操你🐎、操你&x"
	filter.Add([]string{"fuck", "操你*"})
	replaceText, _ := filter.Replace(text)

	fmt.Println("\n===== 进阶用法：通配符、忽略大小写、符号干扰无效")
	fmt.Println("添加:", "fuck", "操你*")
	fmt.Println("测试:", text)
	fmt.Println("替换后:", replaceText)
}

// 设置分类标签
const (
	TagDefault = iota
	TagInsult
	TagTrade
)

func demo3() {
	const text = "你是傻b么？这么傻呢！"
	fails, err := filter.AddByFile(myDemoKeywords, TagInsult, true)
	if err != nil {
		panic(err)
	}
	replaceText, keywords := filter.Replace(text)

	fmt.Println("\n===== 特殊功能：打标签")
	fmt.Printf("添加词典: %s 设置标签: %d 自动替换\n", myDemoKeywords, TagInsult)
	if len(fails) > 0 {
		fmt.Println("添加失败:", fails)
	}
	fmt.Println("测试:", text)
	fmt.Println("替换后:", replaceText)
	fmt.Println("包含敏感词: ", keywords)
}

func demo4() {
	const text = "二手房怎么样"
	filter.AddByFile(myDemoKeywords, TagTrade, false)
	replaceText, keywords := filter.Replace(text)

	fmt.Println("\n===== 特殊功能：不替换 (仅发现，另行处理) =====")
	fmt.Printf("添加词典: %s 设置标签: %d 不替换\n", myDemoKeywords, TagTrade)
	fmt.Println("测试:", text)
	fmt.Println("替换后:", replaceText)
	fmt.Println("包含敏感词: ", keywords)
}

var longText = strings.Repeat("哦😯哈HA", 20) // 100字符

func BenchmarkReplace(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filter.Replace(longText)
	}
}

func BenchmarkSearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		filter.Search(longText)
	}
}
