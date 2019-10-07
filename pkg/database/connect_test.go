package database_test

import (
	"fmt"
	"testing"

	"github.com/inclus/pg-init/pkg/database"
	"github.com/inclus/pg-init/pkg/database/databasefakes"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestSetupPostgisExtensionPropagatesAnyErrors(t *testing.T) {
	handle := &databasefakes.FakeHandle{}
	handle.ExecReturns(nil, fmt.Errorf("something failed"))
	err := database.SetupPostgisExtension(handle)
	assert.EqualError(t, err, "something failed")
}

func TestSetupPostgisExtensionPropagatesExecutesCorrectStatement(t *testing.T) {
	handle := &databasefakes.FakeHandle{}
	handle.ExecReturns(nil, nil)
	err := database.SetupPostgisExtension(handle)
	stmt, _ := handle.ExecArgsForCall(0)
	assert.Equal(t, "CREATE EXTENSION IF NOT EXISTS postgis", stmt)
	assert.NoError(t, err)
}

func TestBuildConnectionString(t *testing.T) {
	viper.Set("db.host", "localhost")
	viper.Set("db.port", "5432")
	viper.Set("db.database", "postgres")
	viper.Set("db.user", "postgres")
	viper.Set("db.password", "secret")
	viper.Set("db.extra", "sslmode=disable")
	result := database.BuildConnectionString()
	assert.Equal(t, "postgres://postgres:secret@localhost:5432/postgres?sslmode=disable", result)
}

func TestWaitForWorkingConnection(t *testing.T) {
	handle := &databasefakes.FakeHandle{}
	viper.Set("retry.attempts", 3)
	handle.QueryReturnsOnCall(0, nil, fmt.Errorf("something failed"))
	handle.QueryReturnsOnCall(1, nil, fmt.Errorf("something failed"))
	handle.QueryReturnsOnCall(2, nil, nil)
	err := database.WaitForWorkingConnection(handle)
	assert.Equal(t, 3, handle.QueryCallCount())
	query, _ := handle.QueryArgsForCall(0)
	assert.Equal(t, "SELECT 1", query)
	assert.Equal(t, 3, handle.QueryCallCount())
	assert.NoError(t, err)
}
