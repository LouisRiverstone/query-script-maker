/**
 * A simplified SQL parser for visualization purposes
 * Optimized for use with @vue-flow/core
 */

export interface SQLTable {
  id: string; // Unique ID for vue-flow
  name: string;
  alias?: string;
  columns: SQLColumn[];
  position?: { x: number; y: number }; // Position for vue-flow
}

export interface SQLColumn {
  id: string; // Unique ID for column
  name: string;
  alias?: string;
  table?: string;
  isPrimaryKey?: boolean;
  isForeignKey?: boolean;
  isSelected?: boolean; // Whether this column is part of the SELECT clause
  value?: string; // Value for INSERT queries
}

export interface SQLJoin {
  id: string; // Unique ID for edge in vue-flow
  table: SQLTable;
  type: 'INNER' | 'LEFT' | 'RIGHT' | 'FULL' | 'CROSS';
  condition?: {
    leftTable: string;
    leftColumn: string;
    rightTable: string;
    rightColumn: string;
    sourceId?: string; // IDs for vue-flow edge
    targetId?: string;
  };
}

export interface SQLParseResult {
  mainTable?: SQLTable;
  tables: SQLTable[];
  joins: SQLJoin[];
  columns: SQLColumn[];
  selectedColumns: SQLColumn[]; // Columns that appear in the SELECT statement
  error?: string;
  queryType: 'SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'OTHER';
  insertValues?: { [column: string]: string }; // Values for INSERT queries
  // Vue Flow specific data
  nodes: any[];
  edges: any[];
}

/**
 * Parses SQL query to extract table and join information for visualization with Vue Flow
 * @param sql The SQL query to parse
 * @returns Parsed SQL information for diagram visualization
 */
export function parse(sql: string): SQLParseResult {
  try {
    // Initialize result
    const result: SQLParseResult = {
      tables: [],
      joins: [],
      columns: [],
      selectedColumns: [],
      queryType: 'OTHER',
      nodes: [],
      edges: []
    };

    if (!sql || !sql.trim()) {
      result.error = "Empty SQL query";
      return result;
    }

    // Normalize SQL for easier parsing
    const normalizedSQL = sql
      .replace(/\s+/g, ' ')
      .replace(/\n/g, ' ')
      .replace(/,/g, ' , ')
      .trim();
    
    // Determine query type
    if (/^\s*SELECT\b/i.test(normalizedSQL)) {
      result.queryType = 'SELECT';
      return parseSelectQuery(normalizedSQL, result);
    } else if (/^\s*INSERT\b/i.test(normalizedSQL)) {
      result.queryType = 'INSERT';
      return parseInsertQuery(normalizedSQL, result);
    } else if (/^\s*UPDATE\b/i.test(normalizedSQL)) {
      result.queryType = 'UPDATE';
      // Handle UPDATE queries in future implementation
    } else if (/^\s*DELETE\b/i.test(normalizedSQL)) {
      result.queryType = 'DELETE';
      // Handle DELETE queries in future implementation
    }

    return result;
  } catch (error) {
    console.error('Error parsing SQL:', error);
    return {
      tables: [],
      joins: [],
      columns: [],
      selectedColumns: [],
      error: error instanceof Error ? error.message : String(error),
      queryType: 'OTHER',
      nodes: [],
      edges: []
    };
  }
}

/**
 * Parses a SELECT query
 */
