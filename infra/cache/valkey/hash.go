package valkey

import (
	"context"
	"time"
)

func (r *Repository) HSet(
	ctx context.Context,
	key string,
	values map[string]string,
	ttl ...time.Duration,
) error {

	exp := r.defaultTTL
	if len(ttl) > 0 {
		exp = ttl[0]
	}

	pipe := r.client.Pipeline()
	pipe.HSet(ctx, key, values)
	pipe.Expire(ctx, key, exp)

	_, err := pipe.Exec(ctx)
	return err
}

func (r *Repository) HGetAll(
	ctx context.Context,
	key string,
) (map[string]string, error) {
	return r.client.HGetAll(ctx, key).Result()
}
