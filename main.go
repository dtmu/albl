package main

import (
	"errors"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var reg *regexp.Regexp = regexp.MustCompile(`[0-9]+\_elasticloadbalancing\_.+-.+-[0-9]\_.+\_[0-9]{8}T[0-9]{4}Z\_.*`)

func main() {
	app := cli.NewApp()
	app.Name = "albl"
	app.Usage = "This is cli tool to convert alb logs to xlsx."
	app.Version = "0.0.1"
	app.Action = action

	cdir, _ := os.Getwd()
	app.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "name",
			Aliases: []string{"n"},
			Value:   "alb_access_log.xlsx",
			Usage:   "specified xlsx file name",
		},
		&cli.StringFlag{
			Name:    "directory",
			Aliases: []string{"d"},
			Value:   cdir,
			Usage:   "specified root directory of log file",
		},
	}
	app.Run(os.Args)
}

func action(c *cli.Context) error {
	ex, err := readALBLog(c.String("directory"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	ex.DeleteSheet("Sheet1")
	err = ex.SaveAs(c.String("name"))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

var albLogHeader []string = []string{
	"type",
	"elb",
	"target:port",
	"request_processing_time",
	"target_processing_time",
	"response_processing_time",
	"elb_status_code",
	"target_status_code",
	"received_bytes",
	"sent_bytes",
	"\"request\"",
	"\"user_agent\"",
	"ssl_cipher",
	"ssl_protocol",
	"target_group_arn",
	"\"trace_id\"",
	"\"domain_name\"",
	"\"chosen_cert_arn\"",
	"matched_rule_priority",
	"request_creation_time",
	"\"actions_executed\"",
	"\"redirect_url\"",
	"\"error_reason\"",
	"\"target:port_list\"",
	"\"target_status_code_list\"",
}

func readALBLog(dir string) (*excelize.File, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	ex := excelize.NewFile()
	sheetsEndLine := make(map[string]int)

	for _, file := range files {
		fileName := file.Name()
		if !file.IsDir() {
			if reg.Match([]byte(fileName)) {
				e, err := ioutil.ReadFile(filepath.Join(dir, fileName))
				if err != nil {
					fmt.Println("ReadFile Error: " + err.Error())
					continue
				}
				entries := strings.Split(string(e), "\n")
				ip := strings.Split(fileName, "_")[5]
				if ex.Sheet[ip] == nil {
					i := ex.NewSheet(ip)
					ex.SetActiveSheet(i)
					sheetsEndLine[ip] = 2
				}
				for _, entry := range entries {
					strLine := strconv.Itoa(sheetsEndLine[ip])
					clmNo := 1
					cellValue := ""
					isDblQuote := false
					for _, v := range strings.Split(entry, " ") {
						cellValue += v
						if c := strings.Count(v, "\""); c > 0 {
							if !isDblQuote && c == 1 {
								isDblQuote = true
								v += " "
								continue
							}
							isDblQuote = false
						} else {
							if isDblQuote {
								v += " "
								continue
							}
						}
						clm, _ := excelize.ColumnNumberToName(clmNo)
						clmNo++
						if ex.SetCellValue(ip, clm+strLine, cellValue) != nil {
							fmt.Println("SetCellValue Error: " + cellValue + " to " + clm + strLine + " in \"" + ip + "\" sheet")
						}
						cellValue = ""
					}
					if sheetsEndLine[ip]-1 == 1 {
						for i, v := range albLogHeader {
							clm, _ := excelize.ColumnNumberToName(i + 1)
							if ex.SetCellValue(ip, clm+"1", v) != nil {
								return nil, errors.New("SetLogHeader Error: setting alb log header to xlsx is faild.")
							}
						}
					}
					sheetsEndLine[ip]++
				}
			}
		}
	}

	for _, v := range sheetsEndLine {
		if v != 1 {
			return ex, nil
		}
	}

	return nil, errors.New("NoLogEntryError: " + "there are not log entories")
}
