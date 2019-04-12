package tidbtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	//mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type tiDB struct {
	client *sql.DB
}

//DbConfig 测试数据库相关信息
type DbConfig struct {
	Host                   string
	Port                   int
	User, Password, DbName string
}

func newTiDB(cfg DbConfig) (*tiDB, error) {
	source := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password,
		cfg.Host, cfg.Port, cfg.DbName)
	db, err := sql.Open("mysql", source)
	if err != nil {
		return nil, fmt.Errorf("open db %s error %s", source, err)
	}
	return &tiDB{
		client: db,
	}, nil
}

func (db *tiDB) exec(ctx context.Context, sql string) error {
	if db == nil {
		return errors.New("db is nil")
	}
	_, err := db.client.ExecContext(ctx, sql)
	if err != nil {
		return fmt.Errorf("execute sql error %s", err)
	}
	return nil
}

func (db *tiDB) close() error {
	if db == nil {
		return errors.New("db is nil")
	}
	return db.client.Close()
}
