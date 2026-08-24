package domain

import "gorm.io/gorm"

type Topic struct {
	gorm.Model
	ChapterID uint  `gorm:"not null;index"`
	UserID    *uint `gorm:"index"` // 작성자 탈퇴 시 SET NULL을 위해 포인터 처리

	Title   string `gorm:"type:varchar(200);not null"`
	Content string `gorm:"type:text;not null"`

	// 조회용 매핑 (주제 입장에서 부모 정보 가져오기)
	Chapter Chapter `gorm:"foreignKey:ChapterID;"`
	User    User    `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL;"`
}