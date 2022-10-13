package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	db2 "github.com/shelestinaa/justparser/external/db"
	"github.com/shelestinaa/justparser/external/parser"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"golang.org/x/net/context"
	"io"
	"mime/multipart"
	"net/http"
	"strings"
	"time"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
}

type response struct {
	code string
	data []byte
}

func New(config *Config) *APIServer {

	return &APIServer{
		config: config,
		logger: logrus.New(),
		router: mux.NewRouter(),
	}
}

func (s *APIServer) Start() error {

	err := s.configureLogger()
	if err != nil {
		return err
	}

	s.configureRouter()

	s.logger.Info("starting api server is going on as fast as u are running to obtain white drugs, sport")

	return http.ListenAndServe(s.config.BindAddr, s.router)
}

func (s *APIServer) configureLogger() error {
	level, err := logrus.ParseLevel(s.config.LogLevel)

	if err != nil {
		return err
	}

	s.logger.SetLevel(level)

	return nil
}

func (s *APIServer) configureRouter() {
	s.router.HandleFunc("/hello", s.handleHello())
	s.router.HandleFunc("/parse", s.handleParse())
	s.router.HandleFunc("/get-list", s.handleGetCollections())
	s.router.HandleFunc("/get-parsed-file", s.handleGetCollectionByName())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *APIServer) handleParse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		err := r.ParseMultipartForm(32 << 20)
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		var buf bytes.Buffer

		file, header, err := r.FormFile("key")
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		defer func(file multipart.File) {
			err := file.Close()
			if err != nil {
				logrus.Fatalf(err.Error())
			}
		}(file)

		name := strings.Split(header.Filename, ".")
		filename := name[0]
		fmt.Printf("File name %s\n", filename)

		io.Copy(&buf, file)

		fmt.Printf("File %s has been copied successfully\n", filename)

		parser.Parse(buf.Bytes(), filename)
		fmt.Printf("File %s has been parsed successfully\n", filename)

	}
}

func (s *APIServer) handleGetCollections() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), time.Minute)

		db, err := db2.NewClient(ctx, "localhost", "27017", "", "", "justparser", "")
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		collectionNames, err := db.ListCollectionNames(ctx, bson.D{{}})
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		_, err = io.WriteString(w, strings.Join(collectionNames, ", "))
		if err != nil {
			logrus.Fatalf(err.Error())
		}
	}
}

func (s *APIServer) handleGetCollectionByName() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, _ := context.WithTimeout(context.Background(), time.Minute)

		db, err := db2.NewClient(ctx, "localhost", "27017", "", "", "justparser", "")
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		filename := r.URL.Query().Get("name")
		currentCollection := db.Collection(filename)

		cursor, err := currentCollection.Find(ctx, bson.D{{}})
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		type Document struct {
			Title string
			Data  []byte
		}

		type TransformedDocument struct {
			Title string
			Data  string
		}

		type Response struct {
			Code       int
			Filename   string
			ParsedFile []TransformedDocument
		}

		var documents []Document
		var transformedDocuments []TransformedDocument

		err = cursor.All(ctx, &documents)

		for _, document := range documents {

			var documentData string
			documentData = string(document.Data)

			if err != nil {
				logrus.Fatalf(err.Error())
			}

			transformedDocuments = append(transformedDocuments, TransformedDocument{
				Title: document.Title,
				Data:  documentData,
			})
		}

		response := Response{
			Code:       200,
			Filename:   filename,
			ParsedFile: transformedDocuments,
		}

		marshalled, err := json.Marshal(response)
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		_, err = w.Write(marshalled)
		if err != nil {
			logrus.Fatalf(err.Error())
		}
	}
}
