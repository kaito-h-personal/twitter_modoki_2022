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

type AddTweet struct {
    Auther     int    `json:"auther"`
    Text       string `json:"text"`
}



type Response struct {
    Time   string `json:"time"`
    Status string `json:"status"`
    Result []Tweet `json:"result"`
}

func main() {
    http.HandleFunc("/tweets", tweetsHandler)
    http.HandleFunc("/add_tweets", addTweetsHandler)

    log.Fatal(http.ListenAndServe(":8007", nil))
}

func fetch_tweets() ([]Tweet, error) {
    query := "SELECT * FROM tweet;"
    //     query := `CREATE tweet SET
// 	id = 1,
//   auther = 1,
//   text = 'テスト内容1',
// 	created_at = time::now()
// ;`
//     query := `CREATE tweet SET
// 	id = 2,
//   auther = 2,
//   text = 'テスト内容2',
// 	created_at = time::now()
// ;`
    jsonString, err := sendQuery(query)
    if err != nil {
        fmt.Println(err)
        return []Tweet{}, err
    }

    // 構造体に変換
    var responses []Response
    err = json.Unmarshal([]byte(jsonString), &responses)
    if err != nil {
        fmt.Println(err)
        return []Tweet{}, err
    }

    response := responses[0] // responsesの要素は1つの想定
    var tweets []Tweet = response.Result
    return tweets, nil
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
    tweets, err := fetch_tweets()
    if err != nil {
        fmt.Println(err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
    json.NewEncoder(w).Encode(tweets)
}

func addTweetsHandler(w http.ResponseWriter, r *http.Request) {
    var addTweet AddTweet
    err := json.NewDecoder(r.Body).Decode(&addTweet)
    if err != nil {
        fmt.Println(err.Error())
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // TODO: SQLインジェクション対策
    query := fmt.Sprintf(`CREATE tweet SET
      auther = %d,
      text = '%s',
    	created_at = time::now()
    ;`, addTweet.Auther, addTweet.Text)

    fmt.Println("query")
    fmt.Println(query)


    // TODO: 名前はerr2で良い？
    add_result, err2 := sendQuery(query)
    if err2 != nil {
        fmt.Println(err2.Error())
        http.Error(w, err.Error(), http.StatusBadRequest)
    }

    fmt.Println("add_result")
    fmt.Println(add_result)

    tweets, err := fetch_tweets()
    if err != nil {
        fmt.Println(err)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")
    json.NewEncoder(w).Encode(tweets)
}
