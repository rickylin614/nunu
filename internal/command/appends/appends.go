package appends

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	"github.com/rickylin614/nunu/internal/pkg/helper"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

type Append struct {
	ProjectName        string
	FilePath           string
	FileName           string
	FileNameTitleLower string
	FileNameFirstChar  string
	FileNameSnakeCase  string
	IsFull             bool
	TemplateFiles      map[string]*template.Template // templates of an external project
	Config             Config
}

func NewAppend() *Append {
	return &Append{}
}

type File struct {
	Path     string `yaml:"path"`
	Regex    string `yaml:"regex"`
	Template string `yaml:"template"`
}

type Config struct {
	Files []File `yaml:"files"`
}

var AppendCmd = &cobra.Command{
	Use:     "append [handler-name]",
	Short:   "Append Template to Regexp path",
	Example: "nunu append user",
	Args:    cobra.ExactArgs(1),
	Run:     runProvider,
}

func runProvider(cmd *cobra.Command, args []string) {
	a := NewAppend()
	a.FilePath, a.FileName = filepath.Split(args[0])
	a.FileName = strings.ReplaceAll(strings.ToUpper(string(a.FileName[0]))+a.FileName[1:], ".go", "")
	a.FileNameTitleLower = strings.ToLower(string(a.FileName[0])) + a.FileName[1:]
	a.FileNameSnakeCase = helper.ToSnakeCase(a.FileName)
	a.InitConfig()
	a.AppendTemplate()
}

func (a *Append) AppendTemplate() {
	for _, file := range a.Config.Files {
		// 讀取檔案內容
		data, err := ioutil.ReadFile(file.Path)
		if err != nil {
			log.Fatalf("\033[33;1mcmd run failed %s\u001B[0m", err)
		}

		// 復原正則
		file.Regex = helper.ReplaceEscapeString(file.Regex)
		file.Template = helper.ReplaceEscapeString(file.Template)

		// 創建一個多行匹配的正則表達式
		re := regexp.MustCompile("(?s)" + file.Regex)

		// 找到所有的匹配
		matches := re.FindAllStringIndex(string(data), -1)

		// 結果初始化為原始數據
		result := string(data)

		if len(matches) == 0 {
			log.Printf("not found any match: %s", file.Path)
			continue
		}

		// 從最後一個匹配開始，反向遍曆所有的匹配
		for i := len(matches) - 1; i >= 0; i-- {
			// 找到匹配的開始和結束
			start := matches[i][0]
			end := matches[i][1]

			// 找到匹配內的所有行
			lines := strings.Split(string(data)[start:end], "\n")

			// 如果有超過一行，則在倒數第二行後添加模板
			if len(lines) > 1 {
				// 根據模板執行資料替換
				tmpl, err := template.New("test").Parse(file.Template)
				if err != nil {
					log.Fatalf("Parse: %v", err)
				}
				var tpl bytes.Buffer
				err = tmpl.Execute(&tpl, a)
				if err != nil {
					log.Fatalf("Execute: %v", err)
				}

				// 插入處理過的模板資料
				lines[len(lines)-2] += tpl.String()

				// 將修改後的匹配替換回結果
				result = result[:start] + strings.Join(lines, "\n") + result[end:]
			}
		}

		// 將修改後的結果寫回檔案
		err = ioutil.WriteFile(file.Path, []byte(result), 0644)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Append template to: %s", file.Path)
	}
}

func (a *Append) InitConfig() {
	file, err := os.Open("./template/nunu/append.yaml")
	if err != nil {
		log.Fatalf("read ./template/nunu/append.yaml error: %v", err)
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("read ./template/nunu/append.yaml error: %v", err)
		return
	}

	a.Config = Config{}
	if err := yaml.Unmarshal(content, &a.Config); err != nil {
		log.Fatalf("Unmarshal append.yaml error :%v", err)
		return
	}
}
