package sqlai

import (
	"fmt"
	"log"
	"testing"
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

// TestWhereEqualCondition testa especificamente consultas com 'é igual a' na condição
func TestWhereEqualCondition(t *testing.T) {
	assistant := GetSQLAssistant()

	// Mock da estrutura do banco de dados para o teste
	mockDB := `{
		"tables": [
			{
				"name": "family",
				"columns": [
					{
						"name": "rfam_id",
						"type": "int",
						"isPrimary": true
					},
					{
						"name": "rfam_acc",
						"type": "varchar",
						"isPrimary": false
					}
				]
			}
		],
		"dbType": "mysql"
	}`

	err := assistant.Init(mockDB)
	if err != nil {
		t.Fatalf("Falha ao inicializar o assistente: %v", err)
	}

	testCases := []struct {
		name     string
		prompt   string
		expected string
	}{
		{
			name:     "Prompt com 'é igual a' numérico",
			prompt:   "selecione a coluna rfam_id da tabela family onde a coluna rfam_acc é igual a 3",
			expected: "SELECT family.rfam_id FROM family WHERE family.rfam_acc = 3",
		},
		{
			name:     "Prompt com 'é igual a' texto",
			prompt:   "selecione a coluna rfam_id da tabela family onde a coluna rfam_acc é igual a ABC123",
			expected: "SELECT family.rfam_id FROM family WHERE family.rfam_acc = 'ABC123'",
		},
		{
			name:     "Prompt com 'é' como operador",
			prompt:   "selecione a coluna rfam_id da tabela family onde a coluna rfam_acc é 3",
			expected: "SELECT family.rfam_id FROM family WHERE family.rfam_acc = 3",
		},
		{
			name:     "Prompt com '=' como operador",
			prompt:   "selecione a coluna rfam_id da tabela family onde a coluna rfam_acc = 3",
			expected: "SELECT family.rfam_id FROM family WHERE family.rfam_acc = 3",
		},
		{
			name:     "Prompt com 'menor ou igual a' numérico",
			prompt:   "selecione a coluna rfam_id da tabela family onde a coluna rfam_acc for menor ou igual a 100",
			expected: "SELECT family.rfam_id FROM family WHERE family.rfam_acc <= 100",
		},
		{
			name:     "Prompt com 'menor que' numérico",
			prompt:   "selecione a coluna rfam_id da tabela family onde a coluna rfam_acc for menor que 100",
			expected: "SELECT family.rfam_id FROM family WHERE family.rfam_acc < 100",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sql, err := assistant.GenerateSQL(tc.prompt)
			if err != nil {
				t.Fatalf("Erro ao gerar SQL: %v", err)
			}

			if sql != tc.expected {
				t.Errorf("SQL gerado incorreto.\nEsperado: %s\nObtido: %s", tc.expected, sql)
			}
		})
	}
}
