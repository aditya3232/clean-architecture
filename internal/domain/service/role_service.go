package service

import (
	"clean-architecture/internal/domain/entity"
	"clean-architecture/internal/port/outbound"
	"context"
)

type RoleServiceInterface interface {
	GetAll(ctx context.Context, search string) ([]entity.RoleEntity, error)
	GetByID(ctx context.Context, id int64) (*entity.RoleEntity, error)
	Create(ctx context.Context, req entity.RoleEntity) error
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, req entity.RoleEntity) error
}

type roleService struct {
	repo outbound.RoleRepositoryInterface
}

func NewRoleService(repo outbound.RoleRepositoryInterface) RoleServiceInterface {
	return &roleService{repo: repo}
}

func (r *roleService) Create(ctx context.Context, req entity.RoleEntity) error {
	return r.repo.Create(ctx, req)
}

func (r *roleService) Delete(ctx context.Context, id int64) error {
	return r.repo.Delete(ctx, id)
}

func (r *roleService) GetAll(ctx context.Context, search string) ([]entity.RoleEntity, error) {
	return r.repo.GetAll(ctx, search)
}

func (r *roleService) GetByID(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	return r.repo.GetByID(ctx, id)
}

func (r *roleService) Update(ctx context.Context, req entity.RoleEntity) error {
	return r.repo.Update(ctx, req)
}
