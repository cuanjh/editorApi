package mgdb

import (
	"bytes"
	"context"
	"editorApi/config"
	"editorApi/init/qmlog"
	"log"
	"time"

	"github.com/mongodb/mongo-go-driver/mongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoClient *mongo.Client
var EditorClient *mongo.Client
var OnlineClient *mongo.Client
var TestClient *mongo.Client

type EnvConfig string

var (
	EnvTest   EnvConfig = "test"
	EnvOnline EnvConfig = "online"
	EnvEditor EnvConfig = "editor"
	DbEditor  string    = "editor"
	DbDict    string    = "dict"
	DbKuyu    string    = "kuyu"
	DbContent string    = "courseContent"
)

func init() {
	MongoClient = GetClient("")
	EditorClient = GetClient("")
	OnlineClient = GetClient("to")
	TestClient = GetClient("test")
}

func getContext() (ctx context.Context) {
	ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
	return
}

//通过连接字符串，连接到MongoDB
func GetClient(dbConfig string) *mongo.Client {
	connString := bytes.NewBufferString("mongodb://")

	if dbConfig == "to" {
		mgdb := config.GinVueAdminconfig.Tomongodb
		if mgdb.User != "" {
			connString.WriteString(mgdb.User + ":" + mgdb.Passwd + "@")
		}
		connString.WriteString(mgdb.Hosts)
		if mgdb.ReplicaSet != "" {
			connString.WriteString("?replicaSet=" + mgdb.ReplicaSet)
		}
	} else if dbConfig == "test" {
		mgdb := config.GinVueAdminconfig.Testmongodb
		if mgdb.User != "" {
			connString.WriteString(mgdb.User + ":" + mgdb.Passwd + "@")
		}
		connString.WriteString(mgdb.Hosts)
		if mgdb.ReplicaSet != "" {
			connString.WriteString("?replicaSet=" + mgdb.ReplicaSet)
		}
	} else {
		mgdb := config.GinVueAdminconfig.Mongodb
		if mgdb.User != "" {
			connString.WriteString(mgdb.User + ":" + mgdb.Passwd + "@")
		}
		connString.WriteString(mgdb.Hosts)
		if mgdb.ReplicaSet != "" {
			connString.WriteString("?replicaSet=" + mgdb.ReplicaSet)
		}

	}

	readconcern.Majority()
	opt := options.Client().ApplyURI(connString.String())
	opt.SetLocalThreshold(3 * time.Second)  //只使用与mongo操作耗时小于3秒的
	opt.SetMaxConnIdleTime(5 * time.Second) //指定连接可以保持空闲的最大毫秒数
	opt.SetMaxPoolSize(500)                 //使用最大的连接数
	client, err := mongo.Connect(getContext(), opt)
	if err != nil {
		qmlog.QMLog.Error("Mongodb链接不可用12", connString.String(), err)
		return nil
	}

	if err := client.Ping(getContext(), readpref.Primary()); err != nil {
		qmlog.QMLog.Error("Mongodb链接不可用22", connString.String(), err)
		return nil
	}

	qmlog.QMLog.Info("链接到MongoDB数据库:" + connString.String())
	return client
}

func Find(
	env EnvConfig,
	dbName,
	colName string,
	filter bson.M,
	sort interface{},
	project interface{},
	skip,
	limit int64,
	rst interface{},
) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	}
	ctx := getContext()
	opts := options.Find()
	if sort != nil {
		opts = opts.SetSort(sort)
	}
	if project != nil {
		opts = opts.SetProjection(project)
	}
	opts = opts.SetSkip(skip)
	if limit == 0 {
		limit = 10000
	}
	opts = opts.SetLimit(limit)

	cusor, err := client.Database(dbName).Collection(colName).Find(
		ctx,
		filter,
		opts,
	)
	defer cusor.Close(ctx)
	if err != nil {
		log.Println(err)
		return
	}
	cusor.All(ctx, rst)
}

func FindOne(
	env EnvConfig,
	dbName,
	colName string,
	filter interface{},
	project interface{},
	rst interface{},
) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	} else if env == "test" {
		client = TestClient
	}
	ctx := getContext()
	opts := options.FindOne()
	if project != nil {
		opts = opts.SetProjection(project)
	}
	single := client.Database(dbName).Collection(colName).FindOne(ctx, filter, opts)
	single.Decode(rst)
}

func FindOneAndUpdate(
	env EnvConfig,
	dbName,
	colName string,
	filter interface{},
	sets interface{},
	project interface{},
	rst interface{},
) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	} else if env == "test" {
		client = TestClient
	}
	ctx := getContext()
	opts := options.FindOneAndUpdate()
	if project != nil {
		opts = opts.SetProjection(project)
	}
	single := client.Database(dbName).Collection(colName).FindOneAndUpdate(ctx, filter, sets, opts)
	single.Decode(rst)
}

func UpdateOne(
	env EnvConfig,
	dbName,
	colName string,
	filter interface{},
	sets interface{},
	upsert bool,
) (*mongo.UpdateResult, error) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	} else if env == "test" {
		client = TestClient
	}
	ctx := getContext()
	opts := options.Update()
	if upsert {
		opts = opts.SetUpsert(true)
	}

	return client.Database(dbName).Collection(colName).UpdateOne(
		ctx,
		filter,
		sets,
		opts,
	)
}

func UpdateMany(
	env EnvConfig,
	dbName,
	colName string,
	filter interface{},
	sets interface{},
	upsert bool,
) (*mongo.UpdateResult, error) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	} else if env == "test" {
		client = TestClient
	}
	ctx := getContext()
	opts := options.Update()
	if upsert {
		opts = opts.SetUpsert(upsert)
	}

	return client.Database(dbName).Collection(colName).UpdateMany(
		ctx,
		filter,
		sets,
		opts,
	)
}

func DeleteOne(env EnvConfig,
	dbName,
	colName string,
	filter interface{},
) (*mongo.DeleteResult, error) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	} else if env == "test" {
		client = TestClient
	}
	ctx := getContext()
	opts := options.Delete()
	return client.Database(dbName).Collection(colName).DeleteOne(
		ctx,
		filter,
		opts,
	)
}

func DeleteMany(
	env EnvConfig,
	dbName,
	colName string,
	filter interface{},
) (*mongo.DeleteResult, error) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	} else if env == "test" {
		client = TestClient
	}
	ctx := getContext()
	opts := options.Delete()
	return client.Database(dbName).Collection(colName).DeleteMany(
		ctx,
		filter,
		opts,
	)
}

func Count(
	env EnvConfig,
	dbName,
	colName string,
	filter interface{},
) (int64, error) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	} else if env == "test" {
		client = TestClient
	}
	ctx := getContext()
	opts := options.Count()

	return client.Database(dbName).Collection(colName).CountDocuments(
		ctx,
		filter,
		opts,
	)
}

func Aggregate(
	env EnvConfig,
	dbName,
	colName string,
	pipeline interface{},
) (cursor *mongo.Cursor, err error) {
	client := EditorClient
	if env == "online" {
		client = OnlineClient
	} else if env == "test" {
		client = TestClient
	}
	ctx := getContext()

	collection := client.Database(dbName).Collection(colName)
	cursor, err = collection.Aggregate(ctx, pipeline)
	return
}
