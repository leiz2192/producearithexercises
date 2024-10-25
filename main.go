package main

import (
	"fmt"
	"math/rand"
	"time"
)

func Hidden() int {
	// 使用当前时间作为种子来创建一个新的随机数源
	source := rand.NewSource(time.Now().UnixNano())
	// 基于新的随机数源创建一个新的随机数生成器
	rng := rand.New(source)
	return rng.Intn(4)
}

func Right() bool {
	// 使用当前时间作为种子来创建一个新的随机数源
	source := rand.NewSource(time.Now().UnixNano())
	// 基于新的随机数源创建一个新的随机数生成器
	rng := rand.New(source)
	return rng.Intn(2) == 0
}

func Equation(lhs, rhs, res int, op string) string {
	idx := Hidden()
	right := Right()
	var equation string
	switch idx {
	case 0:
		if right {
			equation = fmt.Sprintf("%2s %s %2d = %2d", " ", op, rhs, res)
		} else {
			equation = fmt.Sprintf("%2d = %2s %s %2d", res, " ", op, rhs)
		}
	case 1:
		if right {
			equation = fmt.Sprintf("%2d %s %2s = %2d", lhs, op, " ", res)
		} else {
			equation = fmt.Sprintf("%2d = %2d %s %2s", res, lhs, op, " ")
		}
	case 2:
		if right {
			equation = fmt.Sprintf("%2d   %2d = %2d", lhs, rhs, res)
		} else {
			equation = fmt.Sprintf("%2d = %2d   %2d", res, lhs, rhs)
		}
	case 3:
		fallthrough
	default:
		equation = fmt.Sprintf("%2d %s %2d = %2s", lhs, op, rhs, " ")
	}
	return equation
}

func main() {
	problems := make([]string, 0, 256)

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
