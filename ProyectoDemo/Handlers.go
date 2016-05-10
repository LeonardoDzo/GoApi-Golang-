package main
import (
	"encoding/json"
	"fmt"
	"net/http"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome!")
}
func GetAllTodos(w http.ResponseWriter, r *http.Request) {
	c := getSession()
	var todo []Todo

	err := c.DB("Demo").C("Todos").Find(bson.M{}).All(&todo)
	if err != nil {
		w.WriteHeader(404);
		return
	}
	uj, _ := json.Marshal(todo)
	//Write status
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}
func GetTodobyId(w http.ResponseWriter, r *http.Request) {
	//var session *mgo.Session
	//session= getSession()
	vars := mux.Vars(r)
	todoId := vars["todoId"]
	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(todoId) {
		w.WriteHeader(404)
		return
	}
	// Grab id
	id := bson.ObjectIdHex(todoId)
	c := getSession()
	todo := Todo{}

	//err := c.Find(bson.M{"name": id}).One(&todo)
        if err := c.DB("Demo").C("Todos").FindId(id).One(&todo);err != nil {
	        w.Header().Set("Content-Type","application/text")
		w.WriteHeader(404);
	        fmt.Fprintf(w, "%s", "Todo no encontrado")
		return
        }
	uj, _ := json.Marshal(todo)
	//Write status
	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(200)
	fmt.Fprintf(w, "%s", uj)
}
func InsertTodo(w http.ResponseWriter, r *http.Request) {
	var todo Todo
	var session *mgo.Session
	session= getSession()
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}
	json.NewDecoder(r.Body).Decode(&todo)
	todo.Id = bson.NewObjectId()
	session.DB("Demo").C("Todos").Insert(todo)

	uj, _ := json.Marshal(todo)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(201)
	fmt.Fprintf(w, "%s", uj)

}
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	var session *mgo.Session
	session= getSession()
	vars := mux.Vars(r)
	todoId := vars["todoId"]

	// Verify id is ObjectId, otherwise bail
	if !bson.IsObjectIdHex(todoId) {
		w.WriteHeader(404)
		return
	}

	// Grab id
	oid := bson.ObjectIdHex(todoId)

	// Remove Todo
	if err := session.DB("Demo").C("Todos").RemoveId(oid); err != nil {
		w.WriteHeader(404)
		return
	}

	// Write status
	w.WriteHeader(200)

}
func UpdateTodo(w http.ResponseWriter, r *http.Request)  {
	var todo Todo

	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		panic(err)
	}
	if err := r.Body.Close(); err != nil {
		panic(err)
	}
	if err := json.Unmarshal(body, &todo); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			panic(err)
		}
	}

	json.NewDecoder(r.Body).Decode(&todo)
	c := getSession()

	colQuerier := bson.M{"_id": todo.Id}
	change := bson.M{"$set": bson.M{"name": todo.Name, "completed": todo.Completed, "due": todo.Due}}
	err = c.DB("Demo").C("Todos").Update(colQuerier, change)
	if err != nil {
		panic(err)
	}
	w.WriteHeader(200)

}
func getSession() *mgo.Session {
	// Connect to our local mongo
	s, err := mgo.Dial("mongodb://localhost")
	// Check if connection error, is mongo running?
	if err != nil {
		panic(err)
	}
	// Deliver session
	return s
}