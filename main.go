//go:generate goversioninfo
package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

type Config struct {
	PersonalStorage int64 `toml:"personalstorage"`
	CommonStorage   int64 `toml:"commonstorage"`
}

var Cfg Config

// 配置文件初始化
func init() {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("toml")
	v.AddConfigPath(".")
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Printf("配置文件不存在:%v\n", err)
		} else {
			fmt.Printf("配置文件存在,解析失败:%v\n", err)
		}
	}
	if err := v.Unmarshal(&Cfg); err != nil {
		fmt.Println(err)
	}
}

func main() {
	//定义一个G大小
	Gigabyte := Cfg.PersonalStorage << 30
	//定义当前文件夹大小
	CurrentDirTotalSize := int64(0)

	dir, _ := os.Getwd()
	userdir := dir + "\\" + os.Getenv("GIT_AUTHOR_NAME")

	//计算当前文件夹大小
	err := filepath.Walk(userdir, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			CurrentDirTotalSize += info.Size()
		}
		return err
	})
	if err != nil {
		fmt.Printf("计算目录大小异常,请重试\n %s", err)
		time.Sleep(5 * time.Second)
		os.Exit(0)
	}

	if CurrentDirTotalSize > Gigabyte {
		fmt.Printf("目录：%s, 大于4G,请缩减不必要文件后再次尝试提交!\n", userdir)
	} else {
		// 检测用户目录差异
		output, err := exec.Command("cmd", "/c", "git", "diff", userdir).Output()
		if err != nil {
			fmt.Printf("分析目录差异异常：%s,请重试！\n", err)
			time.Sleep(10 * time.Second)
			os.Exit(0)
		}
		if strings.Contains(string(output), "index") {
			fmt.Println("⛏ ⛏ ⛏ 检测本地文件有变更,开始提交...")
			output1, err1 := exec.Command("cmd", "/c", "git", "commit", "-m", fmt.Sprintf("\"用户: %s ,提交文件\"", os.Getenv("GIT_AUTHOR_NAME")), userdir).Output()
			if err1 != nil {
				fmt.Printf("提交文件异常：%s,请重试！\n", err1)
				time.Sleep(10 * time.Second)
				os.Exit(0)
			}
			fmt.Println(string(output1))
			output2, err2 := exec.Command("cmd", "/c", "git", "push").Output()
			if err2 != nil {
				fmt.Printf("提交文件异常：%s,请重试！\n", err2)
				time.Sleep(10 * time.Second)
				os.Exit(0)
			}
			fmt.Println(string(output2))
			fmt.Println("✔ ✔ ✔ 提交成功！")
		} else {
			fmt.Println("⚙ ⚙ ⚙ 本地文件没有变更，请重新打开文件，检查文件内容后再次提交")
		}
	}
	time.Sleep(20 * time.Second)
}
