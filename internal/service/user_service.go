package service

import (
	"context"
	"sync"
	"time"

	"github.com/go-goim/api/errors"
	userv1 "github.com/go-goim/api/user/v1"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/util"
	"github.com/go-goim/gateway/internal/dto"

	"github.com/go-goim/gateway/internal/dao"
)

type UserService struct {
	userDao *dao.UserDao
}

var (
	userService     *UserService
	userServiceOnce sync.Once
)

func GetUserService() *UserService {
	userServiceOnce.Do(func() {
		userService = &UserService{
			userDao: dao.GetUserDao(),
		}
	})
	return userService
}

func (s *UserService) QueryUserInfo(ctx context.Context, req *dto.QueryUserRequest) (*dto.User, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := userv1.NewUserServiceClient(cc).QueryUser(ctx, req.ToPb())
	if err != nil {
		return nil, err
	}

	if err := rsp.GetError().Err(); err != nil {
		return nil, err
	}

	return dto.UserFromPb(rsp.GetUser()), nil
}

// Login check user login status and return user info
func (s *UserService) Login(ctx context.Context, req *dto.UserLoginRequest) (*dto.User, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	var queryReq = &userv1.QueryUserRequest{}
	switch {
	case req.Email != nil:
		queryReq.Field = &userv1.QueryUserRequest_Email{Email: *req.Email}
	case req.Phone != nil:
		queryReq.Field = &userv1.QueryUserRequest_Phone{Phone: *req.Phone}
	default:
		return nil, errors.ErrorCode_InvalidParams.WithMessage("invalid user login request")
	}

	ddl, ok := ctx.Deadline()
	if ok {
		log.Debug("Login ctx deadline", "ddl", ddl)
	}

	ctx2, cancel := context.WithTimeout(ctx, 3*time.Second)
	rsp, err := userv1.NewUserServiceClient(cc).QueryUser(ctx2, queryReq)
	cancel()

	if err != nil {
		return nil, err
	}

	if err := rsp.GetError().Err(); err != nil {
		return nil, err
	}

	user := rsp.GetUser()

	if user.GetPassword() != util.HashString(req.Password) {
		return nil, errors.ErrorCode_InvalidUsernameOrPassword
	}

	agentIP, err := s.userDao.GetUserOnlineAgent(ctx, user.GetUid())
	if err != nil {
		return nil, err
	}

	if len(agentIP) == 0 {
		// not login
		user.LoginStatus = userv1.LoginStatus_LOGIN
		agentIP, err = LoadMatchedPushServer(ctx)
		if err != nil {
			return nil, err
		}
		user.PushServerIp = &agentIP
	} else {
		// already login
		user.LoginStatus = userv1.LoginStatus_ALREADY_LOGIN
		user.PushServerIp = &agentIP
	}

	return dto.UserFromPb(user), nil
}

// Register register user.
func (s *UserService) Register(ctx context.Context, req *dto.CreateUserRequest) (*dto.User, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	// do check user exist and create.
	rsp, err := userv1.NewUserServiceClient(cc).CreateUser(ctx, req.ToPb())
	if err != nil {
		return nil, err
	}

	if err := rsp.GetError().Err(); err != nil {
		return nil, err
	}

	return dto.UserFromPb(rsp.GetUser()), nil
}

// UpdateUser update user info.
func (s *UserService) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*dto.User, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	// do check user exist and update.
	rsp, err := userv1.NewUserServiceClient(cc).UpdateUser(ctx, req.ToPb())
	if err != nil {
		return nil, err
	}

	if err := rsp.GetError().Err(); err != nil {
		return nil, err
	}

	return dto.UserFromPb(rsp.GetUser()), nil
}
