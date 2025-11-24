package repository

import (
	"clean-architecture/internal/adapter/outbound/postgres/model"
	"clean-architecture/internal/domain/entity"
	"clean-architecture/internal/port/outbound"
	"context"
	"errors"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) outbound.RoleRepositoryInterface {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(ctx context.Context, req entity.RoleEntity) error {
	modelRole := model.Role{
		Name: req.Name,
	}

	if err := r.db.WithContext(ctx).Create(&modelRole).Error; err != nil {
		log.Errorf("[RoleRepository-1] Create: %v", err)
		return err
	}

	return nil
}

func (r *roleRepository) Delete(ctx context.Context, id int64) error {
	modelRole := model.Role{}

	if err := r.db.WithContext(ctx).Where("id = ?", id).Preload("Users").First(&modelRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[RoleRepository-1] Delete: Role not found")
			return err
		}
		log.Errorf("[RoleRepository-2] Delete: %v", err)
		return err
	}

	if len(modelRole.Users) > 0 {
		err := errors.New("400")
		log.Infof("[RoleRepository-3] Delete: Role is associated with users")
		return err
	}

	if err := r.db.WithContext(ctx).Delete(&modelRole).Error; err != nil {
		log.Errorf("[RoleRepository-3] Delete: %v", err)
		return err
	}

	return nil
}

func (r *roleRepository) GetAll(ctx context.Context, search string) ([]entity.RoleEntity, error) {
	var (
		modelRoles []model.Role
		entityRole []entity.RoleEntity
	)

	if err := r.db.WithContext(ctx).Where("name ILIKE ?", "%"+search+"%").Find(&modelRoles).Error; err != nil {
		log.Errorf("[RoleRepository-1] GetAll: %v", err)
		return nil, err
	}

	if len(modelRoles) == 0 {
		err := errors.New("404")
		log.Infof("[RoleRepository-2] GetAll: No role found")
		return nil, err
	}

	for _, modelRole := range modelRoles {
		entityRole = append(entityRole, entity.RoleEntity{
			ID:   modelRole.ID,
			Name: modelRole.Name,
		})
	}

	return entityRole, nil
}

func (r *roleRepository) GetByID(ctx context.Context, id int64) (*entity.RoleEntity, error) {
	modelRole := model.Role{}

	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&modelRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[RoleRepository-1] GetByID: Role not found")
			return nil, err
		}
		log.Errorf("[RoleRepository-2] GetAll: %v", err)
		return nil, err
	}

	return &entity.RoleEntity{
		ID:   modelRole.ID,
		Name: modelRole.Name,
	}, nil
}

func (r *roleRepository) Update(ctx context.Context, req entity.RoleEntity) error {
	//modelRole := model.Role{
	//	Name: req.Name,
	//}

	var (
		modelRole = model.Role{}
		updates   = map[string]interface{}{}
	)

	if err := r.db.WithContext(ctx).Where("id = ?", req.ID).First(&modelRole).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[RoleRepository-1] Update: Role not found")
			return err
		}
		log.Errorf("[RoleRepository-2] Update: %v", err)
		return err
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}

	if len(updates) > 0 {
		if err := r.db.WithContext(ctx).Model(&modelRole).Updates(updates).Error; err != nil {
			log.Errorf("[RoleRepository-3] Update: %v", err)
			return err
		}
	}

	return nil
}
