package handlers

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/daniel-macias/instant-quizzer-backend/models"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Handler struct {
	Client *mongo.Client
}

func NewHandler(client *mongo.Client) *Handler {
	return &Handler{Client: client}
}

func (h *Handler) CreateQuiz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var quiz models.Quiz
	_ = json.NewDecoder(r.Body).Decode(&quiz)

	collection := h.Client.Database("instant_quizzer").Collection("Quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := collection.InsertOne(ctx, quiz)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Extract the ID from the InsertOneResult and include it in the response
	id := result.InsertedID.(primitive.ObjectID).Hex()
	json.NewEncoder(w).Encode(map[string]string{"id": id})
}

func (h *Handler) GetAllQuizzes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var quizzes []models.Quiz

	collection := h.Client.Database("instant_quizzer").Collection("Quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer cursor.Close(ctx)

	for cursor.Next(ctx) {
		var quiz models.Quiz
		cursor.Decode(&quiz)
		quizzes = append(quizzes, quiz)
	}

	if err := cursor.Err(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(quizzes)
}

func (h *Handler) GetQuizByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Log the ID for debugging purposes
	log.Printf("Received ID: %s", params["id"])

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var quiz models.Quiz
	collection := h.Client.Database("instant_quizzer").Collection("Quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&quiz)
	if err != nil {
		http.Error(w, "Quiz not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(quiz)
}

func (h *Handler) UpdateQuiz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Log the ID for debugging purposes
	log.Printf("Received ID: %s", params["id"])

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var quiz models.Quiz
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database("instant_quizzer").Collection("Quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a copy of the quiz document without the _id field
	updateDoc := bson.M{
		"quizTitle": quiz.QuizTitle,
		"questions": quiz.Questions,
		"results":   quiz.Results,
	}

	filter := bson.M{"_id": id}
	update := bson.M{
		"$set": updateDoc,
	}

	_, err = collection.UpdateOne(ctx, filter, update)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Quiz updated successfully")
}

func (h *Handler) DeleteQuiz(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Log the ID for debugging purposes
	log.Printf("Received ID: %s", params["id"])

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database("instant_quizzer").Collection("Quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Quiz deleted successfully")
}

func (h *Handler) AddResult(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	// Log the ID for debugging purposes
	log.Printf("Received ID: %s", params["id"])

	id, err := primitive.ObjectIDFromHex(params["id"])
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	var newResult models.Result
	if err := json.NewDecoder(r.Body).Decode(&newResult); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	collection := h.Client.Database("instant_quizzer").Collection("Quizzes")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var quiz models.Quiz
	err = collection.FindOne(ctx, bson.M{"_id": id}).Decode(&quiz)
	if err != nil {
		http.Error(w, "Quiz not found", http.StatusNotFound)
		return
	}

	if len(newResult.Responses) != len(quiz.Questions) {
		http.Error(w, "The number of responses does not match the number of questions in the quiz.", http.StatusBadRequest)
		return
	}

	quiz.Results = append(quiz.Results, newResult)

	// Create a copy of the quiz document without the _id field
	updateDoc := bson.M{
		"quizTitle": quiz.QuizTitle,
		"questions": quiz.Questions,
		"results":   quiz.Results,
	}

	_, err = collection.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": updateDoc})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode("Results added successfully")
}
