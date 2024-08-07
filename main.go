package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("ToDo list")
	})
	mux.HandleFunc("GET /list", checkAuth(getToDoList))
	mux.HandleFunc("POST /task", checkAuth(postTask))
	mux.HandleFunc("PUT /task", checkAuth(editTask))
	mux.HandleFunc("DELETE /task", checkAuth(deleteTask))
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		fmt.Println("Error happened", err.Error())
		return
	}
}

type Authorisation struct {
	UserName string
	Password string
}

var User1 = Authorisation{
	UserName: "Mona",
	Password: "42",
}

var User2 = Authorisation{
	UserName: "Liza",
	Password: "315",
}

func checkAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		if (username != User1.UserName || password != User1.Password) && (username != User2.UserName || password != User2.Password) {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

type TaskManager struct {
	ID          int
	Description string
}

var tasks = []TaskManager{
	{ID: 1, Description: "Open computer"},
	{ID: 2, Description: "Do homework"},
	{ID: 3, Description: "Close computer"},
}

func getToDoList(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func postTask(w http.ResponseWriter, r *http.Request) {
	var newTask TaskManager
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	tasks = append(tasks, newTask)
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func editTask(w http.ResponseWriter, r *http.Request) {

	var updatedTask TaskManager
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, task := range tasks {
		if task.ID == updatedTask.ID {
			tasks[i].Description = updatedTask.Description
			break
		}
	}
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	var taskToDelete TaskManager
	err := json.NewDecoder(r.Body).Decode(&taskToDelete)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	for i, task := range tasks {
		if task.ID == taskToDelete.ID {
			tasks = append(tasks[:i], tasks[i+1:]...)
			break
		}
	}
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
