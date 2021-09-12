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

func run() error {
	var (
		wg     sync.WaitGroup
		bond   *spider.Fundamental
		lock   sync.Mutex
		output = make([]spider.Fundamental, 0)
		// 控制并发数量，10
		ch = make(chan struct{}, 10)
		cs = spider.CsIndex{}
	)

	if len(code) > 0 {
		body, err := cs.Download("", code, fast)
		if err != nil {
			log.Fatal(err)
		}
		fs, err := cs.Parser(body)
		if err != nil {
			return err
		}
		if len(date) > 0 {
			for _, f := range fs.([]spider.Fundamental) {
				if f.Date == date {
					output = append(output, f)
				}
			}
		} else {
			output = fs.([]spider.Fundamental)
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
			bond = &spider.Fundamental{
				Date:      date,
				IndexName: "十年期国债",
				IndexCode: "--",
				EP:        b,
				PE2:       "--",
				DP2:       "--",
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
				body, err := cs.Download(name, code, fast)
				if err != nil {
					log.Fatalf("%v %v %v", name, code, err)
				}
				fs, err := cs.Parser(body)
				if err != nil {
					log.Fatalf("%v %v %v", name, code, err)
				}
				for _, f := range fs.([]spider.Fundamental) {
					if f.Date == date {
						lock.Lock()
						output = append(output, f)
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
	table.Output(output)
	return nil
}

// 计算ep
func epCalc(fs []spider.Fundamental) error {
	for i := range fs {
		if spider.EpCode[fs[i].IndexCode] && fs[i].PE2 != "--" {
			pe2, err := strconv.ParseFloat(fs[i].PE2, 64)
			if err != nil {
				return err
			}
			fs[i].EP = fmt.Sprintf("%.2f", Round((1/pe2)*100, 2))
		} else {
			fs[i].EP = "--"
		}
	}
	return nil
}

// 保留小数点后n位
func Round(f float64, n int) float64 {
	n10 := math.Pow10(n)
	return math.Trunc((f+0.5/n10)*n10) / n10
}
