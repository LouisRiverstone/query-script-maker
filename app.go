package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	rntm "runtime"
	"strings"
	"sync"
	"time"

	"sql_script_maker/sqlai"
	sqlaiModels "sql_script_maker/sqlai/models"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
)

// App struct
type App struct {
	ctx context.Context
}

// Variable struct
type Variable struct {
	Field    string
	Value    string
	Position int
}

type Query struct {
	ID          *int
	Title       string
	Query       string
	Description string
	CreatedAt   *string
	UpdatedAt   *string
	DeletedAt   *string
}

type DatabaseConnection struct {
	ID        *int
	Username  string
	Password  string
	Host      string
	Port      int
	Database  string
	CreatedAt *string
	UpdatedAt *string
	DeletedAt *string
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

	createSqliteTables()
}

func (a *App) getSQLAssistant(structureJSON string) *sqlai.SQLAssistant {
	assistant := sqlai.GetSQLAssistant()
	err := assistant.Init(structureJSON)

	if err != nil {
		log.Fatal(err)
	}

	return assistant
}

func (a *App) ReadXLSXFile() (string, error) {
	// Selecionar o arquivo
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "XLSX (*.xlsx)",
				Pattern:     "*.xlsx",
			},
		},
	})

	if err != nil {
		return "", err
	}

	if selection == "" {
		return "", fmt.Errorf("no file selected")
	}

	// Abre o arquivo diretamente usando a biblioteca excelize
	// Isso é mais eficiente que usar os.Open() seguido de excelize.OpenReader()
	xlsxFile, err := excelize.OpenFile(selection)
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer xlsxFile.Close()

	// Obtém o nome da primeira planilha
	sheetName := xlsxFile.GetSheetName(0)

	// Usa o método Rows para processar em stream, evitando carregar todo o arquivo na memória
	rows, err := xlsxFile.Rows(sheetName)
	if err != nil {
		return "", fmt.Errorf("failed to get rows: %w", err)
	}
	defer rows.Close()

	// Inicialização de variáveis
	var headers []string
	content := make([]map[string]string, 0, 1000) // Pré-aloca com tamanho inicial razoável
	rowIndex := 0

	// Determina quantos CPUs temos disponíveis para uso
	numCPU := rntm.NumCPU()
	batchSize := 1000 // Tamanho do lote para processamento
	rowBatch := make([][]string, 0, batchSize)

	// Processamento em lotes para melhor desempenho
	for rows.Next() {
		rowData, err := rows.Columns()
		if err != nil {
			// Log do erro mas continua processando
			log.Printf("Error reading row: %v", err)
			continue
		}

		// Primeira linha contém os cabeçalhos
		if rowIndex == 0 {
			headers = make([]string, len(rowData))
			copy(headers, rowData)
			rowIndex++
			continue
		}

		// Armazena a linha no lote atual
		rowCopy := make([]string, len(rowData))
		copy(rowCopy, rowData)
		rowBatch = append(rowBatch, rowCopy)

		// Quando o lote atinge o tamanho definido, processa paralelamente
		if len(rowBatch) >= batchSize {
			processBatch(&content, rowBatch, headers, numCPU)
			rowBatch = make([][]string, 0, batchSize)
		}

		rowIndex++
	}

	// Processa o último lote se houver dados restantes
	if len(rowBatch) > 0 {
		processBatch(&content, rowBatch, headers, numCPU)
	}

	// Se não temos dados além dos cabeçalhos, retorna um array vazio
	if len(content) == 0 {
		return "[]", nil
	}

	// Usando a biblioteca padrão, mas com um buffer pré-alocado para melhor desempenho
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	err = encoder.Encode(content)
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}

	// Remove a quebra de linha que o Encode adiciona
	jsonBytes := buffer.Bytes()
	if len(jsonBytes) > 0 && jsonBytes[len(jsonBytes)-1] == '\n' {
		jsonBytes = jsonBytes[:len(jsonBytes)-1]
	}

	return string(jsonBytes), nil
}

