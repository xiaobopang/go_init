## Introduce

    这是一个基于gin搭建的一个包含gorm, goredis,rabbitmq,websocket等操作相关操作的项目结构。

    主要提供一些库和组件的实现案例，以及项目开发部署，发布，执行等流程。纯属个人兴趣，学习整理过程，
    
    如有发现不合理的地方希望大家可以提出建议和指正。

    通过执行 go get -u github.com/xiaobopang/go_init 来下载安装该项目，使用 govendor sync来拉取依赖到vendor目录，

    注意vendor所在的目录一定要在第一个GOPAHT下。

    通过 go build -o go_init main.go 来生成执行文件
    
    启动服务：./run.sh start

    停止服务：./run.sh stop

    注意：有可能出现没有执行权限的情况，执行 sudo chmod +x run.sh来解决



### Router 示例

```
        package routers

        import (
                "net/http"

                "github.com/gin-gonic/gin"
                "github.com/go_init/controllers"
                "github.com/go_init/middleware"
        )

        var indexCtl = new(controllers.IndexController)
        var testCtl = new(controllers.TestController)
        var wsCtl = new(controllers.WsController)
        var mqCtl = new(controllers.MqController)

        func SetupRouter() *gin.Engine {
                router := gin.Default()
                router.Use(gin.Recovery())
                //router.Use(gin.Logger())

                router.GET("/", indexCtl.Welcome)
                router.NoRoute(indexCtl.Handle404)
                router.GET("/redis", testCtl.RedisTest) //redis测试

                router.GET("/redirect", func(c *gin.Context) {
                        c.Redirect(http.StatusMovedPermanently, "https://www.unclepang.com/")
                })
                router.POST("/exchange", func(c *gin.Context) {
                        mqCtl.ExchangeHandler(c.Writer, c.Request)
                })
                router.POST("/queue/bind", func(c *gin.Context) {
                        mqCtl.QueueBindHandler(c.Writer, c.Request)
                })
                router.GET("/queue", func(c *gin.Context) {
                        mqCtl.QueueHandler(c.Writer, c.Request)
                }) //consume queue
                router.POST("/queue", func(c *gin.Context) {
                        mqCtl.QueueHandler(c.Writer, c.Request)
                }) //declare queue
                router.DELETE("/queue", func(c *gin.Context) {
                        mqCtl.QueueHandler(c.Writer, c.Request)
                }) //delete queue
                router.POST("/publish", func(c *gin.Context) {
                        mqCtl.PublishHandler(c.Writer, c.Request)
                })
                router.GET("/ws", func(c *gin.Context) {
                        wsCtl.WsHandler(c.Writer, c.Request)
                })

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

#####  前端通过访问 ws://localhost:7777/ws 即可与服务端建立websocket连接

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

```



### RabbitMQ 示例

