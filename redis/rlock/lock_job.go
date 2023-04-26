package rlock

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
	"github.com/xhigher/hzgo/env"
	"github.com/xhigher/hzgo/logger"
	"github.com/xhigher/hzgo/utils"
)

type LockJobOption struct {
	RedisErrorRetryInterval time.Duration
	ReLockInterval          time.Duration
	LockHoldInterval        time.Duration
	LockRefreshInterval     time.Duration
	TaskRunInterval         time.Duration
}

func NewLockJob(client *redis.Client, key string, call func(), option *LockJobOption) *LockJob {
	job := &LockJob{
		client: client,
		key:    key,
		call:   call,

		text:                    "",
		ctx:                     nil,
		cancel:                  nil,
		redisErrorRetryInterval: 0,
		reLockInterval:          0,
		lockHoldInterval:        0,
		lockRefreshInterval:     0,
		taskRunInterval:         0,
	}
	if option == nil {
		option = &LockJobOption{}
	}
	getDefault := func(d1, d2 time.Duration) time.Duration {
		if d1 != 0 {
			return d1
		} else {
			return d2
		}
	}
	job.redisErrorRetryInterval = getDefault(option.RedisErrorRetryInterval, time.Second*30)
	job.reLockInterval = getDefault(option.ReLockInterval, time.Second*30)
	job.lockHoldInterval = getDefault(option.LockHoldInterval, time.Second*120)
	job.lockRefreshInterval = getDefault(option.LockRefreshInterval, time.Second*20)
	job.taskRunInterval = getDefault(option.TaskRunInterval, time.Second*10)
	return job
}

type LockJob struct {
	client *redis.Client
	key    string
	call   func()

	text   string
	ctx    context.Context
	cancel context.CancelFunc

	redisErrorRetryInterval time.Duration
	reLockInterval          time.Duration
	lockHoldInterval        time.Duration
	lockRefreshInterval     time.Duration
	taskRunInterval         time.Duration
}

func (l *LockJob) Start() {
INIT:
	l.ctx, l.cancel = context.WithCancel(context.Background())
	l.text = fmt.Sprintf("%d:%s", utils.NowTimeNano(), env.GetHostName())
	for {
		res, err := l.client.SetNX(l.ctx, l.key, l.text, l.lockHoldInterval).Result()
		if err != nil {
			logger.Errorf("error call set nx %v", err)
			time.Sleep(l.redisErrorRetryInterval)
			continue
		}
		if !res {
			time.Sleep(l.reLockInterval + time.Millisecond*time.Duration(utils.RandInt64(0, 1000)))
			continue
		}
		go l.holdKey()
		for {
			select {
			case <-l.ctx.Done():
				goto INIT
			default:
			}
			l.call()
			time.Sleep(l.taskRunInterval)
		}
	}
}

func (l *LockJob) holdKey() {
	defer func() {
		l.cancel()
	}()

	var errCount int
	for {
		res, err := l.client.Get(context.Background(), l.key).Result()
		if err == redis.Nil {
			return
		} else if err != nil {
			logger.Errorf("error hold key %v", err)
			time.Sleep(l.redisErrorRetryInterval)
			errCount++
			if errCount > 5 {
				logger.Infof("errorCount %d return", errCount)
				return
			}
			continue
		}
		if res != l.text {
			return
		}
		err = l.client.Expire(l.ctx, l.key, l.lockHoldInterval).Err()
		if err != nil {
			logger.Errorf("error call expire %v", err)
			time.Sleep(l.redisErrorRetryInterval)
			errCount++
			if errCount > 5 {
				logger.Infof("errorCount %d return", errCount)
				return
			}
			continue
		}
		errCount = 0
		time.Sleep(l.lockRefreshInterval)
	}
}

// 针对 key 加锁，然后执行 f， 保持锁， 并间隔 interval 执行一次
func LockDo(client *redis.Client, key string, interval time.Duration, f func(input context.Context) error) {
	var errCount int
	for {
		errCount = 0
		lock, err := NewRetryLock(client, key, time.Minute, time.Second*30, 2)
		if err == ErrNotObtained {
			logger.Errorf("not obtained, continue")
			continue
		} else if err != nil {
			logger.Errorf("error get lock %v", err)
			time.Sleep(time.Minute)
			continue
		}

		// start
		ctx, cancel := context.WithCancel(context.Background())
		go refreshLock(lock, ctx, cancel)
		for errCount < 5 {
			err = f(ctx)
			if err != nil {
				errCount++
				logger.Errorf("error %v", err)
				time.Sleep(time.Second)
				continue
			}
			errCount = 0
			time.Sleep(interval)
		}
		logger.Errorf("error count >= 5, continue")
		cancel()
	}
}

func refreshLock(lock *RetryLock, ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTimer(time.Second * 20)
	defer func() {
		ticker.Stop()
		if cancel != nil {
			cancel()
		}
	}()
	var errCount int
	for {
		select {
		case <-ticker.C:
		case <-ctx.Done():
			logger.Infof("already done")
			return
		}

		own, err := lock.IsOwn()
		if err == ErrNotObtained {
			logger.Errorf("own not obtained return")
			return
		} else if err != nil {
			logger.Errorf("error check own %v", err)
			if errCount > 10 {
				logger.Errorf("error count > 10 return")
				return
			}
			errCount++
			time.Sleep(time.Second)
			continue
		}

		if !own {
			logger.Errorf("not own return")
			return
		}

		err = lock.Refresh(time.Minute)
		if err != nil {
			logger.Errorf("error refresh lock %v", err)
			if errCount > 10 {
				logger.Errorf("error count > 10 return")
				return
			}
			errCount++
			time.Sleep(time.Second)
			continue
		}
		ticker.Reset(time.Minute)
	}
}
