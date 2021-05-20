package els

import (
	"context"
	"os"
	"time"

	"log"

	elastic "github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	_ "github.com/spf13/viper/remote"
)

var client *elastic.Client

func init() {
	setClient()
}

func setClient() {
	v := viper.New()
	v.AddRemoteProvider("consul", "127.0.0.1:8500", "elasticsearch")
	v.SetConfigType("json")
	v.ReadRemoteConfig()

	var setConfig = func() {

		hosts := v.GetStringSlice("hosts")
		user := v.GetString("user")
		pwd := v.GetString("password")

		// log.Println("Init Elastic Client", hosts)
		var err error
		client, err = elastic.NewClient(
			elastic.SetURL(hosts...),
			elastic.SetBasicAuth(user, pwd),
			elastic.SetSniff(false),
			elastic.SetHealthcheckInterval(10*time.Second),
			elastic.SetMaxRetries(5),
			elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
			elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		)
		if err != nil {
			log.Fatalf("Fail to connect elasticsearch :" + err.Error())
		}
	}
	setConfig()

	// //watch changes
	// go func() {
	// 	for {
	// 		time.Sleep(time.Second * 10) // delay after each request
	// 		// log.Println("watch Again")
	// 		err := v.WatchRemoteConfigOnChannel()
	// 		if err != nil {
	// 			log.Fatalf("unable to read remote config: %v", err)
	// 			continue
	// 		}
	// 		setConfig()
	// 	}
	// }()

}

func IndexExists(index string) bool {
	exists, err := client.IndexExists(index).Do(context.Background())
	if err != nil {
		log.Fatalf("Fail to check index :" + err.Error())
	}
	return exists
}

func IndexCreate(index, mapping string) {
	createIndex, err := client.CreateIndex(index).Body(mapping).Do(context.Background())
	if err != nil {
		// Handle error
		log.Fatalf("Fail to create index :" + err.Error())
	}
	if !createIndex.Acknowledged {
		// Not acknowledged
	}
}

func IndexDelete(index string) {
	deleteIndex, err := client.DeleteIndex(index).Do(context.Background())
	if err != nil {
		// Handle error
		panic(err)
	}
	if !deleteIndex.Acknowledged {
		// Not acknowledged
	}
}

func PutJson(index, id string, data interface{}) {
	put, err := client.Index().
		Index(index).
		Id(id).
		BodyJson(data).
		Do(context.Background())
	if err != nil {
		// Handle error
		log.Fatalf("Fail to add doc :" + err.Error())
	}
	_, err = client.Refresh().Index(index).Do(context.Background())
	if err != nil {
		log.Fatalf("Fail to add doc :" + err.Error())
	}
	log.Printf("Indexed doc %s to index %s, type %s\n", put.Id, put.Index, put.Type)
}

func PutStr(index, id, data string) {
	put, err := client.Index().
		Index(index).
		Id(id).
		BodyString(data).
		Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	_, err = client.Refresh().Index(index).Do(context.Background())

	if err != nil {
		handleErr(err)
	}

	log.Printf("Indexed doc %s to index %s, type %s\n", put.Id, put.Index, put.Type)
}

func DelById(index, id string) *elastic.DeleteResponse {
	res, err := client.Delete().Index(index).Id(id).Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	_, err = client.Refresh().Index(index).Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	return res
}

func Update(index, id string, data interface{}) int {

	rsp, err := client.Update().
		Index(index).
		Id(id).
		Doc(data).
		Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	_, err = client.Refresh().Index(index).Do(context.Background())

	if err != nil {
		handleErr(err)
	}

	return rsp.Shards.Total
}

func Upsert(index, id string, data interface{}) int {

	rsp, err := client.Update().
		Index(index).
		Id(id).
		DocAsUpsert(true).
		Doc(data).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		handleErr(err)
	}

	_, err = client.Refresh().Index(index).Do(context.Background())

	if err != nil {
		handleErr(err)
	}

	return rsp.Shards.Total
}

func Exists(index, id string) bool {

	ex, err := client.Exists().
		Index(index).
		Id(id).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	return ex
}

func GetById(index, id string) (*elastic.GetResult, error) {
	get, err := client.Get().
		Index(index).
		Id(id).
		Pretty(true).
		Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	return get, err
}

func Search(
	index string,
	from, size int,
	query elastic.Query,
	sortInfos []elastic.SortInfo,
	includes, excludes []string,
) (*elastic.SearchResult, error) {
	ctx := elastic.NewFetchSourceContext(true)
	if includes != nil {
		ctx = ctx.Include(includes...)
	}
	if excludes != nil {
		ctx = ctx.Exclude(excludes...)
	}
	searchService := client.Search().
		FetchSourceContext(ctx).
		Index(index).
		Query(query).
		From(from).Size(size).
		Pretty(true)
	if sortInfos != nil && len(sortInfos) > 0 {
		for _, s := range sortInfos {
			searchService = searchService.SortWithInfo(s)
		}
	}

	searchResult, err := searchService.
		Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	return searchResult, err
}

//批量增加文档
func PutBulk(index string, ids []string, docs []interface{}) *elastic.BulkResponse {
	if len(ids) == 0 || len(ids) != len(docs) {
		log.Println("No docs Insert!")
		return nil
	}
	bulk := client.Bulk()
	for k, id := range ids {
		req := elastic.NewBulkIndexRequest().Index(index).Id(id).Doc(docs[k])
		bulk = bulk.Add(req)
	}
	rps, err := bulk.Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	return rps
}

//批量更新文档
func UpdateBulk(index string, ids []string, docs []interface{}) *elastic.BulkResponse {
	if len(ids) == 0 || len(ids) != len(docs) {
		log.Println("No docs Update!")
		return nil
	}
	bulk := client.Bulk()
	for k, id := range ids {
		req := elastic.NewBulkUpdateRequest().Index(index).Id(id).Doc(docs[k])
		bulk = bulk.Add(req)
	}
	rps, err := bulk.Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	return rps
}

//批量删除
func DelBulk(index string, ids []string) *elastic.BulkResponse {
	if len(ids) == 0 {
		log.Println("No docs Delete!")
		return nil
	}
	bulk := client.Bulk()
	for _, id := range ids {
		req := elastic.NewBulkDeleteRequest().Index(index).Id(id)
		bulk = bulk.Add(req)
	}
	rps, err := bulk.Do(context.Background())
	if err != nil {
		handleErr(err)
	}
	return rps
}

func handleErr(err error) {
	if err != nil {
		switch {
		case elastic.IsNotFound(err):
			log.Printf("Document not found: %v", err)
		case elastic.IsTimeout(err):
			log.Printf("Timeout retrieving document: %v", err)
		case elastic.IsConnErr(err):
			log.Printf("Connection problem: %v", err)
		default:
			// Some other kind of error
			log.Printf("Unknow problem: %v", err)
		}
	}
}
