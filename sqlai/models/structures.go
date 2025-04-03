package models

import (
	"time"
)

// DatabaseStructureForAI is a copy of the database structure for AI usage
type DatabaseStructureForAI struct {
	Tables      []TableForAI      `json:"tables"`
	QueryCache  map[string]string `json:"-"`
	LastUpdated time.Time         `json:"lastUpdated"`
	DBType      string            `json:"dbType"`
}

// TableForAI represents a database table for AI usage
type TableForAI struct {
	Name             string            `json:"name"`
	Columns          []ColumnForAI     `json:"columns"`
	ForeignKeys      []ForeignKeyForAI `json:"foreignKeys"`
	Description      string            `json:"description"`
	EstimatedRows    int               `json:"estimatedRows"`
	CommonQueryTypes []string          `json:"commonQueryTypes"`
}

// ColumnForAI represents a database column for AI usage
type ColumnForAI struct {
	Name         string   `json:"name"`
	Type         string   `json:"type"`
	Nullable     string   `json:"nullable"`
	Key          string   `json:"key"`
	Default      string   `json:"default"`
	Extra        string   `json:"extra"`
	IsPrimary    bool     `json:"isPrimary"`
	IsUnique     bool     `json:"isUnique"`
	Description  string   `json:"description"`
	HasIndex     bool     `json:"hasIndex"`
	CommonValues []string `json:"commonValues"`
	SampleValues []string `json:"sampleValues"`
	MinValue     string   `json:"minValue"`
	MaxValue     string   `json:"maxValue"`
}

// ForeignKeyForAI represents a foreign key relationship for AI usage
type ForeignKeyForAI struct {
	ColumnName       string `json:"columnName"`
	ReferencedTable  string `json:"referencedTable"`
	ReferencedColumn string `json:"referencedColumn"`
	ConstraintName   string `json:"constraintName"`
	RelationshipType string `json:"relationshipType"`
	CascadeDelete    bool   `json:"cascadeDelete"`
	CascadeUpdate    bool   `json:"cascadeUpdate"`
}

// SQLDialect contains dialect-specific SQL implementations
type SQLDialect struct {
	Name                 string
	LimitSyntax          string
	DateFunctions        map[string]string
	StringFunctions      map[string]string
	AggregationFunctions map[string]string
	Pagination           string
	SupportsCTE          bool
	SupportsWindowFuncs  bool
}

// QueryHistory tracks past queries for learning
type QueryHistory struct {
	Query         string
	ResultCount   int
	ExecutionTime float64
	CreatedAt     time.Time
	Success       bool
}

// TableInfo represents a table with its confidence score
type TableInfo struct {
	Name       string
	Confidence float64
	Alias      string
}

// ColumnInfo represents a column with metadata and confidence score
type ColumnInfo struct {
	Name       string
	TableName  string
	Type       string
	IsPrimary  bool
	Confidence float64
	Function   string
}

// Condition represents a WHERE condition in SQL
type Condition struct {
	ColumnName  string
	TableName   string
	Operator    string
	Value       string
	Conjunction string
	IsComplex   bool
	ComplexExpr string
}

// FeedbackResult allows tracking the quality of generated queries
type FeedbackResult struct {
	Query         string
	WasSuccessful bool
	ErrorMessage  string
	RowCount      int
	ExecutionTime float64
	CreatedAt     time.Time
}
