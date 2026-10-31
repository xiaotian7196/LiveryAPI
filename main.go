package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	_ "modernc.org/sqlite"
)

type Livery struct {
	ID                int64  `json:"id"`
	Name              string `json:"name"`
	Description       string `json:"description,omitempty"`
	Game              string `json:"game,omitempty"`
	Aircraft          string `json:"aircraft,omitempty"`
	ImageURL          string `json:"image_url,omitempty"`
	DownloadURLCN     string `json:"download_url_cn"`
	DownloadURLGlobal string `json:"download_url_global"`
	Version           string `json:"version,omitempty"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

type listResponse struct {
	Items      []Livery `json:"items"`
	Page       int      `json:"page"`
	PageSize   int      `json:"page_size"`
	Total      int      `json:"total"`
	TotalPages int      `json:"total_pages"`
}

type app struct{ db *sql.DB }

func main() {
	db, err := sql.Open("sqlite", getenv("DATABASE_URL", "file:liveries.db"))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	db.SetMaxOpenConns(1)
	if err := migrate(db); err != nil {
		log.Fatal(err)
	}

	port := getenv("PORT", "8080")
	server := &http.Server{Addr: ":" + port, Handler: routes(db), ReadHeaderTimeout: 5 * time.Second}
	log.Printf("livery API listening on http://localhost:%s", port)
	log.Fatal(server.ListenAndServe())
}

func routes(db *sql.DB) http.Handler {
	a := &app{db: db}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", a.health)
	mux.HandleFunc("GET /api/v1/liveries", a.listLiveries)
	mux.HandleFunc("GET /api/v1/liveries/", a.getLivery)
	return cors(logging(mux))
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS liveries (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  name TEXT NOT NULL,
  description TEXT NOT NULL DEFAULT '',
  game TEXT NOT NULL DEFAULT '',
  aircraft TEXT NOT NULL DEFAULT '',
  image_url TEXT NOT NULL DEFAULT '',
  download_url_cn TEXT NOT NULL DEFAULT '',
  download_url_global TEXT NOT NULL DEFAULT '',
  version TEXT NOT NULL DEFAULT '',
  created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS idx_liveries_name ON liveries(name);
CREATE INDEX IF NOT EXISTS idx_liveries_game ON liveries(game);
`)
	if err != nil {
		return err
	}
	// 兼容旧版本数据库：旧版本只有 download_url 字段。
	if err := ensureColumn(db, "download_url_cn", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if err := ensureColumn(db, "download_url_global", "TEXT NOT NULL DEFAULT ''"); err != nil {
		return err
	}
	if hasColumn, _ := tableHasColumn(db, "download_url"); hasColumn {
		_, err = db.Exec(`UPDATE liveries SET download_url_cn = download_url WHERE download_url_cn = ''`)
	}
	return err
}

func ensureColumn(db *sql.DB, name, definition string) error {
	exists, err := tableHasColumn(db, name)
	if err != nil || exists {
		return err
	}
	_, err = db.Exec(`ALTER TABLE liveries ADD COLUMN ` + name + ` ` + definition)
	return err
}

func tableHasColumn(db *sql.DB, name string) (bool, error) {
	rows, err := db.Query(`PRAGMA table_info(liveries)`)
	if err != nil {
		return false, err
	}
	defer rows.Close()
	for rows.Next() {
		var cid int
		var column, columnType string
		var notNull, primaryKey int
		var defaultValue any
		if err := rows.Scan(&cid, &column, &columnType, &notNull, &defaultValue, &primaryKey); err != nil {
			return false, err
		}
		if column == name {
			return true, nil
		}
	}
	return false, rows.Err()
}

func (a *app) health(w http.ResponseWriter, r *http.Request) {
	if err := a.db.PingContext(r.Context()); err != nil {
		writeError(w, http.StatusServiceUnavailable, "database unavailable")
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *app) listLiveries(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	page := positiveInt(q.Get("page"), 1)
	pageSize := positiveInt(q.Get("page_size"), 20)
	if pageSize > 100 {
		pageSize = 100
	}
	search, game := strings.TrimSpace(q.Get("search")), strings.TrimSpace(q.Get("game"))
	where := []string{"1 = 1"}
	args := []any{}
	if search != "" {
		where = append(where, "(name LIKE ? OR description LIKE ? OR aircraft LIKE ?)")
		like := "%" + search + "%"
		args = append(args, like, like, like)
	}
	if game != "" {
		where = append(where, "game = ?")
		args = append(args, game)
	}
	condition := strings.Join(where, " AND ")

	var total int
	if err := a.db.QueryRowContext(r.Context(), "SELECT COUNT(*) FROM liveries WHERE "+condition, args...).Scan(&total); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to count liveries")
		return
	}
	offset := (page - 1) * pageSize
	queryArgs := append(append([]any{}, args...), pageSize, offset)
	rows, err := a.db.QueryContext(r.Context(), `SELECT id, name, description, game, aircraft, image_url, download_url_cn, download_url_global, version, created_at, updated_at
FROM liveries WHERE `+condition+` ORDER BY id DESC LIMIT ? OFFSET ?`, queryArgs...)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query liveries")
		return
	}
	defer rows.Close()
	items := make([]Livery, 0)
	for rows.Next() {
		var item Livery
		if err := rows.Scan(&item.ID, &item.Name, &item.Description, &item.Game, &item.Aircraft, &item.ImageURL, &item.DownloadURLCN, &item.DownloadURLGlobal, &item.Version, &item.CreatedAt, &item.UpdatedAt); err != nil {
			writeError(w, http.StatusInternalServerError, "failed to read livery")
			return
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to read liveries")
		return
	}
	writeJSON(w, http.StatusOK, listResponse{Items: items, Page: page, PageSize: pageSize, Total: total, TotalPages: (total + pageSize - 1) / pageSize})
}

func (a *app) getLivery(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseInt(strings.TrimPrefix(r.URL.Path, "/api/v1/liveries/"), 10, 64)
	if err != nil || id <= 0 {
		writeError(w, http.StatusBadRequest, "invalid livery id")
		return
	}
	var item Livery
	err = a.db.QueryRowContext(r.Context(), `SELECT id, name, description, game, aircraft, image_url, download_url_cn, download_url_global, version, created_at, updated_at FROM liveries WHERE id = ?`, id).
		Scan(&item.ID, &item.Name, &item.Description, &item.Game, &item.Aircraft, &item.ImageURL, &item.DownloadURLCN, &item.DownloadURLGlobal, &item.Version, &item.CreatedAt, &item.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		writeError(w, http.StatusNotFound, "livery not found")
		return
	}
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query livery")
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.RequestURI(), time.Since(start))
	})
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]any{"error": map[string]any{"code": status, "message": message}})
}

func positiveInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	n, err := strconv.Atoi(value)
	if err != nil || n < 1 {
		return fallback
	}
	return n
}

func getenv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
