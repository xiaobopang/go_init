package controllers

/*
 * @Script: test.go
 * @Author: pangxiaobo
 * @Email: 10846295@qq.com
 * @Create At: 2018-11-06 14:50:15
 * @Last Modified By: pangxiaobo
 * @Last Modified At: 2018-12-11 14:39:12
 * @Description: This is description.
 */

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go_init/helpers"
	"github.com/go_init/libs"
	"github.com/go_init/models"
	"strconv"
	"time"
)

type TestController struct{}

func (t *TestController) GetNick(c *gin.Context) {
	nickname := c.DefaultQuery("nick", "guest")
	c.JSON(200, gin.H{
		"code":      200,
		"data":      map[string]string{"nickname": nickname},
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}

//获取用户
func (t *TestController) GetUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Query("id"))
	fmt.Println(id)

	res, _ := models.GetUserById(id)

	c.JSON(200, gin.H{
		"code":      200,
		"data":      res,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}

//新增用户
func (t *TestController) AddUser(c *gin.Context) {

	name := c.PostForm("name")
	password := helpers.EncodeMD5(c.PostForm("password"))
	age, _ := strconv.Atoi(c.DefaultPostForm("age", "20"))
	gender, _ := strconv.Atoi(c.DefaultPostForm("gender", "1"))
	email := c.PostForm("email")

	res := models.AddUser(name, password, age, gender, email)

	c.JSON(200, gin.H{
		"code":      200,
		"data":      res,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}

//删除用户 (硬删除)
func (t *TestController) DelUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.Query("id"))
	fmt.Println(id)

	res := models.DelUser(id)

	c.JSON(200, gin.H{
		"code":      200,
		"data":      res,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}

//更新
func (t *TestController) UptUser(c *gin.Context) {

	id, _ := strconv.Atoi(c.PostForm("id"))
	data := make(map[string]interface{})

	data["username"] = c.PostForm("name")
	data["password"] = helpers.EncodeMD5(c.PostForm("password"))
	data["age"], _ = strconv.Atoi(c.DefaultPostForm("age", "20"))
	data["gender"], _ = strconv.Atoi(c.DefaultPostForm("gender", "1"))
	data["email"] = c.PostForm("email")
	data["updated_at"] = time.Now().Unix()

	res := models.UptUser(id, data)

	c.JSON(200, gin.H{
		"code":      200,
		"data":      res,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}

//Redis 测试
func (t *TestController) RedisTest(c *gin.Context) {
	redisKey := c.Query("redisKey")
	fmt.Println(redisKey)
	userInfo, err := libs.GetKey(redisKey)
	if err != nil {
		data := make(map[string]interface{})
		data["username"] = "jack"
		data["age"] = 22
		data["gender"] = "man"
		data["email"] = "test@test.com"
		data["updated_at"] = time.Now().Unix()
		userInfo, err := json.Marshal(data)
		libs.SetKey(redisKey, userInfo, 3600)
		if err != nil {
			fmt.Println(err)
		}
	}
	c.JSON(200, gin.H{
		"code":      200,
		"data":      userInfo,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}
