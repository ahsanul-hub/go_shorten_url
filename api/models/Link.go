package models

import (
	"errors"
	"log"
	"time"

	"github.com/jinzhu/gorm"
)

type Url struct {
	ID          uint64    `gorm:"primary_key;auto_increment" json:"id"`
	OriginalUrl string    `gorm:"size:255;not null;unique" json:"originalUrl"`
	CustomUrl   string    `gorm:"size:255;not null;" json:"customUrl"`
	UserID      uint32    `gorm:"not null" json:"user_id"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

//
//func (u *Url) Prepare() {
//	u.ID = 0
//	u.OriginalUrl = html.EscapeString(strings.TrimSpace(u.OriginalUrl))
//	u.CustomUrl = html.EscapeString(strings.TrimSpace(u.CustomUrl))
//	u.CreatedAt = time.Now()
//	u.UpdatedAt = time.Now()
//}

func (u *Url) Validate() error {

	if u.OriginalUrl == "" {
		return errors.New("Required URL")
	}
	if u.CustomUrl == "" {
		return errors.New("Required Custom URL")
	}
	return nil
}

func (u *Url) SaveUrl(db *gorm.DB) (*Url, error) {
	var err error
	err = db.Debug().Model(&Url{}).Create(&u).Error
	if err != nil {
		return &Url{}, err
	}
	if u.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", u.UserID).Error
		if err != nil {
			return &Url{}, err
		}
	}
	return u, nil
}

func (u *Url) FindAllUrl(db *gorm.DB) (*[]Url, error) {
	var err error
	url := []Url{}
	err = db.Debug().Model(&Url{}).Limit(100).Find(&url).Error
	if err != nil {
		return &[]Url{}, err
	}
	if len(url) > 0 {
		for i, _ := range url {
			err := db.Debug().Model(&User{}).Where("id = ?", url[i].UserID).Error
			if err != nil {
				return &[]Url{}, err
			}
		}
	}
	return &url, nil
}

func (u *Url) FindUrlByLink(db *gorm.DB, url string) (*Url, error) {
	var err error
	err = db.Debug().Model(&Url{}).Where("custom_url = ?", url).Take(&u).Error
	if err != nil {
		return &Url{}, err
	}
	if u.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", u.UserID).Error
		if err != nil {
			return &Url{}, err
		}
	}
	log.Println(u)
	return u, nil
}

func (u *Url) UpdateAUrl(db *gorm.DB, urlid uint64) (*Url, error) {

	var err error
	db = db.Debug().Model(&Url{}).Where("id = ?", urlid).Take(&Url{}).UpdateColumns(
		map[string]interface{}{
			"originalUrl": u.OriginalUrl,
			"customUrl":   u.CustomUrl,
			"updated_at":  time.Now(),
		},
	)
	err = db.Debug().Model(&Url{}).Where("id = ?", urlid).Take(&u).Error
	if err != nil {
		return &Url{}, err
	}
	if u.ID != 0 {
		err = db.Debug().Model(&User{}).Where("id = ?", u.UserID).Error
		if err != nil {
			return &Url{}, err
		}
	}
	return u, nil
}

func (u *Url) DeleteAUrl(db *gorm.DB, key string, uid uint32) (int64, error) {

	db = db.Debug().Model(&Url{}).Where("custom_url = ? and user_id = ?", key, uid).Take(&Url{}).Delete(&Url{})

	if db.Error != nil {
		if gorm.IsRecordNotFoundError(db.Error) {
			return 0, errors.New("Post not found")
		}
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
