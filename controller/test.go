package controller

/*
 * @Script: test.go
 * @Author: pangxiaobo
 * @Email: 10846295@qq.com
 * @Create At: 2018-11-06 14:50:15
 * @Last Modified By: pangxiaobo
 * @Last Modified At: 2018-12-12 14:25:46
 * @Description: This is description.
 */

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/elastic/go-elasticsearch"
	"github.com/gin-gonic/gin"
	"github.com/xiaobopang/go_init/helper"
	"github.com/xiaobopang/go_init/lib"
	"github.com/xiaobopang/go_init/model"
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

	res, _ := model.GetUserById(id)

	c.JSON(200, gin.H{
		"code":      200,
		"data":      res,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}

//获取用户
func (t *TestController) UserList(c *gin.Context) {
	keyword := c.Query("keyword")
	pageNo := c.GetInt("page_number")
	pageSize := c.GetInt("page_size")

	res := model.UsersList(pageNo, pageSize, "username = ?", keyword)

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
	password := helper.Md5(c.PostForm("password"))
	age, _ := strconv.Atoi(c.DefaultPostForm("age", "20"))
	gender, _ := strconv.Atoi(c.DefaultPostForm("gender", "1"))
	email := c.PostForm("email")

	res := model.AddUser(name, password, age, gender, email)

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

	res := model.DelUser(id)

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
	data["password"] = helper.Md5(c.PostForm("password"))
	data["age"], _ = strconv.Atoi(c.DefaultPostForm("age", "20"))
	data["gender"], _ = strconv.Atoi(c.DefaultPostForm("gender", "1"))
	data["email"] = c.PostForm("email")
	data["updated_at"] = time.Now().Unix()

	res := model.UptUser(id, data)

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

	userInfo, err := lib.GetKey(redisKey)
	if err != nil {
		data := make(map[string]interface{})
		data["username"] = "jack"
		data["age"] = 22
		data["gender"] = "man"
		data["email"] = "test@test.com"
		data["updated_at"] = time.Now().Unix()
		userInfo, err := json.Marshal(data)
		lib.SetKey(redisKey, userInfo, 3600)
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

func (t *TestController) GetToken(c *gin.Context) {
	token, err := lib.GenerateToken(1, "pang@pang.com")
	if err != nil {
		fmt.Println("err: ", err)
	}
	c.JSON(200, gin.H{
		"code":      200,
		"data":      token,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})

}

// es 测试
func (t *TestController) ES(c *gin.Context) {
	cfg := elasticsearch.Config{
		Addresses: []string{
			"http://127.0.0.1:9201",
		},
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		fmt.Println("elasticsearch has error: ", err)
	}
	// 1. Get cluster info
	//
	res, err := es.Info()
	if err != nil {
		fmt.Println("Error getting response: %s", err)
	}
	c.JSON(200, gin.H{
		"code":      200,
		"data":      res,
		"msg":       "success",
		"timestamp": time.Now().Unix(),
	})
}
