package seed

import (
	"log"

	"github.com/jinzhu/gorm"
	"shorterer-link/api/models"
)

var users = models.User{
	Name:     "ahsanul",
	Email:    "ahsanul@gmail.com",
	Password: "password123",
}

var url = models.Url{
	OriginalUrl: "https://www.linkedin.com/in/ahsanulwaladi",
	CustomUrl:   "linked-aw",
}

func Load(db *gorm.DB) {

	err := db.Debug().DropTableIfExists(&models.Url{}, &models.User{}).Error
	if err != nil {
		log.Fatalf("cannot drop table: %v", err)
	}
	err = db.Debug().AutoMigrate(&models.User{}, &models.Url{}).Error
	if err != nil {
		log.Fatalf("cannot migrate table: %v", err)
	}

	err = db.Debug().Model(&models.Url{}).AddForeignKey("user_id", "users(id)", "cascade", "cascade").Error
	if err != nil {
		log.Fatalf("attaching foreign key error: %v", err)
	}

	//for i, _ := range users {
	err = db.Debug().Model(&models.User{}).Create(&users).Error
	if err != nil {
		log.Fatalf("cannot seed users table: %v", err)
	}
	url.UserID = users.ID

	err = db.Debug().Model(&models.Url{}).Create(&url).Error
	if err != nil {
		log.Fatalf("cannot seed posts table: %v", err)
	}
	//}
}
