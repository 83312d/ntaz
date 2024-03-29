package cmd

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

type ResultData struct {
	bytes float64
	time  float64
}

var mbpsCmd = &cobra.Command{
	Use:   "mbps",
	Short: "Read log and perform traffic calculation.",
	Long:  "Read nginx access log of particular format and calculate mbps for all requests.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		var fp string
		var scanner *bufio.Scanner
		dataMap := make(map[string]*ResultData)
		mut := &sync.Mutex{}
		wg := &sync.WaitGroup{}

		if len(args) > 0 {
			fp = args[0]
			file, err := os.Open(fp)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			defer file.Close()

			info, err := file.Stat()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			} else if info.Size() == 0 {
				fmt.Println("File is empty")
				os.Exit(1)
			}

			scanner = bufio.NewScanner(file)
		} else {
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) != 0 {
				fmt.Println("No data provided")
				os.Exit(1)
			} else {
				scanner = bufio.NewScanner(os.Stdin)
			}
		}

		for scanner.Scan() {
			wg.Add(1)

			go func(line string) {
				defer wg.Done()

				parts := strings.Split(line, " - ")
				reqProcessTime, err := strconv.ParseFloat(parts[0], 64)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				splittedIpAndTime := strings.Split(parts[1], " ")
				reqTime := splittedIpAndTime[1]
				bytesSent, err := strconv.ParseFloat(splittedIpAndTime[2], 64)
				if err != nil {
					fmt.Println(err)
					os.Exit(1)
				}

				mut.Lock()
				if data, ok := dataMap[reqTime]; ok {
					data.bytes += bytesSent
					data.time += reqProcessTime
				} else {
					dataMap[reqTime] = &ResultData{bytes: bytesSent, time: reqProcessTime}
				}
				mut.Unlock()
			}(scanner.Text())
		}

		wg.Wait()

		keys := make([]string, 0, len(dataMap))
		for k := range dataMap {
			keys = append(keys, k)
		}
		slices.Sort(keys)

		for _, k := range keys {
			mbps := (dataMap[k].bytes * 8) / (dataMap[k].time * 1024 * 1024)
			fmt.Printf("%s - %f mbps\n", k, mbps)
		}
	},
}

func init() {
	rootCmd.AddCommand(mbpsCmd)
}
