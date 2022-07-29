package controllers

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
	"shorterer-link/api/models"
)

var server = Server{}

var userInstance = models.User{}
var urlInstance = models.Url{}

func TestMain(m *testing.M) {
	var err error
	err = godotenv.Load(os.ExpandEnv("../../.env"))
	if err != nil {
		log.Fatalf("Error getting env %v\n", err)
	}
	Database()

	os.Exit(m.Run())
}

func Database() {

	var err error

	TestDbDriver := os.Getenv("TestDbDriver")

	DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("TestDbUser"), os.Getenv("TestDbPassword"), os.Getenv("TestDbHost"), os.Getenv("TestDbPort"), os.Getenv("TestDbName"))
	server.DB, err = gorm.Open(TestDbDriver, DBURL)
	if err != nil {
		fmt.Printf("Cannot connect to %s database\n", TestDbDriver)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the %s database\n", TestDbDriver)
	}
}

func refreshUserTable() error {
	err := server.DB.DropTableIfExists(&models.User{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed table")
	return nil
}

func seedOneUser() (models.User, error) {

	err := refreshUserTable()
	if err != nil {
		log.Fatal(err)
	}

	user := models.User{
		Name:     "Peter",
		Email:    "pet@gmail.com",
		Password: "password",
	}

	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

//func seedUsers() ([]models.User, error) {
//
//	var err error
//	if err != nil {
//		return nil, err
//	}
//	users := []models.User{
//		models.User{
//			Name:     "Steven victor",
//			Email:    "steven@gmail.com",
//			Password: "password",
//		},
//		models.User{
//			Name:     "Kenny Morris",
//			Email:    "kenny@gmail.com",
//			Password: "password",
//		},
//	}
//	for i, _ := range users {
//		err := server.DB.Model(&models.User{}).Create(&users[i]).Error
//		if err != nil {
//			return []models.User{}, err
//		}
//	}
//	return users, nil
//}

func refreshUserAndUrlTable() error {

	err := server.DB.DropTableIfExists(&models.User{}, &models.Url{}).Error
	if err != nil {
		return err
	}
	err = server.DB.AutoMigrate(&models.User{}, &models.Url{}).Error
	if err != nil {
		return err
	}
	log.Printf("Successfully refreshed tables")
	return nil
}

func seedOneUserAndUrl() (models.Url, error) {

	err := refreshUserAndUrlTable()
	if err != nil {
		return models.Url{}, err
	}
	user := models.User{
		Name:     "Sam Phil",
		Email:    "sam@gmail.com",
		Password: "password",
	}
	err = server.DB.Model(&models.User{}).Create(&user).Error
	if err != nil {
		return models.Url{}, err
	}
	url := models.Url{
		OriginalUrl: "https://www.linkedin.com/in/ahsanulwaladi",
		CustomUrl:   "linkedin-aw",
	}
	err = server.DB.Model(&models.Url{}).Create(&url).Error
	if err != nil {
		return models.Url{}, err
	}
	return url, nil
}

func seedUsersAndUrl() ([]models.User, []models.Url, error) {

	var err error

	if err != nil {
		return []models.User{}, []models.Url{}, err
	}
	var users = []models.User{
		models.User{
			Name:     "Steven victor",
			Email:    "steven@gmail.com",
			Password: "password",
		},
		models.User{
			Name:     "Magu Frank",
			Email:    "magu@gmail.com",
			Password: "password",
		},
	}
	var url = []models.Url{
		models.Url{
			OriginalUrl: "https://www.linkedin.com/in/ahsanulwaladi",
			CustomUrl:   "linkedin-aw",
		},
		models.Url{
			OriginalUrl: "https://www.linkedin.com/in/ahsanulwaladi2",
			CustomUrl:   "linkedin-aw2",
		},
	}

	for i, _ := range users {
		err = server.DB.Model(&models.User{}).Create(&users[i]).Error
		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
		url[i].UserID = users[i].ID

		err = server.DB.Model(&models.Url{}).Create(&url[i]).Error
		if err != nil {
			log.Fatalf("cannot seed url table: %v", err)
		}
	}
	return users, url, nil
}
