package repository

import (
	"clean-architecture/internal/adapter/outbound/postgres/model"
	"clean-architecture/internal/domain/entity"
	outboundport "clean-architecture/internal/port/outbound"
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) outboundport.UserRepositoryInterface {
	return &userRepository{db: db}
}

func (u *userRepository) DeleteCustomer(ctx context.Context, customerID int64) error {
	modelUser := model.User{}
	if err := u.db.WithContext(ctx).Where("id =?", customerID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] DeleteCustomer: User not found")
			return err
		}
		log.Errorf("[UserRepository-2] DeleteCustomer: %v", err)
		return err
	}

	if err := u.db.Delete(&modelUser).Error; err != nil {
		log.Errorf("[UserRepository-3] DeleteCustomer: %v", err)
		return err
	}
	return nil
}

func (u *userRepository) UpdateCustomer(ctx context.Context, req entity.UserEntity) error {
	var (
		modelRole = model.Role{}
		modelUser = model.User{}
		updates   = map[string]interface{}{}
	)

	// Gunakan transaksi
	return u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		// 🔍 1. Cek role
		if err := tx.Where("id = ?", req.RoleID).First(&modelRole).Error; err != nil {
			log.Errorf("[UserRepository-1] UpdateCustomer: Role not found: %v", err)
			return err
		}

		// 🔍 2. Cek user
		if err := tx.Where("id = ?", req.ID).First(&modelUser).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				log.Infof("[UserRepository-2] UpdateCustomer: User not found")
				return errors.New("404")
			}
			log.Errorf("[UserRepository-3] UpdateCustomer: %v", err)
			return err
		}

		// 🧩 3. Siapkan field yang mau diupdate
		if req.Name != "" {
			updates["name"] = req.Name
		}
		if req.Email != "" {
			updates["email"] = req.Email
		}
		if req.Phone != "" {
			updates["phone"] = req.Phone
		}
		if req.Address != "" {
			updates["address"] = req.Address
		}
		if req.Lat != "" {
			updates["lat"] = req.Lat
		}
		if req.Lng != "" {
			updates["lng"] = req.Lng
		}
		if req.Photo != "" {
			updates["photo"] = req.Photo
		}
		if req.Password != "" {
			updates["password"] = req.Password
		}

		// 🚀 4. Jalankan partial update (jika ada field)
		if len(updates) > 0 {
			if err := tx.Model(&modelUser).Where("id = ?", req.ID).Updates(updates).Error; err != nil {
				log.Errorf("[UserRepository-4] UpdateCustomer: %v", err)
				return err
			}
		}

		// 🔗 5. Update relasi Role di pivo table user_role (many2many)
		// Relasi lama user dengan role lain akan dihapus
		// Relasi baru (user ↔ role) akan ditambahkan
		if err := tx.Model(&modelUser).Association("Roles").Replace(&modelRole); err != nil {
			log.Errorf("[UserRepository-5] UpdateCustomer (role): %v", err)
			return err
		}

		// ✅ 6. Commit otomatis jika semua berhasil
		log.Infof("[UserRepository] UpdateCustomer: User %d updated successfully", req.ID)
		return nil
	})
}

// CreateCustomer create user & user_role
func (u *userRepository) CreateCustomer(ctx context.Context, req entity.UserEntity) (int64, error) {
	var (
		modelRole = model.Role{}
		modelUser = model.User{}
	)

	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Cek Role
		if err := tx.Where("id = ?", req.RoleID).First(&modelRole).Error; err != nil {
			log.Errorf("[UserRepository-1] CreateCustomer: %v", err)
			return err
		}

		// Buat User baru
		modelUser = model.User{
			Name:       req.Name,
			Email:      req.Email,
			Password:   req.Password,
			Address:    req.Address,
			Lat:        req.Lat,
			Lng:        req.Lng,
			Phone:      req.Phone,
			Photo:      req.Photo,
			Roles:      []model.Role{modelRole},
			IsVerified: true,
		}

		if err := tx.Create(&modelUser).Error; err != nil {
			log.Errorf("[UserRepository-2] CreateCustomer: %v", err)
			return err
		}

		return nil
	})

	if err != nil {
		log.Errorf("[UserRepository-3] CreateCustomer: %v", err)
		return 0, err
	}

	return modelUser.ID, nil
}

