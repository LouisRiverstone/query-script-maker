package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/xuri/excelize/v2"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) ReadFile() (string, error) {
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

// create sql file
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

func (a *App) CreateSQLiteDatase(data string) error {
	if _, err := os.Stat("database.db"); os.IsNotExist(err) {
		//create the file
		file, err := os.Create("database.db")

		if err != nil {
			return err
		}

		file.Close()

		sqliteDatabase, _ := sql.Open("sqlite3", "./database.db")
		defer sqliteDatabase.Close() // Defer Closing the database
		createTable(sqliteDatabase)
	}

	return nil
}

func createTable(db *sql.DB) {
	// Create Table
	createTableSQL := `CREATE TABLE IF NOT EXISTS queries (
		"id" integer NOT NULL PRIMARY KEY,
		"title" TEXT,
		"query" TEXT
	);`

	_, err := db.Exec(createTableSQL)

	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("employees Table is created successfully")
	}
}

func (a *App) InsertQueryInDatabase(data string) error {
	sqliteDatabase, _ := sql.Open("sqlite3", "./database.db")

	defer sqliteDatabase.Close() // Defer Closing the database

	insertQuery := `INSERT INTO queries(title, query) VALUES(?, ?)`

	stmt, err := sqliteDatabase.Prepare(insertQuery)

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec("Test", data)

	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (a *App) GetQueries() (string, error) {
	sqliteDatabase, _ := sql.Open("sqlite3", "./database.db")

	defer sqliteDatabase.Close() // Defer Closing the database

	rows, err := sqliteDatabase.Query("SELECT * FROM queries")

	if err != nil {
		fmt.Println(err.Error())
	}

	var content []map[string]string

	for rows.Next() {
		var id int
		var title string
		var query string

		rows.Scan(&id, &title, &query)

		rowData := make(map[string]string)
		rowData["id"] = fmt.Sprintf("%d", id)
		rowData["title"] = title
		rowData["query"] = query

		content = append(content, rowData)
	}

	jsonData, err := json.Marshal(content)

	if err != nil {
		return "", err
	}

	return string(jsonData), nil
}

func (a *App) DeleteQuery(id string) error {
	sqliteDatabase, _ := sql.Open("sqlite3", "./database.db")

	defer sqliteDatabase.Close() // Defer Closing the database

	deleteQuery := `DELETE FROM queries WHERE id = ?`

	stmt, err := sqliteDatabase.Prepare(deleteQuery)

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec(id)

	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
}

func (a *App) UpdateQuery(id string, data string) error {
	sqliteDatabase, _ := sql.Open("sqlite3", "./database.db")

	defer sqliteDatabase.Close() // Defer Closing the database

	updateQuery := `UPDATE queries SET query = ? WHERE id = ?`

	stmt, err := sqliteDatabase.Prepare(updateQuery)

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = stmt.Exec(data, id)

	if err != nil {
		fmt.Println(err.Error())
	}

	return nil
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
