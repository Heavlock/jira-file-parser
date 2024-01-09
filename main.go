package main

import (
	"awesomeProject/fileParser"
	"awesomeProject/fileWriter"
	"fmt"
	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	"log"
	"net/url"
	"strings"
)

type MyFileFilter struct{}

func (m *MyFileFilter) Matches(uri fyne.URI) bool {
	if strings.Contains(uri.String(), ".csv") && strings.Contains(uri.String(), "Jira") {
		return true
	}
	return false
}

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("File Chooser Example")
	myWindow.Resize(fyne.NewSize(600, 400))
	var filePath string
	// Create a label to display the selected file path
	filePathLabel := widget.NewLabel("Здесь будет полнуть путь к выбранному файлу: ")
	textResult := &widget.TextGrid{}
	//downloadButton := widget.NewButton("Скачать результат", func() {
	//	//showDialog(myWindow, filePathLabel, &filePath)
	//})
	//downloadContainer := container.NewHBox(textResult, downloadButton)
	//downloadContainer.Hide()
	linkURL, err := url.Parse("https://eda1-lifemart-goulash.atlassian.net/jira/software/c/projects/GT/issues")
	if err != nil {
		panic(err)
	}

	linkToDownload := widget.NewHyperlink("скачайте выгрузка csv а странице в Jira", linkURL)
	// Create a button to trigger the file chooser dialog
	chooseFileButton := widget.NewButton("Выбрать Файл", func() {
		showDialog(myWindow, filePathLabel, &filePath)
	})

	parseFileButton := widget.NewButton("Обработать файл", func() {
		res, err := fileParser.ParseFile(filePath)
		if err != nil {
			textResult.SetText(err.Error())
		} else {
			textResult.SetText(`готово, файл с результатом рядом с вашей программой`)
			fileWriter.WriteToCSV(res, "result.csv")
			//downloadContainer.Show()
		}
	})

	// Set up the layout
	content := container.NewVBox(
		filePathLabel,
		linkToDownload,
		layout.NewSpacer(),
		//downloadContainer,
		chooseFileButton,
		parseFileButton,
	)

	// Set the content and show the window
	myWindow.SetContent(content)
	myWindow.ShowAndRun()
}

func showDialog(window fyne.Window, label *widget.Label, filePath *string) {

	filter := &MyFileFilter{}

	dialogF := dialog.NewFileOpen(func(file fyne.URIReadCloser, err error) {
		if err != nil {
			log.Println("Error opening file:", err)
			return
		}
		if file == nil {
			return
		}

		fileURL := file.URI().String()

		*filePath = strings.TrimPrefix(fileURL, "file://")
		label.SetText(*filePath)
		err = file.Close()
		if err != nil {
			fmt.Println(err.Error())
		}

	}, window)
	dialogF.SetFilter(filter)
	dialogF.Show()
}
