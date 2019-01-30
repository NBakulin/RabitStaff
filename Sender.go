package main

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/streadway/amqp"
	_ "github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	dbsqlite, err := gorm.Open("sqlite3", "C:/SQLiteStudio/FirstFormDB.db")
	if err != nil {
		panic("failed to connect database")
	}

	var firstFormTable []FirstFormTable
	var item []Item
	var language []Language
	var nameNode []NameNode
	var purchase []Purchase
	var user []User
	var userPurchase []UserPurchase
	var desKey = NewDesKey()

	dbsqlite.Find(&firstFormTable)
	fmt.Println(&firstFormTable)

	language = append(language, NewLanguage(0, "English"))
	language = append(language, NewLanguage(1, "Russian"))

	for i := 0; i < len(firstFormTable); i++ {
		item = append(item, NewItem(i, firstFormTable[i].ItemPrice))
		nameNode = append(nameNode, NewNameNode(i, 0, firstFormTable[i].ItemNameEn))
		nameNode = append(nameNode, NewNameNode(i, 1, firstFormTable[i].ItemNameRu))
		purchase = append(purchase, NewPurchase(i, i, firstFormTable[i].ItemAmount))
		var newUser = NewUser(i, firstFormTable[i].UserName, firstFormTable[i].UserEmail)
		var userId = FindUser(user, newUser)
		if userId < 0 {
			user = append(user, newUser)
			userPurchase = append(userPurchase, NewUserPurchase(i, i))
		} else {
			userPurchase = append(userPurchase, NewUserPurchase(i, userId))
		}
	}

	dbsqlite.Close()
	var sx = NewSX(item, language, nameNode, purchase, user, userPurchase)
	var serializedStructs = ToGOB64(sx)
	mytext := []byte(serializedStructs)
	//не робит, сук
	cryptoText, _ := DesEncryption(desKey.Key, desKey.Key, mytext)

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	kq, err := ch.QueueDeclare(
		"key_queue", // name
		false,       // durable
		true,        // delete when unused
		false,       // exclusive
		false,       // no-wait
		nil,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	encryptedRsa := encryptRsa(desKey.Key)
	err = ch.Publish(
		"",      // exchange
		kq.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        encryptedRsa,
		})

	mq, err := ch.QueueDeclare(
		"message_queue", // name
		false,           // durable
		true,            // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	//add
	err = ch.Publish(
		"",      // exchange
		mq.Name, // routing key
		false,   // mandatory
		false,   // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        cryptoText,
		})

	failOnError(err, "Failed to publish a message")
}
