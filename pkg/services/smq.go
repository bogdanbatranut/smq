package services

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"smq/pkg/domains/messagequeue"
	"smq/pkg/domains/simplemq"
	"smq/pkg/smqerrors"

	"github.com/gorilla/mux"
)

type SMQServiceConfiguration func(s *SMQService) error

type SMQService struct {
	configuration IConfig
	messageQueue  messagequeue.IMessageQueueRepository
}

func NewSMQService(cfgs ...SMQServiceConfiguration) (*SMQService, error) {
	s := &SMQService{}
	for _, cfg := range cfgs {
		err := cfg(s)
		if err != nil {
			return nil, err
		}
	}
	return s, nil
}

func WithViperCfg() SMQServiceConfiguration {
	return func(s *SMQService) error {
		cfg, err := NewViperConfig()
		smqerrors.Panic(err)
		s.configuration = cfg
		return nil
	}
}

func WithSimpleMessageQueue() SMQServiceConfiguration {
	return func(s *SMQService) error {
		messageQueue := simplemq.NewSimpleMessageQueueRepository()
		s.messageQueue = messageQueue
		return nil
	}
}

func (s SMQService) Start() {
	port := s.configuration.GetString(AppPort)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/push/{topic}", s.addMessageHandler()).Methods("POST")
	r.HandleFunc("/pop/{topic}", s.popMessageHandler()).Methods("GET")
	r.HandleFunc("/peek/{topic}", s.peekMessageHandler()).Methods("GET")
	r.HandleFunc("/peekall/{topic}", s.peekAllMessageHandler()).Methods("GET")
	r.HandleFunc("/topics", s.getTopics()).Methods("GET")
	r.HandleFunc("/clear/{topic}", s.clear()).Methods("DELETE")

	log.Println(fmt.Sprintf("listening on port %s", port))

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		panic(err)
	}
}

func (s SMQService) getTopics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		keys := s.messageQueue.GetTopics()

		type topicsResults struct {
			Topics []string `json:"topics"`
		}

		res := topicsResults{
			Topics: keys,
		}
		bytes, err := json.Marshal(&res)
		if err != nil {
			panic(err)
		}
		_, err = w.Write(bytes)
		if err != nil {
			panic(err)
		}
	}
}

func (s SMQService) addMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := getTopic(w, r)
		message, err := io.ReadAll(r.Body)
		if err != nil {
			panic(err)
		}
		s.messageQueue.Push(topic, message)
	}
}

func (s SMQService) popMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		topic := getTopic(w, r)
		res := s.messageQueue.Pop(topic)
		if res != nil {
			//log.Println("Pop Message ... ", string(*res))
			_, err := w.Write(*res)
			if err != nil {
				panic(err)
			}
			return
		}
		_, err := w.Write(nil)
		if err != nil {
			panic(err)
		}
	}
}

func (s SMQService) peekMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		topic := getTopic(w, r)
		res := s.messageQueue.Peek(topic)
		if res != nil {
			//log.Println("Peek Message ... ", string(*res))
			_, err := w.Write(*res)
			if err != nil {
				panic(err)
			}
			return
		}
		_, err := w.Write(nil)
		if err != nil {
			panic(err)
		}
	}
}

func (s SMQService) peekAllMessageHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		topic := getTopic(w, r)
		res := s.messageQueue.PeekAll(topic)
		if res != nil {
			_, err := w.Write(*res)
			if err != nil {
				panic(err)
			}
			return
		}
		_, err := w.Write(nil)
		if err != nil {
			panic(err)
		}
	}
}

func (s SMQService) clear() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		topic := getTopic(w, r)
		s.messageQueue.Clear(topic)
		_, err := w.Write(nil)
		if err != nil {
			panic(err)
		}
	}
}

func getTopic(w http.ResponseWriter, r *http.Request) string {
	vars := mux.Vars(r)
	topic, ok := vars["topic"]
	if !ok {
		http.Error(w, "invalid topic", http.StatusBadRequest)
	}
	return topic
}
