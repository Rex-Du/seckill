// Author : rexdu
// Time : 2020-03-22 00:35
package RabbitMQ

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

// url格式： amqp://用户名：密码@rabbitmq服务器地址：端口/virtualhost
const MQURL = "amqp://rexdu:rootroot@192.168.124.129:5672/rexdu"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	// 队列名称
	QueueName string
	// 交换机
	Exchange string
	// key
	Key   string
	Mqurl string
}

func NewRabbitMQ(queueName, exchange, key string) *RabbitMQ {
	var (
		rabbitmq *RabbitMQ
		err      error
	)
	rabbitmq = &RabbitMQ{
		QueueName: queueName,
		Exchange:  exchange,
		Key:       key,
		Mqurl:     MQURL,
	}
	rabbitmq.conn, err = amqp.Dial(rabbitmq.Mqurl)
	rabbitmq.failOnErr(err, "创建连接错误！")
	rabbitmq.channel, err = rabbitmq.conn.Channel()
	rabbitmq.failOnErr(err, "获取channel失败！")
	return rabbitmq
}

// 断开channel和connection
func (r *RabbitMQ) Destory() {
	r.channel.Close()
	r.conn.Close()
}

// 错误处理函数
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s", message, err)
		panic(fmt.Sprintf("%s:%s", message, err))
	}
}

// 创建简单模式下的RabbitMQ实例
func NewRabbitMQSimple(queueName string) *RabbitMQ {
	return NewRabbitMQ(queueName, "", "")
}

// 发布消息的函数
func (r *RabbitMQ) PublishSimple(message string) {
	var (
		err error
	)
	// 申请队列，如果队列不存在就创建队列，保证消息能发送到队列中
	_, err = r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	// 发送消息到队列中
	r.channel.Publish(
		r.Exchange,
		r.QueueName,
		// 如果为true，根据exchange类型和routkey规则，如果无法找到符合条件的队列，会把发送的消息返回给发送都
		false,
		// 如果为true，当exchange发送消息到队列后发现队列上没有绑定消费者，会把发送的消息返回给发送都
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

// 消费消息的函数
func (r *RabbitMQ) ConsumeSimple() {
	var (
		err     error
		msgs    <-chan amqp.Delivery
		forever chan bool
	)
	// 申请队列，如果队列不存在就创建队列，保证消息能发送到队列中
	_, err = r.channel.QueueDeclare(
		r.QueueName,
		// 是否持久化
		false,
		// 是否为自动删除
		false,
		// 是否具有排他性
		false,
		// 是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}
	msgs, err = r.channel.Consume(
		r.QueueName,
		"",
		// 是否自动应答
		true,
		// 是否具有排他性
		false,
		// 如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		// 消费是否阻塞
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		return
	}

	forever = make(chan bool)

	go func() {
		for d := range msgs {
			// 这里写实际处理的逻辑
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Println("[*] Waiting for messages, To exit press Ctrl+C")
	<-forever

}

// 创建订阅模式下的RabbitMQ实例
func NewRabbitMQSub(exchange string) *RabbitMQ {
	return NewRabbitMQ("", exchange, "")
}

// 订阅模式下的生产
func (r *RabbitMQ) PublishPub(message string) {
	var (
		err error
	)
	// 创建交换机
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机的类型：广播类型
		"fanout",
		// 是否持久化
		true,
		false,
		// true表示这个exchange不可以被client用来推送消息，只能用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 发送消息
	err = r.channel.Publish(
		r.Exchange,
		"",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) ReceiveSub() {
	var (
		err      error
		q        amqp.Queue
		messages <-chan amqp.Delivery
	)
	// 创建交换机
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机的类型：广播类型
		"fanout",
		// 是否持久化
		true,
		false,
		// true表示这个exchange不可以被client用来推送消息，只能用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	q, err = r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an queue")

	err = r.channel.QueueBind(
		q.Name,
		"",
		r.Exchange,
		false,
		nil,
	)

	// 消费消息
	messages, err = r.channel.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)
	go func() {
		for m := range messages {
			log.Printf("Received a message: %s", m.Body)
		}
	}()

	log.Println("[*] Waiting for messages, To exit press Ctrl+C")
	<-forever
}

// 创建路由模式下的RabbitMQ实例
func NewRabbitMQRouting(exchangeName, routingKey string) *RabbitMQ {
	return NewRabbitMQ("", exchangeName, routingKey)
}

func (r *RabbitMQ) PublishRouting(message string) {
	var (
		err error
	)
	// 创建交换机
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机的类型：广播类型
		"direct",
		// 是否持久化
		true,
		false,
		// true表示这个exchange不可以被client用来推送消息，只能用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	// 发送消息
	err = r.channel.Publish(
		r.Exchange,
		r.Key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (r *RabbitMQ) ReceiveRouting() {
	var (
		err      error
		q        amqp.Queue
		messages <-chan amqp.Delivery
	)
	// 创建交换机
	err = r.channel.ExchangeDeclare(
		r.Exchange,
		// 交换机的类型：广播类型
		"direct",
		// 是否持久化
		true,
		false,
		// true表示这个exchange不可以被client用来推送消息，只能用来进行exchange和exchange之间的绑定
		false,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an exchange")

	q, err = r.channel.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	r.failOnErr(err, "Failed to declare an queue")

	err = r.channel.QueueBind(
		q.Name,
		r.Key,
		r.Exchange,
		false,
		nil,
	)

	// 消费消息
	messages, err = r.channel.Consume(q.Name, "", true, false, false, false, nil)

	forever := make(chan bool)
	go func() {
		for m := range messages {
			log.Printf("Received a message: %s", m.Body)
		}
	}()

	log.Println("[*] Waiting for messages, To exit press Ctrl+C")
	<-forever
}
