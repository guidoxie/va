package spider

type Spider interface {
	Download(indexName string, indexCode string, isFast bool) ([]byte, error)
	Parser(data []byte) (interface{}, error)
}

const (
	userAgent = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"
)
