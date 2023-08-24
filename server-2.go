package main

import (
  "fmt"
  "strconv"
  "net/http"
)



func mainHandler(w http.ResponseWriter, r *http.Request) {
  // Получаем значение переменной из запроса
    myvar := r.URL.Query().Get("myvar")
  
    // Преобразуем значение в число
    num, err := strconv.Atoi(myvar)
    if err != nil {
    fmt.Fprintf(w, "Ошибка при преобразовании значения переменной: %v", err)
    return
    }
  
    // Прибавляем 1 к значению переменной
    num++
  
    // Отправляем значение обратно на первый сервер
    fmt.Fprintf(w, "%d", num)
  }


func main() {

  http.HandleFunc("/second", mainHandler)
  http.ListenAndServe(":6767", nil)
}