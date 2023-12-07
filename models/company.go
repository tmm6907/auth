package models

import (
	"gorm.io/gorm"
)

type Company struct {
	gorm.Model
	Name        string
	Alias       string
	AddressID   uint
	Address     Address
	Departments []Department
}

func NewCompany(name string, addr Address) *Company {
	return &Company{
		Name:        name,
		Address:     addr,
		Departments: []Department{},
	}
}

func (c *Company) isValidName() error {
	return nil
}

func (c *Company) isValidAddr() error {
	if err := c.Address.IsValidAddr(); err != nil {
		return err
	}
	return nil
}

func (c *Company) isValid() error {
	if err := c.isValidName(); err != nil {
		return err
	}
	if err := c.isValidAddr(); err != nil {
		return err
	}
	return nil
}

func (c *Company) BeforeCreate(tx *gorm.DB) error {
	if err := c.isValid(); err != nil {
		return err
	}
	return nil
}
