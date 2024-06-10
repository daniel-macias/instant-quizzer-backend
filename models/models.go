package models

// Quiz struct represents a quiz document in MongoDB
type Quiz struct {
	ID        string     `bson:"_id,omitempty" json:"id,omitempty"` // Use omitempty for optional export during marshalling
	QuizTitle string     `bson:"quizTitle"`                         // Use bson tags for mapping to MongoDB fields
	Questions []Question `bson:"questions"`                         // Nested struct for questions
	Results   []Result   `bson:"results"`                           // Nested struct for results
}

// Result struct represents a quiz result document in MongoDB
type Result struct {
	PersonName string `bson:"personName"`
	Responses  []bool `bson:"responses"`
}

// Question struct represents a question within a quiz
type Question struct {
	QuestionTitle   string   `bson:"questionTitle"`
	PossibleAnswers []string `bson:"possibleAnswers"`
	CorrectAnswers  []int    `bson:"correctAnswers"`
}
