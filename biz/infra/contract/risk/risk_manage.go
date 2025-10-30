package risk

import (
	"context"
	"errors"
	"strconv"

	"github.com/xh-polaris/synapse/biz/infra/contract/cache"
)

var rm *RiskManager

type RiskManager struct {
	cache cache.Cmdable
}

func New(cli cache.Cmdable) *RiskManager {
	rm = &RiskManager{cache: cli}
	return rm
}

// checkUpperLimit 检验这个key是否得到指定的上限返回是否达到以及实际的值
func (r *RiskManager) checkUpperLimit(ctx context.Context, key string, upper int) (bool, int, error) {
	get := r.cache.Get(ctx, key)
	result, err := get.Val(), get.Err()
	if errors.Is(err, cache.Nil) {
		return 0 >= upper, 0, nil
	} else if err != nil {
		return false, 0, err
	}
	val, err := strconv.Atoi(result)
	if err != nil {
		return false, 0, err
	}
	return val >= upper, val, nil
}

// addOnce 增加一次风险行为
// 如果不存在, 则创建一个key, 设置过期时间为 expire 秒
// 如果存在, 则增加1, 不改变过期时间
func (r *RiskManager) addOnce(ctx context.Context, key string, expire int) error {
	script := `
        local key = KEYS[1]
        local expire = tonumber(ARGV[1])
        local current = redis.call('INCR', key)
        if current == 1 then
            redis.call('EXPIRE', key, expire)
        end
        return current  -- 或者直接返回 OK
    `
	result := r.cache.Eval(ctx, script, []string{key}, expire)
	return result.Err()
}

func CheckUpperLimit(ctx context.Context, key string, upper int) (bool, int, error) {
	return rm.checkUpperLimit(ctx, key, upper)
}

func AddOnce(ctx context.Context, key string, expire int) error {
	return rm.addOnce(ctx, key, expire)
}
