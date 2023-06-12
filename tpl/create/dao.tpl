package dao

import (
	"{{ .ProjectName }}/internal/model"
)

type {{ .FileName }}Dao struct {
	*Dao
}

func New{{ .FileName }}Dao(dao *Dao) *{{ .FileName }}Dao {
	return &{{ .FileName }}Dao{
		Dao: dao,
	}
}

func (d *{{ .FileName }}Dao) FirstById(id int64) (*model.{{ .FileName }}, error) {
	var {{ .FileNameTitleLower }} model.{{ .FileName }}
	if err := d.db.Where("id = ?", id).First(&{{ .FileNameTitleLower }}).Error; err != nil {
		return nil, err
	}
	return &{{ .FileNameTitleLower }}, nil
}