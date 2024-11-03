package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ArithMode int

const (
	Empty ArithMode = iota
	TwoNum
	FillWithTwoNum
	ThreeNum
	Unknown
)

var (
	Options = []string{"", "10以内的加减", "10以内的加减填空", "10以内的连加或连减"}

	// 使用当前时间作为种子来创建一个新的随机数源,基于新的随机数源创建一个新的随机数生成器
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

func Hidden() int {
	return rng.Intn(4)
}

func Right() bool {
	return rng.Intn(2) == 0
}

// Reverse s like "AA - BB = CC" to "CC = AA - BB"
func Reverse(s string) string {
	var sb strings.Builder
	sb.WriteString(s[10:])
	sb.WriteString(s[7:10])
	sb.WriteString(s[:7])
	return sb.String()
}

func Equation(lhs, rhs, res int, op string, hidden bool) string {
	var equation string
	idx := Hidden()
	if !hidden {
		idx = math.MaxInt
	}
	switch idx {
	case 0:
		equation = fmt.Sprintf("%2s %s %2d = %2d", " ", op, rhs, res)
	case 1:
		equation = fmt.Sprintf("%2d %s %2s = %2d", lhs, op, " ", res)
	case 2:
		equation = fmt.Sprintf("%2d   %2d = %2d", lhs, rhs, res)
	case 3:
		fallthrough
	default:
		equation = fmt.Sprintf("%2d %s %2d = %2s", lhs, op, rhs, " ")
	}
	if idx < 3 && !Right() {
		return Reverse(equation)
	}
	return equation
}

func Format(equations []string, colNum int) string {
	rand.Shuffle(len(equations), func(i, j int) {
		equations[i], equations[j] = equations[j], equations[i]
	})

	var ret strings.Builder
	cnt := 0
	for _, p := range equations {
		ret.WriteString(p)
		cnt = cnt + 1
		if cnt%colNum == 0 {
			ret.WriteString("\n")
		} else {
			ret.WriteString(" | ")
		}
	}
	return ret.String()
}

func TwoNumExercises(hidden bool) string {
	equations := make([]string, 0, 100)

	subCnt := 0
	for i := 10; i > 0; i-- {
		for j := i; j > 0; j-- {
			equations = append(equations, Equation(i, j, i-j, "-", hidden))
			subCnt++
		}
	}

	addCnt := 0
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			if i+j > 10 {
				continue
			}
			equations = append(equations, Equation(i, j, i+j, "+", hidden))
			addCnt++
		}
	}

	return fmt.Sprintf("%s\nsubCnt: %d, addCnt: %d\n", Format(equations, 6), subCnt, addCnt)
}

func TriEquation(lhs, mhs, rhs int, op string) string {
	return fmt.Sprintf("%2d %s %2d %s %2d = %2s", lhs, op, mhs, op, rhs, " ")
}

func ThreeNumExercies() string {
	equations := make([]string, 0, 215)

	subCnt := 0
	for i := 10; i > 0; i-- {
		for j := i; j > 0; j-- {
			for k := j; k > 0; k-- {
				if i-j-k < 0 {
					continue
				}
				equations = append(equations, TriEquation(i, j, k, "-"))
				subCnt++
			}
		}
	}

	addCnt := 0
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			for k := 1; k <= 10; k++ {
				if i+j+k > 10 {
					continue
				}
				equations = append(equations, TriEquation(i, j, k, "+"))
				addCnt++
			}
		}
	}

	return fmt.Sprintf("%s\nsubCnt: %d, addCnt: %d\n", Format(equations, 5), subCnt, addCnt)
}

func Produce(mode ArithMode) string {
	switch mode {
	case Empty:
		return ""
	case TwoNum:
		return TwoNumExercises(false)
	case FillWithTwoNum:
		return TwoNumExercises(true)
	case ThreeNum:
		return ThreeNumExercies()
	default:
		return "unsupport this option"
	}
}

// 定义模板数据结构
type TemplateData struct {
	Title    string
	Options  []string
	Selected string
	Content  string
}

func main() {
	// 初始选择和内容
	data := TemplateData{
		Title:   "Dynamic Content Based on Selection",
		Options: Options,
	}

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")
	// 定义路由和处理函数
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	r.POST("/submit", func(c *gin.Context) {
		// 根据选择更新内容
		selected := c.PostForm("options")
		for i, option := range data.Options {
			if selected == option {
				data.Content = "Content for " + selected + "\n" + Produce(ArithMode(i))
			}
		}
		data.Selected = selected
		c.HTML(http.StatusOK, "index.tmpl", data)
	})
	r.Run(":9090")
}
