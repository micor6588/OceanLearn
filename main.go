package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// User 用户的注册信息
type User struct {
	gorm.Model
	Name      string `gorm:"type:varchar(20);not null"`
	Telephone string `gorm:varchar(110);not null;unique*`
	Password  string `gorm:"size:255;not null"`
}

func main() {
	db := InitDB()
	defer db.Close()
	r := gin.Default()
	r.POST("/api/auth/register", func(Ctx *gin.Context) {
		//获取参数
		name := Ctx.PostForm("name")
		telephone := Ctx.PostForm("telephone")
		password := Ctx.PostForm("password")

		//数据验证
		if len(telephone) != 11 {

			Ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "手机号必须是11位"})
			return
		}
		if len(password) < 6 {
			Ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "密码不能少于6位"})
			return
		}
		//如果名称没有传，就给出一个10位的随机字符串
		if len(name) == 0 {
			name = RandomString(10)
			return
		}
		log.Println(name, telephone, password)
		//判断手机号是否存在
		if isTelephoneExitst(db, telephone) {
			Ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
			return
		}
		//创建用户
		newUser := User{
			Name:      name,
			Telephone: telephone,
			Password:  password,
		}
		db.Create(&newUser)
		//返回结果
		Ctx.JSON(200, gin.H{
			"message": "注册成功",
		})
	})
	r.Run(":9090") // listen and serve on 0.0.0.0:8080
}

//与数据库中的数据，验证手机号是否存在
func isTelephoneExitst(db *gorm.DB, telephone string) bool {
	var user User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}

// RandomString 生成10位随机字符串
func RandomString(n int) string {
	var letters = []byte("fhgqewuighweiofwdhiqehqwfbsfiweufowieferpohpkmbncbkxqyyewoevdffovnd")
	result := make([]byte, n)
	rand.Seed(time.Now().Unix())
	for i := range result {
		result[i] = letters[rand.Intn(len(letters))]
	}
	return string(result)
}

// InitDB 开启连接池
func InitDB() *gorm.DB {
	driverName := "mysql"
	host := "localhost"
	port := "3306"
	database := "ginessential"
	username := "root"
	password := "371871"
	charset := "utf8"
	args := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true",

		username,
		password,
		host,
		port,
		database,
		charset,
	)
	db, err := gorm.Open(driverName, args)
	if err != nil {
		panic("failed to connect database err= " + err.Error())
	}
	//自动创建数据表
	db.AutoMigrate(&User{})
	return db
}
