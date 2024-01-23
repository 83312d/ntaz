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

var calcCmd = &cobra.Command{
	Use:   "calc",
	Short: "Read log and perform traffic calculation.",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		var fp string

		fl, _ := cmd.Flags().GetString("path")
		if fl == "" {
			fp = "/home/cdn/log/access.log"
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
		var logArr [][]string

		for scanner.Scan() {
			line := scanner.Text()
			part := strings.Split(line, " - ")
			logArr = append(logArr, part)
		}

		dataMap := make(map[string]*ResultData)

		for _, line := range logArr {
			reqProcessTime,_ := strconv.ParseFloat(line[0], 64)
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
	rootCmd.AddCommand(calcCmd)

	calcCmd.Flags().String("path", "", "path to log file")
}
