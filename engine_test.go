package tidbtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	//初始化测试表
	db, err := newTiDB(testDbConfig)
	assert.Equal(t, nil, err)
	err = db.exec(context.Background(), `
	DROP TABLE IF EXISTS accounts;
	CREATE TABLE IF NOT EXISTS accounts(
		id int primary key AUTO_INCREMENT,
		balance int not null
	);
	INSERT INTO accounts (balance) VALUES (60);
	`)
	assert.Equal(t, nil, err)

	reader, err := NewFileReader([]string{"test1.sql", "test2.sql"})
	assert.Equal(t, nil, err)

	strategy := NewExhaustive()

	err = Run(context.Background(), testDbConfig, reader, strategy)
	assert.Equal(t, nil, err)

	//验证四条语句的执行结果
	var balance int
	err = db.client.QueryRow(`SELECT balance FROM accounts WHERE id=1`).Scan(&balance)
	assert.Equal(t, nil, err)
	//一共会执行6次，最后值应该为0
	assert.Equal(t, 0, balance)
}
