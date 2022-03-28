package tidb_or_mysql

type Resp struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ListData struct {
	Total  int64       `json:"total"`
	Items  interface{} `json:"items"`
	More   bool        `json:"more"`
	Offset int         `json:"offset"`
}
