package main

import (
	"fmt"
	"log"
	"io/ioutil"
	"html/template"
    "net/http"
    "github.com/gorilla/mux"
    "os/exec"
    "os"
    "bytes"
)
var tpl = template.Must(template.ParseFiles("index.html"))

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tpl.Execute(w, nil)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
    fmt.Println("File Upload Endpoint Hit")

    r.ParseMultipartForm(10 << 20)
	
	file, _, err := r.FormFile("myFile")
    if err != nil {
        fmt.Println("Error Retrieving the File")
        fmt.Println(err)
        return
	}

    defer file.Close()
    // fmt.Printf("Uploaded File: %+v\n", handler.Filename)
    // fmt.Printf("File Size: %+v\n", handler.Size)
    // fmt.Printf("MIME Header: %+v\n", handler.Header)

    tempFile, err := ioutil.TempFile("files", "file.*.c")
    if err != nil {
        fmt.Println(err)
    }
    defer tempFile.Close()

    fileBytes, err := ioutil.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
    
    tempFile.Write(fileBytes)
    
    cmd  := exec.Command("cmd.exe","/C","gcc E:\\go\\src\\Complier\\files\\file.*.c -o filename.exe")
    var out bytes.Buffer
    var stderr bytes.Buffer
    cmd.Stdout = &out
    cmd.Stderr = &stderr
    things := cmd.Run()
    if things != nil {
        fmt.Fprintf(w , "Error: " + stderr.String())
        //fmt.Println(fmt.Sprint(things) + ": " + stderr.String())

        tempFile.Close() 
        err =  os.Remove(tempFile.Name())
        if err != nil{
            fmt.Println(err)
        }
        fmt.Println(tempFile.Name()) 

        
        return
    }
    fmt.Println("Result: " + out.String())

    cmd1  := exec.Command("cmd.exe","/C",".\\filename.exe")
    var out1 bytes.Buffer
    var stderr1 bytes.Buffer
    cmd1.Stdout = &out1
    cmd1.Stderr = &stderr1
    things1 := cmd1.Run()
    if things1 != nil {
        fmt.Fprintf(w , "Error: " + stderr1.String())
        return
    }

    
    fmt.Fprintf(w,"Result: " + out1.String())

    //fmt.Fprintf(w, "\nSuccessfully Uploaded File\n")
    tempFile.Close() 
    err =  os.Remove(tempFile.Name())
    if err != nil{
        fmt.Println(err)
    }
    fmt.Println(tempFile.Name()) 
}

func main() {
	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", indexHandler)
	router.HandleFunc("/upload", uploadFile)
	log.Fatal(http.ListenAndServe(":3000", router))

}