function parseSelectQuery(sql: string, result: SQLParseResult): SQLParseResult {
  let columnsMap: Map<string, SQLColumn[]> = new Map();
  let selectedColumns: SQLColumn[] = [];
  
  // Extract columns from SELECT clause
  const selectMatch = /SELECT\s+(.*?)\s+FROM/i.exec(sql);
  if (selectMatch) {
    const columnsSection = selectMatch[1].trim();
    
    // Handle SELECT * special case
    if (columnsSection === '*') {
      // Will be handled after tables are parsed
    } else {
      const columnMatches = columnsSection.split(/\s*,\s*/);
      
      columnMatches.forEach((colStr, colIndex) => {
        // Handle "table.column as alias" pattern
        const colMatch = /(?:([A-Za-z0-9_]+)\.)?([A-Za-z0-9_*]+)(?:\s+[Aa][Ss]\s+([A-Za-z0-9_]+))?/.exec(colStr.trim());
        if (colMatch) {
          const [_, tableName, colName, colAlias] = colMatch;
          
          const column: SQLColumn = {
            id: `col-${colIndex}-${Date.now()}`, // Unique ID for vue-flow
            name: colName,
            alias: colAlias,
            table: tableName,
            isSelected: true
          };
          
          result.columns.push(column);
          selectedColumns.push(column);
          
          // Group columns by table
          if (tableName) {
            if (!columnsMap.has(tableName)) {
              columnsMap.set(tableName, []);
            }
            columnsMap.get(tableName)?.push(column);
          }
        }
      });
    }
  }
  
  // Extract main table (FROM clause)
  const fromMatch = /\s+FROM\s+([^\s,;]+)(?:\s+(?:AS\s+)?([^\s,;]+))?/i.exec(sql);
  if (fromMatch) {
    const tableName = fromMatch[1];
    const tableAlias = fromMatch[2] || tableName;
    
    const mainTable: SQLTable = {
      id: `table-${tableName}-${Date.now()}`, // Unique ID for vue-flow
      name: tableName,
      alias: tableAlias,
      columns: columnsMap.get(tableName) || columnsMap.get(tableAlias) || [],
      position: { x: 50, y: 50 } // Initial position
    };
    
    result.mainTable = mainTable;
    result.tables.push(mainTable);
  } else {
    result.error = "Could not find FROM clause";
    return result;
  }

  // Extract JOIN clauses
  const joinRegex = /\b(INNER|LEFT|RIGHT|FULL|CROSS)?\s*JOIN\s+([^\s]+)(?:\s+(?:AS\s+)?([^\s]+))?\s+ON\s+([^\s.]+)\.([^\s=]+)\s*=\s*([^\s.]+)\.([^\s\s]+)/gi;
  let joinMatch;
  let joinCount = 0;
  
  while ((joinMatch = joinRegex.exec(sql)) !== null) {
    const [_, joinType = 'INNER', tableName, tableAlias, leftTable, leftColumn, rightTable, rightColumn] = joinMatch;
    const alias = tableAlias || tableName;
    
    // Calculate position based on existing tables
    const xOffset = 400 + (result.tables.length % 2) * 100;
    const yOffset = 50 + Math.floor(result.tables.length / 2) * 250;
    
    const joinTable: SQLTable = {
      id: `table-${tableName}-${Date.now()}`, // Unique ID for vue-flow
      name: tableName,
      alias: alias,
      columns: columnsMap.get(tableName) || columnsMap.get(alias) || [],
      position: { x: xOffset, y: yOffset } // Position for vue-flow
    };
    
    if (!result.tables.some(t => t.name === joinTable.name)) {
      result.tables.push(joinTable);
    }
    
    // Mark foreign key columns
    const condition = {
      leftTable,
      leftColumn,
      rightTable,
      rightColumn
    };
    
    // Add IDs for edge connection
    let sourceTableId = '';
    let targetTableId = '';
    
    // Mark foreign key relationships
    result.tables.forEach(table => {
      if (table.name === leftTable || table.alias === leftTable) {
        // Check if column already exists, if not, add it
        let fkColumn = table.columns.find(c => c.name === leftColumn);
        if (!fkColumn) {
          fkColumn = { id: `col-${leftColumn}-${Date.now()}`, name: leftColumn, isForeignKey: true };
          table.columns.push(fkColumn);
          result.columns.push({ ...fkColumn, table: table.name });
        } else {
          fkColumn.isForeignKey = true;
        }
        sourceTableId = table.id;
      }
      if (table.name === rightTable || table.alias === rightTable) {
        // Check if column already exists, if not, add it
        let fkColumn = table.columns.find(c => c.name === rightColumn);
        if (!fkColumn) {
          fkColumn = { id: `col-${rightColumn}-${Date.now()}`, name: rightColumn, isForeignKey: true };
          table.columns.push(fkColumn);
          result.columns.push({ ...fkColumn, table: table.name });
        } else {
          fkColumn.isForeignKey = true;
        }
        targetTableId = table.id;
      }
    });
    
    const joinEdge = {
      id: `join-${joinCount++}-${Date.now()}`,
      table: joinTable,
      type: joinType.toUpperCase() as any,
      condition: {
        ...condition,
        sourceId: sourceTableId,
        targetId: targetTableId
      }
    };
    
    result.joins.push(joinEdge);
  }
  
  // Look for additional columns in WHERE clause
  const whereMatch = /WHERE\s+(.*?)(?:ORDER BY|GROUP BY|HAVING|LIMIT|$)/i.exec(sql);
  if (whereMatch) {
    const whereClause = whereMatch[1];
    const columnMatches = whereClause.matchAll(/([A-Za-z0-9_]+)\.([A-Za-z0-9_]+)/g);
    
    for (const match of columnMatches) {
      const [_, tableName, colName] = match;
      
      // Find the table
      const table = result.tables.find(t => 
        t.name === tableName || t.alias === tableName
      );
      
      if (table && !table.columns.some(c => c.name === colName)) {
        const column = { id: `col-${colName}-${Date.now()}`, name: colName };
        table.columns.push(column);
        result.columns.push({ ...column, table: tableName });
      }
    }
  }
  
  // Handle SELECT * by adding all columns from all tables
  if (selectMatch && selectMatch[1].trim() === '*') {
    result.tables.forEach(table => {
      table.columns.forEach(col => {
        col.isSelected = true;
        selectedColumns.push({
          ...col,
          table: table.name
        });
      });
    });
  }
  
  // If we haven't found any columns in any table, add default ones
  result.tables.forEach((table, tableIndex) => {
    if (table.columns.length === 0) {
      table.columns = [
        { id: `col-id-${tableIndex}-${Date.now()}`, name: 'id', isPrimaryKey: true },
        { id: `col-name-${tableIndex}-${Date.now()}`, name: 'name' },
        { id: `col-created-${tableIndex}-${Date.now()}`, name: 'created_at' }
      ];
    }
  });
  
  // Try to detect primary keys from column names or JOIN conditions
  result.tables.forEach(table => {
    const idColumn = table.columns.find(c => 
      c.name.toLowerCase() === 'id' || 
      c.name.toLowerCase() === `${table.name.toLowerCase()}_id`
    );
    
    if (idColumn) {
      idColumn.isPrimaryKey = true;
    }
  });
  
  result.selectedColumns = selectedColumns;
  
  // Create Vue Flow nodes and edges
  createVueFlowElements(result);

  return result;
}

