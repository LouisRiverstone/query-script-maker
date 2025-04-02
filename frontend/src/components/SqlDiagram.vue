<template>
  <div class="sql-diagram-container">
    <div class="diagram-header">
      <h3 class="text-lg font-semibold text-black dark:text-white">
        SQL Diagram <span v-if="queryType" class="ml-2 text-sm font-normal">({{ queryType }} query)</span>
      </h3>
      <div class="diagram-controls">
        <Button type="button" @click="zoomIn" class="mr-1">+</Button>
        <Button type="button" @click="zoomOut" class="mr-1">-</Button>
        <Button type="button" @click="fitView">Fit</Button>
      </div>
    </div>
    <div ref="diagramContainer" class="diagram-content bg-white dark:bg-gray-800 rounded-lg p-4">
      <VueFlow v-if="nodes.length > 0"
        :default-zoom="1"
        :min-zoom="0.2"
        :max-zoom="4"
        :nodes="nodes"
        :edges="edges"
        class="vue-flow-wrapper"
        fit-view-on-init
      >
        <!-- Custom Node Types -->
        <template #node-custom="nodeProps">
          <div :class="nodeProps.data.className">
            <div v-html="nodeProps.data.label"></div>
          </div>
        </template>
        
        <!-- Add background pattern -->
        <Background 
          :pattern-color="'#aaa'" 
          :gap="24"
          :variant="isDarkMode ? 'dots' : 'lines'" 
        />
        
        <!-- Add Controls -->
        <Controls />
        
        <!-- Add panel with legend -->
        <Panel :position="'top-right'" class="legend-panel">
          <div class="legend bg-white dark:bg-gray-700 p-2 rounded shadow-md">
            <h4 class="text-sm font-semibold text-gray-700 dark:text-white mb-1">Legend</h4>
            <div class="legend-item flex items-center">
              <div class="w-3 h-3 bg-indigo-500 dark:bg-indigo-400 rounded-sm mr-2"></div>
              <span class="text-xs text-gray-600 dark:text-gray-300">Table</span>
            </div>
            <div class="legend-item flex items-center">
              <div class="w-3 h-3 bg-emerald-500 dark:bg-emerald-400 rounded-sm mr-2"></div>
              <span class="text-xs text-gray-600 dark:text-gray-300">Result/Values</span>
            </div>
            <div class="legend-item flex items-center">
              <div class="w-3 h-3 border border-indigo-500 dark:border-indigo-400 mr-2"></div>
              <span class="text-xs text-gray-600 dark:text-gray-300">Join</span>
            </div>
            <div v-if="dbStructure" class="mt-2 pt-2 border-t border-gray-200 dark:border-gray-600">
              <div class="text-xs font-medium text-gray-600 dark:text-gray-400 mb-1">Database Structure</div>
              <div class="legend-item flex items-center">
                <span class="mr-2">🔑</span>
                <span class="text-xs text-gray-600 dark:text-gray-300">Primary Key</span>
              </div>
              <div class="legend-item flex items-center">
                <span class="mr-2">🔗</span>
                <span class="text-xs text-gray-600 dark:text-gray-300">Foreign Key</span>
              </div>
              <div class="legend-item flex items-center">
                <span class="text-xs italic text-gray-400 mr-2">(null)</span>
                <span class="text-xs text-gray-600 dark:text-gray-300">Nullable Field</span>
              </div>
              <div class="legend-item flex items-center">
                <div class="w-8 h-0.5 bg-indigo-500 dark:bg-indigo-400 mr-2 opacity-70" style="border: 1px dashed #818cf8;"></div>
                <span class="text-xs text-gray-600 dark:text-gray-300">Implied FK Relation</span>
              </div>
              <div class="legend-item flex items-center mt-1">
                <span class="relationship-info mr-2 text-xs">→ table.column</span>
                <span class="text-xs text-gray-600 dark:text-gray-300">References</span>
              </div>
              <div class="legend-item flex items-center">
                <span class="relationship-info referenced-by mr-2 text-xs">← table.column</span>
                <span class="text-xs text-gray-600 dark:text-gray-300">Referenced By</span>
              </div>
            </div>
          </div>
        </Panel>
      </VueFlow>
      <div v-else class="empty-diagram">
        <p class="text-gray-500 dark:text-gray-400 text-center">No tables found in the SQL query</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed, onUnmounted } from 'vue';
