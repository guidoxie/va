package spider

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

var chinaBondUrl = "https://www.csindex.com.cn/csindex-home/data/curve/CurveByTradingDateAndCurveCnNameList?tradingDate=%v&bondCurveNameList=cc_ll_gz"

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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	if len(body) == 0 {
		return "", nil
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(body, &m); err != nil {
		return "", err
	}
	data, ok := m["data"].([]interface{})
	if ok {
		for _, d := range data {
			row := d.(map[string]interface{})
			if fmt.Sprintf("%v", row["year"]) == "10" {
				return fmt.Sprintf("%v", row["ytm"]), nil
			}
		}
	}
	return "", nil

}
