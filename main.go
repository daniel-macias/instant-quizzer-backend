package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/daniel-macias/instant-quizzer-backend/handlers"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	ctx := context.Background()

	// Access environment variable using os.Getenv
	uri := os.Getenv("MONGODB_URI")
	if uri == "" {
		log.Fatal("MONGODB_URI environment variable not found!")
	}

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := client.Disconnect(ctx); err != nil {
			log.Println("Error disconnecting from MongoDB:", err)
		}
	}()

	if err := client.Database("admin").RunCommand(ctx, bson.D{{"ping", 1}}).Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	h := handlers.NewHandler(client)

	router := mux.NewRouter()

	router.HandleFunc("/api/quizzes", h.CreateQuiz).Methods("POST")
	router.HandleFunc("/api/quizzes", h.GetAllQuizzes).Methods("GET")
	router.HandleFunc("/api/quizzes/{id}", h.GetQuizByID).Methods("GET")
	router.HandleFunc("/api/quizzes/{id}", h.UpdateQuiz).Methods("PUT")
	router.HandleFunc("/api/quizzes/{id}", h.DeleteQuiz).Methods("DELETE")
	router.HandleFunc("/api/quizzes/{id}/results", h.AddResult).Methods("POST")

	log.Fatal(http.ListenAndServe(":8000", router))
}
