package operations

import (
	"regexp"
	"strings"

	"sql_script_maker/sqlai/models"
	"sql_script_maker/sqlai/util"
)

// SortTablesByConfidence sorts tables by confidence score (descending)
func SortTablesByConfidence(tables []models.TableInfo) {
	for i := 0; i < len(tables)-1; i++ {
		for j := i + 1; j < len(tables); j++ {
			if tables[i].Confidence < tables[j].Confidence {
				tables[i], tables[j] = tables[j], tables[i]
			}
		}
	}
}

// MakeSingular converts a word to its singular form (simplified version)
func MakeSingular(word string) string {
	// Basic implementation for common plurals
	if strings.HasSuffix(word, "s") && len(word) > 2 {
		// Don't singularize words that naturally end in 's'
		exceptions := []string{"status", "analysis", "chassis", "bus"}
		for _, exception := range exceptions {
			if word == exception {
				return word
			}
		}
		return word[:len(word)-1]
	}
	return word
}

// MakePlural converts a word to its plural form (simplified version)
func MakePlural(word string) string {
	// Don't pluralize if already plural or special cases
	if strings.HasSuffix(word, "s") {
		return word
	}

	// Handle some special cases
	if strings.HasSuffix(word, "y") && len(word) > 1 {
		// Check if the letter before 'y' is a consonant
		vowels := "aeiou"
		lastBeforeY := word[len(word)-2]
		if !strings.ContainsRune(vowels, rune(lastBeforeY)) {
			return word[:len(word)-1] + "ies"
		}
	}

	return word + "s"
}

// FindMostRelevantTable determines the most relevant table based on context clues
func FindMostRelevantTable(prompt string, tables []models.TableForAI) string {
	// Convert the prompt to lowercase for better matching
	promptLower := strings.ToLower(prompt)

	// Check for explicit mentions of tables first using regex
	explicitTableRegex := regexp.MustCompile(`(?i)(?:da|de|from)\s+tabela\s+(\w+)`)
	if matches := explicitTableRegex.FindStringSubmatch(prompt); len(matches) > 1 {
		explicitTable := matches[1]
		// Verify this table exists in our schema with exact match or close match
		for _, table := range tables {
			if strings.EqualFold(table.Name, explicitTable) {
				return table.Name // Exact match (case-insensitive)
			}
		}

		// Try partial match if exact match fails
		for _, table := range tables {
			if strings.Contains(strings.ToLower(table.Name), strings.ToLower(explicitTable)) {
				return table.Name
			}
		}
	}

	// Define a scoring system for tables
	// We'll score each table based on multiple criteria and pick the highest score
	tableScores := make(map[string]float64)

	// Initialize all tables with a base score
	for _, table := range tables {
		tableScores[table.Name] = 0.1 // Small base score
	}

	// List of common table names to identify in the prompt
	// This helps with common domains where specific terms correlate to tables
	domainTerms := map[string][]string{
		"user":      {"usuário", "usuários", "user", "users", "cliente", "clientes", "customer", "customers", "pessoa", "pessoas", "person", "people"},
		"product":   {"produto", "produtos", "product", "products", "item", "items", "mercadoria", "mercadorias", "merchandise"},
		"order":     {"pedido", "pedidos", "order", "orders", "compra", "compras", "purchase", "purchases", "venda", "vendas", "sale", "sales"},
		"category":  {"categoria", "categorias", "category", "categories", "classificação", "classificações", "classification"},
		"payment":   {"pagamento", "pagamentos", "payment", "payments", "transação", "transações", "transaction", "transactions"},
		"comment":   {"comentário", "comentários", "comment", "comments", "review", "reviews", "avaliação", "avaliações"},
		"post":      {"postagem", "postagens", "post", "posts", "artigo", "artigos", "article", "articles", "publicação", "publicações", "publication"},
		"address":   {"endereço", "endereços", "address", "addresses", "localização", "localizações", "location", "locations"},
		"inventory": {"estoque", "estoques", "inventory", "inventories", "armazém", "armazéns", "warehouse", "warehouses"},
		"employee":  {"funcionário", "funcionários", "employee", "employees", "trabalhador", "trabalhadores", "worker", "workers"},
	}

	// Score based on domain terminology
	for tableName, terms := range domainTerms {
		for _, term := range terms {
			if strings.Contains(promptLower, term) {
				// Look for matching or related tables
				for _, table := range tables {
					tableLower := strings.ToLower(table.Name)

					// Exact match with domain term
					if tableLower == tableName {
						tableScores[table.Name] += 2.0
					}

					// Partial match with domain term
					// For example, "users" table matches "user" term
					if strings.Contains(tableLower, tableName) {
						tableScores[table.Name] += 1.5
					}

					// Related table by domain (weaker match)
					// For example "user_address" table relates to "address" term
					if strings.Contains(tableLower, term) {
						tableScores[table.Name] += 1.0
					}
				}
			}
		}
	}

	// Score based on singular/plural forms of table names
	for _, table := range tables {
		// Try to match the table name directly
		tableLower := strings.ToLower(table.Name)
		singularForm := MakeSingular(tableLower)
		pluralForm := MakePlural(tableLower)

		// Check for exact matches
		if strings.Contains(promptLower, tableLower) {
			tableScores[table.Name] += 3.0
		}

		// Check for singular form match
		if singularForm != tableLower && strings.Contains(promptLower, singularForm) {
			tableScores[table.Name] += 2.5
		}

		// Check for plural form match
		if pluralForm != tableLower && strings.Contains(promptLower, pluralForm) {
			tableScores[table.Name] += 2.5
		}
	}

	// Check for explicit column mentions
	// This helps when the user mentions a column but not the table
	columnRegex := regexp.MustCompile(`(?i)(?:a coluna|as colunas|column[s]?)\s+(\w+)`)
	if matches := columnRegex.FindStringSubmatch(prompt); len(matches) > 1 {
		explicitColumn := matches[1]

		// Check which tables have this column
		for _, table := range tables {
			for _, col := range table.Columns {
				if strings.EqualFold(col.Name, explicitColumn) {
					tableScores[table.Name] += 2.0
				}
			}
		}
	}

	// Check for operations that make one table more likely
	insertOperations := []string{"insert", "add", "create", "inserir", "adicionar", "criar", "nova", "novo"}

	// Tables that often receive these operations get a boost
	for _, term := range insertOperations {
		if strings.Contains(promptLower, term) {
			// Boost tables that are commonly modified
			for _, table := range tables {
				tableLower := strings.ToLower(table.Name)
				if strings.Contains(tableLower, "user") ||
					strings.Contains(tableLower, "order") ||
					strings.Contains(tableLower, "product") {
					tableScores[table.Name] += 0.5
				}
			}
		}
	}

	// Find the table with the highest score
	var bestTable string
	var highestScore float64

	for tableName, score := range tableScores {
		if score > highestScore {
			highestScore = score
			bestTable = tableName
		}
	}

	// If we have a reasonable confidence (score threshold)
	if highestScore >= 1.0 {
		return bestTable
	}

	// If we have tables but confidence is too low, return the most generic table (often users or products)
	if len(tables) > 0 {
		for _, tableName := range []string{"users", "user", "products", "product", "customers", "customer"} {
			for _, table := range tables {
				if strings.EqualFold(table.Name, tableName) {
					return table.Name
				}
			}
		}

		// If no commonly used table found, return the first table
		return tables[0].Name
	}

	return ""
}

