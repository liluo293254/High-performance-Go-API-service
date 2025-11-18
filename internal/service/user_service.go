package service

import (
	"errors"
	"time"

	"bgame/internal/dao"
	"bgame/internal/model"
	"bgame/internal/util"
)

type UserService struct {
	userDAO        *dao.UserDAO
	userProfileDAO *dao.UserProfileDAO
}

func NewUserService() *UserService {
	return &UserService{
		userDAO:        dao.NewUserDAO(),
		userProfileDAO: dao.NewUserProfileDAO(),
	}
}

type RegAndLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type RegAndLoginResponse struct {
	Token    string      `json:"token"`
	UserInfo *model.User `json:"user_info"`
	UserProfile *model.UserProfile `json:"user_profile"`
}

// RegAndLoginRequest 注册和登录
func (s *UserService) RegAndLogin(req *RegAndLoginRequest) (*RegAndLoginResponse, error) {
	// 检查用户名是否已存在
	_, err := s.userDAO.GetByUsername(req.Username)
	if err == nil {
		return nil, errors.New("用户名已存在")
	}
	// 加密密码
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("密码加密失败")
	}
	// 创建用户
	user := &model.User{
		Username: req.Username,
		Password: hashedPassword,
		Nickname: req.Username,
		Email:    req.Username + "@bgame.com",
		Status:   1,
	}
	if err := s.userDAO.Create(user); err != nil {
		return nil, errors.New("创建用户失败")
	}
	// 创建用户资料
	userProfile := &model.UserProfile{
		UserID:          user.ID,
		Balance:         0,
		ActivityBalance: 0,
		Level:           1,
		Experience:      0,
		RegisterTime:    time.Now(),
	}
	if err := s.userProfileDAO.CreateUserProfile(userProfile); err != nil {
		return nil, errors.New("创建用户资料失败")
	}
	token, err := util.GenerateUserToken(user.ID, user.Username)
	if err != nil {
		return nil, errors.New("生成token失败")
	}
	userProfile, err = s.userProfileDAO.GetUserProfileByUserID(user.ID)
	if err != nil {
		return nil, errors.New("获取用户资料失败")
	}
	return &RegAndLoginResponse{
		Token: token,
		UserInfo: &model.User{
			ID:        user.ID,
			Username:  user.Username,
			Email:     user.Email,
			Nickname:  user.Nickname,
			Status:    user.Status,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
		UserProfile: &model.UserProfile{
			Balance:         userProfile.Balance,
			ActivityBalance: userProfile.ActivityBalance,
			Level:           userProfile.Level,
			Experience:      userProfile.Experience,
			RegisterTime:    userProfile.RegisterTime,
		},
	}, nil
}

// GetUserInfo 获取用户信息
func (s *UserService) GetUserInfo(userID uint) (*model.UserProfile, error) {
	user, err := s.userProfileDAO.GetUserProfileByUserID(userID)
	if err != nil {
		return nil, errors.New("用户不存在")
	}

	return &model.UserProfile{
		UserID:          user.UserID,
		Balance:         user.Balance,
		ActivityBalance: user.ActivityBalance,
		Level:           user.Level,
		Experience:      user.Experience,
		RegisterTime:    user.RegisterTime,
	}, nil
}

// CreateUserProfile 创建用户资料
func (s *UserService) CreateUserProfile(userID uint, userProfile *model.UserProfile) error {
	_, err := s.userDAO.GetByID(userID)
	if err == nil {
		return errors.New("用户不存在")
	}
	return s.userProfileDAO.CreateUserProfile(userProfile)
}

// GetUserProfileByUserID 根据用户ID获取用户资料
func (s *UserService) GetUserProfileByUserID(userID uint) (*model.UserProfile, error) {
	return s.userProfileDAO.GetUserProfileByUserID(userID)
}
