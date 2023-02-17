package main

import "database/sql"
import "log"
import _ "github.com/mattn/go-sqlite3"

import "fmt"

func create_db_table(db *sql.DB) {
// DROP TABLE IF EXISTS cars;
    sts := `CREATE TABLE IF NOT EXISTS cars(id INTEGER PRIMARY KEY, name TEXT, price INT);`
    _, err := db.Exec(sts)
    if err != nil {
        log.Fatal(err)
    }
}


func insert_db(db *sql.DB) {
    sts := `
INSERT INTO cars(name, price) VALUES('Audi', 50000);
INSERT INTO cars(name, price) VALUES('Bugati Veyron', 5000000);
`

    _, err := db.Exec(sts)

    if err != nil {
        log.Fatal(err)
    }

}


func insert_prepared_db(db *sql.DB) {
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
    // db, err := sql.Open("sqlite3", "./test.db")

    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    sts := "SELECT * FROM CARS"
    // sts := "SELECT * FROM BLACKOUTS"

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
    db, err := sql.Open("sqlite3", "./check.db")
    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

    create_db_table(db)
    insert_db(db)
    insert_prepared_db(db)

    read_all_db()
}