func (u *userRepository) GetCustomerByID(ctx context.Context, customerID int64) (*entity.UserEntity, error) {
	modelUser := model.User{}

	if err := u.db.WithContext(ctx).Where("id = ?", customerID).Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] GetCustomerByID: User not found")
			return nil, err
		}
		log.Errorf("[UserRepository-2] GetCustomerByID: %v", err)
		return nil, err
	}

	var roleID int64
	if len(modelUser.Roles) > 0 {
		roleID = modelUser.Roles[0].ID
	}

	return &entity.UserEntity{
		ID:      customerID,
		Name:    modelUser.Name,
		Email:   modelUser.Email,
		RoleID:  roleID,
		Address: modelUser.Address,
		Lat:     modelUser.Lat,
		Lng:     modelUser.Lng,
		Phone:   modelUser.Phone,
		Photo:   modelUser.Photo,
	}, nil
}

func (u *userRepository) GetCustomerAll(ctx context.Context, query entity.QueryStringEntity) ([]entity.UserEntity, int64, int64, error) {
	var (
		modelUsers   []model.User
		respEntities []entity.UserEntity
		countData    int64
	)

	order := fmt.Sprintf("%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	sqlMain := u.db.WithContext(ctx).Preload("Roles", "name = ?", "Customer").
		Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%")

	if err := sqlMain.Model(&modelUsers).Count(&countData).Error; err != nil {
		log.Errorf("[UserRepository-1] GetCustomerAll: %v", err)
		return nil, 0, 0, err
	}

	totalPage := int(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.Order(order).Limit(int(query.Limit)).Offset(int(offset)).Find(&modelUsers).Error; err != nil {
		log.Errorf("[UserRepository-3] GetCustomerAll: %v", err)
		return nil, 0, 0, err
	}

	if len(modelUsers) < 1 {
		err := errors.New("404")
		log.Infof("[UserRepository-4] GetCustomerAll: No Customer found")
		return nil, 0, 0, err
	}

	for _, val := range modelUsers {
		roleName := ""
		for _, role := range val.Roles {
			roleName = role.Name
			break // hanya butuh 1 role name
		}
		respEntities = append(respEntities, entity.UserEntity{
			ID:       val.ID,
			Name:     val.Name,
			Email:    val.Email,
			RoleName: roleName,
			Phone:    val.Phone,
			Photo:    val.Photo,
		})
	}

	return respEntities, countData, int64(totalPage), nil
}

func (u *userRepository) UpdateDataUser(ctx context.Context, req entity.UserEntity) error {
	var (
		modelUser model.User
		updates   = map[string]interface{}{}
	)

	// 🔍 Cek apakah user ditemukan dan sudah terverifikasi
	if err := u.db.WithContext(ctx).
		Where("id = ? AND is_verified = true", req.ID).
		First(&modelUser).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("[UserRepository-1] UpdateDataUser: User not found or not verified")
			return errors.New("404")
		}

		log.Errorf("[UserRepository-2] UpdateDataUser: %v", err)
		return err
	}

	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Email != "" {
		updates["email"] = req.Email
	}
	if req.Address != "" {
		updates["address"] = req.Address
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Photo != "" {
		updates["photo"] = req.Photo
	}
	if req.Lat != "" {
		updates["lat"] = req.Lat
	}
	if req.Lng != "" {
		updates["lng"] = req.Lng
	}

	// 🚀 Jalankan update hanya kalau ada field yang berubah
	if len(updates) > 0 {
		if err := u.db.WithContext(ctx).
			Model(&modelUser).
			Updates(updates).Error; err != nil {
			log.Errorf("[UserRepository-3] UpdateDataUser: %v", err)
			return err
		}
		log.Infof("[UserRepository] UpdateDataUser: User %d updated successfully", req.ID)
	} else {
		log.Infof("[UserRepository] UpdateDataUser: No fields to update for user %d", req.ID)
	}

	return nil
}

func (u *userRepository) GetUserByID(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	modelUser := model.User{}

	if err := u.db.WithContext(ctx).Where("id =? AND is_verified = true", userID).Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Errorf("[UserRepository-1] GetUserByID: %v", err)
			return nil, err
		}
		log.Errorf("[UserRepository-2] GetUserByID: %v", err)
		return nil, err
	}

	var roleName string
	if len(modelUser.Roles) > 0 {
		roleName = modelUser.Roles[0].Name
	}

	return &entity.UserEntity{
		ID:       modelUser.ID,
		Email:    modelUser.Email,
		Name:     modelUser.Name,
		RoleName: roleName,
		Lat:      modelUser.Lat,
		Lng:      modelUser.Lng,
		Address:  modelUser.Address,
		Phone:    modelUser.Phone,
		Photo:    modelUser.Photo,
	}, nil
}

