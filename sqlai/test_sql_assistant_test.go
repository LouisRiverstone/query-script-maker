package sqlai_test

import (
	"sql_script_maker/sqlai"
	"testing"
)

// TestWhereCondition testa condições específicas no WHERE de consultas SQL
func TestWhereCondition(t *testing.T) {
	assistant := sqlai.GetSQLAssistant()

	// Mock da estrutura do banco de dados para os testes
	mockDB := `{
		"tables": [
			{
				"name": "clan",
				"columns": [
					{
						"name": "id",
						"type": "int",
						"isPrimary": true
					},
					{
						"name": "name",
						"type": "varchar",
						"isPrimary": false
					}
				]
			},
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
			},
			{
				"name": "author",
				"columns": [
					{
						"name": "id",
						"type": "int",
						"isPrimary": true
					},
					{
						"name": "name",
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
			name:     "Condição com 'é igual a' numérico",
			prompt:   "selecione a coluna rfam_id da tabela family onde a coluna rfam_acc é igual a 3",
			expected: "SELECT family.rfam_id FROM family WHERE family.rfam_acc = 3",
		},
		{
			name:     "Condição com 'menor ou igual a' numérico",
			prompt:   "selecione a coluna id da tabela clan onde id for menor ou igual a 100",
			expected: "SELECT clan.id FROM clan WHERE clan.id <= 100",
		},
		{
			name:     "Condição com 'menor que' numérico",
			prompt:   "selecione a coluna id da tabela clan onde id for menor que 100",
			expected: "SELECT clan.id FROM clan WHERE clan.id < 100",
		},
		{
			name:     "Condição com 'maior que' numérico",
			prompt:   "selecione a coluna id da tabela clan onde id for maior que 50",
			expected: "SELECT clan.id FROM clan WHERE clan.id > 50",
		},
		{
			name:     "Condição com 'maior ou igual a' numérico",
			prompt:   "selecione a coluna id da tabela clan onde id for maior ou igual a 50",
			expected: "SELECT clan.id FROM clan WHERE clan.id >= 50",
		},
		{
			name:     "Condição com 'diferente de' numérico",
			prompt:   "selecione a coluna id da tabela clan onde id for diferente de 50",
			expected: "SELECT clan.id FROM clan WHERE clan.id != 50",
		},
		{
			name:     "Condição com 'diferente de' texto",
			prompt:   "selecione a coluna name da tabela author onde name for diferente de João",
			expected: "SELECT author.name FROM author WHERE author.name != 'João'",
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
