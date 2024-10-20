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
	return rng.Intn(3)
}

func main() {
	problems := make([]string, 0, 256)

	subCnt := 0
	for i := 10; i > 0; i-- {
		for j := i; j > 0; j-- {
			idx := Hidden()
			switch idx {
			case 0:
				problems = append(problems, fmt.Sprintf("%2s - %2d = %2d", " ", j, i-j))
			case 1:
				problems = append(problems, fmt.Sprintf("%2d - %2s = %2d", i, " ", i-j))
			case 2:
				fallthrough
			default:
				problems = append(problems, fmt.Sprintf("%2d - %2d = %2s", i, j, " "))
			}
			subCnt++
		}
	}

	addCnt := 0
	for i := 1; i <= 10; i++ {
		for j := 1; j <= 10; j++ {
			if i+j > 10 {
				break
			}
			idx := Hidden()
			switch idx {
			case 0:
				problems = append(problems, fmt.Sprintf("%2s + %2d = %2d", " ", j, i+j))
			case 1:
				problems = append(problems, fmt.Sprintf("%2d + %2s = %2d", i, " ", i+j))
			case 2:
				fallthrough
			default:
				problems = append(problems, fmt.Sprintf("%2d + %2d = %2s", i, j, " "))
			}
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
