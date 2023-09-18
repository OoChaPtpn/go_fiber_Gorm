package user

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	_ "github.com/microsoft/go-mssqldb"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error

const DNS = "sqlserver://fx-alnath.cco2.thunghuasinn.com?database=fx_pd_dev" //:9930

type User struct {
	gorm.Model
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"e_mail"`
}

func InitialMigration() {
	DB, err = gorm.Open(sqlserver.Open(DNS), &gorm.Config{})
	if err != nil {
		fmt.Println(err.Error())
		panic("Cannot connect to Database")
	}
	DB.AutoMigrate(&User{})
}

func GetUsers(c *fiber.Ctx) error {
	var users []User
	DB.Find(&users)
	return c.JSON(&users)
}

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	DB.Find(&user, id)
	return c.JSON(&user)
}

func SaveUser(c *fiber.Ctx) error {
	user := new(User)
	if err := c.BodyParser(user); err != nil {
		return c.Status(500).SendString(err.Error())
	} 
	DB.Create(&user)
	return c.JSON(&user)
}

func DeleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	var user User
	//DB.Unscoped().Delete(&user, id)
	DB.First(&user, id)
	if user.Email == "" {
		return c.Status(500).SendString("User not available")
	}

	DB.Delete(&user)
	return c.SendString("User is deleted!!!")
}

func UpdateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user := new(User)
	DB.First(&user, id)
	if user.Email == "" {
		return c.Status(500).SendString("User not available")
	}
	if err := c.BodyParser(&user); err != nil {
		return c.Status(500).SendString(err.Error())
	}
	DB.Save(&user)
	return c.JSON(&user)
}
