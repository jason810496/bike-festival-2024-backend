package service

import (
	"context"
	"time"

	tokensvc "bikefest/internal/token"
	"bikefest/pkg/model"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewUserService(db *gorm.DB, cache *redis.Client) model.UserService {
	return &UserServiceImpl{
		cache: cache,
	}
}

type UserServiceImpl struct {
	db    *gorm.DB
	cache *redis.Client
}

// GetUserByID implements model.UserService.
func (us *UserServiceImpl) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	user := &model.User{}
	err := us.db.WithContext(ctx).Where("id = ?", id).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (us *UserServiceImpl) CreateAccessToken(_ context.Context, user *model.User, secret string, expiry int64) (accessToken string, err error) {
	return tokensvc.CreateAccessToken(user, secret, expiry)
}

func (us *UserServiceImpl) CreateRefreshToken(_ context.Context, user *model.User, secret string, expiry int64) (refreshToken string, err error) {
	return tokensvc.CreateRefreshToken(user, secret, expiry)
}

func (*UserServiceImpl) VerifyRefreshToken(_ context.Context, refreshToken string, secret string) (user *model.User, err error) {
	return tokensvc.VerifyRefreshToken(refreshToken, secret)
}

// Logout implements model.UserService.
func (us *UserServiceImpl) Logout(ctx context.Context, token *string, secrect string) error {
	claims, err := tokensvc.ExtractCustomClaimsFromToken(token, secrect)
	if err != nil {
		return err
	}
	ttl := claims.ExpiresAt.Unix() - time.Now().Unix()
	if ttl < 0 {
		return nil
	}
	return us.cache.Set(ctx, claims.ID, *token, time.Duration(ttl)*time.Second).Err()
}
