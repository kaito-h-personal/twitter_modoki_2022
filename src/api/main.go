package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/surrealdb/surrealdb.go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
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

type TweetResponse struct {
	Id        string `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"created_at"`
	UserName  string `json:"user_name"`
	IconImg   string `json:"icon_img"`
}

// TODO: AddTweetRequestの方が良さそう
type AddTweet struct {
	Text   string `json:"text"`
	UserId string `json:"user_id"`
}

type AuthRequest struct {
	Email   string `json:"email"`
	Password string `json:"password"`
}


// TODO: surrealdb.UnmarshalRawは多めにあっても問題ないので、テーブルに合わせる
type UserAuth struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	EncryptedPassword string `json:"encrypted_password"`
}

var db *surrealdb.DB

func main() {
	// DBとのコネクションを作成
	err := dbSetup()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// デフォルトのtweetをセット
	if err := setDefaultTweets(); err != nil {
		fmt.Println(err.Error())
		return
	}

	http.HandleFunc("/auth", authHandler)
	http.HandleFunc("/user", fetchUserHandler)
	http.HandleFunc("/tweets", fetchTweetsHandler)
	http.HandleFunc("/add_tweets", addTweetHandler)

	log.Fatal(http.ListenAndServe(":8007", nil))
}

func getIconImg(user_id string) (string, error) {
	path := fmt.Sprintf("img/%s.jpeg", user_id)
	// 画像ファイルを読み込む
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}
	defer file.Close()

	// 画像ファイルをバイト配列に変換する
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	// バイト配列をBase64エンコードする
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded, nil
}

func dbSetup() (error) {
	_db, err := surrealdb.New("ws://db:8000/rpc")
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	if _, err := _db.Signin(map[string]interface{}{
		"user": "root",
		"pass": "pasuwado",
	}); err != nil {
		fmt.Println(err.Error())
		return err
	}

	if _, err := _db.Use("test", "test"); err != nil {
		fmt.Println(err.Error())
		return err
	}

	db = _db
	return nil
}

func setDefaultTweets() error {
	query := `
		DELETE user;
		DELETE tweet;
		CREATE user SET
			id = 1
			,name = 'ユーザー1'
			,email = 'user1@example.com'
			,encrypted_password = 'password1'
		;
		CREATE user SET
			id = 3
			,name = 'ユーザー3'
			,email = 'user3@example.com'
			,encrypted_password = 'password3'
		;
		CREATE user SET
			id = 6
			,name = 'ユーザー6'
			,email = 'user6@example.com'
			,encrypted_password = 'password6'
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
			fmt.Println(err.Error())
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

	// リクエストの中身を取得
	var addTweet AddTweet
	err := json.NewDecoder(r.Body).Decode(&addTweet)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	JST := time.FixedZone("Asia/Tokyo", 9*60*60)
	now := time.Now().In(JST).Format("2006/01/02 15:04:05") // Golangは「2006年01月02日 15時04分05秒」でフォーマットを指定する

	query := `
        CREATE tweet SET
            user = $user_id
            ,text = $text
            ,created_at = $created_at
        ;`
	if _, err := db.Query(query, map[string]interface{}{
		"user_id":    addTweet.UserId,
		"text":       addTweet.Text,
		"created_at": now,
	}); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tweets, err := fetchTweets()
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(tweets)
}


func authHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	// リクエストの中身を取得
	var authRequest AuthRequest
	err := json.NewDecoder(r.Body).Decode(&authRequest)
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := "SELECT id, name, encrypted_password FROM user WHERE email = $email AND encrypted_password = $encrypted_password;"
	query_result, err := db.Query(query, map[string]interface{}{
		"email": authRequest.Email,
		"encrypted_password": authRequest.Password,
	})
	if err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// クエリの結果を構造体に変換
	var user_auth []UserAuth
	if _, err := surrealdb.UnmarshalRaw(query_result, &user_auth); err != nil {
		fmt.Println(err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	is_authorized := false

	if len(user_auth) != 0 {
		is_authorized = true
	}

	json.NewEncoder(w).Encode(is_authorized)
}
