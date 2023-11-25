package addmain

import (
    "fmt"
    "encoding/json"
)

func PrettyPrint(v interface{}) {
    b, _ := json.MarshalIndent(v, "", "  ")
    fmt.Println(string(b))
}
