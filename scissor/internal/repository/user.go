package repository

import (
	"context"
	"time"

	"github.com/gorexlv/cabinet/scissor/pkg/ent"
	"github.com/gorexlv/cabinet/scissor/pkg/ent/user"
)

type UserRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) *UserRepository {
	return &UserRepository{client: client}
}

func (r *UserRepository) Create(ctx context.Context, user *ent.User) error {
	_, err := r.client.User.Create().
		SetUsername(user.Username).
		SetPassword(user.Password).
		SetEmail(user.Email).
		SetNickname(user.Nickname).
		SetWxOpenID(user.WxOpenID).
		Save(ctx)
	return err
}

func (r *UserRepository) FindByID(ctx context.Context, id uint) (*ent.User, error) {
	return r.client.User.Get(ctx, int(id))
}

func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*ent.User, error) {
	return r.client.User.Query().
		Where(user.UsernameEQ(username)).
		Only(ctx)
}

func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*ent.User, error) {
	return r.client.User.Query().
		Where(user.Email(email)).
		Only(ctx)
}

func (r *UserRepository) FindByWxOpenID(ctx context.Context, openID string) (*ent.User, error) {
	return r.client.User.Query().
		Where(user.WxOpenIDEQ(openID)).
		Only(ctx)
}

func (r *UserRepository) ExistsByUsername(ctx context.Context, username string) (bool, error) {
	return r.client.User.Query().
		Where(user.UsernameEQ(username)).
		Exist(ctx)
}

func (r *UserRepository) Update(ctx context.Context, id int, fields map[string]interface{}) (*ent.User, error) {
	update := r.client.User.UpdateOneID(id)

	for key, value := range fields {
		switch key {
		case "nickname":
			if v, ok := value.(string); ok {
				update.SetNickname(v)
			}
		case "email":
			if v, ok := value.(string); ok {
				update.SetEmail(v)
			}
		case "password":
			if v, ok := value.(string); ok {
				update.SetPassword(v)
			}
		case "wx_open_id":
			if v, ok := value.(string); ok {
				update.SetWxOpenID(v)
			}
		case "updated_at":
			if v, ok := value.(time.Time); ok {
				update.SetUpdatedAt(v)
			}
		}
	}

	return update.Save(ctx)
}
