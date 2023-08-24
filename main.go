package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"os"
	"path/filepath"
	"time"
)

func ReversCode(code int) string {
	switch code {
	case 32:
		return "无变化"
	case 63:
		return "未提交"
	case 77:
		return "修改文件"
	case 65:
		return "新增文件"
	case 68:
		return "删除文件"
	case 82:
		return "重命名文件"
	default:
		return ""
	}
}

func main() {
	//定义一个管道阻塞程序
	//ch := make(chan int)
	//定义一个G大小
	const Gigabyte = 8 << 30
	//定义当前文件夹大小
	CurrentDirTotalSize := int64(0)

	dir, _ := os.Getwd()
	gitdir := filepath.Dir(dir)

	//计算当前文件夹大小
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			CurrentDirTotalSize += info.Size()
		}
		return err
	})
	if err != nil {
		fmt.Printf("计算目录大小异常,请重试\n %s", err)
		time.Sleep(5 * time.Second)
	}

	if CurrentDirTotalSize > Gigabyte {
		fmt.Printf("当前目录：%s, 大于4G,请缩减不必要文件后再次尝试提交!\n", dir)
	} else {
		// 只获取当前目录

		// 打开当前目录作为仓库
		repo, err := git.PlainOpen(gitdir)
		if err != nil {
			fmt.Printf("获取工作目录异常,请重试\n %s", err)
			time.Sleep(5 * time.Second)
		}

		// 获取工作树(工作目录)和暂存树的差异
		wt, err := repo.Worktree()

		if err != nil {
			fmt.Printf("获取工作目录异常,请重试\n %s", err)
			time.Sleep(5 * time.Second)
		}
		start := time.Now()
		fmt.Println("=====================================================================")
		fmt.Printf("| 【检测当前文件夹文件是否发生变更，开始时间：%s 】|\n", start.Format("2006-01-02 15:04:05"))
		diff, err := wt.Status()
		fmt.Printf("| 【检测完成，耗时：%s 】                                    |\n", time.Since(start))
		fmt.Println("=====================================================================\n\n")

		if err != nil {
			fmt.Printf("检测工作目录异常,请重试\n %s", err)
			time.Sleep(5 * time.Second)
		}

		if !diff.IsClean() {
			// 工作树和暂存区有差异
			fmt.Println("⛏ ⛏ ⛏ 检测本地文件有变更,变更的文件列表：\n")
			flag := 1
			fmt.Println("------------------------------------")
			for key, value := range diff {
				fmt.Printf("%d.%-20s   %-20s\n", flag, key, ReversCode(int(value.Worktree)))
				flag += 1
			}
			fmt.Println("------------------------------------")

			// 添加所有变化到暂存区
			_, err = wt.Add(".")
			if err != nil {
				fmt.Printf("提交异常,请重试\n %s", err)
				time.Sleep(5 * time.Second)
			}

			// 提交
			commit, err := wt.Commit("init", &git.CommitOptions{})
			if err != nil {
				fmt.Printf("提交异常,请重试\n %s", err)
				time.Sleep(5 * time.Second)
			}

			// 使用commit变量
			_, err = repo.CommitObject(commit)
			if err != nil {
				fmt.Printf("提交异常,请重试\n %s", err)
				time.Sleep(5 * time.Second)
			}
			fmt.Println("✔✔✔ 提交成功！\n")

			//// 推送到远程origin的master分支
			//err = repo.Push(&git.PushOptions{
			//	RemoteName: "origin",
			//	RefSpecs: []config.RefSpec{
			//		"refs/heads/master",
			//	},
			//})
			//if err != nil {
			//	fmt.Printf("提交异常,请重试\n %s", err)
			//	time.Sleep(5 * time.Second)
			//}
		} else {
			fmt.Println("☂ ☂ ☂ 本地文件没有变更，请重新打开文件，检查文件内容后再次提交\n")
		}
	}
	time.Sleep(5 * time.Second)
}