import { parse } from '../utils/sqlParser';
import Button from './Button.vue';
import { 
  VueFlow, 
  useVueFlow,
  Panel
} from '@vue-flow/core';
import '@vue-flow/core/dist/style.css';
import '@vue-flow/core/dist/theme-default.css';

// Additional components
import { Background } from '@vue-flow/background';
import { Controls } from '@vue-flow/controls';
import { MiniMap } from '@vue-flow/minimap';

// Additional styles
import '@vue-flow/minimap/dist/style.css';
import '@vue-flow/controls/dist/style.css';

const props = defineProps<{
  query: string;
  databaseStructure?: string;
}>();

const diagramContainer = ref<HTMLElement | null>(null);
const queryType = ref<'SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'OTHER'>('SELECT');
const isDarkMode = ref(false);

// Elements for Vue Flow
const nodes = ref<any[]>([]);
const edges = ref<any[]>([]);

// Get Vue Flow instance
const { zoomIn, zoomOut, fitView } = useVueFlow({
  nodes: [],
  edges: []
});

// Add type definitions for database structure
interface Column {
  name: string;
  type: string;
  nullable: string;
  key: string;
  default: string;
  extra: string;
  isPrimary: boolean;
  isForeign: boolean;
}

interface ForeignKey {
  columnName: string;
  referencedTable: string;
  referencedColumn: string;
  constraintName: string;
}

interface Table {
  name: string;
  columns: Column[];
  foreignKeys: ForeignKey[];
}

interface DatabaseStructure {
  tables: Table[];
}

// Store the parsed database structure
const dbStructure = ref<DatabaseStructure | null>(null);

// Function to check if dark mode is enabled
const updateDarkMode = () => {
  isDarkMode.value = document.documentElement.classList.contains('dark') || 
                   document.body.classList.contains('dark') ||
                   window.matchMedia('(prefers-color-scheme: dark)').matches;
};