// Função auxiliar para processar um lote de linhas em paralelo
func processBatch(content *[]map[string]string, batch [][]string, headers []string, numWorkers int) {
	// Limita o número de workers ao número de itens no lote
	if numWorkers > len(batch) {
		numWorkers = len(batch)
	}

	var wg sync.WaitGroup
	resultChannel := make(chan map[string]string, len(batch))

	// Divide o trabalho entre os workers
	chunkSize := (len(batch) + numWorkers - 1) / numWorkers
	for w := 0; w < numWorkers; w++ {
		wg.Add(1)
		start := w * chunkSize
		end := start + chunkSize
		if end > len(batch) {
			end = len(batch)
		}

		// Processa um intervalo de linhas em uma goroutine
		go func(startIdx, endIdx int) {
			defer wg.Done()

			for i := startIdx; i < endIdx; i++ {
				row := batch[i]
				rowData := make(map[string]string, len(headers))

				for j, cell := range row {
					if j < len(headers) {
						rowData[headers[j]] = cell
					}
				}

				resultChannel <- rowData
			}
		}(start, end)
	}

	// Goroutine para coletar resultados
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// Coleta os resultados do canal
	for rowData := range resultChannel {
		*content = append(*content, rowData)
	}
}

func (a *App) MakeBindedSQL(query string, data []map[string]interface{}, variables []Variable, minify bool) string {
	var result strings.Builder

	re := regexp.MustCompile(`{{ \w+ }}`)

	for index, row := range data {
		modifiedQuery := re.ReplaceAllStringFunc(query, func(match string) string {
			for _, variable := range variables {
				if match == fmt.Sprintf("{{ %s }}", variable.Value) {
					if value, ok := row[variable.Field]; ok {
						return fmt.Sprintf("%v", value)
					}
				}
			}
			return match
		})

		result.WriteString(modifiedQuery)

		if index != len(data)-1 {
			result.WriteString("\n")
		}
	}

	if minify {
		return strings.ReplaceAll(result.String(), "\n", " ")
	}

	return result.String()
}

func (a *App) CreateSQLFile(data string) (string, error) {
	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Save File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "SQL (*.sql)",
				Pattern:     "*.sql",
			},
		},
	})

	if err != nil {
		return "", err
	}

	if selection == "" {
		return "", fmt.Errorf("no file selected")
	}

	file, err := os.Create(selection)

	if err != nil {
		return "", err
	}

	defer file.Close()

	_, err = file.WriteString(data)

	if err != nil {
		return "", err
	}

	return selection, nil
}

