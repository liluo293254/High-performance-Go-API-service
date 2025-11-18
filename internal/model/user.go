package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Username  string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username"`
	Password  string         `gorm:"type:varchar(255);not null" json:"password"`
	Email     string         `gorm:"type:varchar(100);" json:"email"`
	Nickname  string         `gorm:"type:varchar(50)" json:"nickname"`
	Status    int            `gorm:"type:tinyint;default:1" json:"status"` // 1:正常 2:禁用
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

type UserProfile struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	UserID          uint           `gorm:"index;comment:用户ID" json:"user_id"` // 关联到users表的ID
	Balance         float64        `gorm:"type:decimal(10,2);default:0.00;comment:余额" json:"balance"`
	ActivityBalance float64        `gorm:"type:decimal(10,2);default:0.00;comment:活动余额" json:"activity_balance"`
	Level           int            `gorm:"type:int;default:1;comment:等级" json:"level"`
	Experience      int            `gorm:"type:int;default:0;comment:经验值" json:"experience"`
	RegisterTime    time.Time      `gorm:"comment:注册时间" json:"register_time"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

func (UserProfile) TableName() string {
	return "user_profiles"
}
