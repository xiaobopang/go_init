package model

/*
 * @Script: model.go
 * @Author: pangxiaobo
 * @Email: 10846295@qq.com
 * @Create At: 2018-11-07 16:55:04
 * @Last Modified By: pangxiaobo
 * @Last Modified At: 2018-12-12 14:24:45
 * @Description: This is description.
 */
import (
	"time"
)

type User struct {
	ID        int    `json:"id" gorm:"index"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	Age       int    `json:"age"`
	Email     string `json:"email"`
	Gender    int    `json:"gender"`
	CreatedAt int64  `json:"created_at"`
	UpdatedAt int64  `json:"updated_at"`
}

func GetUserById(id int) ([]*User, error) {
	var user []*User
	err := DB.Select("id, username, age, email, gender, created_at").Where("id = ? AND is_deleted = ? ", id, 0).First(&user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func AddUser(name string, password string, age int, gender int, email string) error {
	user := User{
		Username:  name,
		Password:  password,
		Age:       age,
		Gender:    gender,
		Email:     email,
		CreatedAt: time.Now().Unix(),
	}
	if err := DB.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func DelUser(id int) error {
	if err := DB.Where("id = ?", id).Delete(&User{}).Error; err != nil {
		return err
	}

	return nil
}

func UptUser(id int, data interface{}) error {

	if err := DB.Model(&User{}).Where("id = ? AND is_deleted = ? ", id, 0).Updates(data).Error; err != nil {
		return err
	}

	return nil
}
