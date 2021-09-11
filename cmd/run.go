package cmd

import (
	"github.com/guidoxie/va/spider"
	"github.com/modood/table"
	"log"
	"runtime"
	"sort"
	"sync"
	"time"
)

func run() error {
	var (
		wg     sync.WaitGroup
		lock   sync.Mutex
		output = make([]spider.Fundamental, 0)
		// 控制并发数量，为cpu核心数
		ch = make(chan struct{}, runtime.NumCPU())
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
	}

	sort.Slice(output, func(i, j int) bool {
		return output[i].IndexCode < output[j].IndexCode
	})
	table.Output(output)
	return nil
}
