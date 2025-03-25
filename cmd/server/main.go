package main

import (
	"fmt"

	"github.com/restapi/internal/repository"
	"github.com/restapi/internal/service/comment"
	transportHttp "github.com/restapi/internal/transport/http"

	_ "github.com/lib/pq"
	"github.com/restapi/internal/db"
)

func Run() error {
	db, err := db.NewDatabase()
	if err != nil {
		return fmt.Errorf("failed to initialise database: %w", err)
	}
	if err := db.MigrateDB(); err != nil {
		return fmt.Errorf("failed to migrate database: %w", err)
	}
	//comment service
	cmtStore := repository.NewCommentStore(db.Client)
	cmtService := comment.NewService(cmtStore)
	handler := transportHttp.NewHandler(cmtService)

	if err = handler.Serve(); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Go rest api course")
	if err := Run(); err != nil {
		fmt.Println(err)
	}
}
