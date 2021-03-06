package dao

import (
	"context"
	"sync"

	redisv8 "github.com/go-redis/redis/v8"

	"github.com/go-goim/core/pkg/consts"

	"github.com/go-goim/gateway/internal/app"
)

var (
	userDao     *UserDao
	userDaoOnce sync.Once
)

type UserDao struct {
	rdb *redisv8.Client
	// mysql.DB get from context, because we may need use transaction
}

func GetUserDao() *UserDao {
	userDaoOnce.Do(func() {
		userDao = &UserDao{
			rdb: app.GetApplication().Redis,
		}
	})
	return userDao
}

// GetUserOnlineAgent get user online agent from redis
func (u *UserDao) GetUserOnlineAgent(ctx context.Context, uid int64) (string, error) {
	key := consts.GetUserOnlineAgentKey(uid)
	val, err := u.rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redisv8.Nil {
			return "", nil
		}
		return "", err
	}

	return val, nil
}
