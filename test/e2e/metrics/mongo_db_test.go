package metrics

import (
	"context"
	"fmt"
	mongo_driver "github.com/quwan-sre/observability-go-contrib/metrics/mongo-driver"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
	"time"
)

var mongoClient *mongo.Client

func initMongoDBClient() {
	var err error
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://test_user:password@localhost:27017").SetMonitor(mongo_driver.NewCommandMonitor()))
	if err != nil {
		panic(fmt.Sprintf("init mongo db client err: %v", err))
	}
}

func TestMongoDriver(t *testing.T) {
	initMongoDBClient()
	mongoClient.Ping(context.TODO(), nil)
	mongoClient.Database("test").CreateCollection(context.TODO(), "random things")
	//filter := bson.D{{"result", bson.D{{"$exists", true}}},
	//	{"request_id", "123"}}

	doc := struct {
		OK string `bson:"ok,omitempty"`
	}{
		OK: "ok",
	}
	insertResult, err := mongoClient.Database("fake").Collection("random things").InsertOne(context.TODO(), doc)
	fmt.Println(insertResult)
	fmt.Println(err)
	err = mongoClient.Database("fake").Drop(context.TODO())
	fmt.Println(err)

	mongoClient.Database("fake").CreateCollection(context.TODO(), "exist")
	err = mongoClient.Database("fake").CreateCollection(context.TODO(), "exist")
	fmt.Println(err)

	result, err := mongoClient.Database("fake").Collection("yourCollection").UpdateMany(context.Background(), bson.D{{}}, bson.D{{"$set t", bson.D{{"field", "value"}}}})
	fmt.Println(result)
	fmt.Println(err)
	cursor, err := mongoClient.Database("fake").Collection("testset").Find(context.TODO(), bson.D{{"$set t", bson.D{{"field", "value"}}}})
	fmt.Println(cursor)
	fmt.Println(err)
}
