package dogowners

import (
	"time"

	"gorm.io/gorm"
)

type DogDBModel struct {
	gorm.Model
	OwnerID  uint
	Name     string
	Birthday time.Time
}

type OwnerDBModel struct {
	gorm.Model
	Name     string
	Email    *string
	Age      uint8
	Birthday time.Time
	Dogs     []DogDBModel `gorm:"foreignKey:OwnerID"`
}

func RunMigration(db *gorm.DB) error {
	return db.AutoMigrate(&OwnerDBModel{}, &DogDBModel{})
}

type Storage struct {
	db *gorm.DB
}

func NewDatabase(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (d *Storage) SaveOwner(user *OwnerDBModel) error {
	return d.db.Create(user).Error
}

func (d *Storage) SaveDog(dog *DogDBModel) error {
	return d.db.Create(dog).Error
}

func (d *Storage) Owner(id uint) (*OwnerDBModel, error) {
	user := &OwnerDBModel{}
	err := d.db.Preload("Dogs").First(user, id).Error
	return user, err
}
