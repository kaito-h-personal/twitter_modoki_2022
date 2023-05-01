package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
    "encoding/base64"
    "os"
)

type UserId struct {
	UserId   string `json:"user_id"`
}

type DisplayUserInfo struct {
	Name string `json:"name"`
    IconImg   string `json:"icon_img"`
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

type QueryResultUser struct {
	Time   string  `json:"time"`
	Status string  `json:"status"`
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
	_, err := executeQuery(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("デフォルトのtweetをセット")

	// _, _ = fetch_tweets() // 後で消す

	http.HandleFunc("/user", fetchUserHandler)
	http.HandleFunc("/tweets", fetchTweetsHandler)
	http.HandleFunc("/add_tweets", addTweetHandler)

	log.Fatal(http.ListenAndServe(":8007", nil))
}

func fetch_tweets() ([]TweetResponse, error) {
	query := "SELECT * FROM tweet ORDER BY created_at DESC FETCH user;"

	jsonString, err := executeQuery(query)
	if err != nil {
		fmt.Println(err)
		return []TweetResponse{}, err
	}

	// 構造体に変換
	var queryResult []QueryResult
	err = json.Unmarshal([]byte(jsonString), &queryResult)
	if err != nil {
		fmt.Println(err)
		return []TweetResponse{}, err
	}

	var tweets []Tweet = queryResult[0].Result // queryResultの要素は1つの想定

	var tweetResponses []TweetResponse

	for _, t := range tweets {
        // TODO: ここ2回呼ばれている？

        // アイコンの画像を取得
        icon_img, err := LoadImage(t.User.Id)
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

func fetchTweetsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")

	tweets, err := fetch_tweets()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	json.NewEncoder(w).Encode(tweets)
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

	// TODO: SQLインジェクション対策
	query := fmt.Sprintf(`
        CREATE tweet SET
            user = %s
            ,text = '%s'
            ,created_at = '%s'
        ;`, addTweet.UserId, addTweet.Text, now)

	fmt.Println("query")
	fmt.Println(query)

	// TODO: 名前はerr2で良い？
	add_result, err2 := executeQuery(query)
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

	json.NewEncoder(w).Encode(tweets)
}

func LoadImage(user_id string) (string, error) {
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

// TODO: 関数名の規則・関数の順番統一
func fetchUser(user_id string) (DisplayUserInfo, error) {
	// TODO: SQLインジェクション対策
	query := fmt.Sprintf("SELECT id, name FROM user WHERE id=\"%s\";", user_id)

	jsonString, err := executeQuery(query)
	if err != nil {
		fmt.Println(err)
		return DisplayUserInfo{}, err
	}

	// 構造体に変換
	var queryResult []QueryResultUser
	err = json.Unmarshal([]byte(jsonString), &queryResult)
	if err != nil {
		fmt.Println(err)
		return DisplayUserInfo{}, err
	}

	var user User = queryResult[0].Result[0] // queryResult及びResultの要素は1つの想定

	// アイコンの画像を取得
	icon_img, err := LoadImage(user_id)
	if err != nil {
		fmt.Println(err.Error())
		return DisplayUserInfo{}, err
	}

	result := DisplayUserInfo{
		Name:      user.Name,
		IconImg:   icon_img,
	}

	return result, nil
}


func fetchUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("呼ばれた")

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
