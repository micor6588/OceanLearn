package controller

import (
	"log"
	"micor/ginessential/common"
	"micor/ginessential/model"
	"micor/ginessential/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
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
	hasedPassdord, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		Ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 500, "msg": "密码加密错误"})
		return
	}
	newUser := model.User{
		Name:      name,
		Telephone: telephone,
		Password:  string(hasedPassdord),
	}
	DB.Create(&newUser)
	//返回结果
	Ctx.JSON(200, gin.H{
		"msg":  "注册成功",
		"code": 200,
	})
}

// Login 处理客户端登录相关逻辑
func Login(Ctx *gin.Context) {
	DB := common.GetDB()
	//获取参数
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
	//判断手机号是否存在
	var user model.User
	DB.Where("telephone = ?", telephone).First(&user)
	if user.ID == 0 {
		Ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 422, "msg": "用户不存在"})
		return
	}
	//判断密码是否正确,密码不能明文保存，需要加密后保存
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		Ctx.JSON(http.StatusUnprocessableEntity, gin.H{"code": 400, "msg": "密码错误"})
		return
	}
	//发放token
	token := "11"
	//返回结果
	Ctx.JSON(200, gin.H{
		"msg":  "登录成功",
		"code": 200,
		"data": gin.H{"token": token},
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
