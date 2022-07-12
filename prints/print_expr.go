package prints

import (
	"errors"
	"fmt"
	"io/ioutil"
	"runtime"
	"strings"
)

// ExpPrintln 逐行打印参数原来的表达式和参数本身。
func ExpPrintln(a ...any) {
	// 拿到上一个堆栈信息
	_, filename, line, _ := runtime.Caller(1)
	fb, _ := ioutil.ReadFile(filename)
	fbs := string(fb)

	// 解析文件，拿到参数原来的表达式
	lines := strings.SplitN(fbs, "\n", line)
	target := lines[line-1]
	args := parseFuncArgs(target, "ExpPrintln")
	lena := len(a)
	maxLen := 0
	for i := 0; i < lena; i++ {
		tmp := len(args[i])
		if tmp > maxLen {
			maxLen = tmp
		}
	}
	for i := 0; i < lena; i++ {
		fmt.Printf("%d: %s", i+1, args[i])
		fmt.Printf("\t%s->  ", strings.Repeat("    ", (maxLen-len(args[i]))/4))
		fmt.Println(a[i])
	}
}

func parseFuncArgs(source string, funcName string) (res []string) {
	index := strings.Index(source, funcName) + len(funcName)
	if index < 0 {
		panic(errors.New(fmt.Sprintf("function %s not found", funcName)))
	}
	stack := make([]uint8, 20)
	lens := 0
	tmpIndex := index + 1
	for i := index; ; i++ {
		if source[i] == '(' {
			stack[lens] = '('
			lens += 1
			if lens == 20 {
				panic(errors.New("to many func calls"))
			}
		} else if source[i] == ',' && lens == 1 {
			res = append(res, strings.Trim(strings.TrimSpace(source[tmpIndex:i]), "\n"))
			tmpIndex = i + 1
		} else if source[i] == ')' {
			lens -= 1
		}
		if lens == 0 {
			if tmpIndex > index {
				res = append(res, strings.Trim(strings.TrimSpace(source[tmpIndex:i]), "\n"))
			}
			return
		}
	}
}
