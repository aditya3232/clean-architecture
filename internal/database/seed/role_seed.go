package seed

import (
	"clean-architecture/internal/adapter/outbound/postgres/model"

	"github.com/labstack/gommon/log"

	"gorm.io/gorm"
)

func RoleSeed(db *gorm.DB) {
	roles := []model.Role{
		{
			Name: "Super Admin",
		},
		{
			Name: "Customer",
		},
	}

	for _, role := range roles {
		if err := db.FirstOrCreate(&role, model.Role{Name: role.Name}).Error; err != nil {
			log.Errorf("[SeedRole-1]: %v", err)
		} else {
			log.Infof("Role %s created", role.Name)
		}
	}
}
