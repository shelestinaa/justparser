package parser

import (
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
)

func Parse(buf []byte) {

	filename := "files/saved-file.xlsx"
	err := ioutil.WriteFile(filename, buf, 0777)
	if err != nil {
		return
	}

	var reportFile *excelize.File
	reportFile, err = excelize.OpenFile(filename)
	if err != nil {
		return
	}

	defer func() {
		if err := reportFile.Close(); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	sheetMap := reportFile.GetSheetMap()

	//var result []string
	for sheetNumber, sheetTitle := range sheetMap {
		fmt.Sprintf(sheetTitle+": %s", sheetNumber)

	}
	// Тут логика парсера, которую я не придумал
	//todo: узнать про монго, мне кажется коллекции там не просто так, вместо таблиц

}
