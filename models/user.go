package models

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

const MIN_FNAME_SIZE = 2
const MAX_FNAME_SIZE = 50
const MAX_INTITIALS_SIZE = 2
const MIN_LNAME_SIZE = 2
const MAX_LNAME_SIZE = 50
const MIN_USERNAME_SIZE = 6
const MAX_USERNAME_SIZE = 16
const MIN_PASSWORD_SIZE = 8
const MAX_PASSWORD_SIZE = 16
const EMAIL_REGEX = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`

var PHONE_FILTER = []string{"-", "(", ")", " "}

type Role struct {
	gorm.Model
	UserID uint
	Name   string
}

func NewRole(name string) *Role {
	return &Role{
		Name: name,
	}
}

type UserConfig struct {
	FirstName      string
	LastName       string
	MiddleInitials string
	Username       string
	Password       string
	Email          string
	Phone          string
	Role           Role
}

type User struct {
	gorm.Model
	FirstName      string
	LastName       string
	MiddleInitials string
	Username       string `gorm:"unique"`
	Password       string
	Email          string
	Phone          string
	Role           Role
	CompanyID      uint
	Company        Company
	DepartmentID   uint
	Department     Department
	UpdatedAt      time.Time
	CreatedAt      time.Time
	DeletedAt      gorm.DeletedAt
}

func NewUser(config UserConfig, companyID uint, departmentID uint) *User {
	return &User{
		FirstName:      config.FirstName,
		LastName:       config.LastName,
		MiddleInitials: config.MiddleInitials,
		Username:       config.Username,
		Password:       config.Password,
		Email:          config.Email,
		Phone:          config.Phone,
		Role:           config.Role,
		CompanyID:      companyID,
		DepartmentID:   departmentID,
	}
}

func (u *User) isValid() (err error) {
	if err := u.isValidFirstName(); err != nil {
		return err
	}
	if err := u.isValidLastName(); err != nil {
		return err
	}
	if err := u.isValidMiddleInitials(); err != nil {
		return err
	}
	if err := u.isValidUsername(); err != nil {
		return err
	}
	if err := u.isValidPassword(); err != nil {
		return err
	}
	if err := u.isValidEmail(); err != nil {
		return err
	}
	if err := u.isValidPhone(); err != nil {
		return err
	}
	return nil
}

func (u *User) isValidFirstName() (err error) {
	if u.FirstName == "" {
		return fmt.Errorf("user must provide a first name")
	}
	if size := len([]rune(u.FirstName)); size < MIN_FNAME_SIZE {
		return fmt.Errorf("user first name of size %d is too short", size)
	}
	if size := len([]rune(u.FirstName)); size > MAX_FNAME_SIZE {
		return fmt.Errorf("user first name of size %d is too long", size)
	}
	return nil
}

func (u *User) isValidLastName() (err error) {
	if u.LastName == "" {
		return fmt.Errorf("user must provide a last name")
	}
	if size := len([]rune(u.LastName)); size < MIN_LNAME_SIZE {
		return fmt.Errorf("user last name of size %d is too short", size)
	}
	if size := len([]rune(u.LastName)); size > MAX_LNAME_SIZE {
		return fmt.Errorf("user last name of size %d is too long", size)
	}
	return nil
}

func (u *User) isValidMiddleInitials() (err error) {
	if size := len([]rune(u.MiddleInitials)); size > MAX_INTITIALS_SIZE {
		return fmt.Errorf("user middle initials of size %d is greater than max size %d", size, MAX_INTITIALS_SIZE)
	}
	return nil
}

func (u *User) isValidUsername() (err error) {
	if u.Username == "" {
		return fmt.Errorf("user must provide a username")
	}
	if size := len([]rune(u.Username)); size < MIN_USERNAME_SIZE {
		return fmt.Errorf("username of size %d is too short", size)
	}
	if size := len([]rune(u.Username)); size > MAX_USERNAME_SIZE {
		return fmt.Errorf("username of size %d is too long", size)
	}
	return nil
}

func (u *User) isValidPassword() (err error) {
	if u.Password == "" {
		return fmt.Errorf("user must provide a password")
	}
	if size := len([]rune(u.Password)); size < MIN_PASSWORD_SIZE {
		return fmt.Errorf("user password of size %d is too short", size)
	}
	if size := len([]rune(u.Password)); size > MAX_PASSWORD_SIZE {
		return fmt.Errorf("user password of size %d is too long", size)
	}
	return nil
}

func (u *User) isValidEmail() (err error) {
	if u.Email == "" {
		return fmt.Errorf("must provide an user email")
	}
	re := regexp.MustCompile(EMAIL_REGEX)
	// Use the MatchString method to check if the email matches the pattern
	if !re.MatchString(u.Email) {
		return fmt.Errorf("invalid user email: %s", u.Email)
	}
	return nil
}

func (u *User) cleanPhone() {
	for _, pattern := range PHONE_FILTER {
		u.Phone = strings.ReplaceAll(u.Phone, pattern, "")
	}
}

func (u *User) isValidPhone() (err error) {
	if u.Phone == "" {
		return nil
	}
	u.cleanPhone()
	phone := u.Phone
	for _, number := range phone {
		if _, err := strconv.ParseUint(string(number), 10, 64); err != nil {
			return fmt.Errorf("invalid user phone number: %s", phone)
		}
	}
	return nil
}

func (u *User) hashPassword() error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) cleanName() {
	titler := cases.Title(language.English)
	u.FirstName = titler.String(u.FirstName)
	u.LastName = titler.String(u.LastName)
}

func (u *User) BeforeCreate(tx *gorm.DB) error {
	if err := u.isValid(); err != nil {
		return err
	}
	u.cleanName()
	if err := u.hashPassword(); err != nil {
		return err
	}
	return nil
}

func (u *User) CheckPassword(password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)); err != nil {
		return err
	}
	return nil
}
