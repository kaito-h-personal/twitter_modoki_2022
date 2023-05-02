package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/surrealdb/surrealdb.go"
)

type UserId struct {
	UserId string `json:"user_id"`
}

type DisplayUserInfo struct {
	Name    string `json:"name"`
	IconImg string `json:"icon_img"`
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type Tweet struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}

type QueryResult struct {
	Time   string  `json:"time"`
	Status string  `json:"status"`
	Result []Tweet `json:"result"`
}

// TODO: 構造体名
type QueryResultUser struct {
	Time   string `json:"time"`
	Status string `json:"status"`
	Result []User `json:"result"`
}

type TweetResponse struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UserName  string `json:"user_name"`
	IconImg   string `json:"icon_img"`
}

type AddTweet struct {
	Text   string `json:"text"`
	UserId string `json:"user_id"`
}

var db *surrealdb.DB

func main() {
	// DBとのコネクションを作成
	_db, err := surrealdb.New("ws://db:8000/rpc")
	if err != nil {
		fmt.Println(err)
		return
	}

	if _, err := _db.Signin(map[string]interface{}{
		"user": "root",
		"pass": "pasuwado",
	}); err != nil {
		fmt.Println(err)
		return
	}

	if _, err := _db.Use("test", "test"); err != nil {
		fmt.Println(err)
		return
	}

	db = _db

	// デフォルトのtweetをセット
	if err := setDefaultTweets(); err != nil {
		fmt.Println(err)
		return
	}

	http.HandleFunc("/user", fetchUserHandler)
	http.HandleFunc("/tweets", fetchTweetsHandler)
	http.HandleFunc("/add_tweets", addTweetHandler)

	log.Fatal(http.ListenAndServe(":8007", nil))
}

func executeQuery(query string) (string, error) {
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

func getIconImg(user_id string) (string, error) {
	path := fmt.Sprintf("img/%s.jpeg", user_id)
	// 画像ファイルを読み込む
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer file.Close()

	// 画像ファイルをバイト配列に変換する
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	// バイト配列をBase64エンコードする
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

func setDefaultTweets() (error) {
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
	_, err := db.Query(query, nil)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println("デフォルトのtweetをセットしました")
	return nil
}

func fetchUserHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	var user_id UserId
	err := json.NewDecoder(r.Body).Decode(&user_id)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user_info, err := fetchUser(user_id.UserId)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	json.NewEncoder(w).Encode(user_info)
}

func fetchUser(user_id string) (DisplayUserInfo, error) {
	query := "SELECT id, name FROM user WHERE id = $user_id;"
	query_result, err := db.Query(query, map[string]interface{}{
		"user_id": user_id,
	})
	if err != nil {
		fmt.Println(err.Error())
		return DisplayUserInfo{}, err
	}

	// クエリの結果を構造体に変換
	var user_slice []User
	if _, err := surrealdb.UnmarshalRaw(query_result, &user_slice); err != nil {
		fmt.Println(err.Error())
		return DisplayUserInfo{}, err
	}

	var user User = user_slice[0] // 要素は一つの想定

	// アイコンの画像を取得
	icon_img, err := getIconImg(user_id)
	if err != nil {
		fmt.Println(err.Error())
		return DisplayUserInfo{}, err
	}

	result := DisplayUserInfo{
		Name:    user.Name,
		IconImg: icon_img,
	}

	return result, nil
}

func fetchTweetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	tweets, err := fetchTweets()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	json.NewEncoder(w).Encode(tweets)
}

func fetchTweets() ([]TweetResponse, error) {
	query := "SELECT * FROM tweet ORDER BY created_at DESC FETCH user;"
	query_result, err := db.Query(query, nil)
	if err != nil {
		fmt.Println(err.Error())
		return []TweetResponse{}, err
	}

	// クエリの結果を構造体に変換
	var tweets []Tweet
	if _, err := surrealdb.UnmarshalRaw(query_result, &tweets); err != nil {
		fmt.Println(err.Error())
		return []TweetResponse{}, err
	}

	var tweetResponses []TweetResponse
	for _, t := range tweets {
		// アイコンの画像を取得
		icon_img, err := getIconImg(t.User.Id)
		if err != nil {
			fmt.Println(err.Error()) // TODO: 統一
			return []TweetResponse{}, err
		}

		tr := TweetResponse{
			Id:        t.Id,
			Text:      t.Text,
			CreatedAt: t.CreatedAt,
			UserName:  t.User.Name,
			IconImg:   icon_img,
		}
		tweetResponses = append(tweetResponses, tr)

	}

	return tweetResponses, nil

}

func addTweetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	var addTweet AddTweet
	err := json.NewDecoder(r.Body).Decode(&addTweet)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	JST := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(JST).Format("2006/01/02 15:04:05") // Goは「2006/01/02 15:04:05」でフォーマットを指定する

	query := `
        CREATE tweet SET
            user = $user_id
            ,text = $text
            ,created_at = $created_at
        ;`
	if _, err := db.Query(query, map[string]interface{}{
		"user_id": addTweet.UserId,
		"text": addTweet.Text,
		"created_at": now,
	}); err != nil {
		fmt.Println(err.Error())
		return
	}

	tweets, err := fetchTweets()
	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(tweets)
}