// Parse SQL query and create diagram
const parseQuery = (sql: string) => {
  try {
    // Parse SQL using our custom parser
    const result = parse(sql);
    queryType.value = result.queryType;
    
    if (result.error) {
      console.error('SQL parsing error:', result.error);
      return;
    }
    
    // Generate Vue Flow nodes and edges
    const diagramNodes: any[] = [];
    const diagramEdges: any[] = [];
    
    // Process tables as nodes and enhance with DB structure info
    result.tables.forEach((table, index) => {
      // Calculate position
      const position = table.position || {
        x: 100 + (index % 3) * 350,
        y: 100 + Math.floor(index / 3) * 250
      };
      
      // Find table in database structure
      const dbTable = dbStructure.value?.tables.find(t => 
        t.name.toLowerCase() === table.name.toLowerCase()
      );
      
      // Create node for each table
      diagramNodes.push({
        id: table.id,
        type: 'custom',
        data: { 
          label: buildTableLabel(table, result.queryType),
          className: result.queryType === 'INSERT' ? 'insert-table-node' : 'table-node',
          dbTable: dbTable // Store DB structure info for reference
        },
        position,
        style: {
          width: 300,
          borderRadius: '5px',
          backgroundColor: 'white'
        }
      });
    });
    
    // Process explicit joins from the query
    result.joins.forEach(join => {
      if (join.condition && join.condition.sourceId && join.condition.targetId) {
        const condition = join.condition; // Store the condition to prevent undefined checks
        
        // Get table names for better labels
        const sourceTable = result.tables.find(t => t.id === condition.sourceId);
        const targetTable = result.tables.find(t => t.id === condition.targetId);
        
        // Create a more descriptive join label
        const joinType = join.type || 'JOIN';
        let columnInfo = '';
        
        // Add column information if available
        // Most SQL parsers will represent the join condition with leftColumn and rightColumn
        if (condition.leftColumn && condition.rightColumn) {
          const sourceTableName = sourceTable?.name || '';
          const targetTableName = targetTable?.name || '';
          columnInfo = `${sourceTableName}.${condition.leftColumn} = ${targetTableName}.${condition.rightColumn}`;
        }
        
        diagramEdges.push({
          id: join.id,
          source: condition.sourceId,
          target: condition.targetId,
          label: columnInfo ? `${joinType} \n${columnInfo}` : joinType,
          animated: true,
          style: { stroke: '#4f46e5' },
          labelBgStyle: { fill: '#4f46e5', color: '#ffffff', fillOpacity: 0.7 },
          labelStyle: { 
            fill: '#ffffff', 
            fontWeight: 500, 
            fontSize: columnInfo ? 10 : 12 
          },
        });
      }
    });
    
    // Add foreign key relationships from database structure
    // Only do this for SELECT queries without explicit joins
    if (result.queryType === 'SELECT' && dbStructure.value && result.joins.length === 0) {
      // Find all tables in the query
      const tableMap = new Map<string, any>();
      result.tables.forEach(table => {
        tableMap.set(table.name.toLowerCase(), table);
      });
      
      // Check each table for foreign keys
      tableMap.forEach((sourceTable) => {
        const dbSourceTable = dbStructure.value?.tables.find(t => 
          t.name.toLowerCase() === sourceTable.name.toLowerCase()
        );
        
        if (dbSourceTable && dbSourceTable.foreignKeys) {
          dbSourceTable.foreignKeys.forEach(fk => {
            const targetTable = tableMap.get(fk.referencedTable.toLowerCase());
            
            // If both source and target tables are in the query, add implied join
            if (targetTable) {
              const edgeId = `fk-${sourceTable.id}-${fk.columnName}-to-${targetTable.id}-${fk.referencedColumn}`;
              
              // Check if this edge doesn't exist yet
              if (!diagramEdges.some(e => e.id === edgeId)) {
                diagramEdges.push({
                  id: edgeId,
                  source: sourceTable.id,
                  target: targetTable.id,
                  label: `${sourceTable.name}.${fk.columnName} → ${targetTable.name}.${fk.referencedColumn}`,
                  animated: false,
                  style: { stroke: '#4f46e5', strokeDasharray: '5 5', opacity: 0.7 },
                  labelBgStyle: { fill: '#4f46e5', color: '#ffffff', fillOpacity: 0.5 },
                  labelStyle: { fill: '#ffffff', fontSize: 10, fontWeight: 500 },
                });
              }
            }
          });
        }
      });
    }
    
    // Create result node for SELECT queries
    if (result.queryType === 'SELECT' && result.selectedColumns.length > 0) {
      const resultNodeId = `result-${Date.now()}`;
      
      // Calculate position for result node
      const resultY = Math.max(...diagramNodes.map(n => n.position.y)) + 300;
      
      // Create node for query result
      diagramNodes.push({
        id: resultNodeId,
        type: 'custom',
        data: { 
          label: buildResultLabel(result.selectedColumns),
          className: 'result-node'
        },
        position: {
          x: 100 + (result.tables.length % 3) * 200,
          y: resultY
        },
        style: {
          width: 300,
          borderRadius: '5px',
          backgroundColor: 'white'
        }
      });
      
      // Connect tables to result node
      const selectedTables = new Set<string>();
      result.selectedColumns.forEach(col => {
        if (col.table) {
          selectedTables.add(col.table);
        }
      });
      
      selectedTables.forEach(tableName => {
        const sourceTable = result.tables.find(t => 
          t.name === tableName || t.alias === tableName
        );
        
        if (sourceTable) {
          diagramEdges.push({
            id: `select-${tableName}-${Date.now()}`,
            source: sourceTable.id,
            target: resultNodeId,
            animated: true,
            style: { stroke: '#10b981', strokeDasharray: '5 5' },
            label: 'SELECT',
            labelBgStyle: { fill: '#10b981', color: '#ffffff', fillOpacity: 0.7 },
            labelStyle: { fill: '#ffffff' }
          });
        }
      });
    }
    
    // Create values node for INSERT queries
    if (result.queryType === 'INSERT' && result.insertValues && result.tables.length > 0) {
      const valuesNodeId = `values-${Date.now()}`;
      const tableNode = result.tables[0]; // In INSERT, we only have one table
      
      // Create values node
      diagramNodes.push({
        id: valuesNodeId,
        type: 'custom',
        data: { 
          label: buildValuesLabel(result.insertValues),
          className: 'values-node'
        },
        position: {
          x: tableNode.position?.x || 100,
          y: (tableNode.position?.y || 100) - 200
        },
        style: {
          width: 300,
          borderRadius: '5px',
          backgroundColor: 'white'
        }
      });
      
      // Create edge from values to table
      diagramEdges.push({
        id: `insert-${Date.now()}`,
        source: valuesNodeId,
        target: tableNode.id,
        animated: true,
        style: { stroke: '#10b981' },
        label: 'INSERT',
        labelBgStyle: { fill: '#10b981', color: '#ffffff', fillOpacity: 0.7 },
        labelStyle: { fill: '#ffffff', fontWeight: 700 }
      });
    }
    
    // Update nodes and edges
    nodes.value = diagramNodes;
    edges.value = diagramEdges;
    
    // Make sure the view fits all elements
    setTimeout(() => {
      fitView();
    }, 100);
  } catch (error) {
    console.error('Error generating diagram:', error);
  }
};

