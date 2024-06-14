package models

// Quiz struct represents a quiz document in MongoDB
type Quiz struct {
	ID        string     `bson:"_id,omitempty" json:"id,omitempty"` // Use omitempty for optional export during marshalling
	QuizTitle string     `bson:"quizTitle" json:"quizTitle"`        // Use bson and json tags for mapping to MongoDB fields and camelCase
	Questions []Question `bson:"questions" json:"questions"`        // Nested struct for questions
	Results   []Result   `bson:"results" json:"results"`            // Nested struct for results
}

// Result struct represents a quiz result document in MongoDB
type Result struct {
	PersonName string `bson:"personName" json:"personName"`
	Responses  []bool `bson:"responses" json:"responses"`
}

// Question struct represents a question within a quiz
type Question struct {
	QuestionTitle   string   `bson:"questionTitle" json:"questionTitle"`
	PossibleAnswers []string `bson:"possibleAnswers" json:"possibleAnswers"`
	CorrectAnswers  []int    `bson:"correctAnswers" json:"correctAnswers"`
}
