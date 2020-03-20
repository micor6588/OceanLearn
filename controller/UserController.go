package controller

import (
	"log"
	"micor/ginessential/common"
	"micor/ginessential/model"
	"micor/ginessential/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

// Register 处理用户注册
func Register(Ctx *gin.Context) {
	DB := common.GetDB() //获取数据库数据
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
		name = util.RandomString(10)
		return
	}
	log.Println(name, telephone, password)
	//判断手机号是否存在
	if isTelephoneExitst(DB, telephone) {
		Ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户已经存在"})
		return
	}
	//创建用户
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  password,
	}
	DB.Create(&newUser)
	//返回结果
	Ctx.JSON(200, gin.H{
		"message": "注册成功",
	})
}

// isTelephoneExitst 与数据库中的数据，验证手机号是否存在
func isTelephoneExitst(db *gorm.DB, telephone string) bool {
	var user model.User
	db.Where("telephone = ?", telephone).First(&user)
	if user.ID != 0 {
		return true
	}
	return false
}
