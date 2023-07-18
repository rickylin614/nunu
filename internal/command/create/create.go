package create

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rickylin614/nunu/internal/pkg/helper"
	"github.com/rickylin614/nunu/tpl"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

type Create struct {
	ProjectName        string
	CreateType         string
	FilePath           string
	FileName           string
	FileNameTitleLower string
	FileNameFirstChar  string
	FileNameSnakeCase  string
	IsFull             bool
	TemplateFiles      map[string]*template.Template // templates of an external project
	Config             Config
}

func NewCreate() *Create {
	return &Create{}
}

var CreateCmd = &cobra.Command{
	Use:     "create [type] [handler-name]",
	Short:   "Create a new handler/service/repository/model",
	Example: "nunu create handler user",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var CreateHandlerCmd = &cobra.Command{
	Use:     "handler",
	Short:   "Create a new handler",
	Example: "nunu create handler user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CreateServiceCmd = &cobra.Command{
	Use:     "service",
	Short:   "Create a new service",
	Example: "nunu create service user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CreateRepositoryCmd = &cobra.Command{
	Use:     "repository",
	Short:   "Create a new repository",
	Example: "nunu create repository user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CreateModelCmd = &cobra.Command{
	Use:     "model",
	Short:   "Create a new model",
	Example: "nunu create model user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}
var CreateAllCmd = &cobra.Command{
	Use:     "all",
	Short:   "Create a new handler & service & repository & model",
	Example: "nunu create all user",
	Args:    cobra.ExactArgs(1),
	Run:     runCreate,
}

func runCreate(cmd *cobra.Command, args []string) {
	c := NewCreate()
	c.ProjectName = helper.GetProjectName(".")
	c.CreateType = cmd.Use
	c.FilePath, c.FileName = filepath.Split(args[0])
	c.FileName = strings.ReplaceAll(strings.ToUpper(string(c.FileName[0]))+c.FileName[1:], ".go", "")
	c.FileNameTitleLower = strings.ToLower(string(c.FileName[0])) + c.FileName[1:]
	c.FileNameSnakeCase = helper.ToSnakeCase(c.FileName)
	c.InitConfig()

	switch c.CreateType {
	case "handler", "service", "repository", "model":
		c.genFile()
	case "all":

		c.CreateType = "handler"
		c.genFile()

		c.CreateType = "service"
		c.genFile()

		c.CreateType = "repository"
		c.genFile()

		c.CreateType = "model"
		c.genFile()
	default:
		log.Fatalf("Invalid handler type: %s", c.CreateType)
	}
}

func (c *Create) genFile() {
	for _, v := range c.GetPath() {
		// create file
		fileName := strings.ToLower(c.FileName)
		f := createFile(v.Path, fileName+".go")
		if f == nil {
			log.Printf("warn: file %s%s %s", v.Path, fileName+".go", "already exists.")
			return
		}
		defer f.Close()

		// get template
		t := c.GetTemplate(v.TempFile)

		err := t.Execute(f, c)
		if err != nil {
			log.Fatalf("create %s error: %s", c.CreateType, err.Error())
		}
		log.Printf("Created new %s: %s", c.CreateType, v.Path+fileName+".go")
	}
}

func createFile(dirPath string, filename string) *os.File {
	filePath := dirPath + filename
	// 创建文件夹
	err := os.MkdirAll(dirPath, os.ModePerm)
	if err != nil {
		log.Fatalf("Failed to create dir %s: %v", dirPath, err)
	}
	stat, _ := os.Stat(filePath)
	if stat != nil {
		return nil
	}
	// 创建文件
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatalf("Failed to create file %s: %v", filePath, err)
	}

	return file
}

func (c *Create) GetTemplate(tempalteFileName string) *template.Template {
	// check template
	matches, err := filepath.Glob("./template/nunu/*.tpl")
	if err != nil {
		log.Fatal(err)
	}

	if c.TemplateFiles == nil {
		c.TemplateFiles = make(map[string]*template.Template, len(matches))

		for _, path := range matches {
			dataTemp, err := template.ParseFiles(path)
			// data, err := ioutil.ReadFile(path)
			if err != nil {
				log.Fatalf("failed reading data from file: %s", err)
			}

			// Get the file name from the path
			_, filename := filepath.Split(path)
			// Save the file content as string in the map
			c.TemplateFiles[filename] = dataTemp
		}
	}

	if v, ok := c.TemplateFiles[tempalteFileName]; ok {
		return v
	} else {
		t, err := template.ParseFS(tpl.CreateTemplateFS, fmt.Sprintf("create/%s.tpl", c.CreateType))
		if err != nil {
			log.Fatalf("create %s error: %s", c.CreateType, err.Error())
		}
		return t
	}

}

// 確認目標路由
func (c *Create) GetPath() []Path {
	// 設定檔有Mode的流程
	if c.CreateType == "model" && len(c.Config.TargetPath.Model) > 0 {
		return c.Config.TargetPath.Model
	}

	// 沒有Model的流程
	filePath := c.FilePath
	TempFile := fmt.Sprintf("%s.tpl", c.CreateType)

	// 判斷檔案路徑
	if filePath == "" {
		switch c.CreateType {
		case "handler":
			filePath = c.Config.TargetPath.Handler
		case "service":
			filePath = c.Config.TargetPath.Service
		case "repository":
			filePath = c.Config.TargetPath.Repository
		}
	}
	if filePath == "" {
		filePath = fmt.Sprintf("internal/%s/", c.CreateType)
	}
	return []Path{{Path: filePath, TempFile: TempFile}}
}

func (c *Create) InitConfig() {
	file, err := os.Open("./template/nunu/target.yaml")
	if err != nil {
		return
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	c.Config = Config{}
	if err := yaml.Unmarshal(content, &c.Config); err != nil {
		return
	}
}

type Path struct {
	Path     string `yaml:"path"`
	TempFile string `yaml:"temp_file"`
}

type Config struct {
	TargetPath struct {
		Handler    string `yaml:"handler"`
		Service    string `yaml:"service"`
		Repository string `yaml:"repository"`
		Model      []Path `yaml:"model"`
	} `yaml:"target_path"`
}