/**
 * Parses an INSERT query
 */
function parseInsertQuery(sql: string, result: SQLParseResult): SQLParseResult {
  // Extract table name from INSERT INTO clause
  const insertMatch = /INSERT\s+INTO\s+([^\s(]+)(?:\s+\(([^)]+)\))?\s+VALUES\s+\(([^)]+)\)/i.exec(sql);
  
  if (!insertMatch) {
    result.error = "Invalid INSERT query format";
    return result;
  }
  
  const [_, tableName, columnsStr, valuesStr] = insertMatch;
  
  // Create the table object
  const table: SQLTable = {
    id: `table-${tableName}-${Date.now()}`,
    name: tableName,
    columns: [],
    position: { x: 50, y: 50 } // Initial position
  };
  
  result.mainTable = table;
  result.tables.push(table);
  
  // Initialize the insert values object
  result.insertValues = {};
  
  // Parse columns
  let columns: string[] = [];
  if (columnsStr) {
    // Columns explicitly specified
    columns = columnsStr.split(',').map(col => col.trim());
  }
  
  // Parse values - handle string literals, numbers, etc.
  const values: string[] = [];
  let currentValue = '';
  let inString = false;
  
  for (let i = 0; i < valuesStr.length; i++) {
    const char = valuesStr[i];
    
    if (char === "'" && (i === 0 || valuesStr[i-1] !== '\\')) {
      inString = !inString;
      currentValue += char;
    } else if (char === ',' && !inString) {
      values.push(currentValue.trim());
      currentValue = '';
    } else {
      currentValue += char;
    }
  }
  
  // Add the last value
  if (currentValue.trim()) {
    values.push(currentValue.trim());
  }
  
  // Map columns to values
  for (let i = 0; i < columns.length && i < values.length; i++) {
    const column: SQLColumn = {
      id: `col-${columns[i]}-${Date.now()}`,
      name: columns[i],
      value: values[i],
      isSelected: true
    };
    
    table.columns.push(column);
    result.selectedColumns.push(column);
    result.columns.push(column);
    
    // Store in insertValues map
    result.insertValues[columns[i]] = values[i];
  }
  
  // If no columns were specified, create columns based on values positions
  if (columns.length === 0) {
    values.forEach((value, index) => {
      const columnName = `column${index + 1}`;
      const column: SQLColumn = {
        id: `col-${columnName}-${Date.now()}`,
        name: columnName,
        value: value,
        isSelected: true
      };
      
      table.columns.push(column);
      result.selectedColumns.push(column);
      result.columns.push(column);
      
      // Store in insertValues map
      if (result.insertValues) {
        result.insertValues[columnName] = value;
      }
    });
  }
  
  // Try to detect primary keys from column names
  table.columns.forEach(column => {
    if (column.name.toLowerCase() === 'id' || 
        column.name.toLowerCase() === `${table.name.toLowerCase()}_id`) {
      column.isPrimaryKey = true;
    }
  });
  
  // For INSERT queries, create a "values" node
  const valuesNode = {
    id: `values-${Date.now()}`,
    type: 'valuesNode',
    data: { values: result.insertValues },
    position: { x: 50, y: 250 }
  };
  
  result.nodes.push(valuesNode);
  
  // Create Vue Flow nodes and edges
  createVueFlowElements(result);
  
  return result;
}

