package seed

import (
	"clean-architecture/internal/adapter/outbound/postgres/model"
	utilpassword "clean-architecture/utils/password"

	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
)

func AdminSeed(db *gorm.DB) {
	bytes, err := utilpassword.HashPassword("admin123")
	if err != nil {
		log.Fatalf("[SeedAdmin-1]: %v", err)
	}

	modelRole := model.Role{}
	err = db.Where("name = ?", "Super Admin").First(&modelRole).Error
	if err != nil {
		log.Errorf("[SeedAdmin-2]: %v", err)
	}

	admin := model.User{
		Name:       "super admin",
		Email:      "superadmin@mail.com",
		Password:   bytes,
		IsVerified: true,
		Roles:      []model.Role{modelRole},
	}

	if err := db.FirstOrCreate(&admin, model.User{Email: "superadmin@mail.com"}).Error; err != nil {
		log.Errorf("[SeedAdmin-3]: %v", err)
	} else {
		log.Infof("Admin %s created", admin.Name)
	}
}
