package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

var (
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

func Equation(lhs, rhs, res int, op string) string {
	var equation string
	idx := Hidden()
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

func main() {
	problems := make([]string, 0, 100)

	subCnt := 0
	for i := 10; i > 0; i-- {
		for j := i; j > 0; j-- {
			problems = append(problems, Equation(i, j, i-j, "-"))
			subCnt++
		}
	}

	addCnt := 0
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			if i+j > 10 {
				break
			}
			problems = append(problems, Equation(i, j, i+j, "+"))
			addCnt++
		}
	}

	rand.Shuffle(len(problems), func(i, j int) {
		problems[i], problems[j] = problems[j], problems[i]
	})

	cnt := 0
	for _, p := range problems {
		fmt.Print(p)
		cnt = cnt + 1
		if cnt%6 == 0 {
			fmt.Print("\n")
		} else {
			fmt.Print("  |  ")
		}
	}
	fmt.Printf("\nsubCnt: %d, addCnt: %d\n", subCnt, addCnt)
}
