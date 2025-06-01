package main

import (
	"fmt"
	"math"
	"math/rand"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type ArithMode int

const (
	Empty ArithMode = iota
	AddOrSubWithin10
	FillForOneEquationWithin10
	MultiAddOrSubWithin10
	MixAddAndSubWithin10
	FillForTwoEquationWithin10
	FillForMixAddAndSubWithin10
	AddOrSubWithin20
	FillForOneEquationWithin20
	Unknown
)

var (
	Options = []string{
		"",
		"10以内加减",
		"10以内算式填空",
		"10以内连加或连减",
		"10以内加减混合",
		"10以内两边算式填空",
		"10以内两边加减算式填空",
		"20以内加减",
		"20以内算式填空",
	}

	// 使用当前时间作为种子来创建一个新的随机数源,基于新的随机数源创建一个新的随机数生成器
	rng = rand.New(rand.NewSource(time.Now().UnixNano()))
)

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
	idx := rng.Intn(4)
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
	if idx < 3 && rng.Intn(2) != 0 {
		return Reverse(equation)
	}
	return equation
}

func Format(equations []string, colNum int) string {
	rand.Shuffle(len(equations), func(i, j int) { equations[i], equations[j] = equations[j], equations[i] })

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

func AddOrSubWithin10Exercises(hidden bool) string {
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

func MultiAddOrSubWithin10Exercies() string {
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

func MixAddAndSubWithin10Exercies() string {
	equations := make([]string, 0, 128)

	preSubCnt := 0
	for i := 10; i > 0; i-- {
		for j := i - 1; j > 0; j-- {
			for k := 1; k <= 10; k++ {
				if i-j+k > 10 || k == j {
					continue
				}
				equations = append(equations, fmt.Sprintf("%2d - %2d + %2d = %2s", i, j, k, " "))
				preSubCnt++
			}
		}
	}

	preAddCnt := 0
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			if i+j > 10 {
				continue
			}
			for k := 1; k <= 10; k++ {
				if i+j-k < 0 {
					continue
				}
				if k == i || k == j {
					continue
				}
				equations = append(equations, fmt.Sprintf("%2d + %2d - %2d = %2s", i, j, k, " "))
				preAddCnt++
			}
		}
	}
	return fmt.Sprintf("%s\npreSubCnt: %d, preAddCnt: %d\n", Format(equations, 5), preSubCnt, preAddCnt)
}

func ReplaceCharAt(s string, start int, replacement string) string {
	var sb strings.Builder
	sb.WriteString(s[0:start])
	sb.WriteString(replacement)
	sb.WriteString(s[start+len(replacement):])
	return sb.String()
}

func FillForTwoEquationWithin10Exercies() string {
	expressions := make(map[int][]string)
	for i := 10; i > 0; i-- {
		for j := i - 1; j > 0; j-- {
			expressions[i-j] = append(expressions[i-j], fmt.Sprintf("%2d - %2d", i, j))
		}
	}

	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			if i+j > 10 {
				continue
			}
			expressions[i+j] = append(expressions[i+j], fmt.Sprintf("%2d + %2d", i, j))
		}
	}

	totalCnt := 0
	equations := make([]string, 0, 128)
	for _, es := range expressions {
		rand.Shuffle(len(es), func(i, j int) { es[i], es[j] = es[j], es[i] })
		for i := 0; i < len(es)-1; i++ {
			for j := i + 1; j < len(es); j++ {
				fill := rng.Intn(4)
				// AA + BB = CC + DD, replace AA/BB/CC/DD to "  "
				equations = append(equations, ReplaceCharAt(fmt.Sprintf("%s = %s", es[i], es[j]), fill*5, "  "))
				totalCnt++
			}
		}
	}
	return fmt.Sprintf("%s\ntotalCnt: %d\n", Format(equations, 5), totalCnt)
}

func FilleForMixAddAndSubExercies() string {
	expressions := make(map[int][]string)
	for i := 10; i > 1; i-- {
		for j := i - 1; j > 1; j-- {
			for k := 2; k <= 10; k++ {
				if i-j+k > 10 || k == j {
					continue
				}
				expressions[i-j+k] = append(expressions[i-j+k], fmt.Sprintf("%2d - %2d + %2d", i, j, k))
			}
		}
	}
	for i := 2; i <= 10; i++ {
		for j := 2; j <= 10; j++ {
			if i+j > 10 {
				continue
			}
			for k := 2; k <= 10; k++ {
				if i+j-k < 0 {
					continue
				}
				if k == i || k == j {
					continue
				}
				expressions[i+j-k] = append(expressions[i+j-k], fmt.Sprintf("%2d + %2d - %2d", i, j, k))
			}
		}
	}
	totalCnt := 0
	equations := make([]string, 0, 128)
	for s, es := range expressions {
		if s == 0 || s == 1 {
			continue
		}
		rand.Shuffle(len(es), func(i, j int) { es[i], es[j] = es[j], es[i] })
		// for i := range len(es) - 1 {
		// 	for j := i + 1; j < len(es); j++ {
		// 		if es[i][2:7] == es[j][2:7] || es[i][7:12] == es[j][7:12] {
		// 			continue
		// 		}
		// 		fill := rng.Intn(6)
		// 		// AA + BB + CC = DD + EE + FF, replace AA/BB/CC/DD/EE/FF to "  "
		// 		equations = append(equations, ReplaceCharAt(fmt.Sprintf("%s = %s", es[i], es[j]), fill*5, "  "))
		// 		totalCnt++
		// 	}
		// }
		for i := range len(es) {
			j := i + 1
			if j >= len(es) {
				j = 0
			}
			// if es[i][2:7] == es[j][2:7] || es[i][7:12] == es[j][7:12] {
			// 	continue
			// }
			fill := rng.Intn(6)
			// AA + BB + CC = DD + EE + FF, replace AA/BB/CC/DD/EE/FF to "  "
			equations = append(equations, ReplaceCharAt(fmt.Sprintf("%s = %s", es[i], es[j]), fill*5, "  "))
			totalCnt++
		}
	}
	return fmt.Sprintf("%s\ntotalCnt: %d\n", Format(equations, 3), totalCnt)
}

