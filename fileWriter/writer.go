package fileWriter

import (
	"awesomeProject/fileParser"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func WriteToCSV(data map[string]*fileParser.Stat, filename string) error {
	sliceData := convertDataToSlice(data)
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	//BOM
	bomUtf8 := []byte{0xEF, 0xBB, 0xBF}
	file.Write(bomUtf8)

	// Write data to the CSV file
	err = writer.WriteAll(sliceData)
	if err != nil {
		return err
	}
	return nil
}
func convertDataToSlice(data map[string]*fileParser.Stat) [][]string {
	var sliceData [][]string

	sliceData = append(sliceData, []string{
		"Исполнитель",
		"Должность",
		"Количество выполненных задач",
		"Количество поинтов",
		"Количество выполненных задач с типом баг",
		"Количество выполненных задач с 1 поинтом",
		"Количество выполненных задач с 5 поинтами",
		"Количество выполненных задач с 10 поинтами",
	})
	for _, val := range data {
		item := []string{
			val.Name,
			val.EmpType,
			fmt.Sprintf("%d", val.TaskCountComplete),
			fmt.Sprintf("%.2f", val.PointsCount),
			fmt.Sprintf("%d", val.BagsCount),
			fmt.Sprintf("%d", val.SmallTaskCount),
			fmt.Sprintf("%d", val.MiddleTaskCount),
			fmt.Sprintf("%d", val.BigTaskCount),
			fmt.Sprintf("%s", strings.Join(val.TaskTitles[:], "\n===================\n")),
		}
		sliceData = append(sliceData, item)
	}
	return sliceData
}
