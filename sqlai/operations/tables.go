package operations

import (
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

// FindMostRelevantTable determines the most relevant table based on the query context
func FindMostRelevantTable(prompt string, dbTables []models.TableForAI) string {
	// Score each table based on context clues in the prompt
	tableScores := make(map[string]float64)

	// Look for singular/plural forms, word boundaries, etc.
	for _, table := range dbTables {
		tableName := strings.ToLower(table.Name)
		singularName := strings.TrimSuffix(tableName, "s")
		pluralName := tableName + "s"

		// Check for exact matches with word boundaries
		if RegexMatch(`\b`+tableName+`\b`, prompt) {
			tableScores[table.Name] = 1.0
			continue
		}

		// Check for singular/plural forms
		if RegexMatch(`\b`+singularName+`\b`, prompt) {
			tableScores[table.Name] = 0.9
			continue
		}

		if RegexMatch(`\b`+pluralName+`\b`, prompt) {
			tableScores[table.Name] = 0.9
			continue
		}

		// Check for partial matches
		if strings.Contains(prompt, tableName) {
			tableScores[table.Name] = 0.7
		}

		// Check for semantic relevance based on domain terminology
		domainTerms := map[string][]string{
			"user":     {"people", "person", "customer", "member", "client", "profile", "usuario", "usuário", "cliente", "pessoa"},
			"product":  {"item", "goods", "merchandise", "stock", "produto", "mercadoria", "estoque", "item"},
			"order":    {"purchase", "transaction", "sale", "invoice", "pedido", "compra", "venda", "transação", "fatura"},
			"payment":  {"transaction", "invoice", "billing", "receipt", "pagamento", "fatura", "recibo", "cobrança"},
			"category": {"group", "section", "department", "class", "type", "categoria", "grupo", "seção", "departamento", "tipo"},
		}

		for term, synonyms := range domainTerms {
			if strings.Contains(tableName, term) {
				for _, synonym := range synonyms {
					if strings.Contains(prompt, synonym) {
						tableScores[table.Name] += 0.4
					}
				}
			}
		}
	}

	// Find the table with the highest score
	var highestScore float64
	var mostRelevantTable string

	for tableName, score := range tableScores {
		if score > highestScore {
			highestScore = score
			mostRelevantTable = tableName
		}
	}

	return mostRelevantTable
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
