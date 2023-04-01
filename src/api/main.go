package main

import (
    "bytes"
    "fmt"
    "net/http"
    "io/ioutil"
)

func main() {
    query := "INFO FOR DB;"
    body, err := sendQuery(query)
    if err != nil {
        fmt.Println(err)
        return
    }

    fmt.Println("!")
    fmt.Println(body)
}

func sendQuery(query string) (string, error) {
    url := "http://db:8000/sql"
    data := []byte(query)
    req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
    if err != nil {
        fmt.Println(err)
        return "", err
    }
    req.Header.Set("Accept", "application/json")
    req.Header.Set("NS", "test")
    req.Header.Set("DB", "test")
    req.SetBasicAuth("root", "pasuwado")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        fmt.Println(err)
        return "", err
    }
    defer resp.Body.Close()

    fmt.Println(resp.Status)
    fmt.Println("~")
    fmt.Println(resp)

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return "", err
    }

    return string(body), nil
}
