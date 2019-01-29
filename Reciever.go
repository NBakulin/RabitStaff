package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/streadway/amqp"
	_ "github.com/streadway/amqp"
	"log"
)

var item []Item
var language []Language
var nameNode []NameNode
var purchase []Purchase
var user []User
var userPurchase []UserPurchase
var key []byte
var key2 string

func failIfError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func mainн() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failIfError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failIfError(err, "Failed to open a channel")
	defer ch.Close()

	keyMsg, err := ch.Consume(
		"key_queue", // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failOnError(err, "Failed to register a consumer")

	for d := range keyMsg {
		key2 = string(d.Body[:])
		log.Printf(key2)
	}

	conn2, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failIfError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch2, err := conn2.Channel()
	failIfError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch2.Consume(
		"message_queue",    // queue
		"message_consumer", // consumer
		true,               // auto-ack
		false,              // exclusive
		false,              // no-local
		false,              // no-wait
		nil,                // args
	)
	failIfError(err, "Failed to register a consumer")

	//forever2 := make(chan bool)
	//
	//go func() {
	for d := range msgs {
		log.Printf("Received a message: %s", FromGOB64(string(d.Body[:])))
		var sx = FromGOB64(string(d.Body[:]))
		log.Printf("Received a message: %s", sx)
		item = sx.Items
		language = sx.Languages
		nameNode = sx.NameNodes
		purchase = sx.Purchases
		user = sx.Users
		userPurchase = sx.UserPurchases
		//<-forever2
	}
	//}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")

	dbgorm, err := gorm.Open("mysql", "root:root@/normalizeddatabase?charset=utf8")

	dbgorm.Delete(&User{})
	dbgorm.Delete(&Item{})
	dbgorm.Delete(&Language{})
	dbgorm.Delete(&Purchase{})
	dbgorm.Delete(&UserPurchase{})
	dbgorm.Delete(&NameNode{})

	dbgorm.AutoMigrate()

	for i := 0; i < 2; i++ {
		dbgorm.Create(&language[i])
	}

	for i := 0; i < len(purchase); i++ {
		dbgorm.Create(&item[i])
		dbgorm.Create(&purchase[i])
		dbgorm.Create(&userPurchase[i])
	}
	for i := 0; i < len(user); i++ {
		dbgorm.Create(&user[i])
	}

	for i := 0; i < len(nameNode); i++ {
		dbgorm.Create(&nameNode[i])
	}

	dbgorm.Close()
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