// Parse database structure
const parseDatabaseStructure = (structureStr: string) => {
  try {
    if (!structureStr) {
      dbStructure.value = null;
      return;
    }

    const structure: DatabaseStructure = JSON.parse(structureStr);
    
    if (!structure.tables || !structure.tables.length) {
      dbStructure.value = null;
      return;
    }
    
    dbStructure.value = structure;
  } catch (error) {
    console.error('Error parsing database structure:', error);
    dbStructure.value = null;
  }
};

// Enhanced function to build table label with actual column info from database structure
const buildTableLabel = (table: any, queryType: string) => {
  // Try to find this table in the database structure
  const dbTable = dbStructure.value?.tables.find(t => 
    t.name.toLowerCase() === table.name.toLowerCase()
  );
  
  let headerClass = queryType === 'INSERT' ? 'insert-header' : 'table-header';
  let headerTitle = table.name;
  if (table.alias) {
    headerTitle += ` (${table.alias})`;
  }
  
  let columnsHtml = '';
  
  // If we have database structure information for this table, use it
  if (dbTable) {
    // Build columns HTML using actual column information
    columnsHtml = '<div class="table-columns">';
    
    // Find foreign keys for this table
    const foreignKeyColumns = new Set<string>();
    const referencedByColumns = new Set<string>();
    
    // Check if this table has foreign keys
    if (dbTable.foreignKeys) {
      dbTable.foreignKeys.forEach(fk => {
        foreignKeyColumns.add(fk.columnName.toLowerCase());
      });
    }
    
    // Check if this table is referenced by other tables
    dbStructure.value?.tables.forEach(otherTable => {
      if (otherTable.foreignKeys) {
        otherTable.foreignKeys.forEach(fk => {
          if (fk.referencedTable.toLowerCase() === dbTable.name.toLowerCase()) {
            referencedByColumns.add(fk.referencedColumn.toLowerCase());
          }
        });
      }
    });
    
    // First add columns mentioned in the query
    if (table.columns && table.columns.length > 0) {
      table.columns.forEach((col: any) => {
        const dbCol = dbTable.columns.find(c => 
          c.name.toLowerCase() === col.name.toLowerCase()
        );
        
        if (dbCol) {
          const isPrimary = dbCol.isPrimary;
          const isForeign = dbCol.isForeign || foreignKeyColumns.has(dbCol.name.toLowerCase());
          const isReferenced = referencedByColumns.has(dbCol.name.toLowerCase());
          const isNullable = dbCol.nullable === 'YES';
          
          let relationshipInfo = '';
          if (isForeign) {
            // Find which table this foreign key references
            const fk = dbTable.foreignKeys?.find(
              fk => fk.columnName.toLowerCase() === dbCol.name.toLowerCase()
            );
            if (fk) {
              relationshipInfo = `<span class="relationship-info">→ ${fk.referencedTable}.${fk.referencedColumn}</span>`;
            }
          }
          
          if (isReferenced) {
            // Find which tables reference this column
            const referencingTables: string[] = [];
            dbStructure.value?.tables.forEach(otherTable => {
              if (otherTable.foreignKeys) {
                otherTable.foreignKeys.forEach(fk => {
                  if (
                    fk.referencedTable.toLowerCase() === dbTable.name.toLowerCase() &&
                    fk.referencedColumn.toLowerCase() === dbCol.name.toLowerCase()
                  ) {
                    referencingTables.push(`${otherTable.name}.${fk.columnName}`);
                  }
                });
              }
            });
            
            if (referencingTables.length > 0) {
              const truncatedList = referencingTables.length > 2 
                ? referencingTables.slice(0, 2).join(', ') + `, +${referencingTables.length - 2} more`
                : referencingTables.join(', ');
              relationshipInfo += `<span class="relationship-info referenced-by">← ${truncatedList}</span>`;
            }
          }
          
          columnsHtml += `
            <div class="column-row ${isPrimary ? 'primary-key' : ''} ${isForeign ? 'foreign-key' : ''} ${isReferenced ? 'referenced-column' : ''}">
              <div class="column-icon">
                ${isPrimary ? '🔑' : isForeign ? '🔗' : ''}
              </div>
              <div class="column-name">
                ${col.name}${isNullable ? ' <span class="nullable">(null)</span>' : ''}
              </div>
              <div class="column-type">${dbCol.type}</div>
              ${relationshipInfo}
            </div>
          `;
        } else {
          // Column mentioned in query but not found in db structure
          columnsHtml += `
            <div class="column-row">
              <div class="column-name">${col.name}</div>
            </div>
          `;
        }
      });
    } else if (queryType === 'SELECT' && table.columns.length === 0) {
      // For SELECT * case, show all columns from the DB structure
      dbTable.columns.forEach(col => {
        const isPrimary = col.isPrimary;
        const isForeign = col.isForeign || foreignKeyColumns.has(col.name.toLowerCase());
        const isReferenced = referencedByColumns.has(col.name.toLowerCase());
        const isNullable = col.nullable === 'YES';
        
        let relationshipInfo = '';
        if (isForeign) {
          // Find which table this foreign key references
          const fk = dbTable.foreignKeys?.find(
            fk => fk.columnName.toLowerCase() === col.name.toLowerCase()
          );
          if (fk) {
            relationshipInfo = `<span class="relationship-info">→ ${fk.referencedTable}.${fk.referencedColumn}</span>`;
          }
        }
        
        if (isReferenced) {
          // Find which tables reference this column
          const referencingTables: string[] = [];
          dbStructure.value?.tables.forEach(otherTable => {
            if (otherTable.foreignKeys) {
              otherTable.foreignKeys.forEach(fk => {
                if (
                  fk.referencedTable.toLowerCase() === dbTable.name.toLowerCase() &&
                  fk.referencedColumn.toLowerCase() === col.name.toLowerCase()
                ) {
                  referencingTables.push(`${otherTable.name}.${fk.columnName}`);
                }
              });
            }
          });
          
          if (referencingTables.length > 0) {
            const truncatedList = referencingTables.length > 2 
              ? referencingTables.slice(0, 2).join(', ') + `, +${referencingTables.length - 2} more`
              : referencingTables.join(', ');
            relationshipInfo += `<span class="relationship-info referenced-by">← ${truncatedList}</span>`;
          }
        }
        
        columnsHtml += `
          <div class="column-row ${isPrimary ? 'primary-key' : ''} ${isForeign ? 'foreign-key' : ''} ${isReferenced ? 'referenced-column' : ''}">
            <div class="column-icon">
              ${isPrimary ? '🔑' : isForeign ? '🔗' : ''}
            </div>
            <div class="column-name">
              ${col.name}${isNullable ? ' <span class="nullable">(null)</span>' : ''}
            </div>
            <div class="column-type">${col.type}</div>
            ${relationshipInfo}
          </div>
        `;
      });
    }
    
    columnsHtml += '</div>';
  } else {
    // Fallback to original implementation if table not found in db structure
    // Only extract the columns mentioned in the query
    let columnListHtml = '';
    if (table.columns && table.columns.length > 0) {
      table.columns.forEach((col: any) => {
        columnListHtml += `<div class="column-row"><span class="column-name">${col.name}</span></div>`;
      });
    } else {
      columnListHtml = '<div class="column-row"><span class="font-italic text-gray-400">All columns (*)</span></div>';
    }
    
    columnsHtml = `<div class="table-columns">${columnListHtml}</div>`;
  }
  
  return `
    <div class="table-label">
      <div class="${headerClass}">
        <span class="header-title">${headerTitle}</span>
      </div>
      ${columnsHtml}
    </div>
  `;
};

