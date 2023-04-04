package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type TweetResponse struct {
	Id        string `json:"id"`
	UserName  string `json:"user_name"`
	CreatedAt string `json:"created_at"`
	Text      string `json:"text"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Tweet struct {
	User      User   `json:"user"`
	CreatedAt string `json:"created_at"`
	Id        string `json:"id"`
	Text      string `json:"text"`
}

type AddTweet struct {
	UserId int    `json:"user_id"`
	Text   string `json:"text"`
}

type Response struct {
	Time   string  `json:"time"`
	Status string  `json:"status"`
	Result []Tweet `json:"result"`
}

func main() {
	// デフォルトのtweetをセット
	query := `
        DELETE user;
        DELETE tweet;
        CREATE user SET
            id = 1
            ,name = 'ユーザー1'
        ;
        CREATE user SET
            id = 3
            ,name = 'ユーザー3'
        ;
        CREATE tweet SET
            user = user:1
            ,text = 'テスト内容1'
            ,created_at = '2006/01/02 15:04:05'
        ;
        CREATE tweet SET
            user = user:3
            ,text = 'テスト内容2'
            ,created_at = '2009/01/02 15:04:05'
        ;
        CREATE user SET
            id = 6
            ,name = 'ユーザー6'
        ;
    `
	_, err := sendQuery(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("デフォルトのtweetをセット")

	_, _ = fetch_tweets() // 後で消す

	http.HandleFunc("/tweets", tweetsHandler)
	http.HandleFunc("/add_tweets", addTweetsHandler)

	log.Fatal(http.ListenAndServe(":8007", nil))
}

func fetch_tweets() ([]TweetResponse, error) {
	// todo: 順番が変
	query := "SELECT * FROM tweet ORDER BY tweet.created_at DESC FETCH user;"

	jsonString, err := sendQuery(query)
	if err != nil {
		fmt.Println(err)
		return []TweetResponse{}, err
	}

	fmt.Println("jsonString")
	fmt.Println(jsonString)

	// 構造体に変換
	var responses []Response
	err = json.Unmarshal([]byte(jsonString), &responses)
	if err != nil {
		fmt.Println(err)
		return []TweetResponse{}, err
	}

	response := responses[0] // responsesの要素は1つの想定
	var tweets []Tweet = response.Result
	fmt.Println("tweets")
	fmt.Println(tweets)

	var tweetResponses []TweetResponse

	for _, t := range tweets {
		tr := TweetResponse{
			Id:        t.Id,
			UserName:  t.User.Name,
			CreatedAt: t.CreatedAt,
			Text:      t.Text,
		}
		tweetResponses = append(tweetResponses, tr)

	}

	fmt.Println("tweetResponses")
	fmt.Println(tweetResponses)

	return tweetResponses, nil

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
		// TODO: エラーの時CORSになる
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

	JST := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(JST).Format("2006/01/02 15:04:05") // Goは「2006/01/02 15:04:05」でフォーマットを指定する

	// TODO: 将来的に「user = %s」になる予定
	// TODO: SQLインジェクション対策
	query := fmt.Sprintf(`
        CREATE tweet SET
            user = user:%d
            ,text = '%s'
            ,created_at = '%s'
        ;`, addTweet.UserId, addTweet.Text, now)

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
