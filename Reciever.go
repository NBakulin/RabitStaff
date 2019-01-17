package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/streadway/amqp"
	_ "github.com/streadway/amqp"
	"log"
)

func failIfError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failIfError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failIfError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failIfError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failIfError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", FromGOB64(string(d.Body[:])))
			var sx = FromGOB64(string(d.Body[:]))
			log.Printf("Received a message: %s", sx)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

	//dbgorm, err := gorm.Open("mysql", "root:root@/normalizeddatabase?charset=utf8")
	//
	//dbgorm.Delete(&User{})
	//dbgorm.Delete(&Item{})
	//dbgorm.Delete(&Language{})
	//dbgorm.Delete(&Purchase{})
	//dbgorm.Delete(&UserPurchase{})
	//dbgorm.Delete(&NameNode{})
	//
	//dbgorm.AutoMigrate()
	//
	//for i := 0; i < 2; i++ {
	//	dbgorm.Create(&language[i])
	//}
	//
	//for i := 0; i < len(firstFormTable); i++ {
	//	dbgorm.Create(&item[i])
	//	dbgorm.Create(&purchase[i])
	//	dbgorm.Create(&userPurchase[i])
	//}
	//for i := 0; i < len(user); i++ {
	//	dbgorm.Create(&user[i])
	//}
	//
	//for i := 0; i < len(nameNode); i++ {
	//	dbgorm.Create(&nameNode[i])
	//}
	//
	//dbgorm.Close()
}

func FromGOB64(str string) SX {
	m := SX{}
	by, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		fmt.Println(`failed base64 Decode`, err)
	}
	b := bytes.Buffer{}
	b.Write(by)
	d := gob.NewDecoder(&b)
	err = d.Decode(&m)
	if err != nil {
		fmt.Println(`failed gob Decode`, err)
	}
	return m
}