// Helper function to build result HTML label
function buildResultLabel(columns: any[]): string {
  // Generate HTML label for the result node
  let html = `<div class="result-label">`;
  
  // Header
  html += `<div class="result-header">
    <div class="header-title">Query Result</div>
    <div class="column-count">${columns.length} column${columns.length !== 1 ? 's' : ''}</div>
  </div>`;
  
  // Columns section
  html += `<div class="result-columns">
    <div class="result-column-header">
      <div class="column-name-header">COLUMN</div>
      <div class="column-source-header">SOURCE</div>
    </div>`;
  
  // Add all columns
  for (const col of columns) {
    html += `<div class="result-column-row">
      <div class="result-column-name">${col.name || col.alias || 'Unknown'}</div>
      <div class="result-column-source">${col.table || '-'}</div>
    </div>`;
  }
  
  html += `</div></div>`;
  
  return html;
}

// Helper function to build values HTML label
function buildValuesLabel(values: Record<string, string>): string {
  // Generate HTML label for values node
  let html = `<div class="values-label">`;
  
  // Header
  html += `<div class="values-header">
    <div class="header-title">INSERT Values</div>
  </div>`;
  
  // Values section
  html += `<div class="values-list">`;
  
  // Add column header
  html += ``;
  
  // Add all values
  for (const [key, value] of Object.entries(values)) {
    html += `<div class="value-row">
      <div class="value-column">${key}</div>
      <div class="value-data">${formatValue(value)}</div>
    </div>`;
  }
  
  html += `</div></div>`;
  
  return html;
}

