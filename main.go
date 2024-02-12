package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

// NewsFullDetailed представляет полные детали новости
type NewsFullDetailed struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Content string `json:"content"`
}

// NewsShortDetailed представляет краткие детали новости
type NewsShortDetailed struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
}

// Comment представляет комментарий к новости
type Comment struct {
	ID       int    `json:"id"`
	Text     string `json:"text"`
	NewsID   int    `json:"news_id"`
	ParentID int    `json:"parent_id"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./comments.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Создание таблицы комментариев
	createTableSQL := `CREATE TABLE IF NOT EXISTS comments (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        text TEXT,
        news_id INTEGER,
        parent_id INTEGER
    );`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal(err)
	}

	// Обработчик для вывода списка новостей
	http.HandleFunc("/news", func(w http.ResponseWriter, r *http.Request) {
		newsList := []NewsShortDetailed{
			{ID: 1, Title: "Новость 1"},
			{ID: 2, Title: "Новость 2"},
			{ID: 3, Title: "Новость 3"},
		}
		json.NewEncoder(w).Encode(newsList)
	})

	// Обработчик для фильтрации новостей
	http.HandleFunc("/news/filter", func(w http.ResponseWriter, r *http.Request) {
		// Здесь можно добавить логику фильтрации
	})

	// Обработчик для получения детальной новости
	http.HandleFunc("/news/detail", func(w http.ResponseWriter, r *http.Request) {
		// Здесь можно добавить логику получения детальной новости
	})

	// Обработчик для добавления комментария
	http.HandleFunc("/comments/add", func(w http.ResponseWriter, r *http.Request) {
		// Получение данных из запроса
		var comment Comment
		err := json.NewDecoder(r.Body).Decode(&comment)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		// Сохранение комментария в базе данных
		_, err = db.Exec("INSERT INTO comments (text, news_id, parent_id) VALUES (?, ?, ?)",
			comment.Text, comment.NewsID, comment.ParentID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Отправка ответа
		w.WriteHeader(http.StatusCreated)
	})

	// Обработчик для получения всех комментариев по ID новости
	http.HandleFunc("/comments/news", func(w http.ResponseWriter, r *http.Request) {
		// Получение ID новости из параметров запроса
		newsID := r.URL.Query().Get("news_id")

		// Запрос всех комментариев по ID новости из базы данных
		rows, err := db.Query("SELECT id, text, news_id, parent_id FROM comments WHERE news_id = ?", newsID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		// Преобразование результатов запроса в структуру Comment и запись в массив
		var comments []Comment
		for rows.Next() {
			var comment Comment
			if err := rows.Scan(&comment.ID, &comment.Text, &comment.NewsID, &comment.ParentID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			comments = append(comments, comment)
		}

		// Отправка массива комментариев в виде JSON
		json.NewEncoder(w).Encode(comments)
	})

	// Запуск сервера на порту 8080
	log.Fatal(http.ListenAndServe(":8080", nil))
}
