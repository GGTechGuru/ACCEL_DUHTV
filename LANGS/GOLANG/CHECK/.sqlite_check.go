package main

import "database/sql"
import "log"
import _ "github.com/mattn/go-sqlite3"

import "fmt"

func insert_db() {
    db, err := sql.Open("sqlite3", "./check.db")
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    sts := `
DROP TABLE IF EXISTS cars;
CREATE TABLE cars(id INTEGER PRIMARY KEY, name TEXT, price INT);
INSERT INTO cars(name, price) VALUES('Audi', 50000);
INSERT INTO cars(name, price) VALUES('Bugati Veyron', 5000000);
`

    _, err = db.Exec(sts)

    if err != nil {
        log.Fatal(err)
    }

}


func insert_prepared_db() {

    db, err := sql.Open("sqlite3", "./check.db")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    insert_prepared := "INSERT INTO cars(name, price) VALUES(?, ?)"

    query, err := db.Prepare(insert_prepared)
    if err != nil {
        log.Fatal(err)
    }

    _,_ = query.Exec("Maybach", 300000)
    _,_ = query.Exec("Acura NSX", 70000)
}


func read_all_db() {
    db, err := sql.Open("sqlite3", "./check.db")

    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    sts := "SELECT * FROM CARS"

    result_set, err := db.Query(sts)

    if err != nil {
        log.Fatal(err)
    }

    for result_set.Next() {
        var id int
        var name string
        var price int

        result_set.Scan(&id, &name, &price)

        fmt.Printf("Car: %d %s %d\n", id, name, price)
      
    }

}


func main() {
    // insert_db()
    insert_prepared_db()

    read_all_db()
}
