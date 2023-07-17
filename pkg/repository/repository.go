package repository

import (
	"errors"
	"log"
	"service2/pkg/db"
	"service2/pkg/entity"

	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

func init() {
	DB, err = db.ConnectToDB()
	if err != nil {
		log.Fatal(err)
	}
}

func GetAllUsers(offset, limit int) ([]entity.User, error) {
	var users []entity.User
	err := DB.Offset(offset).Limit(limit).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsersByPermission(offset, limit int, permission bool) ([]entity.User, error) {
	var users []entity.User
	err := DB.Offset(offset).Limit(limit).Where("permission = ?", permission).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsersById(userId int) ([]entity.User, error) {
	var users []entity.User
	err := DB.Where("id = ?", userId).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetAllUsersByName(name string) ([]entity.User, error) {
	var users []entity.User
	err := DB.Where("first_name LIKE ?", name).Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func GetByID(id int) (*entity.User, error) {
	var user entity.User
	result := DB.Where(&entity.User{ID: id}).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}
	return &user, nil
}

func Update(user *entity.User) error {
	return DB.Save(user).Error
}

func GetByApparelName(apparelName string) error {

	var existingApparel entity.Ticket
	result := DB.Where(&entity.Apparel{Name: apparelName}).First(&existingApparel)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return result.Error
		}
		return result.Error
	}
	return nil

}

func CreateApparel(apparel *entity.Apparel) (int, error) {
	if err := DB.Create(apparel).Error; err != nil {
		return 0, err
	}
	return apparel.ID, nil
}

func GetApparelByID(id int) (*entity.Apparel, error) {
	var apparel entity.Apparel
	result := DB.First(&apparel, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, result.Error
		}
		return nil, result.Error
	}
	return &apparel, nil
}

func UpdateApparel(apparel *entity.Apparel) error {
	return DB.Save(apparel).Error
}

func CreateCoupon(coupon *entity.Coupon) error {
	if err := DB.Create(coupon).Error; err != nil {
		return err
	}
	return nil
}

func CreateOffer(offer *entity.Offer) error {
	if err := DB.Create(offer).Error; err != nil {
		return err
	}
	return nil
}
