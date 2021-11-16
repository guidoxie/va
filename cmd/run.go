package cmd

import (
	"fmt"
	"github.com/guidoxie/va/spider"
	"github.com/modood/table"
	"log"
	"math"
	"sort"
	"strconv"
	"sync"
	"time"
)

type OutPut struct {
	Date      string `json:"date" table:"日期"`
	IndexName string `json:"index_name" table:"名称"`
	IndexCode string `json:"index_code" table:"代码"`
	EP        string `json:"ep" table:"盈利收益率(%)"`
	PE        string `json:"pe" table:"市盈率"`
	DP        string `json:"dp" table:"股息率(%)"`
}

func run() error {
	var (
		wg     sync.WaitGroup
		bond   *OutPut
		lock   sync.Mutex
		output = make([]OutPut, 0)
		// 控制并发数量，10
		ch = make(chan struct{}, 10)
		cs = spider.CsIndex{}
	)

	if len(code) > 0 {
		body, err := cs.Download("", code)
		if err != nil {
			log.Fatal(err)
		}
		fs, err := cs.Parser(body)
		if err != nil {
			return err
		}
		for _, f := range fs {
			if len(date) > 0 && f.Date == date {
				output = append(output, OutPut{
					Date:      f.Date,
					IndexName: f.IndexName,
					IndexCode: f.IndexCode,
					PE:        f.PE,
					DP:        f.DP,
				})
			} else {
				output = append(output, OutPut{
					Date:      f.Date,
					IndexName: f.IndexName,
					IndexCode: f.IndexCode,
					PE:        f.PE,
					DP:        f.DP,
				})
			}
		}
	} else if vaTable {
		if len(date) == 0 {
			date = time.Now().Add(-24 * time.Hour).Format("2006-01-02")
		}
		wg.Add(1)

		go func() {
			defer wg.Done()
			// 10年国债收益率
			b, err := spider.Get10BondEP(date)
			if err != nil {
				log.Fatal(err)
			}
			bond = &OutPut{
				Date:      date,
				IndexName: "十年期国债",
				EP:        b,
			}
		}()
		for n, c := range spider.IndexCode {
			wg.Add(1)
			ch <- struct{}{}
			go func(name, code string) {
				defer func() {
					<-ch
					wg.Done()
				}()
				body, err := cs.Download(name, code)
				if err != nil {
					log.Fatalf("Download %v %v %v", name, code, err)
				}
				fs, err := cs.Parser(body)
				if err != nil {
					log.Fatalf("Parser %v %v %v", name, code, err)
				}
				for _, f := range fs {
					if f.Date == date {
						lock.Lock()
						output = append(output, OutPut{
							Date:      f.Date,
							IndexName: f.IndexName,
							IndexCode: f.IndexCode,
							PE:        f.PE,
							DP:        f.DP,
						})
						lock.Unlock()
						break
					}
				}
			}(n, c)
		}
		wg.Wait()
		sort.Slice(output, func(i, j int) bool {
			return output[i].IndexCode < output[j].IndexCode
		})

	}
	if err := epCalc(output); err != nil {
		return err
	}
	if bond != nil {
		output = append(output, *bond)
	}
	table.Output(before(output))
	return nil
}

// 计算ep
func epCalc(fs []OutPut) error {
	for i := range fs {
		if spider.EpCode[fs[i].IndexCode] {
			pe2, err := strconv.ParseFloat(fs[i].PE, 64)
			if err != nil {
				return err
			}
			fs[i].EP = fmt.Sprintf("%.2f", Round((1/pe2)*100, 2))
		}
	}
	return nil
}

func before(fs []OutPut) []OutPut {
	for i := range fs {
		for _, s := range []*string{&fs[i].Date, &fs[i].IndexName, &fs[i].IndexCode, &fs[i].EP,
			&fs[i].PE, &fs[i].DP} {
			if len(*s) == 0 {
				*s = "--"
			}
		}
	}
	return fs
}

// 保留小数点后n位
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}
