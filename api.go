package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"final-project/model"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	_ "final-project/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

var host = "127.0.0.1"
var port = "8081"
var connString = "root:root@tcp(127.0.0.1:3308)/todo?charset=utf8&parseTime=True&loc=Local"

// @title Hacktiv8 Final Project - Muhammad Reyhan
// @version 1.0
// @description Final Project Hacktiv8
// @termsOfService http://swagger.io/terms/

// @contact.name Muhammad Reyhan
// @contact.url
// @contact.email muh.reyhan@gmail.com
func main() {
	router := mux.NewRouter().StrictSlash(true)
	apiRouter := router.PathPrefix("/api").Subrouter()
	apiRouter.PathPrefix("/todos").HandlerFunc(GetToDos).Methods(http.MethodGet)
	apiRouter.PathPrefix("/todo").HandlerFunc(GetToDoById).Methods(http.MethodGet)
	apiRouter.PathPrefix("/todo").HandlerFunc(CreateToDo).Methods(http.MethodPost)
	apiRouter.PathPrefix("/todo").HandlerFunc(UpdateToDo).Methods(http.MethodPut)
	apiRouter.PathPrefix("/todo").HandlerFunc(DeleteToDo).Methods(http.MethodDelete)

	router.PathPrefix("/user").HandlerFunc(AddUser).Methods(http.MethodPost)
	router.PathPrefix("/user").HandlerFunc(UpdateUser).Methods(http.MethodPut)
	router.PathPrefix("/user").HandlerFunc(DeleteUser).Methods(http.MethodDelete)

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods(http.MethodGet)
	fmt.Println("Listening on port: ", port)
	http.ListenAndServe(host+":"+port, router)

}

//MARK:- Todos
// Get Todos godoc
// @Summary Show all available todos
// @Description Get All Todos
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ToDo
// @Failure 500
// @Router /api/todos [GET]
func GetToDos(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}
	defer db.Close()

	var tasks []model.ToDo
	rows, err := db.Query("SELECT * from ToDo;")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong.")
		fmt.Println(err.Error())
		return
	}

	defer rows.Close()
	for rows.Next() {
		var eachTodo model.ToDo
		rows.Scan(&eachTodo.ID, &eachTodo.Title, &eachTodo.Desc, &eachTodo.DueDate, &eachTodo.PersonInCharge, &eachTodo.Status)
		tasks = append(tasks, eachTodo)
	}

	respondWithJSON(w, http.StatusOK, tasks)

}

// Get Todo godoc
// @Summary Get specific todo by ID
// @Description Get specific todo by ID
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ToDo
// @Failure 500
// @Param id query int true "todo serarch by id"
// @Router /api/todo [GET]
func GetToDoById(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		return
	}

	defer db.Close()

	id := r.URL.Query().Get("id")
	var todo model.ToDo
	err = db.QueryRow("SELECT Title, `Desc`, DueDate, PersonInCharge, `Status` FROM ToDo WHERE id=?", id).Scan(&todo.Title, &todo.Desc, &todo.DueDate, &todo.PersonInCharge, &todo.Status)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No task found with the id="+id)
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	default:
		todo.ID, _ = strconv.Atoi(id)
		respondWithJSON(w, http.StatusOK, todo)

	}

}

// Create Todo godoc
// @Summary Create Todo
// @Description Create new Todo
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ToDo
// @Failure 500
// @Router /api/todo [POST]
func CreateToDo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var todo model.ToDo
	err = decoder.Decode(&todo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}
	statement, err := db.Prepare("insert into ToDo (Title, `Desc`, DueDate, PersonInCharge, `Status`) values(?,?,?,?,?)")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}
	defer statement.Close()
	res, err := statement.Exec(todo.Title, todo.Desc, todo.DueDate, todo.PersonInCharge, todo.Status)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the task.")
		fmt.Println(err.Error())
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		todo.ID = int(id)
		respondWithJSON(w, http.StatusOK, todo)
	}
}

