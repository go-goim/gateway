package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	responsepb "github.com/go-goim/api/transport/response"
	userv1 "github.com/go-goim/api/user/v1"
	"github.com/go-goim/core/pkg/log"
	"github.com/go-goim/core/pkg/util"

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

func (s *UserService) QueryUserInfo(ctx context.Context, req *userv1.QueryUserRequest) (*userv1.User, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	rsp, err := userv1.NewUserServiceClient(cc).QueryUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.GetResponse().Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetUser().ToUser(), nil
}

// Login check user login status and return user info
func (s *UserService) Login(ctx context.Context, req *userv1.UserLoginRequest) (*userv1.User, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	var queryReq = &userv1.QueryUserRequest{}
	switch {
	case req.GetEmail() != "":
		queryReq.User = &userv1.QueryUserRequest_Email{Email: req.GetEmail()}
	case req.GetPhone() != "":
		queryReq.User = &userv1.QueryUserRequest_Phone{Phone: req.GetPhone()}
	default:
		return nil, fmt.Errorf("invalid user login request")
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

	if !rsp.GetResponse().Success() {
		return nil, rsp.GetResponse()
	}

	user := rsp.GetUser()

	if user.GetPassword() != util.HashString(req.GetPassword()) {
		return nil, responsepb.Code_InvalidUsernameOrPassword.BaseResponse()
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

	return user.ToUser(), nil
}

// Register register user.
func (s *UserService) Register(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.User, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	// do check user exist and create.
	rsp, err := userv1.NewUserServiceClient(cc).CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.GetResponse().Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetUser().ToUser(), nil
}

// UpdateUser update user info.
func (s *UserService) UpdateUser(ctx context.Context, req *userv1.UpdateUserRequest) (*userv1.User, error) {
	cc, err := userServiceConnPool.Get()
	if err != nil {
		return nil, err
	}

	// do check user exist and update.
	rsp, err := userv1.NewUserServiceClient(cc).UpdateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	if !rsp.GetResponse().Success() {
		return nil, rsp.GetResponse()
	}

	return rsp.GetUser().ToUser(), nil
}
