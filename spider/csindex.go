package spider

import (
	"bytes"
	"fmt"
	"github.com/shakinm/xlsReader/xls"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

const (
	CsIndexHost       = "www.csindex.com.cn"
	csIndexStatsFiles = "http://www.csindex.com.cn/zh-CN/stats/files"
)

// 基本面
type Fundamental struct {
	Date      string `json:"日期" table:"日期"`
	IndexName string `json:"index_name" table:"名称"`
	IndexCode string `json:"index_code" table:"代码"`
	EP        string `json:"ep" table:"盈利收益率(%)"`
	//PE1       string `json:"pe_1" table:"(总股本)P/E"`
	PE2 string `json:"pe_2" table:"市盈率"`
	//DP1 string `json:"dp_1" table:"(总股本)D/P"`
	DP2      string `json:"dp_2" table:"股息率(%)"`
	Proposal string `json:"proposal" table:"建议"`
}

type CsIndex struct {
}

func (c CsIndex) Download(indexName string, indexCode string, isFast bool) ([]byte, error) {
	timestamp := time.Now().Unix()
	if !isFast {
		// 模仿人为点击
		req, err := http.NewRequest("POST", csIndexStatsFiles, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Add("Referer", fmt.Sprintf("http://www.csindex.com.cn/zh-CN/indices/index-detail/%s", indexCode))
		req.Header.Add("User-Agent", userAgent)
		req.Header.Add("Host", CsIndexHost)
		req.PostForm = url.Values{}
		req.PostForm.Add("href",
			fmt.Sprintf("http://www.csindex.com.cn/uploads/file/autofile/indicator/%sindicator.xls?t=%d", indexCode, timestamp))
		req.PostForm.Add("title", fmt.Sprintf("%s估值", indexName))

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()
	}
	// 下载文件
	reqFile, err := http.NewRequest("GET",
		fmt.Sprintf("http://www.csindex.com.cn/uploads/file/autofile/indicator/%sindicator.xls?t=%d", indexCode, timestamp), nil)
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

func (c CsIndex) Parser(data []byte) (interface{}, error) {
	var res = make([]Fundamental, 0)
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
			date := time.Date(1900, 1, 1, 0, 0, 0, 0, time.Local).
				Add(time.Duration(day.GetInt64()-1) * 24 * time.Hour).Format("2006-01-02")
			indexCode, _ := row.GetCol(1)
			IndexName, _ := row.GetCol(3)
			//pe1, _ := row.GetCol(6)
			pe2, _ := row.GetCol(7)
			//dp1, _ := row.GetCol(8)
			dp2, _ := row.GetCol(9)

			r := Fundamental{
				Date:      date,
				IndexCode: indexCode.GetString(),
				IndexName: IndexName.GetString(),
				//PE1:       pe1.GetString(),
				PE2: pe2.GetString(),
				//DP1:       dp1.GetString(),
				DP2: dp2.GetString(),
			}
			//if len(r.PE1) == 0 {
			//	r.PE1 = "--"
			//}
			if len(r.PE2) == 0 {
				r.PE2 = "--"
			}
			//if len(r.DP1) == 0 {
			//	r.DP1 = "--"
			//}
			if len(r.DP2) == 0 {
				r.DP2 = "--"
			}
			res = append(res, r)
		}
	}
	return res, nil
}
