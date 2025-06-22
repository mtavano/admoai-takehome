package api

import (
	"context"
	"database/sql"

	"github.com/mtavano/admoai-takehome/internal/store"
	"github.com/stretchr/testify/mock"
)

// MockDatabase implements the Database interface for testing
type MockDatabase struct {
	mock.Mock
}

func (m *MockDatabase) BeginTx(ctx context.Context) (store.Transactioner, error) {
	args := m.Called(ctx)
	return args.Get(0).(store.Transactioner), args.Error(1)
}

func (m *MockDatabase) Exec(query string, params ...interface{}) (sql.Result, error) {
	args := m.Called(query, params)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockDatabase) Query(query string, params ...interface{}) (*sql.Rows, error) {
	args := m.Called(query, params)
	return args.Get(0).(*sql.Rows), args.Error(1)
}

func (m *MockDatabase) Get(dest interface{}, query string, args ...interface{}) error {
	mockArgs := m.Called(dest, query, args)
	return mockArgs.Error(0)
}

func (m *MockDatabase) Select(dest interface{}, query string, args ...interface{}) error {
	mockArgs := m.Called(dest, query, args)
	return mockArgs.Error(0)
}

// MockTransaction implements the Transactioner interface for testing
type MockTransaction struct {
	mock.Mock
}

func (m *MockTransaction) Exec(query string, params ...interface{}) (sql.Result, error) {
	args := m.Called(query, params)
	return args.Get(0).(sql.Result), args.Error(1)
}

func (m *MockTransaction) Query(query string, params ...interface{}) (*sql.Rows, error) {
	args := m.Called(query, params)
	return args.Get(0).(*sql.Rows), args.Error(1)
}

func (m *MockTransaction) Get(dest interface{}, query string, args ...interface{}) error {
	mockArgs := m.Called(dest, query, args)
	return mockArgs.Error(0)
}

func (m *MockTransaction) Select(dest interface{}, query string, args ...interface{}) error {
	mockArgs := m.Called(dest, query, args)
	return mockArgs.Error(0)
}

func (m *MockTransaction) Commit() error {
	args := m.Called()
	return args.Error(0)
}

func (m *MockTransaction) Rollback() error {
	args := m.Called()
	return args.Error(0)
}
