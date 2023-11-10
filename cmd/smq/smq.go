package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"smq/pkg/entity"
	"smq/pkg/init/service"

	"github.com/gorilla/mux"
)

func main() {

	messageQueue, port := service.InitServiceDependencies()

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/push/{topic}", addMessageHandler(messageQueue)).Methods("POST")
	r.HandleFunc("/pop/{topic}", popMessageHandler(messageQueue)).Methods("GET")
	r.HandleFunc("/peek/{topic}", peekMessageHandler(messageQueue)).Methods("GET")
	r.HandleFunc("/peekall/{topic}", peekAllMessageHandler(messageQueue)).Methods("GET")
	r.HandleFunc("/clear/{topic}", clear(messageQueue)).Methods("DELETE")

	log.Println(fmt.Sprintf("Listening on port %s", port))

	err := http.ListenAndServe(fmt.Sprintf(":%s", port), r)
	if err != nil {
		panic(err)
	}

}

func addMessageHandler(mq *entity.MessageQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := getTopic(w, r)
		message, err := io.ReadAll(r.Body)
		log.Println("Received message... ", string(message))
		if err != nil {
			panic(err)
		}
		mq.Push(topic, message)
	}
}

func popMessageHandler(mq *entity.MessageQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := getTopic(w, r)
		res := mq.Pop(topic)
		if res != nil {
			log.Println("Pop Message ... ", string(*res))
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

func peekMessageHandler(mq *entity.MessageQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := getTopic(w, r)
		res := mq.Peek(topic)
		if res != nil {
			log.Println("Peek Message ... ", string(*res))
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

func peekAllMessageHandler(mq *entity.MessageQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := getTopic(w, r)
		res := mq.PeekAll(topic)
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

func clear(mq *entity.MessageQueue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		topic := getTopic(w, r)
		mq.Clear(topic)
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
