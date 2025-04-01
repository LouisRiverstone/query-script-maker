declare module 'sql-parser-mistic' {
  /**
   * Parses SQL query into an abstract syntax tree
   * @param sql The SQL query to parse
   * @returns The parsed AST representation of the SQL query
   */
  export function parse(sql: string): any;
  
  // Add other exports from the package if needed
} 