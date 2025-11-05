// models/user.go
package models

import (
	"fmt"
	"list-of-maldives/internal/database"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID         uint           `gorm:"primaryKey" json:"id"`
	UUID       string         `gorm:"uniqueIndex;not null" json:"uuid"`
	Email      string         `gorm:"uniqueIndex;not null" json:"email"`
	NickName   string         `gorm:"size:100" json:"nickname"`
	Password   string         `json:"-"`
	Provider   string         `gorm:"size:50;default:'email'" json:"provider"`
	ProviderID string         `gorm:"size:255;index" json:"provider_id"`
	IsVerified bool           `gorm:"default:false" json:"is_verified"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`
}

// HashPassword hashes the user's password
func (u *User) HashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

// CheckPassword compares the provided password with the stored hash
func (u *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	u.UUID = uuid.New().String()

	// Hash password if it's provided and not already hashed
	if u.Password != "" && len(u.Password) < 60 {
		if err := u.HashPassword(); err != nil {
			return err
		}
	}
	return nil
}

// FindOrCreateByProvider finds or creates a user by OAuth provider
func FindOrCreateByProvider(s database.Service, provider, providerID, email, name string) (*User, error) {
	db := s.GormDB()

	var user User
	// 1. Try to find the user by provider and providerID
	err := db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error

	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	if err == gorm.ErrRecordNotFound {
		// 2. User not found, create a new one
		user = User{
			Email:      email,
			NickName:   name,
			Provider:   provider,
			ProviderID: providerID,
		}
		if err = db.Create(&user).Error; err != nil {
			return nil, err
		}
	} else {
		// 3. User found, update NickName and Email
		updates := map[string]interface{}{
			"email":     email,
			"nick_name": name,
		}

		// Update the database record
		if err = db.Model(&user).Updates(updates).Error; err != nil {
			return nil, err
		}

		// Update the local struct to reflect the new data before returning
		user.Email = email
		user.NickName = name
	}
	return &user, nil
}

// FindByEmail finds a user by email
func FindByEmail(s database.Service, email string) (*User, error) {
	db := s.GormDB()
	var user User
	err := db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser creates a new email/password user
func CreateUser(s database.Service, email, password, nickname string) (*User, error) {
	db := s.GormDB()

	// Check if user already exists
	var existingUser User
	err := db.Where("email = ? AND provider = 'email'", email).First(&existingUser).Error
	if err == nil {
		return nil, fmt.Errorf("user already exists")
	} else if err != gorm.ErrRecordNotFound {
		return nil, err
	}

	user := User{
		Email:      email,
		Password:   password,
		NickName:   nickname,
		Provider:   "email",
		IsVerified: true,
	}

	if err := db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func MigrateUser(db *gorm.DB) error {
	return db.AutoMigrate(&User{})
}
