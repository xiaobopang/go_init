## Introduce

    这是一个基于gin框架和其他一些项目组织的一个包含mysql, redis,rabbitmq,websocket等操作的一个项目结构。

    主要目的是希望为golang web开发入门者提供一个学习和简单项目组织结构。更加方便进行开发和一些核心组件的研究，

    以及项目开发完成后的部署，发布，执行等流程。

    可通过 go build -o go_init main.go 来生成执行文件，然后执行：./run.sh start 即可启动相关服务。

    要想停止服务可执行：./run.sh stop

    注意⚠️：有可能出现没有执行权限的情况：可运行: sudo chmod +x run.sh来解决



### Router 示例

```
        package routers

        import (
                "github.com/gin-gonic/gin"
                "github.com/go_init/controllers"
                "github.com/go_init/middleware"
                "net/http"
        )

        var indexCtl = new(controllers.IndexController)
        var testCtl = new(controllers.TestController)

        func SetupRouter() *gin.Engine {

                router := gin.Default()
                router.Use(gin.Recovery())
                //router.Use(gin.Logger())

                router.GET("/", indexCtl.Welcome)
                router.GET("/redirect", func(c *gin.Context) {
                        c.Redirect(http.StatusMovedPermanently, "https://www.unclepang.com/")
                })
                router.NoRoute(indexCtl.Handle404)

                v1 := router.Group("/v1")
                v1.Use(middleware.CORS(middleware.CORSOptions{}))
                {
                        v1.GET("/test", testCtl.GetNick)
                }

                v2 := router.Group("/v2")
                v2.Use(middleware.CORS(middleware.CORSOptions{}))
                {
                        v2.GET("/user", testCtl.GetUser)
                        v2.POST("/user", testCtl.AddUser)
                        v2.DELETE("/user", testCtl.DelUser)
                        v2.PATCH("/user", testCtl.UptUser)
                }

                return router
        }

```

### Request and Response 示例

```
        package controllers

        import (
                "fmt"
                "github.com/gin-gonic/gin"
                "github.com/go_init/helpers"
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

                res := models.GetUserById(id)

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
                data["age"] = c.PostForm("age")
                data["gender"] = c.PostForm("gender")
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

```

### model 示例

```
        package models

        import (
                "time"
        )

        type User struct {
                ID        int    `gorm:"primary_key" json:"id"`
                Username  string `json:"username"`
                Password  string `json:"password"`
                Age       int    `json:"age"`
                Email     string `json:"email"`
                Gender    int    `json:"gender"`
                CreatedAt int64  `json:"created_at"`
                UpdatedAt int64  `json:"updated_at"`
        }

        func GetUserById(id int) *User {
                var user User
                DB.First(&user, "id = ?", id)
                return &user
        }

        func AddUser(name string, password string, age int, gender int, email string) error {
                user := User{
                        Username:  name,
                        Password:  password,
                        Age:       age,
                        Gender:    gender,
                        Email:     email,
                        CreatedAt: time.Now().Unix(),
                }
                if err := DB.Create(&user).Error; err != nil {
                        return err
                }
                return nil
        }

        func DelUser(id int) error {
                if err := DB.Where("id = ?", id).Delete(&User{}).Error; err != nil {
                        return err
                }

                return nil
        }

        func UptUser(id int, data interface{}) error {

                if err := DB.Model(&User{}).Where("id = ? AND is_deleted = ? ", id, 0).Updates(data).Error; err != nil {
                        return err
                }

                return nil
        }

```

### websocket 示例

```

        package controllers

        import (
                "fmt"
                "github.com/gorilla/websocket"
                "net/http"
                "time"
        )

        type WsController struct{}

        var wsupgrader = websocket.Upgrader{
                ReadBufferSize:    4096,
                WriteBufferSize:   4096,
                EnableCompression: true,
                HandshakeTimeout:  5 * time.Second,
                // CheckOrigin: 处理跨域问题，线上环境慎用
                CheckOrigin: func(r *http.Request) bool {
                        return true
                },
        }

        func (ws *WsController) WsHandler(w http.ResponseWriter, r *http.Request) {
                conn, err := wsupgrader.Upgrade(w, r, nil)
                if err != nil {
                        http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
                        return
                }
                go echo(conn)
        }

        func echo(conn *websocket.Conn) {
                for {
                        msgType, msg, err := conn.ReadMessage()
                        if err != nil {
                                fmt.Println(err)
                                return
                        }
                        if string(msg) == "ping" {
                                fmt.Println("ping")
                                time.Sleep(time.Second * 2)
                                err = conn.WriteMessage(msgType, []byte("pong"))
                                if err != nil {
                                        fmt.Println(err)
                                        return
                                }
                        } else {
                                conn.Close()
                                fmt.Println(string(msg))
                                return
                        }
                }
        }

```


## 如果你使用的是MacOS,那么Mac下编译Linux, Windows平台的64位可执行程序如下：


#### 编译Linux服务器可执行文件：
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o go_init main.go

#### 编译Windows服务器可执行文件：
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o go_init main.go



## 如果你使用的是Linux系统，那么Linux下编译Mac, Windows平台的64位可执行程序如下：

#### 编译MacOS可执行文件：
        CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o go_init main.go

#### 编译windows下可执行文件：
        CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o go_init main.go



## 如果你使用的是Windows系统，那么Windows下编译Mac, Linux平台的64位可执行程序如下：

#### 编译MacOS可执行文件如下：

        SET CGO_ENABLED=0

        SET GOOS=darwin

        SET GOARCH=amd64

        go build -o go_init main.go

#### 编译Windows可执行文件如下：

        SET CGO_ENABLED=0

        SET GOOS=linux

        SET GOARCH=amd64

        go build -o go_init main.go



## Nginx负载均衡


```

        user nginx;
        worker_processes auto;
        error_log /var/log/nginx/error.log;
        pid /var/run/nginx.pid;

        events {
        worker_connections 1024;
        }

        http {
        log_format  main  '$remote_addr - $remote_user [$time_local] "$request" '
                        '$status $body_bytes_sent "$http_referer" '
                        '"$http_user_agent" "$http_x_forwarded_for"';
        
        access_log  /var/log/nginx/access.log  main;
        
        sendfile            on;
        tcp_nopush          on;
        tcp_nodelay         on;
        keepalive_timeout   65;
        types_hash_max_size 2048;
        
        include             /etc/nginx/mime.types;
        default_type        application/octet-stream;
                
        index   index.html index.htm;
        
        upstream docker_nginx {
                ip_hash; #同一个ip一定时间内负载到一台机器
                server 172.31.0.155:8081;
                server 172.31.0.155:8082;
                server 172.31.0.155:8083;
                server 172.31.0.155:8084;
        }
        
        server {
                # 使用openssl自建的rsa证书
                ssl_certificate /opt/ssl/nginx.unclepang.com.crt;
                ssl_certificate_key /opt/ssl/nginx.unclepang.com.key;
                ssl_session_timeout 5m;
                ssl_ciphers ECDHE-RSA-AES128-GCM-SHA256:ECDHE:ECDH:AES:HIGH:!NULL:!aNULL:!MD5:!ADH:!RC4;
                ssl_protocols TLSv1 TLSv1.1 TLSv1.2;
                ssl_prefer_server_ciphers on;
        
                listen 443;
                ssl on;
                server_name www.unclepang.com;
                
                location / {
                        # 代理到真实机器，如果真实机器也安装了https则使用https
                        # 一般代理集群对流量进行了https后，真实机器可不再使用https
                        proxy_pass http://docker_nginx;
                }
        }
        }

```
