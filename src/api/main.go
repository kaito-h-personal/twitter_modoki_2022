package main

import (
    "bytes"
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "log"
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

    // 構造体に変換
    var responses []Response
    err = json.Unmarshal([]byte(jsonString), &responses)
    if err != nil {
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

    // TODO: 上記で取得したものを返す

    http.HandleFunc("/tweets", tweetsHandler)
    http.HandleFunc("/add_tweets", addTweetsHandler)

    log.Fatal(http.ListenAndServe(":8007", nil))
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

func tweetsHandler(w http.ResponseWriter, r *http.Request) {
    tweets := []Tweet{
        {
            Auther:    1,
            CreatedAt: "2023-04-01T16:18:18.419644996Z",
            ID:        "tweet:1",
            Text:      "I got it.",
        },
        {
            Auther:    2,
            CreatedAt: "2023-04-01T16:19:07.287544979Z",
            ID:        "tweet:2",
            Text:      "I'm sleepy.",
        },
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
    json.NewEncoder(w).Encode(tweets)
}

func addTweetsHandler(w http.ResponseWriter, r *http.Request) {
    tweets := []Tweet{
        {
            Auther:    1,
            CreatedAt: "2023-04-01T16:18:18.419644996Z",
            ID:        "tweet:1",
            Text:      "I got it.",
        },
        {
            Auther:    2,
            CreatedAt: "2023-04-01T16:19:07.287544979Z",
            ID:        "tweet:2",
            Text:      "I'm sleepy.",
        },
        {
            Auther:    3,
            CreatedAt: "2023-04-01T16:19:07.287544979Z",
            ID:        "tweet:3",
            Text:      "テスト3",
        },
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
    json.NewEncoder(w).Encode(tweets)
}
