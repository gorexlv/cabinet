package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorexlv/cabinet/scissor/internal/domain"
	"github.com/gorexlv/cabinet/scissor/internal/repository"
	"github.com/gorexlv/cabinet/scissor/pkg/ent"
	"github.com/gorexlv/cabinet/scissor/pkg/jwt"
	"github.com/gorexlv/cabinet/scissor/pkg/wechat"
)

type UserService struct {
	repo      *repository.UserRepository
	jwtSecret string
	wxClient  *wechat.Client
}

func NewUserService(repo *repository.UserRepository, jwtSecret string, wxClient *wechat.Client) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
		wxClient:  wxClient,
	}
}

func (s *UserService) Create(ctx context.Context, req *domain.CreateUserRequest) (*domain.LoginResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.repo.ExistsByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &ent.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Email:    req.Email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := jwt.GenerateToken(uint(user.ID), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token: token,
		User: &domain.User{
			ID:        uint(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			WxOpenID:  user.WxOpenID,
			Nickname:  user.Nickname,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *UserService) Login(ctx context.Context, req *domain.LoginRequest) (*domain.LoginResponse, error) {
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, errors.New("invalid password")
	}

	token, err := jwt.GenerateToken(uint(user.ID), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token: token,
		User: &domain.User{
			ID:        uint(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			WxOpenID:  user.WxOpenID,
			Nickname:  user.Nickname,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}

func (s *UserService) WxLogin(ctx context.Context, req *domain.WxLoginRequest) (*domain.LoginResponse, error) {
	// 查找或创建用户
	user, err := s.repo.FindByWxOpenID(ctx, req.OpenID)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}

		// 用户不存在，创建新用户
		// 生成随机密码
		password, err := bcrypt.GenerateFromPassword([]byte(time.Now().String()), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("生成密码失败: %v", err)
		}

		user = &ent.User{
			Username:  fmt.Sprintf("wx_%s", req.OpenID[:8]),
			Password:  string(password),
			WxOpenID:  req.OpenID,
			Nickname:  req.Nickname,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		if err := s.repo.Create(ctx, user); err != nil {
			return nil, err
		}
	} else {
		// 更新用户信息
		user, err = s.repo.Update(ctx, int(user.ID), map[string]interface{}{
			"nickname":   req.Nickname,
			"updated_at": time.Now(),
		})
		if err != nil {
			return nil, err
		}
	}

	token, err := jwt.GenerateToken(uint(user.ID), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	return &domain.LoginResponse{
		Token: token,
		User: &domain.User{
			ID:        uint(user.ID),
			Username:  user.Username,
			Email:     user.Email,
			WxOpenID:  user.WxOpenID,
			Nickname:  user.Nickname,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		},
	}, nil
}
