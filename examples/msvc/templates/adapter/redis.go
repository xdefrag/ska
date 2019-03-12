{{ if has "redis" .adapters }}
package adapter

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
	"github.com/spf13/viper"
)

type Redis struct {
	pool *redis.Pool
}

func NewRedis(vpr *viper.Viper) Redis {
	{{if .redis.defaults }}
	vpr.SetDefault("redis.host", "{{.redis.defaults.host}}")
	vpr.SetDefault("redis.port", "{{.redis.defaults.port}}")
	vpr.SetDefault("redis.password", "{{.redis.defaults.password}}")
	vpr.SetDefault("redis.database", "{{.redis.defaults.database}}")
	{{end}}

	addr := fmt.Sprintf("%s:%d", vpr.GetString("redis.host"), vpr.GetString("redis.port"))

	var opts []redis.DialOption

	if password := vpr.GetString("redis.password"); password != "" {
		opts = append(opts, redis.DialPassword(password))
	}

	if database := vpr.GetInt("redis.database"); database != 0 {
		opts = append(opts, redis.DialDatabase(database))
	}

	pool := &redis.Pool{
		Dial:            func() (redis.Conn, error) { return redis.Dial("tcp", addr, opts...) },
		MaxIdle:         vpr.GetInt("redis.pool.max.idle"),
		MaxActive:       vpr.GetInt("redis.pool.max.active"),
		IdleTimeout:     vpr.GetDuration("redis.pool.idle.timeout"),
		MaxConnLifetime: vpr.GetDuration("redis.pool.max.conn.lifetime"),
		Wait:            false,
	}

	return Redis{pool}
}
{{end}}
