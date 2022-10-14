package model

import (
	"gorm.io/gorm"
	"time"
)

type User struct {
	gorm.Model
	Username string `gorm:"unique"`
	Password string
}

type Memo struct {
	gorm.Model
	UserID  uint
	User    User
	Title   string
	Content string
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
}

type MemoInfo struct {
	ID           uint      `json:"id"`
	Title        string    `json:"title"`
	LastModified time.Time `json:"lastModified"`
}

type MemoDetail struct {
	MemoInfo
	Content string `json:"content"`
}

func (u *User) Info() UserInfo {
	return UserInfo{ID: u.ID, Username: u.Username}
}

func (m *Memo) Info() MemoInfo {
	return MemoInfo{
		ID:           m.ID,
		Title:        m.Title,
		LastModified: m.UpdatedAt,
	}
}

func (m *Memo) Detail() MemoDetail {
	return MemoDetail{
		MemoInfo: m.Info(),
		Content:  m.Content,
	}
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&User{},
		&Memo{},
	)
}
