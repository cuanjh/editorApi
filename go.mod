module editorApi

go 1.14

require (
	github.com/DataDog/zstd v1.4.5 // indirect
	github.com/EDDYCJY/go-gin-example v0.0.0-20191007083155-a98c25f2172a
	github.com/PuerkitoBio/goquery v1.5.1
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/boombuler/barcode v1.0.1-0.20190219062509-6c824513bacc
	github.com/casbin/casbin v1.9.1
	github.com/casbin/gorm-adapter v1.0.0
	github.com/chenhg5/collection v0.0.0-20191118032303-cb21bccce4c3
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/facebookgo/ensure v0.0.0-20200202191622-63f1cf65ac4c // indirect
	github.com/facebookgo/inject v0.0.0-20180706035515-f23751cae28b
	github.com/facebookgo/stack v0.0.0-20160209184415-751773369052 // indirect
	github.com/facebookgo/structtag v0.0.0-20150214074306-217e25fb9691 // indirect
	github.com/facebookgo/subset v0.0.0-20200203212716-c811ad88dec4 // indirect
	github.com/fastly/go-utils v0.0.0-20180712184237-d95a45783239 // indirect
	github.com/fsnotify/fsnotify v1.4.9
	github.com/gin-contrib/cors v1.3.1
	github.com/gin-gonic/gin v1.5.0
	github.com/go-openapi/spec v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.6 // indirect
	github.com/go-playground/locales v0.13.0
	github.com/go-playground/universal-translator v0.17.0
	github.com/go-playground/validator/v10 v10.3.0
	github.com/go-redis/redis v6.15.8+incompatible
	github.com/jehiah/go-strftime v0.0.0-20171201141054-1d33003b3869 // indirect
	github.com/jinzhu/copier v0.0.0-20190924061706-b57f9002281a
	github.com/jinzhu/gorm v1.9.10
	github.com/juju/ratelimit v1.0.1
	github.com/klauspost/compress v1.10.10 // indirect
	github.com/lestrrat/go-envload v0.0.0-20180220120943-6ed08b54a570 // indirect
	github.com/lestrrat/go-file-rotatelogs v0.0.0-20180223000712-d3151e2a480f
	github.com/lestrrat/go-strftime v0.0.0-20180220042222-ba3bf9c1d042 // indirect
	github.com/metaverse/truss v0.1.0 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/micro/go-plugins v1.5.1
	//github.com/micro/go-plugins/registry/consul v0.0.0-20200119172437-4fe21aa238fd
	github.com/mongodb/mongo-go-driver v1.2.1
	github.com/natefinch/lumberjack v2.0.0+incompatible
	github.com/nats-io/nats.go v1.10.0
	github.com/pkg/errors v0.9.1
	github.com/qiniu/api.v7 v7.2.5+incompatible
	github.com/qiniu/x v7.0.8+incompatible // indirect
	github.com/rifflock/lfshook v0.0.0-20180920164130-b9218ef580f5
	github.com/robfig/cron v1.2.0 // indirect
	github.com/robfig/cron/v3 v3.0.1 // indirect
	github.com/satori/go.uuid v1.2.0
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/viper v1.4.0
	github.com/swaggo/gin-swagger v1.2.0
	github.com/swaggo/swag v1.6.5
	github.com/tealeg/xlsx v1.0.5
	github.com/tebeka/strftime v0.1.3 // indirect
	github.com/tencentcloud/tencentcloud-sdk-go v3.0.158+incompatible
	github.com/tidwall/gjson v1.6.0 // indirect
	github.com/unrolled/secure v1.0.6
	go.mongodb.org/mongo-driver v1.2.1
	go.uber.org/zap v1.12.0
	golang.org/x/crypto v0.0.0-20200604202706-70a84ac30bf9 // indirect
	golang.org/x/tools v0.0.0-20200407041343-bf15fae40dea // indirect
	gopkg.in/go-playground/validator.v9 v9.31.0
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/natefinch/lumberjack.v2 v2.0.0 // indirect
	grpcSrv/proto/chatroom v0.0.0-00010101000000-000000000000
	grpcSrv/proto/im v0.0.0-00010101000000-000000000000
	honnef.co/go/tools v0.0.1-2020.1.3 // indirect
	qiniupkg.com/x v7.0.8+incompatible // indirect
	tkCommon v0.0.0-00010101000000-000000000000
)

replace grpcSrv/proto/chatroom => ../grpcSrv/proto/chatroom

replace grpcSrv/proto/im => ../grpcSrv/proto/im

replace tkCommon => ../tkCommon

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/micro/go-micro v1.18.0 => github.com/micro/go-micro v1.18.0

replace github.com/micro/go-plugins => github.com/micro/go-plugins v1.5.1

replace github.com/golang/protobuf => github.com/golang/protobuf v1.3.2