```
        package controllers

        import (
                "encoding/json"
                "flag"
                "fmt"
                "io/ioutil"
                "log"
                "net/http"

                "github.com/streadway/amqp"
        )

        var (
                amqpUri = flag.String("amqp", "amqp://guest:guest@127.0.0.1:5672/", "amqp uri")
        )

        func init() {
                flag.Parse()
        }

        type MqController struct{}

        // Entity for HTTP Request Body: Message/Exchange/Queue/QueueBind JSON Input
        type MessageEntity struct {
                Exchange     string `json:"exchange"`
                Key          string `json:"key"`
                DeliveryMode uint8  `json:"deliverymode"`
                Priority     uint8  `json:"priority"`
                Body         string `json:"body"`
        }

        type ExchangeEntity struct {
                Name       string `json:"name"`
                Type       string `json:"type"`
                Durable    bool   `json:"durable"`
                AutoDelete bool   `json:"autodelete"`
                NoWait     bool   `json:"nowait"`
        }

        type QueueEntity struct {
                Name       string `json:"name"`
                Durable    bool   `json:"durable"`
                AutoDelete bool   `json:"autodelete"`
                Exclusive  bool   `json:"exclusive"`
                NoWait     bool   `json:"nowait"`
        }

        type QueueBindEntity struct {
                Queue    string   `json:"queue"`
                Exchange string   `json:"exchange"`
                NoWait   bool     `json:"nowait"`
                Keys     []string `json:"keys"` // bind/routing keys
        }

        // RabbitMQ Operate Wrapper
        type RabbitMQ struct {
                conn    *amqp.Connection
                channel *amqp.Channel
                done    chan error
        }

        func (r *RabbitMQ) Connect() (err error) {
                r.conn, err = amqp.Dial(*amqpUri)
                if err != nil {
                        log.Printf("[amqp] connect error: %s\n", err)
                        return err
                }
                r.channel, err = r.conn.Channel()
                if err != nil {
                        log.Printf("[amqp] get channel error: %s\n", err)
                        return err
                }
                r.done = make(chan error)
                return nil
        }

        func (r *RabbitMQ) Publish(exchange, key string, deliverymode, priority uint8, body string) (err error) {
                err = r.channel.Publish(exchange, key, false, false,
                        amqp.Publishing{
                                Headers:         amqp.Table{},
                                ContentType:     "text/plain",
                                ContentEncoding: "",
                                DeliveryMode:    deliverymode,
                                Priority:        priority,
                                Body:            []byte(body),
                        },
                )
                if err != nil {
                        log.Printf("[amqp] publish message error: %s\n", err)
                        return err
                }
                return nil
        }

        func (r *RabbitMQ) DeclareExchange(name, typ string, durable, autodelete, nowait bool) (err error) {
                err = r.channel.ExchangeDeclare(name, typ, durable, autodelete, false, nowait, nil)
                if err != nil {
                        log.Printf("[amqp] declare exchange error: %s\n", err)
                        return err
                }
                return nil
        }

        func (r *RabbitMQ) DeleteExchange(name string) (err error) {
                err = r.channel.ExchangeDelete(name, false, false)
                if err != nil {
                        log.Printf("[amqp] delete exchange error: %s\n", err)
                        return err
                }
                return nil
        }

        func (r *RabbitMQ) DeclareQueue(name string, durable, autodelete, exclusive, nowait bool) (err error) {
                _, err = r.channel.QueueDeclare(name, durable, autodelete, exclusive, nowait, nil)
                if err != nil {
                        log.Printf("[amqp] declare queue error: %s\n", err)
                        return err
                }
                return nil
        }

        func (r *RabbitMQ) DeleteQueue(name string) (err error) {
                // TODO: other property wrapper
                _, err = r.channel.QueueDelete(name, false, false, false)
                if err != nil {
                        log.Printf("[amqp] delete queue error: %s\n", err)
                        return err
                }
                return nil
        }

        func (r *RabbitMQ) BindQueue(queue, exchange string, keys []string, nowait bool) (err error) {
                for _, key := range keys {
                        if err = r.channel.QueueBind(queue, key, exchange, nowait, nil); err != nil {
                                log.Printf("[amqp] bind queue error: %s\n", err)
                                return err
                        }
                }
                return nil
        }

        func (r *RabbitMQ) UnBindQueue(queue, exchange string, keys []string) (err error) {
                for _, key := range keys {
                        if err = r.channel.QueueUnbind(queue, key, exchange, nil); err != nil {
                                log.Printf("[amqp] unbind queue error: %s\n", err)
                                return err
                        }
                }
                return nil
        }

        func (r *RabbitMQ) ConsumeQueue(queue string, message chan []byte) (err error) {
                deliveries, err := r.channel.Consume(queue, "", true, false, false, false, nil)
                if err != nil {
                        log.Printf("[amqp] consume queue error: %s\n", err)
                        return err
                }
                go func(deliveries <-chan amqp.Delivery, done chan error, message chan []byte) {
                        for d := range deliveries {
                                message <- d.Body
                        }
                        done <- nil
                }(deliveries, r.done, message)
                return nil
        }

        func (r *RabbitMQ) Close() (err error) {
                err = r.conn.Close()
                if err != nil {
                        log.Printf("[amqp] close error: %s\n", err)
                        return err
                }
                return nil
        }

        // HTTP Handlers
        func (m *MqController) QueueHandler(w http.ResponseWriter, r *http.Request) {
                if r.Method == "POST" || r.Method == "DELETE" {
                        if r.Body == nil {
                                fmt.Println("missing form body")
                                return
                        }

                        body, err := ioutil.ReadAll(r.Body)
                        if err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }

                        entity := new(QueueEntity)
                        if err = json.Unmarshal(body, entity); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }

                        rabbit := new(RabbitMQ)
                        if err = rabbit.Connect(); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                        defer rabbit.Close()

                        if r.Method == "POST" {
                                if err = rabbit.DeclareQueue(entity.Name, entity.Durable, entity.AutoDelete, entity.Exclusive, entity.NoWait); err != nil {
                                        http.Error(w, err.Error(), http.StatusInternalServerError)
                                        return
                                }
                                w.Write([]byte("declare queue ok"))
                        } else if r.Method == "DELETE" {
                                if err = rabbit.DeleteQueue(entity.Name); err != nil {
                                        http.Error(w, err.Error(), http.StatusInternalServerError)
                                        return
                                }
                                w.Write([]byte("delete queue ok"))
                        }
                } else if r.Method == "GET" {
                        r.ParseForm()
                        rabbit := new(RabbitMQ)
                        if err := rabbit.Connect(); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                        defer rabbit.Close()

                        message := make(chan []byte)

                        for _, name := range r.Form["name"] {
                                if err := rabbit.ConsumeQueue(name, message); err != nil {
                                        fmt.Println("Receive message ", message)
                                        http.Error(w, err.Error(), http.StatusInternalServerError)
                                        return
                                }
                        }

                        w.Write([]byte(""))
                        w.(http.Flusher).Flush()

                        for {
                                fmt.Fprintf(w, "%s\n", <-message)
                                w.(http.Flusher).Flush()
                        }
                } else {
                        w.WriteHeader(http.StatusMethodNotAllowed)
                }
        }

        func (m *MqController) QueueBindHandler(w http.ResponseWriter, r *http.Request) {
                if r.Method == "POST" || r.Method == "DELETE" {
                        body, err := ioutil.ReadAll(r.Body)
                        if err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }

                        entity := new(QueueBindEntity)
                        if err = json.Unmarshal(body, entity); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }

                        rabbit := new(RabbitMQ)
                        if err = rabbit.Connect(); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                        defer rabbit.Close()

                        if r.Method == "POST" {
                                if err = rabbit.BindQueue(entity.Queue, entity.Exchange, entity.Keys, entity.NoWait); err != nil {
                                        http.Error(w, err.Error(), http.StatusInternalServerError)
                                        return
                                }
                                w.Write([]byte("bind queue ok"))
                        } else if r.Method == "DELETE" {
                                if err = rabbit.UnBindQueue(entity.Queue, entity.Exchange, entity.Keys); err != nil {
                                        http.Error(w, err.Error(), http.StatusInternalServerError)
                                        return
                                }
                                w.Write([]byte("unbind queue ok"))
                        }
                } else {
                        w.WriteHeader(http.StatusMethodNotAllowed)
                }
        }

        func (m *MqController) PublishHandler(w http.ResponseWriter, r *http.Request) {
                if r.Method == "POST" {
                        body, err := ioutil.ReadAll(r.Body)
                        if err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }

                        entity := new(MessageEntity)
                        if err = json.Unmarshal(body, entity); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }

                        rabbit := new(RabbitMQ)
                        if err = rabbit.Connect(); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                        defer rabbit.Close()

                        if err = rabbit.Publish(entity.Exchange, entity.Key, entity.DeliveryMode, entity.Priority, entity.Body); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                        w.Write([]byte("publish message ok"))
                } else {
                        w.WriteHeader(http.StatusMethodNotAllowed)
                }
        }

        func (m *MqController) ExchangeHandler(w http.ResponseWriter, r *http.Request) {
                if r.Method == "POST" || r.Method == "DELETE" {
                        body, err := ioutil.ReadAll(r.Body)
                        if err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }

                        entity := new(ExchangeEntity)
                        if err = json.Unmarshal(body, entity); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }

                        rabbit := new(RabbitMQ)
                        if err = rabbit.Connect(); err != nil {
                                http.Error(w, err.Error(), http.StatusInternalServerError)
                                return
                        }
                        defer rabbit.Close()

                        if r.Method == "POST" {
                                if err = rabbit.DeclareExchange(entity.Name, entity.Type, entity.Durable, entity.AutoDelete, entity.NoWait); err != nil {
                                        http.Error(w, err.Error(), http.StatusInternalServerError)
                                        return
                                }
                                w.Write([]byte("declare exchange ok"))
                        } else if r.Method == "DELETE" {
                                if err = rabbit.DeleteExchange(entity.Name); err != nil {
                                        http.Error(w, err.Error(), http.StatusInternalServerError)
                                        return
                                }
                                w.Write([]byte("delete exchange ok"))
                        }
                } else {
                        w.WriteHeader(http.StatusMethodNotAllowed)
                }
        }

```

#### 声明queue，注意：你的request格式必须是 Content-Type:application/json

        http://127.0.0.1:7777/queue POST

        {
                "name":"test",
                "durable":false,
                "autodelete":true,
                "exclusive":false,
                "nowait":false
        }

#### 声明exchange,注意：你的request格式必须是 Content-Type:application/json

        http://127.0.0.1:7777/exchange POST

       {
                "name":"test",
                "type":"direct",
                "durable":false,
                "autodelete":true,
                "nowait":false
        }

#### 绑定queue，注意：你的request格式必须是 Content-Type:application/json

        http://127.0.0.1:7777/queue/bind POST

        {
                "queue":"test",
                "exchange":"test",
                "nowait":false,
                "keys": ["router"]
        }

#### pusblish messge,注意：你的request格式必须是 Content-Type:application/json

        http://127.0.0.1:7777/publish POST

        {
                "exchange":"test",
                "key":"router",
                "deliverymode":2,  //1:不用持久化，2:消息持久化
                "priority":3,  //优先级 （一般情况下1-10之间）
                "body":"this is a message from test."
        }



#### 消费queue

        http://127.0.0.1:7777/queue?name=test   GET





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
