package config

// fonted setting
var FontedSet = &__webSetting{
	Host: "127.0.0.1",
	Port: "8080",
}

// backend setting
var BackendSet = &__webSetting{
	Host: "127.0.0.1",
	Port: "8081",
}

var ValidateSet = &__webSetting{
	Host: "127.0.0.1",
	Port: "8082",
}

var (
	//集群地址
	ClusterHostArray = []string{"172.18.105.29", "172.18.105.29"}
	// 每个用户抢购间隔时间
	Interval = 5
)

// rabbitmq setting
// url：amqp://user:password @RabbitMQ_host:port/Vhost
var RMQSet = &__rmqSetting{
	Host:     "172.29.125.42",
	Port:     "15672",
	Vhost:    "productshop",
	User:     "admin",
	Password: "rabbitmq..2233",
}

//redis setting
//user:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local
var RedisSet = &__mysqlSetting{
	Host:     "172.29.125.42",
	Port:     "6379",
	Password: "123",
}

//mysql setting
//user:password@tcp(host:port)/dbname?charset=utf8&parseTime=True&loc=Local
var SQLSet = &__mysqlSetting{
	Host:     "172.29.125.42",
	Port:     "3306",
	Dbname:   "productshop",
	User:     "user_win",
	Password: "123",
}

type __rmqSetting struct {
	Host     string
	Port     string
	Vhost    string
	User     string
	Password string
}

type __mysqlSetting struct {
	Host     string
	Port     string
	Dbname   string
	User     string
	Password string
}

type __webSetting struct {
	Host string
	Port string
}
