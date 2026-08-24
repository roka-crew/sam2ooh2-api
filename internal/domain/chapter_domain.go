package domain

import (
	"time"

	"gorm.io/gorm"
)

type Chapter struct {
	gorm.Model
	GroupID    uint      `gorm:"not null;index"`
	UserID     uint      `gorm:"not null;index"`
	TargetPage int       `gorm:"type:integer;not null;check:target_page > 0"`
	Deadline   time.Time `gorm:"not null"`

	// Chapter < 1:N > Topic : 이 챕터 안에서 논의될 주제들
	Topics []Topic `gorm:"foreignKey:ChapterID;constraint:OnDelete:CASCADE;"`
	
	// 조회용 매핑 (챕터 입장에서 부모 정보 가져오기)
	Group Group `gorm:"foreignKey:GroupID;"`
	User  User  `gorm:"foreignKey:UserID;"`
}