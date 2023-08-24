package main

import (
  "fmt"
  "io/ioutil"
  "net/http"
)



func mainHandler(w http.ResponseWriter, r *http.Request) {
    // Получаем значение переменной из запроса
    myvar := r.URL.Query().Get("myvar")

    // Создаем URL для второго сервера
    url := "http://localhost:6767/second?myvar=" + myvar

    // Отправляем GET-запрос на второй сервер
    resp, err := http.Get(url)

    if err != nil {
      fmt.Fprintf(w, "Ошибка при отправке запроса на второй сервер: %v", err)
      return
    }
    defer resp.Body.Close()

    // Читаем ответ от второго сервера
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
      fmt.Fprintf(w, "Ошибка при чтении ответа от второго сервера: %v", err)
      return
    }

    fmt.Fprintf(w, "Ответ от второго сервера: %s", body)
  }


func main() {

  http.HandleFunc("/main", mainHandler)
  http.ListenAndServe(":5656", nil)
}