func AddOrSubWithin20Exercies() string {
	equations := make([]string, 0, 100)

	subCnt := 0
	for i := 11; i <= 20; i++ {
		for j := 1; j < 20; j++ {
			if i <= j || j == 10 {
				continue
			}
			equations = append(equations, fmt.Sprintf("%2d - %2d = %2s", i, j, " "))
			subCnt++
		}
	}

	addCnt := 0
	for i := 1; i < 20; i++ {
		for j := 1; j < 20; j++ {
			if i+j <= 10 || i+j > 20 {
				continue
			}
			equations = append(equations, fmt.Sprintf("%2d + %2d = %2s", i, j, " "))
			addCnt++
		}
	}

	return fmt.Sprintf("%s\nsubCnt: %d, addCnt: %d\n", Format(equations, 7), subCnt, addCnt)
}

func FillForOneEquationWithin20Exercies() string {
	expressions := make(map[int][]string)

	for i := 11; i <= 20; i++ {
		for j := 1; j < 20; j++ {
			if i <= j || j == 10 {
				continue
			}
			expressions[i-j] = append(expressions[i-j], fmt.Sprintf("%2d - %2d", i, j))
		}
	}

	for i := 1; i < 20; i++ {
		for j := 1; j < 20; j++ {
			if i+j <= 10 || i+j > 20 {
				continue
			}
			expressions[i+j] = append(expressions[i+j], fmt.Sprintf("%2d + %2d", i, j))
		}
	}

	totalCnt := 0
	equations := make([]string, 0, 128)
	for s, es := range expressions {
		if s < 2 {
			continue
		}
		rand.Shuffle(len(es), func(i, j int) { es[i], es[j] = es[j], es[i] })
		for i := range len(es) - 1 {
			for j := i + 1; j < len(es); j++ {
				if es[i][3] == es[j][3] && es[i][0:1] == es[j][5:6] {
					continue
				}
				fill := rng.Intn(3)
				// AA + BB = CC + DD, replace AA/BB/CC/DD to "  "
				equations = append(equations, ReplaceCharAt(fmt.Sprintf("%s = %s", es[i], es[j]), fill*5, "  "))
				totalCnt++
			}
		}
	}
	return fmt.Sprintf("%s\ntotalCnt: %d\n", Format(equations, 5), totalCnt)
}

func Produce(mode ArithMode) string {
	switch mode {
	case Empty:
		return ""
	case AddOrSubWithin10:
		return AddOrSubWithin10Exercises(false)
	case FillForOneEquationWithin10:
		return AddOrSubWithin10Exercises(true)
	case MultiAddOrSubWithin10:
		return MultiAddOrSubWithin10Exercies()
	case MixAddAndSubWithin10:
		return MixAddAndSubWithin10Exercies()
	case FillForTwoEquationWithin10:
		return FillForTwoEquationWithin10Exercies()
	case FillForMixAddAndSubWithin10:
		return FilleForMixAddAndSubExercies()
	case AddOrSubWithin20:
		return AddOrSubWithin20Exercies()
	case FillForOneEquationWithin20:
		return FillForOneEquationWithin20Exercies()
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
	r.StaticFile("/favicon.ico", "./resources/favicon.ico")
	// 定义路由和处理函数
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", data)
	})

	r.POST("/submit", func(c *gin.Context) {
		// 根据选择更新内容
		selected := c.PostForm("options")
		for i, option := range data.Options {
			if selected == option {
				data.Content = Produce(ArithMode(i))
			}
		}
		data.Selected = selected
		c.HTML(http.StatusOK, "index.tmpl", data)
	})
	// r.Run(":9090")

	// 创建一个TCP监听器，使用 ":0" 表示随机端口
	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		panic(err)
	}
	defer ln.Close()

	// 获取随机分配的端口号
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening and serving HTTP on", port)

	if err := r.RunListener(ln); err != nil {
		panic(err)
	}
}
