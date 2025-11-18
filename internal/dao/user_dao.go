package dao

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"bgame/internal/model"
	"bgame/pkg/mysql"
	"bgame/pkg/redis"
)

const (
	userCachePrefix = "user:"
	userCacheTTL    = 3600 * time.Second
)

type UserDAO struct{}
type UserProfileDAO struct{}

func NewUserDAO() *UserDAO {
	return &UserDAO{}
}

func NewUserProfileDAO() *UserProfileDAO {
	return &UserProfileDAO{}
}

// Create 创建用户
func (d *UserDAO) Create(user *model.User) error {
	return mysql.DB.Create(user).Error
}

// GetByID 根据ID获取用户（带缓存）
func (d *UserDAO) GetByID(id uint) (*model.User, error) {
	// 先查缓存
	cacheKey := fmt.Sprintf("%s%d", userCachePrefix, id)
	cached, err := redis.Client.Get(context.Background(), cacheKey).Result()
	if err == nil {
		var user model.User
		if json.Unmarshal([]byte(cached), &user) == nil {
			return &user, nil
		}
	}

	// 查数据库排除软删除和password字段
	var user model.User
	if err := mysql.DB.Where("id = ? AND status = 1").Select("id, username, email, nickname, status, created_at, updated_at").First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByUsername 根据用户名获取用户
func (d *UserDAO) GetByUsername(username string) (*model.User, error) {
	var user model.User
	if err := mysql.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// GetByEmail 根据邮箱获取用户
func (d *UserDAO) GetByEmail(email string) (*model.User, error) {
	var user model.User
	if err := mysql.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (d *UserDAO) Update(user *model.User) error {
	err := mysql.DB.Save(user).Error
	if err == nil {
		// 清除缓存
		cacheKey := fmt.Sprintf("%s%d", userCachePrefix, user.ID)
		redis.Client.Del(context.Background(), cacheKey)
	}
	return err
}

// DeleteCache 删除用户缓存
func (d *UserDAO) DeleteCache(userID uint) {
	cacheKey := fmt.Sprintf("%s%d", userCachePrefix, userID)
	redis.Client.Del(context.Background(), cacheKey)
}

// CreateUserProfile 创建用户资料
func (d *UserProfileDAO) CreateUserProfile(userProfile *model.UserProfile) error {
	return mysql.DB.Create(userProfile).Error
}

// GetUserProfileByUserID 根据用户ID获取用户资料
func (d *UserProfileDAO) GetUserProfileByUserID(userID uint) (*model.UserProfile, error) {
	var userProfile model.UserProfile
	if err := mysql.DB.Where("user_id = ?", userID).First(&userProfile).Error; err != nil {
		return nil, err
	}
	return &userProfile, nil
}

// UpdateUserProfileByUserID 根据用户ID更新用户资料
func (d *UserProfileDAO) UpdateUserProfileByUserID(userID uint, userProfile *model.UserProfile) error {
	return mysql.DB.Model(&model.UserProfile{}).Where("user_id = ?", userID).Updates(userProfile).Error
}
