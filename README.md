## 完全自己写的，喜欢的同学给个🌟，我需要你

### 介绍
- 使用DFA算法实现关键词/敏感词检测。有问题可以联系我 QQ:`772532526`
- 优点：
	- 快速。复杂度O(1)
	- 忽略英文大小写。
	- 忽略常见中英文符号。例如设置关键字「苹果」，文本中包含「苹& 果」也可以检测到，无视符号。（关键词第一位为符号则关键词无效）
	- 通配符*。一个* 匹配且必定匹配任何一个非符号字符。例如「苹果*」，可以匹配到「苹果醋」、「苹果&奶」，但是不能匹配到「苹果」。（关键词第一位为通配符则关键词无效）
	- 打标签。对关键词分类，更好地定制化处理。例如人工分类多个种类的关键词：政治、色情、辱骂、广告 ...
	- 选择性替换。检测但不处理。例如你希望对疑似营销、广告的关键词不进行替换，仅灰字提示消息接受者注意。

- 缺点：
	- 英文字符的效率比中文低，在关键词和文本都含有英文的情况下。

- 举例关键词
	- 色情。中文前后语意复杂，「色情」如果设置为关键词，那会影响「白色情人节」，此类属于特例，不建议设置关键词。(其实也可以通过标签的方式再增加白名单二次过滤，但是那样无疑增加了复杂度)
	- 傻逼、傻B、煞笔、傻X。 这种不应该使用通配符，太简洁，只能设置多个。不要设置「sha B」。
	- 操你x、操你🐎。 这种变化多样，但是可以设置「操你*」。（但是「我操！你干啥」会被误杀）
	- 迷魂药、迷情药。这种可以设置「迷*药」
	- VX: 2341443， wx:43545 。此类我也没什么好办法，它数字还可以改为全角，圆圈字符。（➕qq、➕vx可以，但是不要直接添加vx这种常见的字母开头）


### 函数说明
- [x] `New() *Sieve` 创建新的实例
- [x] `Add(keywords []string) int` 添加关键词，返回添加成功数。
- [x] `AddByFile(filename string, tag uint8, canReplace bool) (int, error)` 从文本添加关键词，附带标签和是否需要被替换，返回添加成功数和错误 。
- [x] `AddWithTag(keywords []string, tag uint8, canReplace bool) int` 添加关键词，附带标签和是否需要被替换，返回添加成功数
- [x] `Remove(keywords []string) int` 移除关键词
- [x] `Len() int` 实例中关键词数量
- [x] `Has(text string) bool` 是否包含关键词
- [x] `Replace(text string) string` 替换文本中的所有关键词 (设置canReplace为false的除外)
- [x] `Search(text string) (string, tag)` 搜索文本中出现的第一个关键词和其标签


### 使用

```sh
go get -u github.com/tomatocuke/sieve@latest
```

```go
package main

import (
	"fmt"

	"github.com/tomatocuke/sieve"
)

var (
	filter = sieve.New()
	text   = "A：‘我想要个苹果手机，iphone-13或iphone14都行。’ B：‘行，买斤苹果吧。’"
)

func main() {

	// =================== 基础用法 ===================

	// 测试：添加、删除、长度
	filter1 := New()
	filter1.Add([]string{"foo", "bar"})
	fmt.Println(filter1.Len() == 2)
	filter1.Remove([]string{"bar"})
	fmt.Println(filter1.Len() == 1)

	// 测试：中文、英文、忽略大小写、忽略符号、优先检测长关键词
	filter2 := New()
	filter2.Add([]string{"苹果", "苹果手机", "iphone 13"}) // 只13不行噢？
	has := filter2.Has(text)
	fmt.Println("是否含有关键词:", has)
	s1 := filter2.Replace(text)
	fmt.Println("普通替换：", s1)

	// 测试：通配符。哪款iphone都不行！
	filter3 := New()
	filter3.Add([]string{"苹果", "苹果手机", "iphone**"}) // 管你iphone多少
	s2 := filter3.Replace(text)
	fmt.Println("通配符替换：", s2)

	// =================== 高阶用法 ===================

	// 测试：打标签、不替换
	const (
		TAG_FRUIT = iota + 1
		TAG_PHONE
	)
	var tags = [...]string{"", "水果", "手机"}
	filter4 := New()
	filter4.AddWithTag([]string{"苹果"}, TAG_FRUIT, false)              // 水果OK，不替换
	filter4.AddWithTag([]string{"苹果手机", "iphone**"}, TAG_PHONE, true) // 手机不行，替换

	s3 := filter4.Replace(text)
	fmt.Println("标签不替换:", s3)
	s4, tag := filter4.Search(text)
	fmt.Println("搜索第一个关键词:", s4, " 标签:", tags[tag])

}


```
打印结果
```sh
# A：‘我想要个苹果手机，iphone-13或iphone14都行。’ B：‘行，买斤苹果吧。’
是否含有关键词: true
普通替换： A：‘我想要个****，*********或iphone14都行。’ B：‘行，买斤**吧。’
通配符替换： A：‘我想要个****，*********或********都行。’ B：‘行，买斤**吧。’
标签不替换: A：‘我想要个****，*********或********都行。’ B：‘行，买斤苹果吧。’
搜索第一个关键词: 苹果手机  标签: 手机
```
