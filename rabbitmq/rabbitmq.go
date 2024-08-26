package rabbitmq

import (
	"encoding/json"
	"fmt"
	"log"
	"mygoshop/config"
	"mygoshop/datamodels"
	"mygoshop/services"
	"sync"

	"github.com/streadway/amqp"
)

// url：amqp://user:password @RabbitMQ_host:port/Vhost
var R = config.RMQSet
var MQURL = fmt.Sprintf("amqp://%s:%s@%s:%s/%s", R.User, R.Password, R.Host, R.Port, R.Vhost)

type RabbitMQ struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	QueueName string
	Exchange  string
	Key       string
	MqUrl     string
	sync.Mutex
}

func NewRabbitMQ(queueName string, exchange string, key string) *RabbitMQ {
	rabbitmq := &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		MqUrl:     MQURL,
	}
	var err error
	rabbitmq.conn, err = amqp.Dial(rabbitmq.MqUrl)
	failOnErr(err, "创建 RabbitMQ 连接失败")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	failOnErr(err, "获取 channel 失败")
	return rabbitmq
}

func failOnErr(err error, msg string) {
	if err != nil {
		log.Printf("%s:%s", msg, err)
		panic(fmt.Sprintf("%s:%s", msg, err))
	}
}

func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

func (r *RabbitMQ) PublishSimple(msg string) error {
	r.Lock()
	defer r.Unlock()
	//申请队列
	_, err := r.channel.QueueDeclare(
		r.QueueName, false, false, false, false, nil,
	)
	if err != nil {
		return err
	}

	r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		},
	)
	return nil
}

func (r *RabbitMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	r.Lock()
	defer r.Unlock()
	// 1.申请队列 如果队列不存在则会自动创建，如果存在则跳过创建
	_, err := r.channel.QueueDeclare(
		// 队列名称
		r.QueueName,
		// 是否持久化
		false,
		// 是否自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		// 额外属性
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	// 2.接收消息
	msgs, err := r.channel.Consume(
		// 队列名称
		r.QueueName,
		// 用来区分多个消费者
		"",
		// 是否自动应答
		false,
		// 是否具有排他性
		false,
		// 如果设置为 true，表示不能将同一个 connection 中发送的消息传递给这个 connection 中的消费
		false,
		// 队列消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
	}

	forever := make(chan bool)
	// 启用协程处理消息
	r.channel.Qos(1, 0, false)

	go func() {
		for d := range msgs {
			log.Printf("接收到消息 Received a message: %s", d.Body)
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(d.Body), message)
			if err != nil {
				fmt.Println(err)
			}

			// 插入订单
			_, err = orderService.InsertOrderByMessage(message)
			if err != nil {
				fmt.Println(err)
				continue
			}
			// 扣除商品数量
			err = productService.SubProductNum(message.ProductID, message.ProductNum)
			if err != nil {
				fmt.Println(err)
			}
			// true表示确认所有未确认的消息，false表示确认当前消息
			d.Ack(false)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

// 断开 channel 和 connection 连接
func (r *RabbitMQ) Destroy() {
	r.channel.Close()
	r.conn.Close()
}
