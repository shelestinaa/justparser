package parser

import (
	"context"
	"fmt"
	db2 "github.com/shelestinaa/justparser/external/db"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"io/ioutil"
	"strings"
	"time"
)

func Parse(buf []byte, filename string) {

	filenameToSave := "files/saved-file.xlsx"
	err := ioutil.WriteFile(filenameToSave, buf, 0777)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	var reportFile *excelize.File

	reportFile, err = excelize.OpenFile(filenameToSave)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	defer func() {
		if err := reportFile.Close(); err != nil {
			logrus.Fatalf(err.Error())
		}
	}()

	ctx, _ := context.WithTimeout(context.Background(), time.Minute)

	db, err := db2.NewClient(ctx, "localhost", "27017", "", "", "justparser", "")
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	sheetMap := reportFile.GetSheetMap()

	err = db.Collection(filename).Drop(ctx)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	err = db.CreateCollection(ctx, filename)
	if err != nil {
		logrus.Fatalf(err.Error())
	}

	for _, sheetTitle := range sheetMap {

		var rowsToPersist []byte

		rows, err := reportFile.GetRows(sheetTitle)
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		for _, row := range rows {
			implodedStringRow := strings.Join(row, ",")
			rowsToPersist = append(rowsToPersist, []byte(implodedStringRow)...)

		}
		type Document struct {
			Title string
			Data  []byte
		}

		newDocument := &Document{sheetTitle, rowsToPersist}

		documentToPersist, err := bson.Marshal(&newDocument)
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		insertResult, err := db.Collection(filename).InsertOne(ctx, documentToPersist)
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		fmt.Printf(fmt.Sprintf("File have been parsed and persisted successfully. Its ID is: %s ", insertResult.InsertedID))
	}
}
