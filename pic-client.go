package main

import (
    "log"
    "net"
    "os"
    "bytes"
    "encoding/binary"
    "fmt"
)

const connAddr = "192.168.0.185:4000"

// Error handling not done
// How to close the connection when C-c is pressed?
func main() {
    conn, err := net.Dial("tcp", connAddr)
    if err != nil {
        log.Fatal(err)
    }

    var (
        buf bytes.Buffer
        b []byte
        size int64
	fileArg string
    )

    fileArg = os.Args[1]
    fmt.Println(fileArg)

    // Request file
    conn.Write([]byte(fileArg))

    b = make([]byte, 8)
    nr, err := conn.Read(b)
    if err != nil {
        fmt.Println("failed read")
        log.Fatal(err)
    }

    // From []byte to int64
    fmt.Println(nr)
    fmt.Println(b)
    size, _ = binary.Varint(b)
    fmt.Printf("Filesize : %d bytes\n", size)

    // Read file to buffer
    for int64(buf.Len()) < size {
        conn.Read(b)
        buf.Write(b)
    }
    //fmt.Println(buf)
    
    f, err := os.Create(fileArg)
    if err != nil {
        log.Fatal(err)
    }

    _, err = f.Write(buf.Bytes())
    if err != nil {
        log.Fatal(err)
    }
    f.Close()
    fmt.Println("DONE")
}
