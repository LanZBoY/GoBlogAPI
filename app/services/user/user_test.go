package user_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	UserModel "wentee/blog/app/model/mongodb/user"
	"wentee/blog/app/schema/apperror"
	"wentee/blog/app/schema/apperror/errcode"
	"wentee/blog/app/schema/basemodel"
	UserSchema "wentee/blog/app/schema/user"
	UserSvc "wentee/blog/app/services/user"
	"wentee/blog/app/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MockUserRepo struct {
	mock.Mock
}

// CountUsers implements user.IUserRepository.
func (m *MockUserRepo) CountUsers(ctx context.Context) (int64, error) {
	args := m.Called(ctx)

	return args.Get(0).(int64), args.Error(1)
}

// CreateUser implements user.IUserRepository.
func (m *MockUserRepo) CreateUser(ctx context.Context, userCreate *UserModel.UserDocument) error {
	args := m.Called(ctx, userCreate)

	return args.Error(0)
}

// DeleteUserById implements user.IUserRepository.
func (m *MockUserRepo) DeleteUserById(ctx context.Context, id primitive.ObjectID, opts ...*options.DeleteOptions) error {
	args := m.Called(ctx, id, opts)
	return args.Error(0)
}

// GetUserByEmail implements user.IUserRepository.
func (m *MockUserRepo) GetUserByEmail(ctx context.Context, email string, opts ...*options.FindOneOptions) (*UserModel.UserDocument, error) {

	args := m.Called(testutils.AppendCallArgs([]any{ctx, email}, opts)...)

	return args.Get(0).(*UserModel.UserDocument), args.Error(1)
}

// GetUserById implements user.IUserRepository.
func (m *MockUserRepo) GetUserById(ctx context.Context, id primitive.ObjectID, opts ...*options.FindOneOptions) (*UserModel.UserDocument, error) {
	args := m.Called(testutils.AppendCallArgs([]any{ctx, id}, opts)...)
	return args.Get(0).(*UserModel.UserDocument), args.Error(1)
}

// QueryUsers implements user.IUserRepository.
func (m *MockUserRepo) QueryUsers(ctx context.Context, query *basemodel.BaseQuery) ([]UserModel.UserDocument, error) {
	args := m.Called(ctx, query)
	return args.Get(0).([]UserModel.UserDocument), args.Error(1)
}

// UpdateUserById implements user.IUserRepository.
func (m *MockUserRepo) UpdateUserById(ctx context.Context, id primitive.ObjectID, updateData *UserSchema.UserUpdate, opts ...*options.UpdateOptions) error {
	args := m.Called(testutils.AppendCallArgs([]any{ctx, id, updateData}, opts)...)
	return args.Error(0)
}

func TestNewUserService(t *testing.T) {
	mockRepo := new(MockUserRepo)
	mockPasswordUtils := new(testutils.MockPasswordUtils)
	mockObjectIdCreator := new(testutils.MockObjectIDCreator)
	userSvc := UserSvc.NewUserService(mockRepo, mockPasswordUtils, mockObjectIdCreator)
	assert.NotNil(t, userSvc)
}

func TestCountUsers(t *testing.T) {
	mockRepo := new(MockUserRepo)
	ctx := context.TODO()
	mockRepo.On("CountUsers", ctx).Return(int64(10), nil)
	mockPasswordUtils := new(testutils.MockPasswordUtils)
	mockObjectIdCreator := new(testutils.MockObjectIDCreator)
	userSvc := UserSvc.NewUserService(mockRepo, mockPasswordUtils, mockObjectIdCreator)

	total, _ := userSvc.CountUsers(ctx)

	assert.Equal(t, int64(10), total)
}

