package api

import (
	"One.2/api/middleware"
	"One.2/dao"
	"One.2/model"
	"One.2/utils"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"time"
)

var db *sqlx.DB

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}

func initDB() {
	var err error

	dsn := "root:kuangivan05@tcp(127.0.0.1:3306)/d1215?charset=utf8mb4&parseTime=True&loc=Local"

	db, err = sqlx.Connect("mysql", dsn)
	if err != nil {
		fmt.Printf("connect DB failed, err:%v\n", err)
		return
	}

	fmt.Println("connecting to MySQL...")
	return
}
func reset(c *gin.Context) {
	initDB()
	form := model.User{}
	if err := c.ShouldBind(&form); err != nil {
		fmt.Println(err)
		utils.RespSuccess(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否重复
	flag := dao.SelectUser(username, db)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}

	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username, db)
	// 若不正确则传出错误
	if selectPassword != password {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "wrong password")
		return
	}
	password = c.PostForm("newpassword")
	dao.AddUser(username, password, db)
	utils.RespSuccess(c, "Change password successful")
}

func register(c *gin.Context) {
	initDB()
	form := model.User{}
	if err := c.ShouldBind(&form); err != nil {
		fmt.Println(err)
		utils.RespSuccess(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否重复
	if dao.SelectUser(username, db) {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user already exists")
		return
	}

	dao.AddUser(username, password, db)
	// 以 JSON 格式返回信息
	utils.RespSuccess(c, "add user successful")
}

func login(c *gin.Context) {
	initDB()
	if err := c.ShouldBind(&model.User{}); err != nil {
		utils.RespFail(c, "verification failed")
		return
	}
	// 传入用户名和密码
	username := c.PostForm("username")
	password := c.PostForm("password")

	// 验证用户名是否存在
	flag := dao.SelectUser(username, db)
	// 不存在则退出
	if !flag {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "user doesn't exists")
		return
	}

	// 查找正确的密码
	selectPassword := dao.SelectPasswordFromUsername(username, db)
	// 若不正确则传出错误
	if selectPassword != password {
		// 以 JSON 格式返回信息
		utils.RespFail(c, "wrong password")
		return
	}

	// 正确则登录成功
	// 创建一个我们自己的声明
	claim := model.MyClaims{
		Username: username, // 自定义字段
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 2).Unix(), // 过期时间
			Issuer:    "Yxh",                                // 签发人
		},
	}
	// 使用指定的签名方法创建签名对象
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	// 使用指定的secret签名并获得完整的编码后的字符串token
	tokenString, _ := token.SignedString(middleware.Secret)
	utils.RespSuccess(c, tokenString)
}

func getUsernameFromToken(c *gin.Context) {
	initDB()
	username, _ := c.Get("username")
	utils.RespSuccess(c, username.(string))
}