// SelectRelevantColumns chooses appropriate columns based on the query context
func SelectRelevantColumns(tables []models.TableInfo, operation string, prompt string, dbTables []models.TableForAI) []models.ColumnInfo {
	var selectedColumns []models.ColumnInfo

	if len(tables) == 0 {
		return []models.ColumnInfo{{Name: "*", TableName: "", Confidence: 1.0}}
	}

	// For count operations, use COUNT(*) by default
	if operation == "count" {
		return []models.ColumnInfo{{Name: "*", TableName: tables[0].Name, Confidence: 1.0, Function: "COUNT"}}
	}

	// Look for aggregation functions in the prompt
	aggregationFunctions := map[string]string{
		"average":            "AVG(%s)",
		"mean":               "AVG(%s)",
		"sum":                "SUM(%s)",
		"total":              "SUM(%s)",
		"maximum":            "MAX(%s)",
		"highest":            "MAX(%s)",
		"minimum":            "MIN(%s)",
		"lowest":             "MIN(%s)",
		"count":              "COUNT(%s)",
		"count distinct":     "COUNT(DISTINCT %s)",
		"unique count":       "COUNT(DISTINCT %s)",
		"standard deviation": "STDDEV(%s)",
		"variance":           "VARIANCE(%s)",
		"median":             "PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY %s)",
		// Portuguese terms
		"média":             "AVG(%s)",
		"soma":              "SUM(%s)",
		"máximo":            "MAX(%s)",
		"maior":             "MAX(%s)",
		"mínimo":            "MIN(%s)",
		"menor":             "MIN(%s)",
		"contagem":          "COUNT(%s)",
		"contagem distinta": "COUNT(DISTINCT %s)",
		"desvio padrão":     "STDDEV(%s)",
		"variância":         "VARIANCE(%s)",
		"mediana":           "PERCENTILE_CONT(0.5) WITHIN GROUP (ORDER BY %s)",
	}

	for funcName := range aggregationFunctions {
		if strings.Contains(prompt, funcName) {
			// Find a numeric column to apply the function to
			for _, table := range dbTables {
				if table.Name == tables[0].Name {
					for _, col := range table.Columns {
						if strings.Contains(col.Type, "int") ||
							strings.Contains(col.Type, "float") ||
							strings.Contains(col.Type, "decimal") ||
							strings.Contains(col.Type, "double") {
							selectedColumns = append(selectedColumns, models.ColumnInfo{
								Name:       col.Name,
								TableName:  tables[0].Name,
								Type:       col.Type,
								Confidence: 0.8,
								Function:   strings.ToUpper(funcName),
							})
							break
						}
					}
				}
			}
		}
	}

	// If we found aggregate functions, return those columns
	if len(selectedColumns) > 0 {
		return selectedColumns
	}

	// Otherwise, use meaningful columns from tables
	for _, tableInfo := range tables {
		for _, table := range dbTables {
			if table.Name == tableInfo.Name {
				// Add ID column
				for _, col := range table.Columns {
					if col.IsPrimary || strings.ToLower(col.Name) == "id" {
						selectedColumns = append(selectedColumns, models.ColumnInfo{
							Name:       col.Name,
							TableName:  table.Name,
							Type:       col.Type,
							IsPrimary:  col.IsPrimary,
							Confidence: 0.9,
						})
						break
					}
				}

				// Add descriptive columns (name, title, etc.)
				for _, col := range table.Columns {
					colLower := strings.ToLower(col.Name)
					if util.ContainsAny(colLower, []string{"name", "title", "description", "label", "nome", "titulo", "descrição", "etiqueta"}) {
						selectedColumns = append(selectedColumns, models.ColumnInfo{
							Name:       col.Name,
							TableName:  table.Name,
							Type:       col.Type,
							Confidence: 0.8,
						})
					}
				}

				// Add date/time columns
				for _, col := range table.Columns {
					if strings.Contains(strings.ToLower(col.Type), "date") ||
						strings.Contains(strings.ToLower(col.Type), "time") {
						selectedColumns = append(selectedColumns, models.ColumnInfo{
							Name:       col.Name,
							TableName:  table.Name,
							Type:       col.Type,
							Confidence: 0.7,
						})
						break
					}
				}
			}
		}
	}

	// If we still don't have columns, use all columns from primary table or *
	if len(selectedColumns) == 0 {
		return []models.ColumnInfo{{Name: "*", TableName: tables[0].Name, Confidence: 1.0}}
	}

	return selectedColumns
}

