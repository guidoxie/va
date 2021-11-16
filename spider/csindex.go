package spider

import (
	"bytes"
	"fmt"
	"github.com/shakinm/xlsReader/xls"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	CsIndexHost = "www.csindex.com.cn"
	userAgent   = "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/92.0.4515.159 Safari/537.36"
)

// 基本面
type CsFundamental struct {
	Date      string `json:"date"`
	IndexName string `json:"index_name"`
	IndexCode string `json:"index_code"`
	PE        string `json:"pe"`
	DP        string `json:"dp"`
}

type CsIndex struct {
}

func (c CsIndex) Download(indexName string, indexCode string) ([]byte, error) {
	// 下载文件
	reqFile, err := http.NewRequest("GET",
		fmt.Sprintf("https://csi-web-dev.oss-cn-shanghai-finance-1-pub.aliyuncs.com/static/html/csindex/public/uploads/file/autofile/indicator/%sindicator.xls", indexCode), nil)
	if err != nil {
		return nil, err
	}
	reqFile.Header.Add("Referer", fmt.Sprintf("http://www.csindex.com.cn/zh-CN/indices/index-detail/%s", indexCode))
	reqFile.Header.Add("User-Agent", userAgent)
	reqFile.Header.Add("Host", CsIndexHost)

	resp, err := http.DefaultClient.Do(reqFile)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		log.Fatal(err)
	}
	resp.Body.Close()

	return body, nil
}

func (c CsIndex) Parser(data []byte) ([]CsFundamental, error) {
	var res = make([]CsFundamental, 0)
	f, err := xls.OpenReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	sheet, err := f.GetSheet(0)
	if err != nil {
		return nil, err
	}
	for i := 1; i < sheet.GetNumberRows(); i++ {
		row, err := sheet.GetRow(i)
		if err != nil {
			return nil, err
		}
		if row != nil {
			day, _ := row.GetCol(0)
			//  Excel uses Julian dates prior to March 1st 1900
			//date := time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local).
			//	Add(time.Duration(day.GetInt64()-1) * 24 * time.Hour).Format("2006-01-02")
			d, _ := time.Parse("20060102", day.GetString())
			date := d.Format("2006-01-02")
			indexCode, _ := row.GetCol(1)
			IndexName, _ := row.GetCol(3)
			pe2, _ := row.GetCol(7)
			dp2, _ := row.GetCol(9)

			r := CsFundamental{
				Date:      date,
				IndexCode: indexCode.GetString(),
				IndexName: IndexName.GetString(),
				PE:        pe2.GetString(),
				DP:        dp2.GetString(),
			}
			res = append(res, r)
		}
	}
	return res, nil
}
