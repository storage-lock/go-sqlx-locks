package sqlx_locks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobalSqlxLockFactory_InitiallyNil(t *testing.T) {
	assert.Nil(t, GlobalSqlxLockFactory)
}
