package models

import (
	"html"
	"log"
	"strings"
	"time"

	"github.com/MinhL4m/blogs/helpers"
	"github.com/jinzhu/gorm"
)

type User struct {
	ID        uint32    `gorm:"primary_key;auto_increment" json:"id"`
	Username  string    `gorm:"size:255;not null" json:"username"`
	Email     string    `gorm:"size:100;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

func (u *User) BeforeSave() error {
	if u.Password == "" {
		hashedPassword, err := helpers.Hash(u.Password)

		if err != nil {
			return err
		}

		u.Password = string(hashedPassword)
	}
	return nil
}

func (u *User) Prepare() {
	u.ID = 0
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	u.Email = html.EscapeString(strings.TrimSpace(u.Email))
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validation(action string) error {
	switch strings.ToLower(action) {
	case "update":
		return helpers.UpdateValid(u.Username, u.Password, u.Email)
	case "login":
		return helpers.LoginValid(u.Password, u.Username)
	default:
		return helpers.UpdateValid(u.Username, u.Password, u.Email)
	}
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	err := db.Debug().Create(&u).Error

	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) UpdateUser(db *gorm.DB, uid uint32) (*User, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).UpdateColumns(
		map[string]interface{}{
			"username":  u.Username,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}

	err := db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error

	if err != nil {
		return &User{}, db.Error
	}

	return u, nil
}

func (u *User) UpdateUserPassword(db *gorm.DB, uid uint32) (*User, error) {
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"update_at": time.Now(),
		},
	)

	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Debug().Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
