package main

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
)

type Item struct {
	Id    int     `gorm:"id"`
	Price float32 `gorm:"price"`
}

type Language struct {
	Id       int    `gorm:"id"`
	Language string `gorm:"language"`
}

type NameNode struct {
	ItemId     int    `gorm:"itemId"`
	LanguageId int    `gorm:"languageId"`
	Content    string `gorm:"content"`
}

type Purchase struct {
	Id          int `gorm:"id"`
	ItemId      int `gorm:"itemId"`
	ItemsAmount int `gorm:"itemsAmount"`
}

type User struct {
	Id    int    `gorm:"id"`
	Name  string `gorm:"name"`
	Email string `gorm:"email"`
}

type UserPurchase struct {
	PurchaseId int `gorm:"purchaseId"`
	User       int `gorm:"userId"`
}

type FirstFormTable struct {
	Id         int     `gorm:"id"`
	UserName   string  `gorm:"userName"`
	UserEmail  string  `gorm:"userEmail"`
	ItemAmount int     `gorm:"itemAmount"`
	ItemPrice  float32 `gorm:"itemPrice"`
	ItemNameEn string  `gorm:"itemNameEn"`
	ItemNameRu string  `gorm:"itemNameRu"`
}

type ExcelResult struct {
	id           int
	name         string
	email        string
	items_amount int
	price        float32
	content      string
	language     string
}

func NewUser(id int, name string, email string) User {
	return User{
		Id:    id,
		Name:  name,
		Email: email,
	}
}

func NewItem(id int, price float32) Item {
	return Item{
		Id:    id,
		Price: price,
	}
}

func NewLanguage(id int, language string) Language {
	return Language{
		Id:       id,
		Language: language,
	}
}

func NewNameNode(id int, languageId int, content string) NameNode {
	return NameNode{
		ItemId:     id,
		LanguageId: languageId,
		Content:    content,
	}
}

func NewPurchase(id int, itemId int, itemsAmount int) Purchase {
	return Purchase{
		Id:          id,
		ItemId:      itemId,
		ItemsAmount: itemsAmount,
	}
}

func NewUserPurchase(purchase int, user int) UserPurchase {
	return UserPurchase{
		PurchaseId: purchase,
		User:       user,
	}
}

func FindUser(userArray []User, user User) int {
	for _, v := range userArray {
		if v.Name == user.Name && v.Email == user.Email {
			return v.Id
		}
	}
	return -1
}

func ToGOB64(m SX) string {
	b := bytes.Buffer{}
	e := gob.NewEncoder(&b)
	err := e.Encode(m)
	if err != nil {
		fmt.Println(`failed gob Encode`, err)
	}
	return base64.StdEncoding.EncodeToString(b.Bytes())
}

type SX struct {
	Items         []Item
	Languages     []Language
	NameNodes     []NameNode
	Purchases     []Purchase
	Users         []User
	UserPurchases []UserPurchase
}

func NewSX(item []Item, language []Language, nameNode []NameNode, purchase []Purchase, user []User, userPurchase []UserPurchase) SX {
	return SX{
		Items:         item,
		Languages:     language,
		NameNodes:     nameNode,
		Purchases:     purchase,
		Users:         user,
		UserPurchases: userPurchase,
	}
}

func getItems(sx SX) []Item {
	return sx.Items
}

func getLanguages(sx SX) []Language {
	return sx.Languages
}

func getNameNodes(sx SX) []NameNode {
	return sx.NameNodes
}

func getPurchases(sx SX) []Purchase {
	return sx.Purchases
}

func getUsers(sx SX) []User {
	return sx.Users
}

func getUserPurchases(sx SX) []UserPurchase {
	return sx.UserPurchases
}
