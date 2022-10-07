package parser

import (
	"context"
	"fmt"
	db2 "github.com/shelestinaa/justparser/external/db"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"strings"
	"time"
)

func Parse(buf []byte, filename string) error {

	filenameToSave := "files/saved-file.xlsx"
	err := ioutil.WriteFile(filenameToSave, buf, 0777)
	if err != nil {
		return err
	}

	var reportFile *excelize.File
	reportFile, err = excelize.OpenFile(filenameToSave)
	if err != nil {
		return err
	}

	defer func() {
		if err := reportFile.Close(); err != nil {
			fmt.Printf(err.Error())
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	db, err := db2.NewClient(ctx, "localhost", "27017", "", "", "justparser", "")
	if err != nil {
		return err
	}

	sheetMap := reportFile.GetSheetMap()

	//var result []string
	for _, sheetTitle := range sheetMap {

		var rowsToPersist []string

		err := db.Collection(filename).Drop(ctx)
		if err != nil {
			return err
		}

		err = db.CreateCollection(ctx, filename)
		if err != nil {
			return err
		}

		rows, _ := reportFile.GetRows(sheetTitle)

		for _, row := range rows {
			implodedStringRow := strings.Join(row, ",")
			//encodedRow, arr, err := bson.MarshalValue(implodedStringRow)
			//if err != nil {
			//	return err
			//}
			rowsToPersist = append(rowsToPersist, implodedStringRow)

		}
		type Document struct {
			title string
			data  []byte
		}

		newDocument := &Document{sheetTitle, []byte(rowsToPersist)}
		documentToPersist, _ := bson.Marshal(newDocument)
		insertResult, err := db.Collection(filename).InsertOne(ctx, documentToPersist)
		if err != nil {
			return err
		}
		fmt.Printf(fmt.Sprintf("File have been parsed and persisted successfully. Its ID is: %s ", insertResult.InsertedID))
	}

	// Тут логика парсера, которую я не придумал
	//todo: узнать про монго, мне кажется коллекции там не просто так, вместо таблиц

	/*
		Двигаемся слева направо по столбцам *1, как только теряем значение, конец
		Затем двигаемся по буквам
	*/

	return nil
}
