package main

import "encoding/json"
import "fmt"
import "log"
import "io/ioutil"
import "net/http"
import "github.com/gorilla/mux"

import "database/sql"
import _ "github.com/mattn/go-sqlite3"

import "os"


type Blackout struct {
    ID		string	`json:"ID"`
    STARTS_UTC	string	`json:"STARTS_UTC"`
    ENDS_UTC	string	`json:"ENDS_UTC"`
    PROGRAM_ID	string	`json:"PROGRAM_ID"`
    REGION_ID	string	`json:"REGION_ID"`
}

var Blackouts []Blackout

var sqliteDb string


func delBlackoutFromDb(blackout Blackout) {
    del_prepared := fmt.Sprintf("DELETE FROM BLACKOUTS WHERE ID=%s", blackout.ID)

    db, err := sql.Open("sqlite3", sqliteDb)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    query, err := db.Prepare(del_prepared)
    if err != nil {
        log.Fatal(err)
    }

    _,err = query.Exec()
    if err != nil {
        log.Fatal(err)
    }
}


func blackoutIntoDb(blackout Blackout) {
    insert_prepared :=`
INSERT INTO BLACKOUTS(ID, STARTS_UTC, ENDS_UTC, PROGRAM_ID, REGION_ID) VALUES(?, ?, ?, ?, ?)`

    db, err := sql.Open("sqlite3", sqliteDb)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    query, err := db.Prepare(insert_prepared)
    if err != nil {
        log.Fatal(err)
    }

    _,err = query.Exec(blackout.ID, blackout.STARTS_UTC, blackout.ENDS_UTC, blackout.PROGRAM_ID,
                       blackout.REGION_ID)
    if err != nil {
        log.Fatal(err)
    }
}


func updateBlackoutInDb(id string, blackout Blackout) {
    update_prepared :=fmt.Sprintf(`
UPDATE BLACKOUTS SET STARTS_UTC='%s', ENDS_UTC='%s', PROGRAM_ID='%s', REGION_ID='%s' WHERE ID = '%s';`,
blackout.STARTS_UTC, blackout.ENDS_UTC, blackout.PROGRAM_ID, blackout.REGION_ID, id)


    fmt.Printf("Executing SQL:%s\n", update_prepared)

    db, err := sql.Open("sqlite3", sqliteDb)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    query, err := db.Prepare(update_prepared)
    if err != nil {
        log.Fatal(err)
    }

    _,err = query.Exec()
    if err != nil {
        log.Fatal(err)
    }
}

func blackoutsIntoDb(blackouts []Blackout) {
    for _, blackout := range blackouts {
        blackoutIntoDb(blackout)
    }
}

func read_all_db() ([]Blackout) {

    db, err := sql.Open("sqlite3", sqliteDb)
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()
    sts := "SELECT ID, STARTS_UTC, ENDS_UTC, PROGRAM_ID, REGION_ID FROM BLACKOUTS"
    result_set, err := db.Query(sts)
    if err != nil {
        log.Fatal(err)
    }

    var db_blackouts []Blackout
    for result_set.Next() {
        var bo Blackout
        result_set.Scan(&(bo.ID),  &(bo.STARTS_UTC), &(bo.ENDS_UTC), &(bo.PROGRAM_ID), &(bo.REGION_ID))
        db_blackouts = append(db_blackouts, bo)
    }
    return db_blackouts
}


func addBlackout(respWriter http.ResponseWriter, req *http.Request) {
    var blackout Blackout
    reqBody, _ := ioutil.ReadAll(req.Body)
    json.Unmarshal(reqBody, &blackout)
    blackoutIntoDb(blackout)
    Blackouts = append(Blackouts, blackout)
    json.NewEncoder(respWriter).Encode(blackout)
}

func putBlackout(respWriter http.ResponseWriter, req *http.Request) {
    var new_blackout Blackout
    reqBody, _ := ioutil.ReadAll(req.Body)

    new_blackout.PROGRAM_ID = "FILLER"

    json.Unmarshal(reqBody, &new_blackout)
    fmt.Printf("New blackout is: ")
    fmt.Println(new_blackout)

    vars := mux.Vars(req)
    id := vars["ID"]
    for _, blackout := range Blackouts {
        if blackout.ID == id {
            delBlackoutFromDb(blackout)
        }
    }

    blackoutIntoDb(new_blackout)
    Blackouts = read_all_db()
    json.NewEncoder(respWriter).Encode(new_blackout)
}


func patchBlackout(respWriter http.ResponseWriter, req *http.Request) {
    var blackoutUpdates Blackout
    reqBody, _ := ioutil.ReadAll(req.Body)

    vars := mux.Vars(req)
    id := vars["ID"]
    for _, blackout := range Blackouts {
        if blackout.ID == id {
            json.Unmarshal(reqBody, &blackout)
            updateBlackoutInDb(id, blackout)
        }
    }

    Blackouts = read_all_db()
    json.Unmarshal(reqBody, &blackoutUpdates)
    json.NewEncoder(respWriter).Encode(blackoutUpdates)
}

func getAllBlackouts(respWriter http.ResponseWriter, req *http.Request) {
    Blackouts = read_all_db()
    json.NewEncoder(respWriter).Encode(Blackouts)
}

func getBlackout(respWriter http.ResponseWriter, req *http.Request) {
    vars := mux.Vars(req)

    Blackouts = read_all_db()
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
            delBlackoutFromDb(blackout)
            Blackouts = append(Blackouts[:index], Blackouts[index+1:]...) // Not needed.
        }
    }
}



func setUpApis() {

    muxRouter := mux.NewRouter().StrictSlash(true)

    muxRouter.HandleFunc("/blackout", addBlackout).Methods("POST")

    muxRouter.HandleFunc("/blackout", getAllBlackouts)

    muxRouter.HandleFunc("/blackout/{ID}", deleteBlackout).Methods("DELETE")

    muxRouter.HandleFunc("/blackout/{ID}", putBlackout).Methods("PUT")

    muxRouter.HandleFunc("/blackout/{ID}", patchBlackout).Methods("PATCH")

    muxRouter.HandleFunc("/blackout/{ID}", getBlackout)

    log.Fatal(http.ListenAndServe(":4221", muxRouter))
}


func main() {
    argn := len(os.Args[1:])

    sqliteDb = ""
    if (argn != 1) {
        fmt.Printf("ARGUMENTS ERROR: %s: Need 1 args: SQLite DB path\n", os.Args[0])
        os.Exit(1)
    } else {
        sqliteDb = os.Args[1]
    } 

    Blackouts = []Blackout{
        Blackout{ID: "1", STARTS_UTC: "2001", ENDS_UTC: "3001", PROGRAM_ID: "Prog1",
                 REGION_ID: "Kentucky1"},
        Blackout{ID: "2", STARTS_UTC: "6001", ENDS_UTC: "10001", PROGRAM_ID: "Prog1",
                 REGION_ID: "Texas0"},
    }

    dbConn, err := sql.Open("sqlite3", sqliteDb)
    if err != nil {
        log.Fatal(err)
    }
    defer dbConn.Close()

    blackoutsIntoDb(Blackouts)

    Blackouts = read_all_db()

    setUpApis()

}
