package main

import "encoding/json"
import "fmt"
import "log"
import "io/ioutil"
import "net/http"
import "github.com/gorilla/mux"


type Blackout struct {
    ID		string	`json:"ID"`
    STARTS_UTC	string	`json:"STARTS_UTC"`
    ENDS_UTC	string	`json:"ENDS_UTC"`
    PROGRAM_ID	string	`json:"PROGRAM_ID"`
    REGION_ID	string	`json:"REGION_ID"`
}

var Blackouts []Blackout

func addBlackout(respWriter http.ResponseWriter, req *http.Request) {
    reqBody, _ := ioutil.ReadAll(req.Body)
    var blackout Blackout

    fmt.Println(fmt.Printf("%s", req.Body))

    json.Unmarshal(reqBody, &blackout)
    Blackouts = append(Blackouts, blackout)
    json.NewEncoder(respWriter).Encode(blackout)
}


func getAllBlackouts(respWriter http.ResponseWriter, req *http.Request) {
    json.NewEncoder(respWriter).Encode(Blackouts)
}

func getBlackout(respWriter http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    key := vars["ID"]
    for _, blackout := range Blackouts {
        if blackout.ID == key {
            json.NewEncoder(respWriter).Encode(blackout)
        }
    }
}

func deleteBlackout(respWriter http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)
    id := vars["ID"]
    for index, blackout := range Blackouts {
        if blackout.ID == id {
            Blackouts = append(Blackouts[:index], Blackouts[index+1:]...)
        }
    }
}


func setUpApis() {

    muxRouter := mux.NewRouter().StrictSlash(true)

    muxRouter.HandleFunc("/blackout", addBlackout).Methods("POST")
    muxRouter.HandleFunc("/blackout", getAllBlackouts)

    muxRouter.HandleFunc("/blackout/{ID}", deleteBlackout).Methods("DELETE")
    muxRouter.HandleFunc("/blackout/{ID}", getBlackout)

    log.Fatal(http.ListenAndServe(":4221", muxRouter))
}


func main() {
    Blackouts = []Blackout{
        Blackout{ID: "1", STARTS_UTC: "2001", ENDS_UTC: "3001", PROGRAM_ID: "Prog1",
                 REGION_ID: "Region1"},
    }

    setUpApis()
}
