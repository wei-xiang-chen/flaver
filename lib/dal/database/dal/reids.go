package dal

import (
	"context"
	"time"
)

type IRedis interface {
	RedisIncr(key string, ttl *time.Duration) error
	RedisDecr(key string, ttl *time.Duration) error

	WaitRedisLockRelease(key string) error
	RedisGet(key string) (string, error)

	RedisLock(key string, ttl time.Duration) (bool, error)
	ReleaseRedisLock(key string) error
	RedisSet(key string, val interface{}, ttl time.Duration) error
	RedisDelByPattern(keyPatterns []string) error
	RedisScriptKill() error
}

func (this *Dal) RedisIncr(key string, ttl *time.Duration) error {
	pipe := this.redis.TxPipeline()
	pipe.Incr(key)
	if ttl != nil {
		pipe.Expire(key, *ttl)
	}

	_, err := pipe.Exec()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Dal) RedisDecr(key string, ttl *time.Duration) error {
	pipe := this.redis.TxPipeline()
	pipe.Decr(key)
	if ttl != nil {
		pipe.Expire(key, *ttl)
	}

	_, err := pipe.Exec()
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (this *Dal) WaitRedisLockRelease(key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return this.waitRedisLockRelease(ctx, key)
}

func (this *Dal) waitRedisLockRelease(ctx context.Context, key string) error {

	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			if exist, err := this.redis.Exists(key).Result(); err != nil {
				return err
			} else if exist <= 0 {
				return nil
			}
			// 降低打redis次數 怕把redis打爆
			time.Sleep(time.Millisecond * 100)
		}
	}
}

func (this *Dal) RedisGet(key string) (string, error) {
	return this.redis.Get(key).Result()
}

func (this *Dal) RedisLock(key string, ttl time.Duration) (bool, error) {
	result, err := this.redis.SetNX(key, "", ttl).Result()

	if err != nil {
		return false, nil
	}

	return result, nil
}

func (this *Dal) ReleaseRedisLock(key string) error {
	return this.redis.Del(key).Err()
}

func (this *Dal) RedisSet(key string, val interface{}, ttl time.Duration) error {
	return this.redis.Set(key, val, ttl).Err()
}

func (this *Dal) RedisScriptKill() error {
	return this.redis.ScriptKill().Err()
}

func (this *Dal) RedisDelByPattern(keyPatterns []string) error {

	// ref: http://events.jianshu.io/p/5a95a8209e5b
	script := `
	local count = 0;
	local cursor= 0;
	repeat
		local scanResult = redis.call("SCAN", cursor, "MATCH", KEYS[1], "COUNT", 100);
		local keys = scanResult[2];
		if(scanResult ~= nil and #scanResult > 0)then
            count = count + #keys;
			cursor = tonumber(scanResult[1]);
			for i = 1, #keys do
				redis.replicate_commands()
				local key = keys[i];
				redis.call("UNLINK", key);
			end;
		end;
	until (cursor <= 0);

	return count;
	`

	for _, keyPattern := range keyPatterns {
		if err := this.redis.Eval(script, []string{keyPattern}).Err(); err != nil {
			return err
		}
	}

	return nil
}
