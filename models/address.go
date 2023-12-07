package models

import (
	"fmt"

	"gorm.io/gorm"
)

const MAX_STREETNUM_SIZE = 8
const MAX_STREETNAME_SIZE = 120
const MAX_SUITE_SIZE = 8
const MAX_CITY_SIZE = 32
const MAX_STATE_SIZE = 32
const MAX_ZIP_SIZE = 5

type Address struct {
	gorm.Model
	StreetNumber string
	StreetName   string
	Suite        string
	City         string
	State        string
	ZipCode      string
}

func (a *Address) IsValidAddr() error {
	if a.StreetName == "" {
		return fmt.Errorf("must provide a street name")
	}
	if len([]rune(a.StreetName)) > MAX_STREETNAME_SIZE {
		return fmt.Errorf("street name too long")
	}
	if len([]rune(a.Suite)) > MAX_STREETNAME_SIZE {
		return fmt.Errorf("suite too long")
	}
	if a.City == "" {
		return fmt.Errorf("must provide a city")
	}
	if a.State == "" {
		return fmt.Errorf("must provide a state")
	}
	if a.ZipCode == "" {
		return fmt.Errorf("must provide a zipcode")
	}
	if len([]rune(a.City)) > MAX_CITY_SIZE {
		return fmt.Errorf("invalid city")
	}
	if len([]rune(a.State)) > MAX_STATE_SIZE {
		return fmt.Errorf("invalid state")
	}
	if len([]rune(a.ZipCode)) > MAX_ZIP_SIZE {
		return fmt.Errorf("invalid zip code")
	}
	return nil
}
