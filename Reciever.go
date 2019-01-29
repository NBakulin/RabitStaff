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

func failIfError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failIfError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failIfError(err, "Failed to open a channel")

	keyMsg, err := ch.Consume(
		"key_queue", // queue
		"",          // consumer
		true,        // auto-ack
		false,       // exclusive
		false,       // no-local
		false,       // no-wait
		nil,         // args
	)
	failIfError(err, "Failed to register a consumer")

	for d := range keyMsg {
		key = d.Body[:]
		log.Printf("Received a message: %s", string(key))
		ch.Close()
	}

	ch2, err := conn.Channel()
	failIfError(err, "Failed to open a channel")

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

	for d := range msgs {
		var decryptedSx, _ = DesDecryption(key, key, d.Body[:])
		var sx = FromGOB64(string(decryptedSx))
		log.Printf("%s", sx)
		item = sx.Items
		language = sx.Languages
		nameNode = sx.NameNodes
		purchase = sx.Purchases
		user = sx.Users
		userPurchase = sx.UserPurchases
		ch2.Close()
	}
	conn.Close()

	dbgorm, err := gorm.Open("mysql", "root:root@/normalizeddatabase?charset=utf8")

	dbgorm.Delete(&User{})
	dbgorm.Delete(&Item{})
	dbgorm.Delete(&Language{})
	dbgorm.Delete(&Purchase{})
	dbgorm.Delete(&UserPurchase{})
	dbgorm.Delete(&NameNode{})

	dbgorm.AutoMigrate()

	for i := 0; i < len(language); i++ {
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