func (a *App) InsertQueryInDatabase(data Query) error {
	db := openSqliteConnection()

	defer db.Close()

	tx, err := db.Begin()

	if err != nil {
		fmt.Println(err.Error())
	}

	insertQuery := `INSERT INTO queries(title, query, description) VALUES(?, ?, ?)`

	stmt, err := tx.Prepare(insertQuery)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer stmt.Close()

	_, err = stmt.Exec(data.Title, data.Query, data.Description)

	if err != nil {
		fmt.Println(err.Error())
	}

	err = tx.Commit()

	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (a *App) GetQueriesList(withTrashed bool) ([]Query, error) {
	db := openSqliteConnection()
	defer db.Close()

	var queries []Query = make([]Query, 0)

	var selectQuery string

	if withTrashed {
		selectQuery = `SELECT * FROM queries`
	} else {
		selectQuery = `SELECT * FROM queries WHERE deleted_at IS NULL`
	}

	rows, err := db.Query(selectQuery)

	if err != nil {
		fmt.Println(err.Error())
	}

	defer rows.Close()

	for rows.Next() {
		var query Query

		err = rows.Scan(&query.ID, &query.Title, &query.Query, &query.Description, &query.CreatedAt, &query.UpdatedAt, &query.DeletedAt)

		if err != nil {
			fmt.Println(err.Error())
			continue
		}

		queries = append(queries, query)
	}

	err = rows.Err()

	if err != nil {
		log.Fatal(err)
	}

	return queries, nil
}

func (a *App) DeleteQuery(id int) error {
	db := openSqliteConnection()

	defer db.Close()

	deleteQuery := `UPDATE queries SET deleted_at = CURRENT_TIMESTAMP WHERE id = ?`

	stmt, err := db.Prepare(deleteQuery)

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (a *App) UpdateQuery(id int, data Query) error {
	db := openSqliteConnection()

	defer db.Close()

	updateQuery := `UPDATE queries SET query = ?, description = ? WHERE id = ?`

	stmt, err := db.Prepare(updateQuery)

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec(data.Query, data.Description, id)

	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (a *App) GetDatabaseConnection() (DatabaseConnection, error) {
	db := openSqliteConnection()

	defer db.Close()

	databaseConnectionQuery := `SELECT * FROM database_connections LIMIT 1`

	rows, err := db.Query(databaseConnectionQuery)

	if err != nil {
		return DatabaseConnection{}, err
	}

	defer rows.Close()

	var databaseConnection DatabaseConnection

	for rows.Next() {
		err = rows.Scan(&databaseConnection.ID, &databaseConnection.Username, &databaseConnection.Password, &databaseConnection.Host, &databaseConnection.Port, &databaseConnection.Database, &databaseConnection.CreatedAt, &databaseConnection.UpdatedAt, &databaseConnection.DeletedAt)

		if err != nil {
			return DatabaseConnection{}, err
		}
	}

	err = rows.Err()

	if err != nil {
		return DatabaseConnection{}, err
	}

	return databaseConnection, nil
}

func (a *App) CreateOrUpdateDatabaseConnection(input DatabaseConnection) (DatabaseConnection, error) {
	db := openSqliteConnection()

	defer db.Close()

	databaseConnectionQuery := `SELECT * FROM database_connections LIMIT 1`

	rows, err := db.Query(databaseConnectionQuery)

	if err != nil {
		return DatabaseConnection{}, err
	}

	defer rows.Close()

	var databaseConnection DatabaseConnection

	for rows.Next() {
		err = rows.Scan(&databaseConnection.ID, &databaseConnection.Username, &databaseConnection.Password, &databaseConnection.Host, &databaseConnection.Port, &databaseConnection.Database, &databaseConnection.CreatedAt, &databaseConnection.UpdatedAt, &databaseConnection.DeletedAt)

		if err != nil {
			return DatabaseConnection{}, err
		}
	}

	err = rows.Err()

	if err != nil {
		return DatabaseConnection{}, err
	}

	if input.ID != nil {
		updateQuery := `UPDATE database_connections SET username = ?, password = ?, database = ?, host = ?, port = ? WHERE id = ?`

		stmt, err := db.Prepare(updateQuery)

		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = stmt.Exec(input.Username, input.Password, input.Database, input.Host, input.Port, *input.ID)

		if err != nil {
			fmt.Println(err.Error())
		}

	} else {
		insertQuery := `INSERT INTO database_connections(username, password, database, host, port) VALUES(?, ?, ?, ?, ?)`

		stmt, err := db.Prepare(insertQuery)

		if err != nil {
			fmt.Println(err.Error())
		}

		_, err = stmt.Exec(input.Username, input.Password, input.Database, input.Host, input.Port)

		if err != nil {
			fmt.Println(err.Error())
		}
	}

	databaseConnection.Username = input.Username
	databaseConnection.Password = input.Password
	databaseConnection.Host = input.Host
	databaseConnection.Port = input.Port
	databaseConnection.Database = input.Database

	return databaseConnection, nil
}

func (a *App) ImportDatabaseFile() error {
	selection, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Select File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Sqlite 3 (*.db)",
				Pattern:     "*.db",
			},
		},
	})

	if err != nil {
		return err
	}

	if selection == "" {
		return fmt.Errorf("no file selected")
	}

	_, err = os.Stat("database.db")

	if err == nil {
		err = os.Remove("database.db")

		if err != nil {
			return err
		}
	}

	err = os.Rename(selection, "database.db")

	if err != nil {
		return err
	}

	return nil
}

