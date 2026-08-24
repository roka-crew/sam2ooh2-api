package domain

import (
	"gorm.io/gorm"
)

type Gruops []Group

type Group struct {
	gorm.Model
	Title       string `gorm:"type:varchar(200);index"`
	Author      string `gorm:"type:varchar(100)"`
	PageCount   int    `gorm:"type:integer;check:page_count > 0"`
	Publisher   string `gorm:"type:varchar(100)"`
	Description string `gorm:"type:varchar(1024)"`

	// Group < N:M > User
	Users []User `gorm:"many2many:users_groups;"`

	// Group < 1:N > Chapter : 그룹 내에 생성된 챕터들
	Chapters []Chapter `gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE;"`
}
