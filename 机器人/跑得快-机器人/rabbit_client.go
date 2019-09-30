package main

import (
	"encoding/json"
	"fmt"
	"logs"
	"net"
	"time"

	"github.com/streadway/amqp"
)

type AmqpConf struct {
	Host      string
	Port      int
	User      string
	Pass      string
	QueueName string
	Vhost     string
	Exchange  string
}

type RabbitClient struct {
	conf    AmqpConf
	conn    *amqp.Connection
	channel *amqp.Channel

	SendMsg chan []byte
}

func (this *RabbitClient) init(conf AmqpConf) bool {
	this.conf = conf
	this.conn = nil
	this.channel = nil
	this.SendMsg = make(chan []byte, 1000)

	dialaddr := fmt.Sprintf("%s://%s:%s@%s:%d/%s", "amqp", conf.User, conf.Pass, conf.Host, conf.Port, conf.Vhost)

	fmt.Println("连接地址", dialaddr)
	var dialConfig amqp.Config
	dialConfig.Heartbeat = 10 * time.Second
	dialConfig.Locale = "en_US"
	dialConfig.Dial = func(network, addr string) (net.Conn, error) {
		return net.DialTimeout(network, addr, 5*time.Second)
	}

	amqpConn, err := amqp.DialConfig(dialaddr, dialConfig)
	if err != nil {
		fmt.Println("连接失败", err)
		return false
	}
	this.conn = amqpConn

	return true
}

func (this *RabbitClient) AddMsgNative(d interface{}) {
	b, _ := json.Marshal(d)

	select {
	case this.SendMsg <- b:

	default:
		ErrorLog("丢失信息")
	}
}

func (this *RabbitClient) getChannel() (*amqp.Channel, error) {
	if this.channel != nil {
		return this.channel, nil
	}

	amqpChan, err := this.conn.Channel()
	if err != nil {
		fmt.Println("获取通道失败", err)
	}
	this.channel = amqpChan
	return amqpChan, err
}

func (this *RabbitClient) SendMessageThread() {
	for s := range this.SendMsg {
		msg := amqp.Publishing{
			ContentType: "application/json",
			Body:        s,
		}
		chann, err := this.getChannel()
		if err != nil {
			ErrorLog("获取通道失败", err)
			break
		}

		err = chann.Publish("", this.conf.QueueName, false, false, msg)
		if err != nil {
			ErrorLog("发送消息失败", err)
			break
		}

		DebugLog("发送归还机器人队列消息成功", this.conf.QueueName, string(msg.Body))
	}
}

func (this *RabbitClient) Clean() {
	if this.channel != nil {
		this.channel.Close()
		this.channel = nil
	}
	if this.conn != nil {
		this.conn.Close()
		this.conn = nil
	}
}

var GRabbitClient *RabbitClient = new(RabbitClient)

func init() {
	conf := AmqpConf{
		Host:      gameConfig.GCRabbitClient.Host,
		Port:      gameConfig.GCRabbitClient.Port,
		User:      gameConfig.GCRabbitClient.User,
		Pass:      gameConfig.GCRabbitClient.Pass,
		QueueName: gameConfig.GCRabbitClient.QueueName,
		Vhost:     "",
		Exchange:  "",
	}

	go func() {
		for {
			logs.Debug("开始连接消息队列服务器")
			succ := GRabbitClient.init(conf)
			if succ {
				logs.Debug("消息队列服务器连接成功")
				GRabbitClient.SendMessageThread()
			}
			logs.Debug("消息队列服务器连接失败")
			GRabbitClient.Clean()
			time.Sleep(time.Second)
		}
	}()
}