func (a *App) ExportDatabaseFile() error {
	selection, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title: "Save File",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Sqlite 3 (*.db)",
				Pattern:     "*.db",
			},
		},
	})

	if err != nil {
		return err
	}

	if selection == "" {
		return fmt.Errorf("no file selected")
	}

	_, err = os.Stat("database.db")

	if err != nil {
		return err
	}

	err = os.Rename("database.db", selection)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) TestQueryInDatabase(input DatabaseConnection, query string, useTransaction bool) ([]map[string]interface{}, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", input.Username, input.Password, input.Host, input.Port, input.Database))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	// Begin transaction if requested
	var tx *sql.Tx
	if useTransaction {
		tx, err = db.Begin()
		if err != nil {
			return nil, fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback() // Always rollback to ensure no changes are committed
	}

	// Execute query based on whether we're in a transaction
	var rows *sql.Rows
	if useTransaction {
		rows, err = tx.Query(query)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result []map[string]interface{}
	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range columns {
			valuePtrs[i] = &values[i]
		}

		err = rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		row := make(map[string]interface{})
		for i, col := range columns {
			val := values[i]
			switch v := val.(type) {
			case []byte:
				row[col] = string(v)
			default:
				row[col] = v
			}
		}

		result = append(result, row)
	}

	return result, nil
}

func (a *App) TestBatchQueryInDatabase(input DatabaseConnection, queries []string, useTransaction bool) ([][]map[string]interface{}, error) {
	if len(queries) == 0 {
		return [][]map[string]interface{}{}, nil
	}

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", input.Username, input.Password, input.Host, input.Port, input.Database))
	if err != nil {
		return nil, err
	}
	defer db.Close()

	var results [][]map[string]interface{}

	// Begin transaction if requested
	if useTransaction {
		tx, err := db.Begin()
		if err != nil {
			return nil, fmt.Errorf("failed to start transaction: %w", err)
		}
		defer tx.Rollback() // Always rollback to ensure no changes are committed

		// Execute each query within the same transaction
		for _, query := range queries {
			rows, err := tx.Query(query)
			if err != nil {
				return nil, err
			}

			columns, err := rows.Columns()
			if err != nil {
				rows.Close()
				return nil, err
			}

			var result []map[string]interface{}
			for rows.Next() {
				values := make([]interface{}, len(columns))
				valuePtrs := make([]interface{}, len(columns))

				for i := range columns {
					valuePtrs[i] = &values[i]
				}

				err = rows.Scan(valuePtrs...)
				if err != nil {
					rows.Close()
					return nil, err
				}

				row := make(map[string]interface{})
				for i, col := range columns {
					val := values[i]
					switch v := val.(type) {
					case []byte:
						row[col] = string(v)
					default:
						row[col] = v
					}
				}

				result = append(result, row)
			}
			rows.Close()
			results = append(results, result)
		}
	} else {
		// Run each query individually without transaction
		for _, query := range queries {
			rows, err := a.TestQueryInDatabase(input, query, false)
			if err != nil {
				return nil, err
			}
			results = append(results, rows)
		}
	}

	return results, nil
}

func (a *App) TestDatabaseConnection(input DatabaseConnection) bool {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", input.Username, input.Password, input.Host, input.Port, input.Database))

	if err != nil {
		return false
	}

	defer db.Close()

	err = db.Ping()

	if err != nil {
		return false
	}

	return true
}

func openSqliteConnection() *sql.DB {
	db, err := sql.Open("sqlite3", "./database.db")

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func createSqliteTables() {
	db := openSqliteConnection()

	defer db.Close()

	createTableSQL := `
		CREATE TABLE IF NOT EXISTS queries (
			id integer NOT NULL PRIMARY KEY,
			title TEXT,
			query TEXT,
			description TEXT DEFAULT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP DEFAULT NULL
		);

		CREATE TABLE IF NOT EXISTS database_connections (
			id integer NOT NULL PRIMARY KEY,
			username TEXT,
			password TEXT,
			host TEXT,
			port INTEGER,
			database TEXT,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			deleted_at TIMESTAMP DEFAULT NULL
		);

		CREATE TABLE IF NOT EXISTS database_structure (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			structure TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`
	_, err := db.Exec(createTableSQL)

	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) GetBuildParams() map[string]interface{} {
	buildParams := make(map[string]interface{})

	buildParams["version"] = "0.0.2"

	return buildParams
}

func (a *App) CheckHasUpdate() bool {
	buildParams := a.GetBuildParams()

	if buildParams["version"] == nil {
		return false
	}

	version := buildParams["version"].(string)

	resp, err := http.Get("https://raw.githubusercontent.com/LouisRiverstone/query-script-maker/refs/heads/master/build_params.json")

	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var versionData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&versionData)
	if err != nil {
		return false
	}

	latestVersion := versionData["version"].(string)

	return latestVersion > version
}

func (a *App) GetDatabaseStructure(input DatabaseConnection) (string, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", input.Username, input.Password, input.Host, input.Port, input.Database))

	if err != nil {
		return "", err
	}

	defer db.Close()

	// Get all tables
	tableRows, err := db.Query("SHOW TABLES")
	if err != nil {
		return "", err
	}
	defer tableRows.Close()

	type Column struct {
		Name      string `json:"name"`
		Type      string `json:"type"`
		Nullable  string `json:"nullable"`
		Key       string `json:"key"`
		Default   string `json:"default"`
		Extra     string `json:"extra"`
		IsPrimary bool   `json:"isPrimary"`
	}

	type ForeignKey struct {
		ColumnName       string `json:"columnName"`
		ReferencedTable  string `json:"referencedTable"`
		ReferencedColumn string `json:"referencedColumn"`
		ConstraintName   string `json:"constraintName"`
	}

	type Table struct {
		Name        string       `json:"name"`
		Columns     []Column     `json:"columns"`
		ForeignKeys []ForeignKey `json:"foreignKeys"`
	}

	type DatabaseStructure struct {
		Tables []Table `json:"tables"`
	}

	structure := DatabaseStructure{
		Tables: []Table{},
	}

	// Collect all tables
	var tables []string
	for tableRows.Next() {
		var tableName string
		if err := tableRows.Scan(&tableName); err != nil {
			return "", err
		}
		tables = append(tables, tableName)
	}

	// For each table, get columns and foreign keys
	for _, tableName := range tables {
		table := Table{
			Name:        tableName,
			Columns:     []Column{},
			ForeignKeys: []ForeignKey{},
		}

		// Get columns
		columnRows, err := db.Query(fmt.Sprintf("DESCRIBE `%s`", tableName))
		if err != nil {
			return "", err
		}

		for columnRows.Next() {
			var col Column
			var field, colType, null, key, defaultVal, extra sql.NullString
			if err := columnRows.Scan(&field, &colType, &null, &key, &defaultVal, &extra); err != nil {
				columnRows.Close()
				return "", err
			}

			col.Name = field.String
			col.Type = colType.String
			col.Nullable = null.String
			col.Key = key.String
			col.Default = defaultVal.String
			col.Extra = extra.String
			col.IsPrimary = key.String == "PRI"

			table.Columns = append(table.Columns, col)
		}
		columnRows.Close()

		// Get foreign keys
		fkQuery := `
		SELECT 
			COLUMN_NAME, 
			REFERENCED_TABLE_NAME, 
			REFERENCED_COLUMN_NAME,
			CONSTRAINT_NAME
		FROM 
			INFORMATION_SCHEMA.KEY_COLUMN_USAGE 
		WHERE 
			TABLE_SCHEMA = ? 
			AND TABLE_NAME = ? 
			AND REFERENCED_TABLE_NAME IS NOT NULL
		`
		fkRows, err := db.Query(fkQuery, input.Database, tableName)
		if err != nil {
			return "", err
		}

		for fkRows.Next() {
			var fk ForeignKey
			if err := fkRows.Scan(&fk.ColumnName, &fk.ReferencedTable, &fk.ReferencedColumn, &fk.ConstraintName); err != nil {
				fkRows.Close()
				return "", err
			}
			table.ForeignKeys = append(table.ForeignKeys, fk)
		}
		fkRows.Close()

		structure.Tables = append(structure.Tables, table)
	}

	// Convert to JSON
	structureJSON, err := json.Marshal(structure)
	if err != nil {
		return "", err
	}

	// Store in SQLite
	sqliteDB := openSqliteConnection()
	defer sqliteDB.Close()

	// Create table if not exists
	_, err = sqliteDB.Exec(`
		CREATE TABLE IF NOT EXISTS database_structure (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			structure TEXT NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return "", err
	}

	// Insert new structure
	_, err = sqliteDB.Exec("INSERT INTO database_structure (structure) VALUES (?)", string(structureJSON))
	if err != nil {
		return "", err
	}

	return string(structureJSON), nil
}

func (a *App) GetLatestDatabaseStructure() (string, error) {
	db := openSqliteConnection()
	defer db.Close()

	var structure string
	err := db.QueryRow("SELECT structure FROM database_structure ORDER BY created_at DESC LIMIT 1").Scan(&structure)
	if err != nil {
		if err == sql.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return structure, nil
}

// GenerateSQLFromPrompt is the app method that interfaces with the SQL assistant
func (a *App) GenerateSQLFromPrompt(prompt string) (string, error) {
	// Get the database structure from SQLite
	structureJSON, err := a.GetLatestDatabaseStructure()

	if err != nil {
		return "", fmt.Errorf("error getting database structure: %w", err)
	}

	if structureJSON == "" {
		return "", fmt.Errorf("database structure not found, please scan your database first")
	}

	// Initialize the assistant with the database structure
	assistant := a.getSQLAssistant(structureJSON)

	// Generate SQL from the prompt
	sql, err := assistant.GenerateSQL(prompt)
	if err != nil {
		return "", fmt.Errorf("error generating SQL: %w", err)
	}

	return sql, nil
}

// ResetSQLAssistant resets the assistant's state
func (a *App) ResetSQLAssistant() {
	assistant := a.getSQLAssistant("")

	assistant.Reset()
}

// RecordQueryFeedback records query execution feedback for learning
func (a *App) RecordQueryFeedback(query string, successful bool, errorMsg string, rowCount int, execTime float64) {
	assistant := a.getSQLAssistant("")

	assistant.RecordQueryFeedback(sqlaiModels.FeedbackResult{
		Query:         query,
		WasSuccessful: successful,
		ErrorMessage:  errorMsg,
		RowCount:      rowCount,
		ExecutionTime: execTime,
		CreatedAt:     time.Now(),
	})
}
