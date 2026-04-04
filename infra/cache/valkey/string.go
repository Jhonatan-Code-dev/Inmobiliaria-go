// infra/cache/valkey/string.go
package valkey

import (
	"context"
	"time"
)

func (r *Repository) Set(
	ctx context.Context,
	key string,
	value string,
	ttl ...time.Duration,
) error {

	exp := r.defaultTTL
	if len(ttl) > 0 {
		exp = ttl[0]
	}

	return r.client.Set(ctx, key, value, exp).Err()
}

func (r *Repository) Get(
	ctx context.Context,
	key string,
) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *Repository) Delete(
	ctx context.Context,
	key string,
) error {
	return r.client.Del(ctx, key).Err()
}
