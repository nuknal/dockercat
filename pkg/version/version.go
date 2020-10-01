package version

import (
	"flag"
	"fmt"
)

// 版本信息
var (
	BuildVersion string
	BuildTime    string
	BuildName    string
	CommitID     string
	GoVersion    string
)

// Show 显示版本信息
func Show() {
	fmt.Printf("Commit ID  : %s\n", CommitID)
	fmt.Printf("Build  Name: %s\n", BuildName)
	fmt.Printf("Build  Time: %s\n", BuildTime)
	fmt.Printf("Build  Vers: %s\n", BuildVersion)
	fmt.Printf("Golang Vers: %s\n", GoVersion)
}

// ShowVersion 命令行参数处理
var ShowVersion bool

func init() {
	flag.BoolVar(&ShowVersion, "v", false, "show version and exit")
	flag.BoolVar(&ShowVersion, "version", false, "show version and exit")
}
