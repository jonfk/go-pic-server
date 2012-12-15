package main

import (
    "log"
    "net"
    "os"
    "bytes"
    "encoding/binary"
    "fmt"
    "time"
    "strconv"
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
    t0 := time.Now()

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
    t1 := time.Now()
    time_taken := t1.Sub(t0)
    time_f, err := os.OpenFile("send_time.csv", os.O_APPEND | os.O_RDWR | os.O_CREATE, 0666)
    if err != nil {
        log.Fatal(err)
    }
    //time_taken_buf := make([]byte, 256)
    //binary.PutVarint(time_taken_buf, time_taken.Nanoseconds())
    _, err = time_f.Write([]byte( strconv.FormatFloat(time_taken.Seconds(), 'e', 10, 64) ))
    if err != nil {
        log.Fatal(err)
    }
    _, err = time_f.Write( []byte("\n") )
    if err != nil {
        log.Fatal(err)
    }
    time_f.Close()
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
