package spider

import (
	"bytes"
	"fmt"
	"github.com/shakinm/xlsReader/xls"
	"io/ioutil"
	"net/http"
)

var chinaBondUrl = "http://www.csindex.com.cn/zh-CN/bond-valuation/bond-yield-curve?type=2&line_id=1&line_date=1&line_type=2&start_date=%s&end_date=undefined&download=1"

// 获取10年国债收益率
func Get10BondEP(date string) (string, error) {

	req, err := http.NewRequest("GET", fmt.Sprintf(chinaBondUrl, date), nil)
	if err != nil {
		return "", err
	}
	req.Header.Add("Referer", "http://www.csindex.com.cn/zh-CN/bond-valuation/bond-yield-curve")
	req.Header.Add("User-Agent", userAgent)
	req.Header.Add("Host", CsIndexHost)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(data) == 0 {
		return "", nil
	}
	f, err := xls.OpenReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	sheet, err := f.GetSheet(0)
	if err != nil {
		return "", err
	}
	for i := 1; i < sheet.GetNumberRows(); i++ {
		row, err := sheet.GetRow(i)
		if err != nil {
			return "", err
		}
		if row != nil {
			y, _ := row.GetCol(2)
			if y.GetString() == "10" {
				v, _ := row.GetCol(4)
				return v.GetString(), nil
			}

		}

	}
	return "", nil

}
