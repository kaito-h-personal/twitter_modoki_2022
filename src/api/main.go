package main

import (
    "bytes"
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
)

type Tweet struct {
    Auther     int    `json:"auther"`
    CreatedAt  string `json:"created_at"`
    ID         string `json:"id"`
    Text       string `json:"text"`
}



type Response struct {
    Time   string `json:"time"`
    Status string `json:"status"`
    Result []Tweet `json:"result"`
}

func main() {
    jsonString, err := fetch_tweets()
    if err != nil {
        fmt.Println(err)
        return
    }

    // fmt.Println("!")
    fmt.Println(jsonString)

    var responses []Response
    err = json.Unmarshal([]byte(jsonString), &responses)
    if err != nil {
        fmt.Println("@@")
        fmt.Println(err)
        return
    }

    for _, response := range responses {
        fmt.Println("Time:", response.Time)
        fmt.Println("Status:", response.Status)
        for _, tweet := range response.Result {
            fmt.Println("Auther:", tweet.Auther)
            fmt.Println("Created At:", tweet.CreatedAt)
            fmt.Println("ID:", tweet.ID)
            fmt.Println("Text:", tweet.Text)
        }
    }
    fmt.Println("Fin.")
}

func fetch_tweets() (string, error) {
//     query := `CREATE tweet SET
// 	id = 2,
//   auther = 2,
//   text = 'I\'m sleepy.',
// 	created_at = time::now()
// ;`
    query := "SELECT * FROM tweet;"
    body, err := sendQuery(query)
    if err != nil {
        fmt.Println(err)
        return "", err
    }

    return body, nil
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

    // fmt.Println(resp.Status)
    // fmt.Println("~")
    // fmt.Println(resp)

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
        return "", err
    }

    return string(body), nil
}