// FindJoinCondition determines how two tables should be joined
func FindJoinCondition(table1, table2 string, dbTables []models.TableForAI) string {
	// First check for foreign key relationships
	for _, table := range dbTables {
		if table.Name == table1 {
			for _, fk := range table.ForeignKeys {
				if fk.ReferencedTable == table2 {
					return strings.Join([]string{
						table1,
						fk.ColumnName,
						"=",
						table2,
						fk.ReferencedColumn,
					}, ".")
				}
			}
		} else if table.Name == table2 {
			for _, fk := range table.ForeignKeys {
				if fk.ReferencedTable == table1 {
					return strings.Join([]string{
						table2,
						fk.ColumnName,
						"=",
						table1,
						fk.ReferencedColumn,
					}, ".")
				}
			}
		}
	}

	// Next, try common naming patterns
	table1IdCol := table1 + "_id"
	table2IdCol := table2 + "_id"

	// Check if table1 has a column referencing table2
	for _, table := range dbTables {
		if table.Name == table1 {
			for _, col := range table.Columns {
				if col.Name == table2IdCol || col.Name == strings.TrimSuffix(table2, "s")+"_id" {
					return table1 + "." + col.Name + " = " + table2 + ".id"
				}
			}
		}
	}

	// Check if table2 has a column referencing table1
	for _, table := range dbTables {
		if table.Name == table2 {
			for _, col := range table.Columns {
				if col.Name == table1IdCol || col.Name == strings.TrimSuffix(table1, "s")+"_id" {
					return table2 + "." + col.Name + " = " + table1 + ".id"
				}
			}
		}
	}

	// If no FK or naming pattern found, join on IDs as a last resort
	return table1 + ".id = " + table2 + ".id"
}
