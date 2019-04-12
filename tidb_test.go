package tidbtest

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testDbConfig = DbConfig{
	Host:     "127.0.0.1",
	Port:     4000,
	User:     "root",
	Password: "",
	DbName:   "test_db",
}

func TestTiDB(t *testing.T) {
	db, err := newTiDB(testDbConfig)
	assert.Equal(t, nil, err)
	defer db.close()

	err = db.exec(context.Background(), "show tables;")
	assert.Equal(t, nil, err)
}
