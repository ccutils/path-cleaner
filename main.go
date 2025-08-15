package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"strings"
)

func main() {
	// 解析命令行参数
	pathFlag := flag.String("path", "", "Custom PATH string (optional)")
	flag.Parse()
	removePrefixes := flag.Args()

	// 获取系统PATH
	inputPath := getInputPath(*pathFlag)

	// 处理路径分隔符
	sep := getPathSeparator()
	paths := splitPath(inputPath, sep)

	// 过滤需要移除的前缀
	filtered := filterPaths(paths, removePrefixes)

	// 生成最终PATH
	result := strings.Join(filtered, sep)
	fmt.Println(result)
}

func getInputPath(customPath string) string {
	if customPath != "" {
		return customPath
	}
	return os.Getenv("PATH")
}

func getPathSeparator() string {
	switch runtime.GOOS {
	case "windows":
		return ";"
	default:
		return ":"
	}
}

func splitPath(path string, sep string) []string {
	if path == "" {
		return []string{}
	}
	return strings.Split(path, sep)
}

func filterPaths(paths []string, prefixes []string) []string {
	var result []string
	for _, p := range paths {
		if p == "" {
			continue
		}
		if shouldRemove(p, prefixes) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func shouldRemove(path string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}