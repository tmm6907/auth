package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tmm6907/auth/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initDB(path string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(path), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		return nil, err
	}

	// synchronous updates
	syncResult := db.Exec("PRAGMA synchronous = NORMAL")
	if syncResult.Error != nil {
		return nil, syncResult.Error
	}

	err = db.AutoMigrate(
		&models.User{},
		&models.Role{},
	)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {

	db, err := initDB("db/auth.db")

	if err != nil {
		log.Fatalln(err)
	}
	server := gin.Default()

	server.GET("/", func(ctx *gin.Context) {
		var count int64
		addr := models.Address{
			StreetNumber: "34",
			StreetName:   "Aspen St.",
			Suite:        "300",
			City:         "Washingon",
			State:        "DC",
			ZipCode:      "20030",
		}
		company := models.NewCompany("The Company", addr)

		res := db.Model(&models.Company{}).Create(&company)
		if res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  fmt.Sprint(res.Error),
			})
			return
		}

		department1 := models.NewDepartment("Human Resources", company.ID)
		department2 := models.NewDepartment("Legal", company.ID)
		department3 := models.NewDepartment("IT", company.ID)
		department4 := models.NewDepartment("Facilities", company.ID)
		company.Departments = []models.Department{
			*department1, *department2, *department3, *department4,
		}
		for _, dep := range company.Departments {
			res = db.Model(&models.Department{}).Create(&dep)
			if res.Error != nil {
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"status": http.StatusBadRequest,
					"error":  fmt.Sprint(res.Error),
				})
				return
			}
		}

		user := models.NewUser(
			models.UserConfig{
				FirstName: "Test",
				LastName:  "User3",
				Username:  "tuser003",
				Password:  "hello world!!",
				Email:     "tmm6907@gmail.com",
			},
			company.ID,
			department1.ID,
		)

		res = db.Model(&models.User{}).Create(&user)
		if res.Error != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"status": http.StatusBadRequest,
				"error":  fmt.Sprint(res.Error),
			})
			return
		}
		res1 := user.CheckPassword("hello world")
		res2 := user.CheckPassword("hello world!")
		res3 := user.CheckPassword("hello world!!")
		var check1 bool
		var check2 bool
		var check3 bool
		if res1 == nil {
			check1 = true
		}
		if res2 == nil {
			check2 = true
		}
		if res3 == nil {
			check3 = true
		}
		db.Model(&models.User{}).Count(&count)
		ctx.JSON(http.StatusOK, gin.H{
			"message":         "hello world! " + strconv.Itoa(int(count)) + " count",
			"validPassCheck1": check1,
			"validPassCheck2": check2,
			"validPassCheck3": check3,
		})
	})
	server.Run(":3030")
}
