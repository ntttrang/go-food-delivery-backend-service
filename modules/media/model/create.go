package mediamodel

import "fmt"

type ImageCreateDTO struct {
	Filename  string `json:"filename" gorm:"column:filename;"`
	CloudName string `json:"cloudName" gorm:"column:cloud_name;"`
	Size      int64  `json:"size" gorm:"column:size;"`
	Ext       string `json:"ext" gorm:"column:ext;"`
	Url       string `json:"url" gorm:"column:url;"`
}

func (m *ImageCreateDTO) Fulfill(domain string) {
	m.Url = fmt.Sprintf("%s/%s", domain, m.Filename)
}

func (ImageCreateDTO) TableName() string {
	return Image{}.TableName()
}

func (m *ImageCreateDTO) Validate() error {
	return nil
}
