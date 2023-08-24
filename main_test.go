package main

import (
	"fmt"
	"testing"
	"time"
)

func TestM(t *testing.T) {
	start := time.Now()
	fmt.Println("==============================================================")
	fmt.Printf("| 【检测当前文件夹文件是否发生变更，开始时间：%s 】 |\n", start.Format("2006-01-02 15:04:05"))
	fmt.Printf("| 【检测完成，耗时：%s 】                                 |\n", time.Since(start))
	fmt.Println("==============================================================")
}
