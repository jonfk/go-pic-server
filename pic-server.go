package main

import (
    "log"
    "net"
    "io/ioutil"
    "os"
    "fmt"
    "bytes"
    "encoding/binary"
)

const listenAddr = ":4000"

func main() {
    l, err := net.Listen("tcp", listenAddr)
    if err != nil {
        log.Fatal(err)
    }
    for {
        c, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go serve(c)
    }
}

func serve(c net.Conn) {
    var (
        filenameBuf *bytes.Buffer
        filename string
        f *os.File
        info os.FileInfo
        b []byte
        sizeRead int
        filesize int64
    )

    // Prepare the buffer
    tempb := make([]byte, 256)
    filenameBuf = bytes.NewBuffer(tempb)

    // Read Request
    c.Read(filenameBuf.Bytes())

    // Never forget to splice the string or the comparison will fail
    // Dumb thing in my opinion. Really annoying!
    fmt.Println(filenameBuf.Bytes())
    filename, _ = filenameBuf.ReadString(0)
    filename = filename[:len(filename) - 1]
    fmt.Print("length : ")
    fmt.Println(len(filename))

    if !isValid(filename) {
        fmt.Println("cannot serve request for file: " + filename)
        os.Exit(1)
    }

    f, err := os.Open("/var/www/images/" + filename)//("/home/faiz/andrea.jpg")
    if err != nil {
        log.Fatal(err)
    }
    defer f.Close()

    // Get file size
    info, err = f.Stat()
    if err != nil {
        log.Fatal(err)
    }
    filesize = info.Size()


    fmt.Printf("Filesize : %d bytes\n" , filesize)

    // Write filesize to client
    buf := make([]byte, 8)
    binary.PutVarint(buf, filesize)
    _, err = c.Write(buf)
    if err != nil {
            log.Fatal(err)
    }
    fmt.Println(">>Wrote Size to client<<")

    // Write file to client
    b, err = ioutil.ReadAll(f)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("1")
    sizeRead = 0
    for sizeRead < len(b) {
        fmt.Println("2")
        addSize, err := c.Write(b)
        fmt.Println("3")
        if err != nil {
            log.Fatal(err)
        }
	fmt.Println("size written: "+string(addSize))
        sizeRead += addSize
    }
    fmt.Println("---File Transfer Complete---")
}

func isValid(filename string) bool {
    fileList := []string{"andrea.jpg", "erinstern.jpg", "hamlet.txt", "girl.bmp"}

    for _, f := range fileList {
        if filename == f {
            return true
        }
    }
    return false
}