// Get column type with reasonable defaults
function getColumnType(column: any): string {
  if (column.type) return column.type;
  
  // Try to infer type from column name
  const name = column.name.toLowerCase();
  if (name === 'id' || name.endsWith('_id') || name.includes('id_')) {
    return 'int';
  } else if (name.includes('date') || name.includes('time')) {
    return 'datetime';
  } else if (name.includes('price') || name.includes('amount') || name.includes('total')) {
    return 'decimal';
  } else if (name === 'active' || name === 'enabled' || name === 'status') {
    return 'boolean';
  } else {
    return 'varchar';
  }
}

// Format value for display (truncate if too long, format strings)
function formatValue(value: string): string {
  // If it's a string (starts with quote), format it and truncate if needed
  if (value.startsWith("'") || value.startsWith('"')) {
    const cleaned = value.substring(1, value.length - 1);
    return cleaned.length > 15 ? `'${cleaned.substring(0, 15)}...'` : value;
  }
  
  // For other values (numbers, etc.)
  return value.length > 15 ? `${value.substring(0, 15)}...` : value;
}

onMounted(() => {
  // Parse initial query
  if (props.query) {
    parseQuery(props.query);
  }
  
  // Initialize dark mode detection
  updateDarkMode();
  
  // Add listener for theme changes
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', updateDarkMode);
  
  // Check for Tailwind/custom dark mode toggle
  const observer = new MutationObserver(updateDarkMode);
  observer.observe(document.documentElement, { attributes: true, attributeFilter: ['class'] });
  observer.observe(document.body, { attributes: true, attributeFilter: ['class'] });
});

// Clean up event listeners
onUnmounted(() => {
  // Remove dark mode listeners
  window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', updateDarkMode);
});

// Update the watch for query changes
watch(() => props.query, (newQuery) => {
  if (newQuery) {
    parseQuery(newQuery);
  }
});

// Watch for database structure changes
watch(() => props.databaseStructure, (newValue) => {
  if (newValue) {
    parseDatabaseStructure(newValue);
    
    // Re-parse query if we already have one to update the diagram
    if (props.query) {
      parseQuery(props.query);
    }
  }
}, { immediate: true });
</script>

<style>
.sql-diagram-container {
  display: flex;
  flex-direction: column;
  width: 100%;
  height: 100%;
  min-height: 500px;
}

.diagram-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.diagram-content {
  flex: 1;
  overflow: hidden;
  border: 1px solid #e2e8f0;
  border-radius: 0.5rem;
  min-height: 500px;
}

.vue-flow-wrapper {
  width: 100%;
  height: 100%;
  min-height: 500px;
}

.empty-diagram {
  display: flex;
  justify-content: center;
  align-items: center;
  height: 100%;
  min-height: 300px;
}

