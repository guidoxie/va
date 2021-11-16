package cmd

import (
	"fmt"
	"github.com/guidoxie/va/spider"
	"github.com/spf13/cobra"
	"log"
)

var desc = fmt.Sprintf("指数估值表，不构成投资建议，投资有风险，入市需谨慎\n\n数据来源:中证指数官网（%s）", spider.CsIndexHost)

var (
	date    string // -d
	vaTable bool   // -t
	code    string // -c
)

func init() {
	root.Flags().StringVarP(&date, "date", "d", "", "日期，格式：2006-12-28")
	root.Flags().BoolVarP(&vaTable, "table", "t", true, "常见指数估值表，日期默认为昨天")
	root.Flags().StringVarP(&code, "code", "c", "", "指数代码")
}

var root = &cobra.Command{
	Use:  "",
	Long: desc,
	Run: func(cmd *cobra.Command, args []string) {
		if err := run(); err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() error {
	return root.Execute()
}