/**
 * Creates Vue Flow nodes and edges from the parsed SQL data
 */
function createVueFlowElements(result: SQLParseResult): void {
  // Create nodes for tables
  result.tables.forEach((table, index) => {
    // Ensure position is set
    if (!table.position) {
      table.position = {
        x: 50 + (index % 3) * 350,
        y: 50 + Math.floor(index / 3) * 300
      };
    }
    
    // Create table node
    result.nodes.push({
      id: table.id,
      type: 'tableNode',
      data: {
        label: table.name,
        alias: table.alias,
        columns: table.columns,
        isMainTable: result.mainTable?.id === table.id,
        queryType: result.queryType
      },
      position: table.position
    });
  });
  
  // Create edges for joins
  result.joins.forEach(join => {
    if (join.condition && join.condition.sourceId && join.condition.targetId) {
      result.edges.push({
        id: join.id,
        source: join.condition.sourceId,
        target: join.condition.targetId,
        type: 'joinEdge',
        data: {
          joinType: join.type,
          sourceColumn: join.condition.leftColumn,
          targetColumn: join.condition.rightColumn
        },
        animated: true,
        markerEnd: {
          type: 'arrowclosed',
          width: 20,
          height: 20
        }
      });
    }
  });
  
  // Create result node for SELECT queries
  if (result.queryType === 'SELECT' && result.selectedColumns.length > 0) {
    const resultNode = {
      id: `result-${Date.now()}`,
      type: 'resultNode',
      data: {
        columns: result.selectedColumns
      },
      position: {
        x: 50 + (result.tables.length % 3) * 350,
        y: Math.max(...result.tables.map(t => (t.position?.y || 0) + 300)) + 50
      }
    };
    
    result.nodes.push(resultNode);
    
    // Connect selected columns to result node
    result.selectedColumns.forEach((col, index) => {
      if (col.table) {
        const sourceTable = result.tables.find(t => 
          t.name === col.table || t.alias === col.table
        );
        
        if (sourceTable) {
          result.edges.push({
            id: `select-edge-${index}-${Date.now()}`,
            source: sourceTable.id,
            target: resultNode.id,
            type: 'selectEdge',
            data: {
              column: col.name,
              alias: col.alias
            },
            style: { stroke: '#10b981' },
            animated: true
          });
        }
      }
    });
  }
  
  // Create edges for INSERT queries
  if (result.queryType === 'INSERT' && result.nodes.find(n => n.type === 'valuesNode')) {
    const valuesNode = result.nodes.find(n => n.type === 'valuesNode');
    const tableNode = result.nodes.find(n => n.type === 'tableNode');
    
    if (valuesNode && tableNode) {
      result.edges.push({
        id: `insert-edge-${Date.now()}`,
        source: valuesNode.id,
        target: tableNode.id,
        type: 'insertEdge',
        data: {
          values: result.insertValues
        },
        style: { stroke: '#10b981' },
        animated: true,
        markerEnd: {
          type: 'arrowclosed',
          width: 20,
          height: 20
        }
      });
    }
  }
} 