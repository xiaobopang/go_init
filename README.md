## Introduce

    这是一个基于gin搭建的一个包含gorm, goredis,rabbitmq,websocket等操作相关操作的项目结构。

    主要提供一些库和组件的实现案例，以及项目开发部署，发布，执行等流程。纯属个人兴趣，学习整理过程
    
    如有发现不合理的地方希望可以大家可以提出建议和指正。

    通过 go build -o go_init main.go 来生成执行文件
    
    启动服务：./run.sh start

    停止服务：./run.sh stop

    注意：有可能出现没有执行权限的情况，执行 sudo chmod +x run.sh来解决



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
                "errors"
                "fmt"
                "github.com/gorilla/websocket"
                "net/http"
                "sync"
                "time"
        )

        var wsUpgrader = websocket.Upgrader{
                ReadBufferSize:    4096,
                WriteBufferSize:   4096,
                EnableCompression: true,
                HandshakeTimeout:  5 * time.Second,
                // CheckOrigin: 处理跨域问题，线上环境慎用
                CheckOrigin: func(r *http.Request) bool {
                        return true
                },
        }

        // 客户端读写消息
        type wsMessage struct {
                messageType int
                data        []byte
        }
        type Wscontroller struct{}

        // 客户端连接
        type wsConnection struct {
                wsSocket *websocket.Conn // 底层websocket
                inChan   chan *wsMessage // 读队列
                outChan  chan *wsMessage // 写队列

                mutex     sync.Mutex // 避免重复关闭管道
                isClosed  bool
                closeChan chan byte // 关闭通知
        }

        func (wsConn *wsConnection) wsReadLoop() {
                for {
                        // 读一个message
                        msgType, data, err := wsConn.wsSocket.ReadMessage()
                        if err != nil {
                                goto error
                        }
                        req := &wsMessage{
                                msgType,
                                data,
                        }
                        // 放入请求队列
                        select {
                        case wsConn.inChan <- req:
                        case <-wsConn.closeChan:
                                goto closed
                        }
                }
        error:
                wsConn.wsClose()
        closed:
        }
        func (wsConn *wsConnection) wsWriteLoop() {
                for {
                        select {
                        // 取一个应答
                        case msg := <-wsConn.outChan:
                                // 写给websocket
                                if err := wsConn.wsSocket.WriteMessage(msg.messageType, msg.data); err != nil {
                                        goto error
                                }
                        case <-wsConn.closeChan:
                                goto closed
                        }
                }
        error:
                wsConn.wsClose()
        closed:
        }

        func (wsConn *wsConnection) procLoop() {
                // 启动一个gouroutine发送心跳
                go func() {
                        for {
                                time.Sleep(2 * time.Second)
                                if err := wsConn.wsWrite(websocket.TextMessage, []byte("heartbeat from server")); err != nil {
                                        fmt.Println("heartbeat fail")
                                        wsConn.wsClose()
                                        break
                                }
                        }
                }()

                // 这是一个同步处理模型（只是一个例子），如果希望并行处理可以每个请求一个gorutine，注意控制并发goroutine的数量!!!
                for {
                        msg, err := wsConn.wsRead()
                        if err != nil {
                                fmt.Println("read fail")
                                break
                        }
                        fmt.Println(string(msg.data))
                        err = wsConn.wsWrite(msg.messageType, msg.data)
                        if err != nil {
                                fmt.Println("write fail")
                                break
                        }
                }
        }

        func (w *Wscontroller) WsHandler(resp http.ResponseWriter, req *http.Request) {
                // 应答客户端告知升级连接为websocket
                wsSocket, err := wsUpgrader.Upgrade(resp, req, nil)
                if err != nil {
                        return
                }
                wsConn := &wsConnection{
                        wsSocket:  wsSocket,
                        inChan:    make(chan *wsMessage, 1000),
                        outChan:   make(chan *wsMessage, 1000),
                        closeChan: make(chan byte),
                        isClosed:  false,
                }

                // 处理器
                go wsConn.procLoop()
                // 读协程
                go wsConn.wsReadLoop()
                // 写协程
                go wsConn.wsWriteLoop()
        }

        func (wsConn *wsConnection) wsWrite(messageType int, data []byte) error {
                select {
                case wsConn.outChan <- &wsMessage{messageType, data}:
                case <-wsConn.closeChan:
                        return errors.New("websocket closed")
                }
                return nil
        }

        func (wsConn *wsConnection) wsRead() (*wsMessage, error) {
                select {
                case msg := <-wsConn.inChan:
                        return msg, nil
                case <-wsConn.closeChan:
                }
                return nil, errors.New("websocket closed")
        }

        func (wsConn *wsConnection) wsClose() {
                wsConn.wsSocket.Close()

                wsConn.mutex.Lock()
                defer wsConn.mutex.Unlock()
                if !wsConn.isClosed {
                        wsConn.isClosed = true
                        close(wsConn.closeChan)
                }
        }

        前端通过访问 ws://localhost:7777/ws 即可与服务端建立websocket连接。

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
