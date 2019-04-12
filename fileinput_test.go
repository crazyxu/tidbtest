package tidbtest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileInput(t *testing.T) {
	r, err := NewFileReader([]string{"test1.sql"})
	assert.Equal(t, nil, err)

	ss, err := r.Read()
	assert.Equal(t, nil, err)

	assert.Equal(t, map[string]string{
		"test1.sql": `UPDATE accounts SET balance=balance-1 WHERE id=1;
UPDATE accounts SET balance=balance-2 WHERE id=1;`,
	}, ss)
}

func TestFileInput_Error(t *testing.T) {
	_, err := NewFileReader([]string{"test_not_found.sql"})
	assert.NotEqual(t, nil, err)
}
