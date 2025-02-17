package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

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

func (a *App) ReadXLSXFile() (string, error) {
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

	file, err := os.Open(selection)

	if err != nil {
		return "", err
	}

	defer file.Close()

	xlsxFile, err := excelize.OpenReader(file)

	if err != nil {
		return "", err
	}

	rows, err := xlsxFile.GetRows(xlsxFile.GetSheetName(0))

	if err != nil {
		return "", err
	}

	var content []map[string]string

	headers := rows[0]

	for _, row := range rows[1:] {
		rowData := make(map[string]string)

		for i, cell := range row {
			rowData[headers[i]] = cell
		}

		content = append(content, rowData)
	}

	jsonData, err := json.Marshal(content)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (a *App) MakeBindedSQL(query string, data []map[string]interface{}, variables []Variable) string {
	var result strings.Builder

	re := regexp.MustCompile(`{{\w+}}`)

	for index, row := range data {
		modifiedQuery := re.ReplaceAllStringFunc(query, func(match string) string {
			for _, variable := range variables {
				if match == fmt.Sprintf("{{%s}}", variable.Value) {
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

func (a *App) TestQueryInDatabase(input DatabaseConnection, query string) ([]map[string]interface{}, error) {
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", input.Username, input.Password, input.Host, input.Port, input.Database))

	if err != nil {
		return nil, err
	}

	defer db.Close()

	rows, err := db.Query(query)

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
	`
	_, err := db.Exec(createTableSQL)

	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) GetBuildParams() map[string]interface{} {
	buildParams := make(map[string]interface{})

	buildParams["version"] = "0.0.1"

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
