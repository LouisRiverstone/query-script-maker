package sqlai

import (
	"fmt"
	"sql_script_maker/sqlai/models"
	"time"
)

// GenerateSQLFromPrompt is a wrapper function for apps to use
func GenerateSQLFromPrompt(prompt string, structureJSON string) (string, error) {
	// Initialize the assistant with the database structure
	assistant := GetSQLAssistant()
	err := assistant.Init(structureJSON)
	if err != nil {
		return "", fmt.Errorf("error initializing SQL assistant: %w", err)
	}

	// Generate SQL from the prompt
	sql, err := assistant.GenerateSQL(prompt)
	if err != nil {
		return "", fmt.Errorf("error generating SQL: %w", err)
	}

	return sql, nil
}

// ResetSQLAssistant resets the assistant's state
func ResetSQLAssistant() {
	assistant := GetSQLAssistant()
	assistant.Reset()
}

// RecordQueryFeedback records query execution feedback for learning
func RecordQueryFeedback(query string, successful bool, errorMsg string, rowCount int, execTime float64) {
	assistant := GetSQLAssistant()
	assistant.RecordQueryFeedback(models.FeedbackResult{
		Query:         query,
		WasSuccessful: successful,
		ErrorMessage:  errorMsg,
		RowCount:      rowCount,
		ExecutionTime: execTime,
		CreatedAt:     time.Now(),
	})
}