/* Legend styling */
.legend-panel {
  z-index: 10;
}

.legend {
  font-size: 12px;
}

.legend-item {
  margin-top: 4px;
}

/* Node styling */
.vue-flow__node {
  border-radius: 5px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  font-family: ui-sans-serif, system-ui, sans-serif;
}

.vue-flow__node.table-node {
  border-left: 4px solid #4f46e5;
}

.vue-flow__node.insert-table-node {
  border-left: 4px solid #10b981;
}

.vue-flow__node.result-node {
  border-left: 4px solid #10b981;
}

.vue-flow__node.values-node {
  border-left: 4px solid #10b981;
}

/* Edge styling */
.vue-flow__edge-path {
  stroke-width: 2;
}

.vue-flow__edge-text {
  font-size: 12px;
}

/* Handle styling */
.vue-flow__handle {
  width: 8px;
  height: 8px;
}

/* Table styling inside nodes */
.table-label, .result-label, .values-label {
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.table-header, .insert-header, .result-header, .values-header {
  padding: 10px;
  color: white;
  font-weight: bold;
  border-top-left-radius: 3px;
  border-top-right-radius: 3px;
}

.table-header {
  background-color: #4f46e5;
}

.insert-header, .result-header, .values-header {
  background-color: #10b981;
}

.header-title {
  font-size: 14px;
  font-weight: bold;
}

.header-subtitle {
  font-size: 12px;
  opacity: 0.8;
}

.column-count {
  font-size: 12px;
  font-weight: normal;
  opacity: 0.9;
}

.columns-list, .result-columns, .values-list {
  background-color: white;
  border: 1px solid #e5e7eb;
  border-top: none;
  padding: 4px 0;
}

.dark .columns-list, .dark .result-columns, .dark .values-list {
  background-color: #1f2937;
  border-color: #374151;
  color: #e5e7eb;
}

.column-row, .result-column-row, .value-row {
  display: flex;
  padding: 8px 10px;
  border-bottom: 1px solid #f3f4f6;
  font-size: 13px;
  align-items: center;
}

.dark .column-row, .dark .result-column-row, .dark .value-row {
  border-bottom-color: #374151;
}

.column-row:last-child, .result-column-row:last-child, .value-row:last-child {
  border-bottom: none;
}

.column-name, .result-column-name, .value-name, .value-column {
  flex: 1.5;
  display: flex;
  align-items: center;
  font-weight: 600;
  color: #1f2937;
}

.column-type {
  flex: 1;
  color: #6b7280;
}

.dark .column-type {
  color: #9ca3af;
}

.column-value, .value-data {
  flex: 1.5;
  color: #10b981;
}

.dark .column-value, .dark .value-data {
  color: #34d399;
}

.column-badge {
  font-size: 10px;
  padding: 1px 4px;
  border-radius: 3px;
  font-weight: bold;
  margin-left: 5px;
}

.pk-badge {
  background-color: rgba(147, 51, 234, 0.2);
  color: #7e22ce;
}

.dark .pk-badge {
  background-color: rgba(147, 51, 234, 0.3);
  color: #a855f7;
}

.fk-badge {
  background-color: rgba(37, 99, 235, 0.2);
  color: #2563eb;
}

.dark .fk-badge {
  background-color: rgba(37, 99, 235, 0.3);
  color: #3b82f6;
}

/* Result columns styling */
.result-column-header {
  display: flex;
  padding: 6px 10px;
  background-color: #f9fafb;
  border-bottom: 1px solid #f3f4f6;
  font-size: 11px;
  font-weight: 500;
}

.dark .result-column-header {
  background-color: #374151;
  border-bottom-color: #4b5563;
}

.column-name-header, .column-source-header {
  color: black;
}

.dark .column-name-header, .dark .column-source-header {
  color: #9ca3af;
}

.column-name-header {
  flex: 1.5;
}

.column-source-header {
  flex: 1;
  text-align: right;
}

.result-column-source {
  flex: 1;
  text-align: right;
  color: #6b7280;
}

.dark .result-column-source {
  color: #9ca3af;
}

.values-header {
  background-color: #10b981;
  padding: 10px;
  color: white;
  font-weight: bold;
  border-top-left-radius: 3px;
  border-top-right-radius: 3px;
}

.value-column {
  font-weight: 600;
  color: #1f2937;
  padding-right: 8px;
}

.dark .value-column {
  color: #f3f4f6;
}

.value-row {
  display: flex;
  padding: 8px 10px;
  border-bottom: 1px solid #f3f4f6;
  font-size: 13px;
  align-items: center;
}
</style>

<style scoped>
/* Table styling inside nodes */
.table-label, .result-label, .values-label {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.table-header, .insert-header, .result-header, .values-header {
  padding: 10px;
  color: white;
  font-weight: bold;
  border-top-left-radius: 3px;
  border-top-right-radius: 3px;
}

.table-header {
  background-color: #4f46e5;
}

.insert-header, .result-header, .values-header {
  background-color: #10b981;
}

.header-title {
  font-size: 14px;
  font-weight: bold;
}

.table-columns {
  padding: 8px;
  background-color: white;
  border: 1px solid #e2e8f0;
  border-bottom-left-radius: 3px;
  border-bottom-right-radius: 3px;
  overflow-y: auto;
  max-height: 300px;
}

.dark .table-columns {
  background-color: #1f2937;
  border-color: #374151;
}

.column-row {
  display: flex;
  align-items: center;
  padding: 4px 8px;
  border-bottom: 1px solid #f1f5f9;
  font-size: 12px;
}

.dark .column-row {
  border-color: #2d3748;
}

.column-row:last-child {
  border-bottom: none;
}

.column-row.primary-key {
  background-color: #fff7ed;
}

.dark .column-row.primary-key {
  background-color: #422006;
}

.column-row.foreign-key {
  background-color: #f0f9ff;
}

.dark .column-row.foreign-key {
  background-color: #082f49;
}

.column-icon {
  margin-right: 8px;
  width: 16px;
  display: flex;
  justify-content: center;
}

.column-name {
  flex: 1;
  font-weight: 500;
}

.column-type {
  color: #64748b;
  font-size: 10px;
  background-color: #f1f5f9;
  padding: 2px 6px;
  border-radius: 10px;
}

.dark .column-type {
  background-color: #334155;
  color: #94a3b8;
}

.nullable {
  font-style: italic;
  opacity: 0.6;
  font-size: 10px;
}

/* Fix any dark mode specific styling */
.dark .empty-diagram {
  color: #94a3b8;
}

/* Enhance node display */
:deep(.vue-flow__node-custom) {
  background: white;
  border-radius: 5px;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.dark :deep(.vue-flow__node-custom) {
  background: #1f2937;
  box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.3), 0 2px 4px -1px rgba(0, 0, 0, 0.2);
}

/* Enhance legends */
.legend-panel {
  font-size: 12px;
  opacity: 0.9;
}

/* Add database structure legend items */
:deep(.legend) {
  padding: 10px !important;
}

:deep(.legend-item) {
  display: flex;
  align-items: center;
  margin-bottom: 4px;
}

/* Add these new styles for relationship info */
.relationship-info {
  font-size: 10px;
  padding: 1px 4px;
  margin-left: 4px;
  border-radius: 3px;
  color: #4f46e5;
  background-color: rgba(79, 70, 229, 0.1);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
  display: inline-block;
}

.dark .relationship-info {
  color: #818cf8;
  background-color: rgba(129, 140, 248, 0.1);
}

.relationship-info.referenced-by {
  color: #10b981;
  background-color: rgba(16, 185, 129, 0.1);
}

.dark .relationship-info.referenced-by {
  color: #34d399;
  background-color: rgba(52, 211, 153, 0.1);
}

.referenced-column {
  border-left: 2px solid #10b981;
}

.dark .referenced-column {
  border-left: 2px solid #34d399;
}

/* Make the table column slightly wider to accommodate relationship info */
.table-columns {
  min-width: 280px;
}

.column-row {
  flex-wrap: wrap;
}

/* Add icons to relationship indicators */
.relationship-info::before {
  content: "";
  display: inline-block;
  width: 10px;
  height: 10px;
  margin-right: 2px;
  background-size: contain;
  background-repeat: no-repeat;
  background-position: center;
}

/* Style the table nodes to be larger to accommodate the relationship info */
:deep(.vue-flow__node-custom) {
  min-width: 300px;
  width: auto !important;
}
</style> 