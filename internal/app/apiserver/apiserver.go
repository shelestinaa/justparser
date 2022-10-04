package apiserver

import (
	"bytes"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/shelestinaa/justparser/external/parser"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

type APIServer struct {
	config *Config
	logger *logrus.Logger
	router *mux.Router
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
	s.router.HandleFunc("/get-list", s.handleGetFilesList())
	s.router.HandleFunc("/get-parsed-file", s.handleGetFileById())
}

func (s *APIServer) handleHello() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}

func (s *APIServer) handleParse() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseMultipartForm(32 << 20)
		var buf bytes.Buffer
		file, header, err := r.FormFile("key")
		if err != nil {
			logrus.Fatalf(err.Error())
			io.WriteString(w, err.Error())
			panic(err)
		}
		defer file.Close()
		name := strings.Split(header.Filename, ".")
		fmt.Printf("File name %s\n", name[0])

		//тут забираю данные в свой буфер
		io.Copy(&buf, file)
		fmt.Printf("File name %s has been copied successfully\n", name[0])

		//тут вызываем парселку, отдаём в неё buf
		parser.Parse(buf.Bytes())

	}
}

func (s *APIServer) handleGetFilesList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "БАРЭФ")
	}
}

func (s *APIServer) handleGetFileById() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "SALAM")
	}
}
