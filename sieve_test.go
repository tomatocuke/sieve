package sieve

import (
	"fmt"
	"testing"
)

// 自定义关键词分类
var categories = [...]string{"", "分类A", "分类B"}

const (
	categoryYellow uint8 = iota + 1
	categoryForce
)

func TestNew(t *testing.T) {
	filter := New()

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
		"买个苹果手机": true,  // 检测到苹果手机
		"买个苹果派":  false, // 检测不到苹果
	}
	check(filter, test2)

	text := "苹果和苹果手机我都要！"
	fmt.Printf("\n\n==== 替换「%s」=====\n", text)
	str := filter.Replace(text, '*')
	fmt.Printf("「%s」 => 「%s」\n\n", text, str)
}

func check(filter *Sieve, test map[string]bool) {
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
