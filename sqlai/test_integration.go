package sqlai

import (
	"fmt"
	"log"
)

func main() {
	// Simple sample database structure for testing
	structureJSON := `{
		"tables": [
			{
				"name": "users",
				"columns": [
					{"name": "id", "type": "int", "isPrimary": true},
					{"name": "name", "type": "varchar"},
					{"name": "email", "type": "varchar"},
					{"name": "created_at", "type": "datetime"}
				],
				"foreignKeys": []
			},
			{
				"name": "products",
				"columns": [
					{"name": "id", "type": "int", "isPrimary": true},
					{"name": "name", "type": "varchar"},
					{"name": "price", "type": "decimal"},
					{"name": "category_id", "type": "int"}
				],
				"foreignKeys": [
					{
						"columnName": "category_id",
						"referencedTable": "categories",
						"referencedColumn": "id"
					}
				]
			}
		],
		"dbType": "sqlite"
	}`

	// Test the direct integration
	fmt.Println("Testing SQL AI integration:")
	fmt.Println("==========================")

	// Test the integration
	sql, err := GenerateSQLFromPrompt("Show me all users", structureJSON)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	fmt.Println("Generated SQL:", sql)

	// Test feedback
	RecordQueryFeedback(sql, true, "", 5, 0.1)
	fmt.Println("Feedback recorded successfully")

	// Test reset
	ResetSQLAssistant()
	fmt.Println("Assistant reset successfully")
}
