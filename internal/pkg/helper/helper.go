package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

func GetProjectName(dir string) string {

	// 打开 go.mod 文件
	modFile, err := os.Open(dir + "/go.mod")
	if err != nil {
		fmt.Println("go.mod does not exist", err)
		return ""
	}
	defer modFile.Close()

	var moduleName string
	_, err = fmt.Fscanf(modFile, "module %s", &moduleName)
	if err != nil {
		fmt.Println("read go mod error: ", err)
		return ""
	}
	return moduleName
}

func SplitArgs(cmd *cobra.Command, args []string) (cmdArgs, programArgs []string) {
	dashAt := cmd.ArgsLenAtDash()
	if dashAt >= 0 {
		return args[:dashAt], args[dashAt:]
	}
	return args, []string{}
}

// 查詢所有的
func FindMain(base string) (map[string]string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	if !strings.HasSuffix(wd, "/") {
		wd += "/"
	}
	cmdPath := make(map[string]string)
	err = filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".go" {
			content, err := os.ReadFile(path)
			if err != nil {
				return err
			}
			if !strings.Contains(string(content), "package main") {
				return nil
			}
			re := regexp.MustCompile(`func\s+main\s*\(`)
			if re.Match(content) {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				d, _ := filepath.Split(absPath)
				cmdPath[strings.TrimPrefix(absPath, wd)] = d

			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return cmdPath, nil
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func ToKebabCase(str string) string {
	kebab := matchFirstCap.ReplaceAllString(str, "${1}-${2}")
	kebab = matchAllCap.ReplaceAllString(kebab, "${1}-${2}")
	return strings.ToLower(kebab)
}

func ReplaceEscapeString(str string) string {
	output := strings.ReplaceAll(str, "\\t", "\t")
	output = strings.ReplaceAll(output, "\\n", "\n")
	return output
}
