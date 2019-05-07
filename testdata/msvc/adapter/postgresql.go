
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
	
	vpr.SetDefault("postgres.username", "testing")
	vpr.SetDefault("postgres.password", "testing")
	vpr.SetDefault("postgres.database", "testing")
	vpr.SetDefault("postgres.host", "0.0.0.0")
	vpr.SetDefault("postgres.port", "5432")
	

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

