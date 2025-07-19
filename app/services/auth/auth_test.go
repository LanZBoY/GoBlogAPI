package auth

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	UserModel "wentee/blog/app/model/mongodb/user"
	"wentee/blog/app/schema/auth"
	"wentee/blog/app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mockUserRepo implements userRepo interface
type mockUserRepo struct{ mock.Mock }

func (m *mockUserRepo) GetUserByEmail(ctx context.Context, email string, opts ...*options.FindOneOptions) (*UserModel.UserDocument, error) {
	args := m.Called(testutils.AppendCallArgs([]any{ctx, email}, opts)...)
	return args.Get(0).(*UserModel.UserDocument), args.Error(1)
}

func TestTryLogin(t *testing.T) {
	repo := new(mockUserRepo)
	pass := new(testutils.MockPasswordUtils)
	svc := &AuthService{userRepo: repo, passwordUtils: pass}

	user := &UserModel.UserDocument{Id: primitive.NewObjectID(), Email: "a@b.c", Username: "u", Password: "hash", Salt: "s"}
	repo.On("GetUserByEmail", mock.Anything, "a@b.c").Return(user, nil)
	pass.On("VerifyPassword", user.Password, "pw", user.Salt).Return(true)

	token, err := svc.TryLogin(context.TODO(), &auth.LoginInfo{Email: "a@b.c", Password: "pw"})
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestTryLogin_WrongPassword(t *testing.T) {
	repo := new(mockUserRepo)
	pass := new(testutils.MockPasswordUtils)
	svc := &AuthService{userRepo: repo, passwordUtils: pass}

	user := &UserModel.UserDocument{Id: primitive.NewObjectID(), Password: "hash", Salt: "s"}
	repo.On("GetUserByEmail", mock.Anything, "a@b.c").Return(user, nil)
	pass.On("VerifyPassword", user.Password, "pw", user.Salt).Return(false)

	_, err := svc.TryLogin(context.TODO(), &auth.LoginInfo{Email: "a@b.c", Password: "pw"})
	assert.Error(t, err)
}
