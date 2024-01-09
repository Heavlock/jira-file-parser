package fileParser

import (
	"bufio"
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
	"sync"
)

type result struct {
	titles            map[string]int
	data              []map[int]string
	statusTitleNumber int
}

type Stat struct {
	Name              string
	EmpType           string
	TaskCountComplete int
	PointsCount       float64
	BagsCount         int
	SmallTaskCount    int //story point 1
	MiddleTaskCount   int //story point 5
	BigTaskCount      int //story point 10
}

// empTypes
var qa = "QA"
var developer = "developer"

// issue types
var bagType = "Баг"

var statusFieldName = "Status"
var qaFieldName = "Custom field (QA)"
var storyPointFieldName = "Custom field (Story Points)"
var assigneeFieldName = "Assignee"
var issueTypeFieldName = "Issue Type"
var taskTitle = "Summary"

var titleNames = []string{
	assigneeFieldName,
	storyPointFieldName,
	qaFieldName,
	statusFieldName,
	issueTypeFieldName,
	taskTitle,
}
var taskStatusesToAdd = []string{
	"Публикация",
	"Готово",
	"Приёмка ПОСЛЕ публикации",
}

func ParseFile(filePath string) (map[string]*Stat, error) {
	var statMap = make(map[string]*Stat)

	file, err := os.Open(filePath)
	if err != nil {
		//showErrorDialog(w, "Error", fmt.Sprintf("Error opening file: %v", err))
		return statMap, err
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))

	res := &result{
		titles: make(map[string]int),
		data:   []map[int]string{},
	}
	// парсим csv
	record, err := reader.Read()
	findTitle(res, record, titleNames)
	for {
		record, err = reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}

		setResData(res, record)
	}

	collectQAStat(res, statMap)
	collectDevelopersStat(res, statMap)
	return statMap, nil
}

func collectQAStat(res *result, statMap map[string]*Stat) {
	collectStat(res, qa, statMap)
}

func collectDevelopersStat(res *result, statMap map[string]*Stat) {
	collectStat(res, developer, statMap)
}

func collectStat(res *result, employeeType string, statMap map[string]*Stat) {
	var wg = sync.WaitGroup{}
	var mu = sync.Mutex{}

	var storyPointFieldNumber = res.titles[storyPointFieldName]
	var issueTypeFieldNumber = res.titles[issueTypeFieldName]
	var taskTitleFieldNumber = res.titles[taskTitle]
	var fieldName string
	if employeeType == qa {
		fieldName = qaFieldName
	} else {
		fieldName = assigneeFieldName
	}
	var fieldNameNumber = res.titles[fieldName]

	for _, val := range res.data {
		if val[fieldNameNumber] == "" {
			continue
		}

		wg.Add(1)
		go func(val map[int]string) {
			defer wg.Done()
			storyPoints, err := strconv.ParseFloat(val[storyPointFieldNumber], 64)
			if err != nil {
				storyPoints = 1
			}
			mu.Lock()
			if statObj, _ := statMap[val[fieldNameNumber]]; statObj == nil {
				statVar := &Stat{}
				statMap[val[fieldNameNumber]] = statVar
				statVar.EmpType = employeeType
				statVar.Name = val[fieldNameNumber]
			}
			statMap[val[fieldNameNumber]].TaskCountComplete++
			statMap[val[fieldNameNumber]].PointsCount += storyPoints

			//считаем количество задач с типом баг
			if val[issueTypeFieldNumber] == bagType || strings.Contains(strings.ToLower(val[taskTitleFieldNumber]), "rg") {
				statMap[val[fieldNameNumber]].BagsCount++
			}

			//считаем количество маленьких, больших и средних задач
			if storyPoints >= 5 {
				statMap[val[fieldNameNumber]].MiddleTaskCount++
			} else if storyPoints >= 10 {
				statMap[val[fieldNameNumber]].BigTaskCount++
			} else {
				statMap[val[fieldNameNumber]].SmallTaskCount++
			}
			mu.Unlock()
		}(val)
	}
	wg.Wait()
}

func findTitle(res *result, record []string, titles []string) {
	wg := &sync.WaitGroup{}
	mu := &sync.Mutex{}
	for ind1, val1 := range record {
		if val1 == "" {
			continue
		}
		wg.Add(1)
		go func(wg *sync.WaitGroup, ind1 int, val1 string) {
			defer wg.Done()
			for _, val2 := range titles {
				if !strings.Contains(strings.ToLower(val2), strings.ToLower(val1)) {
					continue
				}
				mu.Lock()
				if strings.Contains(strings.ToLower(statusFieldName), strings.ToLower(val1)) {
					res.statusTitleNumber = ind1
				} else {
					res.titles[val2] = ind1
				}
				mu.Unlock()
			}
		}(wg, ind1, val1)
	}
	wg.Wait()
}

func setResData(res *result, record []string) {
	//имеет ли задача необходимый статус из taskStatusesToAdd
	isMatch := false
	for _, val := range taskStatusesToAdd {
		if strings.Contains(strings.ToLower(val), strings.ToLower(record[res.statusTitleNumber])) {
			isMatch = true
			break
		}
	}
	if isMatch == false {
		return
	}

	dataMap := make(map[int]string)
	for _, ind := range res.titles {
		dataMap[ind] = record[ind]
	}
	res.data = append(res.data, dataMap)
}
