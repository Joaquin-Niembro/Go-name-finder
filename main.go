package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func fetchGender(names []string, quit <-chan struct{}) <-chan string {
	genders := make(chan string)

	go func() {
		defer close(genders)
		for _, name := range names {
			resp, err := http.Get("https://api.genderize.io?name=" + name)
			body, _ := io.ReadAll(resp.Body)
			if err != nil {
				panic(err)
			}
			var result struct {
				Gender string `json:"gender"`
			}
			if err := json.Unmarshal(body, &result); err != nil {
				panic(err)
			}
			select {
			case genders <- result.Gender:
			case <-quit:
				return
			}
		}
	}()
	return genders
}

func main() {
	names := []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Hannah", "Ivan", "John", "Katie", "Liam", "Mia", "Nathan", "Olivia", "Peter", "Quinn", "Rachel", "Steve", "Tina", "Ursula", "Victor", "Wendy", "Xavier", "Yvonne", "Zach"}
	quit := make(chan struct{})
	defer close(quit)
	genders := fetchGender(names, quit)
	for name := range genders {
		fmt.Println(name)
	}
}
