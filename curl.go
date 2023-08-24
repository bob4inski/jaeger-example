package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func getPageContent() (string, error) {
	url := "http://localhost:656/main?myvar=1"

	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("Ошибка при выполнении запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("Ошибка при чтении ответа: %v", err)
		}
		return string(body), nil
	}

	return "", fmt.Errorf("Страница недоступна. Пожалуйста, сообщите пользователю, что не удалось получить страницу.")
}

func main() {
	content, err := getPageContent()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(content)
}

    
    
