package test_sqlai

import (
	"fmt"
	"log"
	"sql_script_maker/sqlai"
)

func TestSQLAI() {
	// Test both English and Portuguese queries
	testQueries := []struct {
		name   string
		prompt string
		lang   string
	}{
		{"English Simple Select", "Show me all users", "English"},
		{"English Join", "Find all orders along with customer information", "English"},
		{"English Count with Condition", "How many products have a price greater than 100", "English"},
		{"English Group By", "Count products by category", "English"},

		{"Portuguese Simple Select", "Mostre todos os usuários", "Portuguese"},
		{"Portuguese Join", "Encontre todos os pedidos junto com as informações do cliente", "Portuguese"},
		{"Portuguese Count with Condition", "Quantos produtos têm preço maior que 100", "Portuguese"},
		{"Portuguese Group By", "Contar produtos por categoria", "Portuguese"},
	}

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
			},
			{
				"name": "categories",
				"columns": [
					{"name": "id", "type": "int", "isPrimary": true},
					{"name": "name", "type": "varchar"}
				],
				"foreignKeys": []
			},
			{
				"name": "orders",
				"columns": [
					{"name": "id", "type": "int", "isPrimary": true},
					{"name": "user_id", "type": "int"},
					{"name": "total", "type": "decimal"},
					{"name": "created_at", "type": "datetime"}
				],
				"foreignKeys": [
					{
						"columnName": "user_id",
						"referencedTable": "users",
						"referencedColumn": "id"
					}
				]
			},
			{
				"name": "order_items",
				"columns": [
					{"name": "id", "type": "int", "isPrimary": true},
					{"name": "order_id", "type": "int"},
					{"name": "product_id", "type": "int"},
					{"name": "quantity", "type": "int"},
					{"name": "price", "type": "decimal"}
				],
				"foreignKeys": [
					{
						"columnName": "order_id",
						"referencedTable": "orders",
						"referencedColumn": "id"
					},
					{
						"columnName": "product_id",
						"referencedTable": "products",
						"referencedColumn": "id"
					}
				]
			}
		],
		"dbType": "sqlite"
	}`

	// Initialize the SQL assistant
	assistant := sqlai.GetSQLAssistant()
	err := assistant.Init(structureJSON)
	if err != nil {
		log.Fatalf("Error initializing SQL assistant: %v", err)
	}

	// Run tests
	fmt.Println("Testing SQL AI with both English and Portuguese queries:")
	fmt.Println("======================================================")

	for _, test := range testQueries {
		fmt.Printf("\n[%s - %s]\n", test.lang, test.name)
		fmt.Printf("Prompt: %s\n", test.prompt)

		sql, err := assistant.GenerateSQL(test.prompt)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}

		fmt.Printf("Generated SQL: %s\n", sql)
	}
}
