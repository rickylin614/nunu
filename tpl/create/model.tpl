package model

import "gorm.io/gorm"

type {{ .FileName }} struct {
	gorm.Model
}

func (*{{ .FileName }}) TableName() string {
	return "{{ .FileNameSnakeCase }}"
}
