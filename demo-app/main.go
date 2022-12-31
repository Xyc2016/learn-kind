package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
)

func getClient(uri, userName, userPassword string) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	var connectURI = fmt.Sprintf("mongodb://%s:%s@%s:27017",  userName, userPassword, uri)
	logrus.Infof("connectURI %s", connectURI)
	return mongo.Connect(ctx, options.Client().ApplyURI(connectURI))
}

func main() {
	var MONGO_URL = os.Getenv("DB_URL")
	var MONGO_USER_NAME = os.Getenv("USER_NAME")
	var MONGO_USER_PWD = os.Getenv("USER_PWD")
	client, err := getClient(MONGO_URL, MONGO_USER_NAME, MONGO_USER_PWD)
	if err != nil {
		logrus.Infof("get MongoClient failed :%s", err)
	}
	collection := client.Database("testing").Collection("numbers")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := collection.InsertOne(ctx, bson.D{{"name", "pi"}, {"value", 3.14159}})
	if err != nil {
		logrus.Infof("insert err %s", err)
	}
	fmt.Printf("id", res.InsertedID)
	result := collection.FindOne(ctx, bson.D{})
	var resultDoc bson.D
	result.Decode(&resultDoc)
	logrus.Infof("result %#v", resultDoc)

	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong??",
		})
	})
	r.GET("/envs", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"envs": os.Environ(),
		})
	})
	logrus.Info("start server aaa")
	r.Run("0.0.0.0:8000") // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
