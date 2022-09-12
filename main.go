package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"path"
	"time"
)

var states = make(map[string]bool)

func main() {
	mux := http.NewServeMux()

	port := "8080"

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/dashboard", dashboardHandler)
	mux.HandleFunc("/redirect", redirectHandler)

	log.Println("Server listening on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	publicDir := path.Join(wd, "public")

	tmpl := template.Must(template.ParseFiles(path.Join(publicDir, "index.html")))

	state := genereateRandomString(5)

	result := map[string]string{"Rand": state}

	states[state] = true

	if err := tmpl.Execute(w, result); err != nil {
		log.Fatal(err)
	}
}

func dashboardHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	clientID, ok := os.LookupEnv("CLIENT_ID")
	if !ok {
		log.Fatal("CLIENT_ID env is required.")
	}

	clientSecret, ok := os.LookupEnv("CLIENT_SECRET")
	if !ok {
		log.Fatal("CLIENT_SECRET env is required.")
	}

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodPost, "https://github.com/login/oauth/access_token", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "application/json")

	q := req.URL.Query()
	q.Add("client_id", clientID)
	q.Add("client_secret", clientSecret)
	q.Add("code", code)

	req.URL.RawQuery = q.Encode()

	response, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
		Scope       string `json:"scope"`
		TokenType   string `json:"token_type"`
	}

	if response.StatusCode == http.StatusOK {
		json.NewDecoder(response.Body).Decode(&token)
		log.Printf("%+v\n", token)
	} else {
		w.WriteHeader(response.StatusCode)
		resData := map[string]any{"message": "Something went wrong."}
		res, err := json.Marshal(resData)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
		return
	}

	req, err = http.NewRequest(http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	response, err = client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	if response.StatusCode == http.StatusOK {
		var user map[string]any
		json.NewDecoder(response.Body).Decode(&user)
		res, err := json.MarshalIndent(user, " ", "\t")
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
	} else {
		w.WriteHeader(response.StatusCode)
		resData := map[string]any{"message": "Something went wrong."}
		res, err := json.Marshal(resData)
		if err != nil {
			log.Fatal(err)
		}
		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
	}
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if _, ok := states[state]; !ok {
		errorResponse := make(map[string]any)
		errorResponse["message"] = "State not equal."
		w.WriteHeader(http.StatusUnauthorized)

		res, err := json.Marshal(errorResponse)
		if err != nil {
			log.Fatal(err)
		}

		w.Header().Add("Content-Type", "application/json")
		w.Write(res)
	}

	http.Redirect(w, r, fmt.Sprintf("/dashboard?code=%s", code), http.StatusTemporaryRedirect)
}

func genereateRandomString(size int) string {
	rand.Seed(time.Now().UnixNano())
	res := make([]byte, size)

	for i := range res {
		randomInt := rand.Intn(26)
		res[i] = 'a' + byte(randomInt)
	}

	return string(res)
}
