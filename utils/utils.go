package utils

import (
	"fmt"
	"os"
	"strings"
)

func ExitProgram()  {
	fmt.Println("回车结束本程序...")
	var inputText string
	fmt.Scanln(&inputText)
	os.Exit(0)
}

func ReplaceRight(s, old, new string, n int) string {
	if old == new || n == 0 {
		return s // avoid allocation
	}
	// Compute number of replacements.
	m := strings.Count(s, old)
	if m == 0 {
		return s // avoid allocation
	} else if n < 0 || m < n {
		n = m
	}
	// Apply replacements to buffer.
	var b strings.Builder
	b.Grow(len(s) + n*(len(new)-len(old)))
	var preIndex int
	for i := 0; i < m; i++ {
		index := strings.Index(s[preIndex:], old)
		if index < 0 {
			b.WriteString(s[preIndex:])
			break
		}
		if i < m-n {
			b.WriteString(s[preIndex : index+preIndex+len(old)])
			preIndex += index + len(old)
			continue
		}
		b.WriteString(s[preIndex : index+preIndex])
		b.WriteString(new)
		preIndex += index + len(old)
	}
	b.WriteString(s[preIndex:])
	return b.String()
}