func (u *userRepository) UpdatePasswordByID(ctx context.Context, req entity.UserEntity) error {
	var modelUser model.User
	if err := u.db.WithContext(ctx).Where("id = ?", req.ID).First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("[UserRepository-1] UpdatePasswordByID: User not found")
			return errors.New("404")
		}
		log.Errorf("[UserRepository-2] UpdatePasswordByID: %v", err)
		return err
	}

	// 🚀 Update hanya kolom password
	if err := u.db.WithContext(ctx).
		Model(&modelUser).
		Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"password": req.Password,
		}).Error; err != nil {
		log.Errorf("[UserRepository-3] UpdatePasswordByID: %v", err)
		return err
	}

	return nil
}

func (u *userRepository) UpdateUserVerified(ctx context.Context, userID int64) (*entity.UserEntity, error) {
	var modelUser model.User

	if err := u.db.WithContext(ctx).Where("id = ?", userID).Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Infof("[UserRepository-1] UpdateUserVerified: user not found")
			return nil, errors.New("404")
		}
		log.Errorf("[UserRepository-2] UpdateUserVerified: %v", err)
		return nil, err
	}

	// 🧩 Siapkan data yang mau di-update (hanya jika terisi)
	updateData := map[string]interface{}{}
	if !modelUser.IsVerified { // hanya update kalau belum verified
		updateData["is_verified"] = true
	}

	// ⚙️ Jalankan update hanya jika ada kolom diupdate
	if len(updateData) > 0 {
		if err := u.db.WithContext(ctx).
			Model(&model.User{}).
			Where("id = ?", userID).
			Updates(updateData).Error; err != nil {
			log.Errorf("[UserRepository-3] UpdateUserVerified: %v", err)
			return nil, err
		}
	}

	// 🧾 Ambil ulang data role (jika sudah ada)
	var roleName string
	if len(modelUser.Roles) > 0 {
		roleName = modelUser.Roles[0].Name
	}

	return &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Email:      modelUser.Email,
		RoleName:   roleName,
		Address:    modelUser.Address,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		IsVerified: true, // sudah di-update
	}, nil
}

func (u *userRepository) CreateUserAccount(ctx context.Context, req entity.UserEntity) (int64, error) {
	var (
		roleID    int64
		modelUser model.User
	)

	err := u.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1 Ambil role "Customer"
		if err := tx.Model(&model.Role{}).
			Select("id").
			Where("name = ?", "Customer").
			Scan(&roleID).Error; err != nil {
			log.Errorf("[UserRepository-1] CreateUserAccount: failed to get role 'Customer': %v", err)
			return err
		}

		if roleID == 0 {
			err := errors.New("role 'Customer' not found")
			log.Errorf("[UserRepository-1b] CreateUserAccount: %v", err)
			return err
		}

		// 2 Buat user baru
		modelUser = model.User{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
			Roles:    []model.Role{{ID: roleID}},
		}

		if err := tx.Create(&modelUser).Error; err != nil {
			log.Errorf("[UserRepository-2] CreateUserAccount: failed to create user: %v", err)
			return err
		}

		// 3 Buat verification token
		modelVerify := model.VerificationToken{
			UserID:    modelUser.ID,
			Token:     req.Token,
			TokenType: "email_verification",
			ExpiresAt: time.Now().Add(time.Hour),
		}

		if err := tx.Create(&modelVerify).Error; err != nil {
			log.Errorf("[UserRepository-3] CreateUserAccount: failed to create verification token: %v", err)
			return err
		}

		// ✅ Semua sukses
		log.Infof("[UserRepository-4] CreateUserAccount: user '%s' created successfully (ID=%d, RoleID=%d)", modelUser.Email, modelUser.ID, roleID)
		return nil
	})

	if err != nil {
		log.Errorf("[UserRepository-5] CreateUserAccount: %v", err)
		return 0, err
	}

	return modelUser.ID, nil
}

func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	modelUser := model.User{}

	if err := u.db.WithContext(ctx).Where("email = ? AND is_verified = ?", email, true).
		Preload("Roles").First(&modelUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err = errors.New("404")
			log.Infof("[UserRepository-1] GetUserByEmail: User not found")
			return nil, err
		}
		log.Errorf("[UserRepository-1] GetUserByEmail: %v", err)
		return nil, err
	}

	var roleName string
	if len(modelUser.Roles) > 0 {
		roleName = modelUser.Roles[0].Name
	}

	return &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Email:      email,
		Password:   modelUser.Password,
		RoleName:   roleName,
		Address:    modelUser.Address,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		IsVerified: modelUser.IsVerified,
	}, nil
}
