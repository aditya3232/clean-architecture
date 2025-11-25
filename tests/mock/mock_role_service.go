package mock

import (
	"context"

	"clean-architecture/internal/domain/entity"

	"github.com/stretchr/testify/mock"
)

// MockRoleService adalah mock implementasi dari service.RoleServiceInterface
// harus ada interface service nya untuk memanggil NewRoleHandler
type MockRoleService struct {
	mock.Mock
}

func (m *MockRoleService) GetAll(ctx context.Context, search string) ([]entity.RoleEntity, error) {
	args := m.Called(ctx, search)
	return args.Get(0).([]entity.RoleEntity), args.Error(1)
}

func (m *MockRoleService) GetByID(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*entity.RoleEntity), args.Error(1)
}

func (m *MockRoleService) Create(ctx context.Context, req entity.RoleEntity) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockRoleService) Delete(ctx context.Context, id int64) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockRoleService) Update(ctx context.Context, req entity.RoleEntity) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
