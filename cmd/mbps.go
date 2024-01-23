package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var Path string

type ResultData struct {
	bytes float64
	time  float64
}

var mbpsCmd = &cobra.Command{
	Use:   "mbps",
	Short: "Read log and perform traffic calculation.",
	Long:  "Read nginx access log of particular format and calculate mbps for all requests. Default log path path is /var/log/nginx/access.log",
	Run: func(cmd *cobra.Command, args []string) {
		var fp string
		var logArr [][]string
		dataMap := make(map[string]*ResultData)

		fl, _ := cmd.Flags().GetString("path")
		if fl == "" {
			fp = "/var/log/nginx/access.log"
		} else {
			fp = fl
		}

		file, err := os.Open(fp)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		for scanner.Scan() {
			line := scanner.Text()
			part := strings.Split(line, " - ")
			logArr = append(logArr, part)
		}

		for _, line := range logArr {
			reqProcessTime, _ := strconv.ParseFloat(line[0], 64)
			splittedIpAndTime := strings.Split(line[1], " ")
			reqTime := splittedIpAndTime[1]
			bytesSent, _ := strconv.ParseFloat(splittedIpAndTime[2], 64)

			if data, ok := dataMap[reqTime]; ok {
				data.bytes += bytesSent
				data.time += reqProcessTime
			} else {
				dataMap[reqTime] = &ResultData{bytes: bytesSent, time: reqProcessTime}
			}
		}

		for k, v := range dataMap {
			mbps := (v.bytes * 8) / (v.time * 1024 * 1024)
			fmt.Printf("%s - %f mbps\n", k, mbps)
		}
	},
}

func init() {
	rootCmd.AddCommand(mbpsCmd)

	mbpsCmd.Flags().String("path", "", "path to log file")
}
