package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"log"
	"os"
	"path/filepath"
)

func main() {
	CurrentDirTotalSize := int64(0)
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			CurrentDirTotalSize += info.Size()
		}
		return err
	})
	if err != nil {
		log.Fatal("计算文件大小异常：", err)
	}

	fmt.Println(CurrentDirTotalSize)

	//time.Sleep(10 * time.Second)

	// 打开当前目录作为仓库
	repo, err := git.PlainOpen(".")
	if err != nil {
		log.Fatal(err)
	}

	// 获取工作树(工作目录)和暂存树的差异
	wt, err := repo.Worktree()
	diff, err := wt.Status()

	if !diff.IsClean() {
		// 工作树和暂存区有差异
		fmt.Println("有变化")
		// 添加所有变化到暂存区
		_, err = wt.Add(".")
		if err != nil {
			log.Fatal(err)
		}

		// 提交
		commit, err := wt.Commit("init", &git.CommitOptions{})
		if err != nil {
			log.Fatal(err)
		}

		// 使用commit变量
		_, err = repo.CommitObject(commit)
		if err != nil {
			log.Fatal(err)
		}

		//// 推送到远程origin的master分支
		//err = repo.Push(&git.PushOptions{
		//	RemoteName: "origin",
		//	RefSpecs: []config.RefSpec{
		//		"refs/heads/master",
		//	},
		//})
		//if err != nil {
		//	log.Fatal(err)
		//}
	} else {
		fmt.Println("无变化")
	}

}
