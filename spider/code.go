package spider

/*
todo
H股指数
恒生指数
红利机会
深证成指
香港中小
50AH优选
深圳100
纳斯达克100
标普500
创业板: "399006"
*/

var (
	IndexCode = map[string]string{}
	EpCode    = map[string]bool{}
	PeCode    = map[string]bool{}
	PbCode    = map[string]bool{}
)

func init() {
	for t, m := range map[string]map[string]string{
		"ep": epIndicator,
		"pe": peIndicator,
		"pb": pbIndicator,
	} {
		for k, v := range m {
			switch t {
			case "ep":
				EpCode[v] = true
			case "pe":
				PeCode[v] = true
			case "pb":
				PbCode[v] = true
			}
			IndexCode[k] = v
		}
	}

}

// 适用盈利收益率法估值
var epIndicator = map[string]string{
	"300价值": "000919",
	"上证50":  "000016",
	"上证红利":  "000015",
	"中证红利":  "000922",
	"基本面50": "000925",
	"上证180": "000010",
}

// 适用博格公式法（市盈率）估值
var peIndicator = map[string]string{
	"沪深300":  "000300",
	"中证500":  "000905",
	"养老产业":   "399812",
	"可选消费":   "000989",
	"医药100":  "000978",
	"中证消费":   "000932",
	"500低波动": "930782",
	"科创50":   "000688",
	"深证F120": "399702",
	"深证F60":  "399701",
}

// 适用博格公式的变种（市净率）
var pbIndicator = map[string]string{
	"证券公司": "399975",
	"中证银行": "399986",
}
