package models

import (
	"fmt"

	"gorm.io/gorm"
)

const MAX_NAME_SIZE = 255

type Department struct {
	gorm.Model
	Name      string
	Members   []User
	CompanyID uint
	Company   Company
}

func NewDepartment(name string, companyID uint) *Department {
	return &Department{
		Name:      name,
		CompanyID: companyID,
		Members:   []User{},
	}
}

func (d *Department) isValid() error {
	if len([]rune(d.Name)) > MAX_NAME_SIZE {
		return fmt.Errorf("department name '%s' too long", d.Name)
	}
	return nil
}

func (d *Department) BeforeCreate(tx *gorm.DB) error {
	if err := d.isValid(); err != nil {
		return err
	}
	return nil
}
