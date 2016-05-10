package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"TodoIndex",
		"GET",
		"/todos",
		GetAllTodos,
	},
	Route{
		"TodoShow",
		"GET",
		"/todos/{todoId}",
		GetTodobyId,
	},
	Route{
		"TodoCreate",
		"POST",
		"/todos",
		InsertTodo,
	},
	Route{
		"TodoDelete",
		"DELETE",
		"/todos/{todoId}",
		DeleteTodo,
	},
	Route{
		"TodoUpdate",
		"PUT",
		"/todos",
		UpdateTodo,
	},
}