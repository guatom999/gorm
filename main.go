package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

type SqlLogger struct {
	logger.Interface
}

var db *gorm.DB

func (l SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	sql, _ := fc()
	fmt.Printf("%v \n ==================== \n", sql)
}

func main() {
	dsn := "root:Bossza555@tcp(127.0.0.1:3306)/goorm?parseTime=true"
	disl := mysql.Open(dsn)
	var err error
	db, err = gorm.Open(disl, &gorm.Config{
		Logger: &SqlLogger{},
		DryRun: false,
	})
	if err != nil {
		panic(err)
	}
	// Normal Migrate
	// db.Migrator().CreateTable(Customer{})

	// db.AutoMigrate(Gender{}, Test{}, Customer{})

	// CreateGender("xxxx")
	// UpdateGender2(6, "")
	// GetGender(10)
	// GetGenderByName("Male")
	// DeleteGender(6)
	// CreateTest(0, "Boss")
	// CreateTest(0, "Chon")
	// CreateTest(0, "Toy")
	// CreateTest(1, "Tong")
	// DeleteTest(2)
	// GetTests()

	// CreateCustomer("Boss", 1)
	// CreateCustomer("Chon", 1)
	// CreateCustomer("Toy", 1)
	// GetCustomers()
	UpdateGender2(1, "ชายแท้")
}

func CreateCustomer(name string, genderID uint) {
	customer := Customer{
		Name:     name,
		GenderID: genderID,
	}
	tx := db.Create(&customer)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(customer)
}

func GetCustomers() {
	customers := []Customer{}
	tx := db.Preload(clause.Associations).Find(&customers)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	for _, customer := range customers {
		fmt.Printf("%v | %v | %v \n", customer.ID, customer.Name, customer.Gender.Name)
	}
	// fmt.Println(customers)
}

type Customer struct {
	ID       uint
	Name     string
	Gender   Gender
	GenderID uint
}

type Gender struct {
	ID   uint
	Name string `gorm:"size(10);unique;default:Hello"`
}

func CreateTest(code uint, name string) {
	test := Test{Code: code, Name: name}
	db.Create(&test)
}

func GetTests() {
	tests := []Test{}
	db.Find(&tests)
	for _, t := range tests {
		fmt.Printf("%v|%v", t.ID, t.Name)
	}
}

func DeleteTest(id uint) {
	db.Delete(&Test{}, id)
}

func DeleteGender(id uint) {
	tx := db.Delete(&Gender{}, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println("Deleted")
	GetGender(id)
}

func UpdateGender(id uint, name string) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	gender.Name = name
	fmt.Println(gender.Name)
	tx = db.Save(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)
	// fmt.Println(gender)
}

func UpdateGender2(id uint, name string) {
	gender := Gender{Name: name}
	tx := db.Model(&Gender{}).Where("id=@getID", sql.Named("getID", id)).Updates(gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	GetGender(id)
}

func GetGenders() {
	gender := []Gender{}
	tx := db.Order("id").Find(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	// result := db.First(&gender)
	fmt.Println(gender)
}

func GetGender(id uint) {
	gender := Gender{}
	tx := db.First(&gender, id)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func GetGenderByName(name string) {
	gender := []Gender{}
	tx := db.Where("name=?", name).Find(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}
	fmt.Println(gender)
}

func CreateGender(name string) {
	gender := Gender{
		Name: name,
	}
	tx := db.Create(&gender)
	if tx.Error != nil {
		fmt.Println(tx.Error)
		return
	}

	fmt.Println(gender)

}

type Test struct {
	gorm.Model
	Code uint   `gorm:"comment:This is Code"`
	Name string `gorm:"column:myname;size:20;default:Hello;not null"`
}

// func (l Test) TableName() string {
// 	return "BossTest"
// }
