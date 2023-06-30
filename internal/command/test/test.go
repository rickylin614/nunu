package test

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/spf13/cobra"
)

type Test struct {
	ProjectName        string
	CreateType         string
	FilePath           string
	FileName           string
	FileNameTitleLower string
	FileNameFirstChar  string
	IsFull             bool
	TemplateFiles      map[string]string // templates of an external project
}

func NewTest() *Test {
	return &Test{}
}

var TestCmd = &cobra.Command{
	Use:     "create [type] [handler-name]",
	Short:   "Create a new handler/service/repository/model",
	Example: "nunu create handler user",
	Args:    cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {

	},
}
var PrintCmd = &cobra.Command{
	Use:     "print",
	Short:   "Create a new handler",
	Example: "nunu create handler user",
	Args:    cobra.ExactArgs(1),
	Run:     print,
}

func print(cmd *cobra.Command, args []string) {
	matches, err := filepath.Glob("./create/*.tpl")
	if err != nil {
		log.Fatal(err)
	}

	templates := make(map[string]string)

	for _, file := range matches {
		data, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatalf("failed reading data from file: %s", err)
		}
		// Get the file name from the path
		_, filename := filepath.Split(file)
		// Save the file content as string in the map
		templates[filename] = string(data)
	}

	// Now templates map contains file names as keys and file content as values
	for name, content := range templates {
		fmt.Printf("\nFile: %s", name)
		fmt.Printf("\nData: %s", content)
	}
}
