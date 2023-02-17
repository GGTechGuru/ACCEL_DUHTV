package main

import "fmt"
import "os"
import "strings"

import "io/ioutil"
import "github.com/buger/jsonparser"

import "encoding/json"


func check_json(data []byte, level int) {

    fmt.Println(string(data))

    jsonparser.ObjectEach(
        data,
        func( 
            key []byte, value []byte, dataType jsonparser.ValueType, 
            offset int) error {

            fmt.Printf(
                "Key: '%s'\n Value: '%s'\n Type: %s\n", string(key),
                string(value), dataType)

            dts := fmt.Sprintf("%s", dataType)
            if strings.Compare(dts, "object") == 0 {
                check_json(value, level+1)

            // } else if (level == 0) && strings.Compare(dts, "array") == 0 {
            } else if strings.Compare(dts, "array") == 0 {

                var val_ary []string

                // _ = json.Unmarshal(fmt.Sprintf("%s", value), &val_ary)
                _ = json.Unmarshal(value, &val_ary)

                fmt.Println(val_ary)
                // fmt.Println(string(value))

                for _, el  := range val_ary {

                    // check_json(data[index], level+1)
                    fmt.Println(string(el))
                }

            } else {
                fmt.Printf("dataType(%s) != 'object'\n", dts)
            }

	    return nil

        })

}


func main() {

    argn := len(os.Args[1:])

    json_file := "./test.json"
    if (argn > 0) {
        json_file = os.Args[1]
    } 

    // data, _ := ioutil.ReadFile("./test.json")
    data, err := ioutil.ReadFile(json_file)

    if (err != nil) {
        fmt.Println(err)

    } else {
        check_json(data, 0)
    }
}
