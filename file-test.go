package main

import (
    "fmt"
    "log"
    "os"
    "bytes"
    "io/ioutil"
)
func main() {
    f, err := os.Open("/home/jfokka/hamlet.txt")
    if err != nil {
        log.Fatal(err)
    }
    info, err := f.Stat()
    if err != nil {
        fmt.Println("It's getting in here....")
        log.Fatal(err)
    }
    size := info.Size()
    fmt.Println(size)

    var buf bytes.Buffer

    b, err := ioutil.ReadAll(f)
    if err != nil {
        fmt.Println("It's getting in here....")
        log.Fatal(err)
    }

    fmt.Println(len(b))

    buf.Write( b )
    //buf.WriteTo(os.Stdout)
    fmt.Println(buf.Len())
}