func TestRegistryUser(t *testing.T) {

	tests := []struct {
		name       string
		ctx        context.Context
		createUser *UserSchema.UserCreate
		mockSetup  func(ctx context.Context, createUser *UserSchema.UserCreate, mockRepo *MockUserRepo, mockPasswordUtils *testutils.MockPasswordUtils)
		wantErr    bool
	}{
		{
			name: "Normal",
			ctx:  context.TODO(),
			createUser: &UserSchema.UserCreate{
				Email:    "superjimalex@gmail.com.tw",
				Username: "WenTee",
				Password: "WenTee",
			},
			mockSetup: func(ctx context.Context, createUser *UserSchema.UserCreate, mockRepo *MockUserRepo, mockPasswordUtils *testutils.MockPasswordUtils) {
				mockRepo.On("GetUserByEmail", ctx, createUser.Email).Return((*UserModel.UserDocument)(nil), nil)
				mockPasswordUtils.On("GenerateSalt", 32).Return("FakeSalt", nil)
				mockPasswordUtils.On("HashPassword", mock.Anything, "FakeSalt").Return("HashedPassword", nil)
				mockRepo.On("CreateUser", ctx, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "FoundError",
			ctx:  context.TODO(),
			createUser: &UserSchema.UserCreate{
				Email:    "superjimalex@gmail.com.tw",
				Username: "WenTee",
				Password: "WenTee",
			},
			mockSetup: func(ctx context.Context, createUser *UserSchema.UserCreate, mockRepo *MockUserRepo, mockPasswordUtils *testutils.MockPasswordUtils) {
				mockRepo.On("GetUserByEmail", ctx, createUser.Email).Return(&UserModel.UserDocument{}, nil)
				// mockPasswordUtils.On("GenerateSalt", 32).Return("FakeSalt", nil)
				// mockPasswordUtils.On("HashPassword", mock.Anything, "FakeSalt").Return("HashedPassword", nil)
				// mockRepo.On("CreateUser", ctx, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "Salt Error",
			ctx:  context.TODO(),
			createUser: &UserSchema.UserCreate{
				Email:    "superjimalex@gmail.com.tw",
				Username: "WenTee",
				Password: "WenTee",
			},
			mockSetup: func(ctx context.Context, createUser *UserSchema.UserCreate, mockRepo *MockUserRepo, mockPasswordUtils *testutils.MockPasswordUtils) {
				mockRepo.On("GetUserByEmail", ctx, createUser.Email).Return((*UserModel.UserDocument)(nil), nil)
				mockPasswordUtils.On("GenerateSalt", 32).Return("", errors.New("Fail to Generate Salt"))
				// mockPasswordUtils.On("HashPassword", mock.Anything, "FakeSalt").Return("HashedPassword", nil)
				// mockRepo.On("CreateUser", ctx, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "Hash Fail Error",
			ctx:  context.TODO(),
			createUser: &UserSchema.UserCreate{
				Email:    "superjimalex@gmail.com.tw",
				Username: "WenTee",
				Password: "WenTee",
			},
			mockSetup: func(ctx context.Context, createUser *UserSchema.UserCreate, mockRepo *MockUserRepo, mockPasswordUtils *testutils.MockPasswordUtils) {
				mockRepo.On("GetUserByEmail", ctx, createUser.Email).Return((*UserModel.UserDocument)(nil), nil)
				mockPasswordUtils.On("GenerateSalt", 32).Return("FakeSalt", nil)
				mockPasswordUtils.On("HashPassword", mock.Anything, "FakeSalt").Return("", errors.New("Hash Fail"))
				// mockRepo.On("CreateUser", ctx, mock.Anything).Return(nil)
			},
			wantErr: true,
		},
		{
			name: "Create User Fail",
			ctx:  context.TODO(),
			createUser: &UserSchema.UserCreate{
				Email:    "superjimalex@gmail.com.tw",
				Username: "WenTee",
				Password: "WenTee",
			},
			mockSetup: func(ctx context.Context, createUser *UserSchema.UserCreate, mockRepo *MockUserRepo, mockPasswordUtils *testutils.MockPasswordUtils) {
				mockRepo.On("GetUserByEmail", ctx, createUser.Email).Return((*UserModel.UserDocument)(nil), nil)
				mockPasswordUtils.On("GenerateSalt", 32).Return("FakeSalt", nil)
				mockPasswordUtils.On("HashPassword", mock.Anything, "FakeSalt").Return("HashedPassword", nil)
				mockRepo.On("CreateUser", ctx, mock.Anything).Return(errors.New("Create User Fail"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockPasswordUtils := new(testutils.MockPasswordUtils)
			mockObjectIdCreator := new(testutils.MockObjectIDCreator)
			userSvc := UserSvc.NewUserService(mockRepo, mockPasswordUtils, mockObjectIdCreator)
			tt.mockSetup(tt.ctx, tt.createUser, mockRepo, mockPasswordUtils)

			err := userSvc.RegistryUser(tt.ctx, tt.createUser)

			if tt.wantErr {
				appErr := new(apperror.AppError)
				if errors.As(err, appErr) {
					assert.Equal(t, http.StatusBadRequest, appErr.Status)
					assert.Equal(t, errcode.USER_EXIST, appErr.Code)
				} else {
					assert.Error(t, err)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetUserById(t *testing.T) {
	oid := primitive.NewObjectID()
	tests := []struct {
		name      string
		ctx       context.Context
		id        string
		mockSetup func(mockRepo *MockUserRepo, mockObjIdCreator *testutils.MockObjectIDCreator, ctx context.Context, id string)
		wantRet   *UserSchema.UserInfo
		wantErr   error
	}{
		{
			name: "Normal",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID().Hex(),
			mockSetup: func(mockRepo *MockUserRepo, mockObjIdCreator *testutils.MockObjectIDCreator, ctx context.Context, id string) {
				mockObjIdCreator.On("ObjectIDFromHex", id).Return(oid, nil)
				mockRepo.On("GetUserById", ctx, oid, mock.Anything).Return(&UserModel.UserDocument{Id: oid, Email: "XXX@gmail.com", Username: "XX_Man"}, nil)
			},
			wantRet: &UserSchema.UserInfo{Id: oid, Email: "XXX@gmail.com", Username: "XX_Man"},
			wantErr: nil,
		},
		{
			name: "Bad Id String",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID().Hex(),
			mockSetup: func(mockRepo *MockUserRepo, mockObjIdCreator *testutils.MockObjectIDCreator, ctx context.Context, id string) {
				mockObjIdCreator.On("ObjectIDFromHex", id).Return(primitive.NilObjectID, errors.New(""))
				// mockRepo.On("GetUserById", ctx, oid, mock.Anything).Return(&UserModel.UserDocument{Id: oid, Email: "XXX@gmail.com", Username: "XX_Man"}, nil)
			},
			wantRet: nil,
			wantErr: errors.New(""),
		},
		{
			name: "Bad Get User From Repo",
			ctx:  context.TODO(),
			id:   primitive.NewObjectID().Hex(),
			mockSetup: func(mockRepo *MockUserRepo, mockObjIdCreator *testutils.MockObjectIDCreator, ctx context.Context, id string) {
				mockObjIdCreator.On("ObjectIDFromHex", id).Return(oid, nil)
				mockRepo.On("GetUserById", ctx, oid, mock.Anything).Return((*UserModel.UserDocument)(nil), errors.New(""))
			},
			wantRet: nil,
			wantErr: errors.New(""),
		},
	}

	for _, tt := range tests {
		mockRepo := new(MockUserRepo)
		mockPasswordUtils := new(testutils.MockPasswordUtils)
		mockObjIdCreator := new(testutils.MockObjectIDCreator)
		tt.mockSetup(mockRepo, mockObjIdCreator, tt.ctx, tt.id)
		userSvc := UserSvc.NewUserService(mockRepo, mockPasswordUtils, mockObjIdCreator)

		userInfo, err := userSvc.GetUserById(tt.ctx, tt.id)
		if tt.wantRet != nil {
			assert.Equal(t, tt.wantRet.Id, userInfo.Id)
			assert.Equal(t, tt.wantRet.Email, userInfo.Email)
			assert.Equal(t, tt.wantRet.Username, userInfo.Username)
		}

		if tt.wantErr != nil {
			assert.Error(t, err)
		}
	}
}

func TestListUsers(t *testing.T) {
	query := &basemodel.BaseQuery{Skip: 0, Limit: 10}
	oid := primitive.NewObjectID()
	tests := []struct {
		name      string
		repoUsers []UserModel.UserDocument
		repoErr   error
		wantErr   bool
	}{
		{
			name:      "Normal",
			repoUsers: []UserModel.UserDocument{{Id: oid, Email: "e", Username: "u"}},
			repoErr:   nil,
			wantErr:   false,
		},
		{
			name:      "RepoError",
			repoUsers: nil,
			repoErr:   errors.New("err"),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepo)
			mockPasswordUtils := new(testutils.MockPasswordUtils)
			mockObjIdCreator := new(testutils.MockObjectIDCreator)
			mockRepo.On("QueryUsers", mock.Anything, query).Return(tt.repoUsers, tt.repoErr)

			svc := UserSvc.NewUserService(mockRepo, mockPasswordUtils, mockObjIdCreator)
			users, err := svc.ListUsers(context.TODO(), query)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, users, len(tt.repoUsers))
				if len(tt.repoUsers) > 0 {
					assert.Equal(t, tt.repoUsers[0].Id, users[0].Id)
				}
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestUpdateUserById(t *testing.T) {
	oid := primitive.NewObjectID()
	tests := []struct {
		name      string
		id        string
		mockSetup func(*MockUserRepo, *testutils.MockObjectIDCreator, context.Context, primitive.ObjectID)
		wantErr   bool
	}{
		{
			name: "Normal",
			id:   primitive.NewObjectID().Hex(),
			mockSetup: func(repo *MockUserRepo, oc *testutils.MockObjectIDCreator, ctx context.Context, id primitive.ObjectID) {
				oc.On("ObjectIDFromHex", mock.Anything).Return(id, nil)
				repo.On("UpdateUserById", ctx, id, mock.Anything, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "BadId",
			id:   "bad",
			mockSetup: func(repo *MockUserRepo, oc *testutils.MockObjectIDCreator, ctx context.Context, id primitive.ObjectID) {
				oc.On("ObjectIDFromHex", mock.Anything).Return(primitive.NilObjectID, errors.New(""))
			},
			wantErr: true,
		},
		{
			name: "RepoError",
			id:   primitive.NewObjectID().Hex(),
			mockSetup: func(repo *MockUserRepo, oc *testutils.MockObjectIDCreator, ctx context.Context, id primitive.ObjectID) {
				oc.On("ObjectIDFromHex", mock.Anything).Return(id, nil)
				repo.On("UpdateUserById", ctx, id, mock.Anything, mock.Anything).Return(errors.New("repo"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepo)
			oc := new(testutils.MockObjectIDCreator)
			pass := new(testutils.MockPasswordUtils)
			tt.mockSetup(repo, oc, context.TODO(), oid)
			svc := UserSvc.NewUserService(repo, pass, oc)

			err := svc.UpdateUserById(context.TODO(), tt.id, &UserSchema.UserUpdate{})
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
			oc.AssertExpectations(t)
		})
	}
}

func TestDeleteUserById(t *testing.T) {
	oid := primitive.NewObjectID()
	tests := []struct {
		name      string
		id        string
		mockSetup func(*MockUserRepo, *testutils.MockObjectIDCreator, context.Context, primitive.ObjectID)
		wantErr   bool
	}{
		{
			name: "Normal",
			id:   primitive.NewObjectID().Hex(),
			mockSetup: func(repo *MockUserRepo, oc *testutils.MockObjectIDCreator, ctx context.Context, id primitive.ObjectID) {
				oc.On("ObjectIDFromHex", mock.Anything).Return(id, nil)
				repo.On("DeleteUserById", ctx, id, mock.Anything).Return(nil)
			},
			wantErr: false,
		},
		{
			name: "BadId",
			id:   "bad",
			mockSetup: func(repo *MockUserRepo, oc *testutils.MockObjectIDCreator, ctx context.Context, id primitive.ObjectID) {
				oc.On("ObjectIDFromHex", mock.Anything).Return(primitive.NilObjectID, errors.New(""))
			},
			wantErr: true,
		},
		{
			name: "RepoError",
			id:   primitive.NewObjectID().Hex(),
			mockSetup: func(repo *MockUserRepo, oc *testutils.MockObjectIDCreator, ctx context.Context, id primitive.ObjectID) {
				oc.On("ObjectIDFromHex", mock.Anything).Return(id, nil)
				repo.On("DeleteUserById", ctx, id, mock.Anything).Return(errors.New("repo"))
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := new(MockUserRepo)
			oc := new(testutils.MockObjectIDCreator)
			pass := new(testutils.MockPasswordUtils)
			tt.mockSetup(repo, oc, context.TODO(), oid)
			svc := UserSvc.NewUserService(repo, pass, oc)

			err := svc.DeleteUserById(context.TODO(), tt.id)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
			repo.AssertExpectations(t)
			oc.AssertExpectations(t)
		})
	}
}
