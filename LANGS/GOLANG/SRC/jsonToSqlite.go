package main

// Get entry for each table as outer key => TABLE NAME
// Collect entry for each COLUMN as keys in outer JSON (value of TABLE NAME key)
// Collect entries for each COLUMN'S TYPE as either string, or first key of inner JSON
// [Defer] Get constraints from inner JSON {value in other table, or in inner JSON list}
// Build SQLITE TABLE CREATION QUERY from TABLE NAME, and list of COLUMN-TYPE pairs
// Execute SQLITE TABLE CREATION QUERY.

import "fmt"
import "os"
import "strings"

import "io/ioutil"
import "github.com/buger/jsonparser"

import "database/sql"
import "log"
import _ "github.com/mattn/go-sqlite3"

// import "encoding/json"


func matchType(origType string) (string, error) {

    var localType string = ""

    switch(strings.ToUpper(origType)) {
       case "COMMENT", "COMMENTS": { localType = "COMMENT" }
       case "INTEGER", "LONG", "NUMBER": { localType = "INTEGER" }
       case "STRING", "TEXT", "VARCHAR": { localType = "TEXT" }

       default: {
           fmt.Printf("ERROR: matchType(): No local type matching (%s).\n", origType)
           os.Exit(3)
       }
    }

    return localType, nil
}





func makeTables(dbPath string, tableName string, colTypeMap map[string]string) {

    db, err := sql.Open("sqlite3", dbPath)

    if err != nil {
        log.Fatal(err)
    }

    defer db.Close()

   var count int = 0
   crTableStmt := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s ( ", tableName)
   for colName, typeStr := range colTypeMap {

       if strings.Compare(strings.ToUpper(colName), "COMMENT") == 0 { continue }

       localType, _ := matchType(typeStr)

       if (count > 0) { crTableStmt += ", " }

       crTableStmt += fmt.Sprintf("%s %s", colName, localType)

       count++
   }
   crTableStmt += " );"


    // _, err = db.Exec(sts)
    _, err = db.Exec(crTableStmt)

    if err != nil {
        log.Fatal(err)
    }

}


/////////////////////////////////////////////////////////////////////////////////////////////////

func getFirstKey(data []byte) (string, error) {

    var firstKey string = ""

    jsonparser.ObjectEach(
        data,
        func( 
            key []byte, value []byte, dataType jsonparser.ValueType, 
            offset int) error {

            firstKey = string(key)

            return nil
    })

    return firstKey, nil
}


func getKeyTypeMap(data []byte, level int) (map[string]string, error) {

    var colName string = ""
    var typeString string = ""

    var colTypeMap = make (map[string]string)

    jsonparser.ObjectEach(
        data,
        func( 
            key []byte, value []byte, dataType jsonparser.ValueType, 
            offset int) error {

            colName = string(key)
    
            dts := fmt.Sprintf("%s", dataType)
            if strings.Compare(dts, "object") == 0 {
                typeString, _ = getFirstKey(value)

            } else if strings.Compare(dts, "string") == 0 {
                typeString = string(value)

            } else {
                fmt.Printf("ERROR: getKeyTypeMap(): Badly formatted JSON in %s\n",
                           string(value))

                os.Exit(2)
            }

            colTypeMap[colName] = typeString

            return nil
    } )

    fmt.Println(colTypeMap)

    return colTypeMap, nil
}



func makeSqliteTables(sqliteDb string, data []byte, level int) {

    fmt.Println(sqliteDb)
    fmt.Println(string(data))

    jsonparser.ObjectEach(
        data,
        func( 
            key []byte, value []byte, dataType jsonparser.ValueType, 
            offset int) error {

            var tableName string = string(key)
            fmt.Printf("Table name == %s\n", tableName)

            dts := fmt.Sprintf("%s", dataType)
            if strings.Compare(dts, "object") == 0 {
                keyTypeMap, _ := getKeyTypeMap(value, level+1)
                fmt.Println(keyTypeMap)

                makeTables(sqliteDb, tableName, keyTypeMap)
            }


	    return nil

        })
}


func main() {

    argn := len(os.Args[1:])

    var sqliteDb string = ""
    var tablesJson string = ""

    if (argn != 2) {
        fmt.Printf("ARGUMENTS ERROR: %s: Need 2 args: SQLite DB path, Path to JSON tables descriptor\n",
                    os.Args[0])
        os.Exit(1)
    }

    if (argn == 2) {
        sqliteDb = os.Args[1]
        tablesJson = os.Args[2]
    } 

    data, err := ioutil.ReadFile(tablesJson)

    if (err != nil) {
        fmt.Println(err)

    } else {
        makeSqliteTables(sqliteDb, data, 0)
    }
}
