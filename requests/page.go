package requests

type Page struct {
	PageIndex int64  `json:"page_index"` //分页第几页
	PageSize  int64  `json:"page_size"`  //分页数据条数
	TextField string `json:"text_field"` //排序字段
	SortType  int    `json:"sort_type"`  //排序类型 1、倒叙；-1、顺序
}