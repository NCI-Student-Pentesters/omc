package helpers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"omc/models"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
	"k8s.io/client-go/util/jsonpath"
)

// TYPES
type Contexts []models.Context

// CONSTS
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// VARS
var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

// FUNCS
func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func RandString(length int) string {
	return StringWithCharset(length, charset)
}

func PrintTable(headers []string, data [][]string) {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(headers)
	table.SetAutoWrapText(false)
	table.SetAutoFormatHeaders(true)
	table.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	table.SetAlignment(tablewriter.ALIGN_LEFT)
	table.SetCenterSeparator("")
	table.SetColumnSeparator("")
	table.SetRowSeparator("")
	table.SetHeaderLine(false)
	table.SetBorder(false)
	table.SetTablePadding("   ")
	table.SetNoWhiteSpace(true)
	table.AppendBulk(data)
	table.Render()
}

func FormatDiffTime(diff time.Duration) string {
	if diff.Hours() > 24 {
		if diff.Hours() > 2160000 {
			return "Unknown"
		}
		return strconv.Itoa(int(diff.Hours()/24)) + "d"
	}
	if diff.Minutes() > 60 {
		var hours float64
		hours = diff.Minutes() / 60
		remainMinutes := int(diff.Minutes()) % 60
		if remainMinutes > 0 {
			return strconv.Itoa(int(hours)) + "h" + strconv.Itoa(remainMinutes) + "m"
		}
		return strconv.Itoa(int(hours)) + "h"

	}
	if diff.Seconds() > 60 {
		var minutes float64
		minutes = diff.Seconds() / 60
		remainSeconds := int(diff.Seconds()) % 60
		if remainSeconds > 0 && diff.Minutes() < 4 {
			return strconv.Itoa(int(minutes)) + "m" + strconv.Itoa(remainSeconds) + "s"
		}
		return strconv.Itoa(int(minutes)) + "m"

	}
	return strconv.Itoa(int(diff.Seconds())) + "s"
}

func ExecuteJsonPath(data interface{}, jsonPathTemplate string) {
	buf := new(bytes.Buffer)
	jPath := jsonpath.New("out")
	jPath.AllowMissingKeys(false)
	jPath.EnableJSONOutput(false)
	err := jPath.Parse(jsonPathTemplate)
	if err != nil {
		fmt.Println("error: error parsing jsonpath " + jsonPathTemplate + ", " + err.Error())
		os.Exit(1)
	}
	jPath.Execute(buf, data)
	fmt.Print(buf)
}

func CreateConfigFile(homePath string) {
	config := models.Config{}
	file, _ := json.MarshalIndent(config, "", " ")
	cfgFilePath := homePath + "/.omc.json"
	_ = ioutil.WriteFile(cfgFilePath, file, 0644)
}

func GetData(data [][]string, allNamespacesFlag bool, showLabels bool, labels string, outputFlag string, column int32, _list []string) [][]string {
	var toAppend []string
	if allNamespacesFlag == true {
		if outputFlag == "" {
			toAppend = _list[0:column] // -A
		}
		if outputFlag == "wide" {
			toAppend = _list // -A -o wide
		}
	} else {
		if outputFlag == "" {
			toAppend = _list[1:column]
		}
		if outputFlag == "wide" {
			toAppend = _list[1:] // -o wide
		}
	}

	if showLabels {
		toAppend = append(toAppend, labels)
	}
	data = append(data, toAppend)
	return data
}

func ExtractLabels(_labels map[string]string) string {
	labels := ""
	for k, v := range _labels {
		labels += k + "=" + v + ","
	}
	if labels == "" {
		labels = "<none>"
	} else {
		labels = strings.TrimRight(labels, ",")
	}
	return labels
}