// Update Todo godoc
// @Summary Update Todo
// @Description Update Todo
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ToDo
// @Failure 500
// @Router /api/todo [PUT]
func UpdateToDo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var todo model.ToDo
	err = decoder.Decode(&todo)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}
	statement, err := db.Prepare("UPDATE ToDo set Title=?, `Desc`=?, DueDate=?, PersonInCharge=?, `Status` = ? where ID=?")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}
	defer statement.Close()
	res, err := statement.Exec(todo.Title, todo.Desc, todo.DueDate, todo.PersonInCharge, todo.Status, todo.ID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem updating the task.")
		fmt.Println(err.Error())
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		respondWithJSON(w, http.StatusOK, todo)
	}
}

// Delete Todo godoc
// @Summary Delete Todo
// @Description Delete Todo
// @Accept  json
// @Produce  json
// @Success 200 {object} model.ToDo
// @Failure 500
// @Router /api/todo [DELETE]
func DeleteToDo(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	id := r.URL.Query().Get("id")
	var todo model.ToDo
	err = db.QueryRow("SELECT Title, `Desc`, DueDate, PersonInCharge, `Status` FROM ToDo WHERE ID=?", id).Scan(&todo.Title, &todo.Desc, &todo.DueDate, &todo.PersonInCharge, &todo.Status)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No tasks found with the id="+id)
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	default:
		res, err := db.Exec("DELETE from ToDo where id=?", id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			fmt.Println(err.Error())
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			fmt.Println(err.Error())
			return
		}
		if count == 1 {
			respondWithJSON(w, http.StatusOK, todo)
		}
	}

}

//Mark: - User
// Add User godoc
// @Summary Add User
// @Description Add User
// @Accept  json
// @Produce  json
// @Success 200 {object} model.User
// @Failure 500
// @Router /user [POST]
func AddUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var user model.User
	err = decoder.Decode(&user)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}
	statement, err := db.Prepare("insert into User (Name) values(?)")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}
	defer statement.Close()
	res, err := statement.Exec(user.Name)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem entering the task.")
		fmt.Println(err.Error())
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		user.UserID = int(id)
		respondWithJSON(w, http.StatusOK, user)
	}

}

// Update User godoc
// @Summary Update User
// @Description Update User
// @Accept  json
// @Produce  json
// @Success 200 {object} model.User
// @Failure 500
// @Router /user [PUT]
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	decoder := json.NewDecoder(r.Body)
	var user model.User
	err = decoder.Decode(&user)

	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}
	statement, err := db.Prepare("UPDATE User set Name = ? where userID=?")
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}
	defer statement.Close()
	res, err := statement.Exec(user.Name, user.UserID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "There was problem updating the task.")
		fmt.Println(err.Error())
		return
	}
	if rowsAffected, _ := res.RowsAffected(); rowsAffected == 1 {
		id, _ := res.LastInsertId()
		user.UserID = int(id)
		respondWithJSON(w, http.StatusOK, user)
	}
}

// Delete User godoc
// @Summary Delete User
// @Description Delete User by id
// @Accept  json
// @Produce  json
// @Success 200 {object} model.User
// @Failure 500
// @Router /user [DELETE]
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not connect to the database")
		fmt.Println(err.Error())
		return
	}

	defer db.Close()

	id := r.URL.Query().Get("id")
	var user model.User
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	}

	err = db.QueryRow("SELECT Name FROM User WHERE userID=?", id).Scan(&user.Name)
	switch {
	case err == sql.ErrNoRows:
		respondWithError(w, http.StatusBadRequest, "No User found with the id="+id)
		fmt.Println(err.Error())
		return
	case err != nil:
		respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
		fmt.Println(err.Error())
		return
	default:
		res, err := db.Exec("DELETE from User where userID=?", id)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			fmt.Println(err.Error())
			return
		}
		count, err := res.RowsAffected()
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "Some problem occurred.")
			fmt.Println(err.Error())
			return
		}
		if count == 1 {
			user.UserID, _ = strconv.Atoi(id)
			respondWithJSON(w, http.StatusOK, user)
		}
	}

}

//MARK:- Util Func

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
