{{if has "postgresql" .adapters}}
package adapter

import (
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

type Postgresql struct {
	db *sqlx.DB
	qb sq.StatementBuilderType
}

func NewPostgresql(vpr *viper.Viper) Postgresql {
	{{if .postgresql.defaults}}
	vpr.SetDefault("postgres.username", "{{.postgresql.defaults.username}}")
	vpr.SetDefault("postgres.password", "{{.postgresql.defaults.password}}")
	vpr.SetDefault("postgres.database", "{{.postgresql.defaults.database}}")
	vpr.SetDefault("postgres.host", "{{.postgresql.defaults.host}}")
	vpr.SetDefault("postgres.port", "{{.postgresql.defaults.port}}")
	{{end}}

	vpr.SetDefault("postgres.maxidle", 20)
	vpr.SetDefault("postgres.maxopen", 200)

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		vpr.GetString("postgres.username"),
		vpr.GetString("postgres.password"),
		vpr.GetString("postgres.host"),
		vpr.GetString("postgres.port"),
		vpr.GetString("postgres.database"),
		strings.Join(vpr.GetStringSlice("postgres.params"), "&"),
	)

	db := sqlx.MustConnect("postgres", dsn)

	db.DB.SetMaxIdleConns(vpr.GetInt("postgres.maxidle"))
	db.DB.SetMaxOpenConns(vpr.GetInt("postgres.maxopen"))

	qb := sq.StatementBuilder.PlaceholderFormat(sq.Dollar).RunWith(sq.NewStmtCacher(db))
	return Postgresql{db: db, qb: qb}
}
{{end}}
