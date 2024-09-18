package rocketmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mygoshop/config"
	"mygoshop/datamodels"
	"mygoshop/services"
	"time"

	"github.com/apache/rocketmq-client-go/v2"
	"github.com/apache/rocketmq-client-go/v2/consumer"
	"github.com/apache/rocketmq-client-go/v2/primitive"
	"github.com/apache/rocketmq-client-go/v2/producer"
)

var NAME_SERVER = fmt.Sprintf("%s:%s", config.ROCKSet.Host, config.ROCKSet.Port)

type RocketMQ struct {
	GroupName string
	producer  rocketmq.Producer
	consumer  rocketmq.PushConsumer
}

func NewRocketMQ(groupName string) *RocketMQ {
	__rock := &RocketMQ{
		GroupName: groupName,
	}
	var err error
	__rock.producer, err = rocketmq.NewProducer(producer.WithNameServer([]string{NAME_SERVER}), producer.WithRetry(2))
	if err != nil {
		log.Fatalf("Create NewDefaultProducer error: %s", err.Error())
	}
	__rock.consumer, _ = rocketmq.NewPushConsumer(
		consumer.WithNameServer([]string{NAME_SERVER}),
		consumer.WithGroupName(groupName),
		consumer.WithRetry(2),
	)
	__rock.producer.Start()
	return __rock
}

func (r *RocketMQ) PublishSimple(msg string) error {
	__msg := &primitive.Message{
		Topic: r.GroupName,
		Body:  []byte(msg),
	}
	err := r.producer.SendAsync(context.Background(), func(ctx context.Context, result *primitive.SendResult, err error) {
		if err != nil {
			fmt.Printf("发送失败: %v\n", err)
		} else {
			fmt.Printf("发送成功: result=%s\n", result.String())
		}
	}, __msg)
	if err != nil {
		fmt.Printf("异步发送消息失败: %v\n", err)
	}
	time.Sleep(1 * time.Second)
	return err
}

func (r *RocketMQ) ConsumeSimple(orderService services.IOrderService, productService services.IProductService) {
	err := r.consumer.Subscribe(r.GroupName, consumer.MessageSelector{}, func(ctx context.Context, msgs ...*primitive.MessageExt) (consumer.ConsumeResult, error) {
		for i, msg := range msgs {
			fmt.Printf("收到消息[%d]: %s\n", i, string(msg.Body))
			message := &datamodels.Message{}
			err := json.Unmarshal([]byte(msg.Body), message)
			if err != nil {
				fmt.Println(err)
			}
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
		}
		return consumer.ConsumeSuccess, nil
	})
	if err != nil {
		fmt.Println("订阅消息失败:", err)
		return
	}
	err = r.consumer.Start()
	if err != nil {
		fmt.Println("启动消费者失败:", err)
		return
	}
	defer r.consumer.Shutdown()
	forever := make(chan bool)
	<-forever
}

func (r *RocketMQ) Destroy() {
	r.producer.Shutdown()
	r.consumer.Shutdown()
}
