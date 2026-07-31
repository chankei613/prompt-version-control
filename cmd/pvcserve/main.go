// cmd/pvcserve はPrompt Version Control APIをlocalhostで提供する単体サーバー。
//
//	go run ./cmd/pvcserve -addr :8425 -db prompt-version-control.db
package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/chankei613/prompt-version-control/internal/api"
	"github.com/chankei613/prompt-version-control/internal/db"
)

func main() {
	addr := flag.String("addr", ":8425", "待ち受けアドレス")
	dbPath := flag.String("db", "prompt-version-control.db", "SQLiteファイル")
	flag.Parse()

	conn, err := db.Init(*dbPath)
	if err != nil {
		log.Fatalf("db init failed: %v", err)
	}

	router := api.NewRouter(conn)
	log.Printf("prompt-version-control backend listening on %s", *addr)
	if err := http.ListenAndServe(*addr, router); err != nil {
		log.Fatal(err)
	}
}
