package main

import "encoding/json"
import "fmt"
import "log"
import "io/ioutil"
import "net/http"
import "github.com/gorilla/mux"


type Employee struct {
    Id		string	`json:"Id"`
    Name	string	`json:"Name"`
    Address	string	`json:"Address"`
    Salary	string	`json:"Salary"`
}

var Employees []Employee

func homePage(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Welcome to the API Home page!")
}


func getAllEmployees(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(Employees)
}

func getEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    key := vars["id"]
    for _, employee := range Employees {
        if employee.Id == key {
            json.NewEncoder(w).Encode(employee)
        }
    }
}

func createEmployee(w http.ResponseWriter, r *http.Request) {
    reqBody, _ := ioutil.ReadAll(r.Body)
    var employee Employee
    json.Unmarshal(reqBody, &employee)
    Employees = append(Employees, employee)
    json.NewEncoder(w).Encode(employee)
}


func deleteEmployee(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]
    for index, employee := range Employees {
        if employee.Id == id {
            Employees = append(Employees[:index], Employees[index+1:]...)
        }
    }
}


func apiRequests() {
    route := mux.NewRouter().StrictSlash(true)
    route.HandleFunc("/", homePage)
    route.HandleFunc("/employees", getAllEmployees)
    route.HandleFunc("/employee", createEmployee).Methods("POST")
    route.HandleFunc("/employee/{id}", deleteEmployee).Methods("DELETE")
    route.HandleFunc("/employee/{id}", getEmployee)

    log.Fatal(http.ListenAndServe(":3000", route))
}


func main() {
    Employees = []Employee{
        Employee{Id: "1", Name: "Gerard Anthony", Address: "New Jersey USA", Salary: "200000"},
        Employee{Id: "2", Name: "Anthony Denga", Address: "Wellingdon Island India", Salary: "1000000"},
        Employee{Id: "3", Name: "Gerard Gold", Address: "Londonderry Ireland", Salary: "800000"},
    }

    apiRequests()
}
