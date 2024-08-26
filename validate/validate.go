package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"mygoshop/common"
	"mygoshop/config"
	"mygoshop/datamodels"
	"mygoshop/rabbitmq"
	"net/http"
	"net/url"
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
)

type AccessControl struct {
	sourceArray map[int]time.Time
	// 保证 sourceArray 的并发安全
	sync.RWMutex
}

var (
	localHost        = ""
	rabbitMqValidate *rabbitmq.RabbitMQ
	hashRing         *common.HashRing
	accessControl    = &AccessControl{
		sourceArray: make(map[int]time.Time),
	}
	// redis 全局句柄
	rdb *redis.Client
)

// 获取接入用户的信息
func (m *AccessControl) Get(uid int) time.Time {
	m.RLock()
	defer m.RUnlock()
	return m.sourceArray[uid]
}

// 设置接入用户的信息
func (m *AccessControl) Put(uid int) {
	m.Lock()
	defer m.Unlock()
	m.sourceArray[uid] = time.Now()
}

func Auth(res http.ResponseWriter, req *http.Request) error {
	return cookieCheck(req)
}

func cookieCheck(req *http.Request) error {
	uidCookie, err := req.Cookie("uid")
	if err != nil {
		return errors.New("user is not logged in")
	}
	signCookie, err := req.Cookie("sign")
	if err != nil {
		return errors.New("failed to obtain the user encryption string")
	}
	__signString, err := common.EnPwdCode([]byte(uidCookie.Value))
	if err != nil {
		return errors.New("the encryption string has been tampered with")
	}
	if CheckIdInfo(signCookie.Value, __signString) {
		return nil
	}
	return errors.New("identity verification failed")
}

// 自定义逻辑判断
func CheckIdInfo(checkStr string, signStr string) bool {
	return checkStr == signStr
}

// 抢购处理
// http://localhost:8083/onsale?productID=1 cookie
func OnSaleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("running OnSale")
	// 获得 productID
	queryForm, err := url.ParseQuery(r.URL.RawQuery)
	if err != nil || len(queryForm["productID"]) <= 0 {
		w.Write([]byte("query false"))
		return
	}

	productString := queryForm["productID"][0]

	// 获取用户 cookie
	userCookie, err := r.Cookie("uid")
	if err != nil {
		w.Write([]byte("cookie false"))
		return
	}

	// 获取商品ID
	productID, err := strconv.ParseInt(productString, 10, 64)
	if err != nil {
		w.Write([]byte("Parse productID false"))
		return
	}

	// 获取用户ID
	userID, err := strconv.ParseInt(userCookie.Value, 10, 64)
	if err != nil {
		w.Write([]byte("Parse userID false"))
		return
	}

	// 1.分布式权限验证
	right := accessControl.CheckPermission(r)
	if !right {
		w.Write([]byte("GetDistributedRight false"))
		return
	}

	// 2.获取数量控制权限，防止秒杀出现超卖现象
	result := BuyProduct(userID, productID, r)
	w.Write(result)
}
func BuyProduct(userID int64, productID int64, r *http.Request) []byte {
	Success := GetProductFromRedis(strconv.Itoa(int(productID)), strconv.Itoa(1))
	if Success {

		// 1.创建消息体
		message := datamodels.NewMessage(userID, productID, 1)
		// 类型转化
		byteMessage, err := json.Marshal(message)
		if err != nil {
			return []byte("json Marshal false")
		}

		// 2.生产消息
		err = rabbitMqValidate.PublishSimple(string(byteMessage))
		if err != nil {
			return []byte("PublishSimple false")

		}
		return []byte("true")
	}
	return []byte("false")
}

func CheckRightHandler(w http.ResponseWriter, r *http.Request) {
	right := accessControl.CheckPermission(r)
	if !right {
		w.Write([]byte("false"))
		return
	}
	w.Write([]byte("true"))
}

func StartHTTPServer() {
	filter := common.NewFilter()
	filter.RegisterUriFilter("/onsale", Auth)
	filter.RegisterUriFilter("/checkRight", Auth)

	// 2,启动服务
	// 用于验证和访问数量控制服务
	http.HandleFunc("/onsale", filter.Handler(OnSaleHandler))
	// 用于分布式验证
	http.HandleFunc("/checkRight", filter.Handler(CheckRightHandler))
	http.ListenAndServe(config.ValidateSet.Host+":"+config.ValidateSet.Port, nil)
}

func main() {
	hashRing = common.NewHashRing()
	for _, v := range config.ClusterHostArray {
		hashRing.Add(v)
	}

	localIP, err := common.GetLocalIP()
	if err != nil {
		fmt.Println(err)
	}
	localHost = localIP
	fmt.Printf("localHost: %s\n", localHost)

	// rabbitmq
	rabbitMqValidate = rabbitmq.NewRabbitMQSimple("productshop")
	defer rabbitMqValidate.Destroy()

	rdb = redis.NewClient(&redis.Options{
		Addr:     config.RedisSet.Host + ":" + config.RedisSet.Port,
		Password: config.RedisSet.Password,
		DB:       0, // use default DB
	})

	StartHTTPServer()
}
