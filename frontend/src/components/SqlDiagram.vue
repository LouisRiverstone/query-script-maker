<template>
  <div class="sql-diagram-container">
    <div class="diagram-header">
      <h3 class="text-lg font-semibold text-black dark:text-white">
        SQL Diagram <span v-if="queryType" class="ml-2 text-sm font-normal">({{ queryType }} query)</span>
      </h3>
      <div class="diagram-controls">
        <Button type="button" @click="zoomIn" class="mr-1">+</Button>
        <Button type="button" @click="zoomOut" class="mr-1">-</Button>
        <Button type="button" @click="resetZoom">Reset</Button>
      </div>
    </div>
    <div ref="diagramContainer" class="diagram-content bg-white dark:bg-gray-800 rounded-lg p-4">
      <svg ref="svgElement" :width="width" :height="height" @wheel.prevent="handleZoom">
        <defs>
          <!-- Arrow markers for join relationships -->
          <marker id="join-arrow-light" viewBox="0 0 20 20" refX="8" refY="5" 
                  markerWidth="16" markerHeight="16" orient="auto">
            <path d="M 0 0 L 10 5 L 0 10 z" fill="#4f46e5" stroke="#ffffff" stroke-width="1.5"></path>
          </marker>
          <marker id="join-arrow-dark" viewBox="0 0 20 20" refX="8" refY="5" 
                  markerWidth="16" markerHeight="16" orient="auto">
            <path d="M 0 0 L 10 5 L 0 10 z" fill="#818cf8" stroke="#1f2937" stroke-width="1.5"></path>
          </marker>
          
          <!-- Glowing effects for join lines -->
          <filter id="glow-light" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur result="blur" stdDeviation="3"></feGaussianBlur>
            <feComposite in="SourceGraphic" in2="blur" operator="over"></feComposite>
          </filter>
          <filter id="glow-dark" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur result="blur" stdDeviation="4"></feGaussianBlur>
            <feComposite in="SourceGraphic" in2="blur" operator="over"></feComposite>
          </filter>
          
          <!-- Glow effect for text highlights -->
          <filter id="text-glow" x="-50%" y="-50%" width="200%" height="200%">
            <feGaussianBlur result="blur" stdDeviation="1"></feGaussianBlur>
            <feComposite in="SourceGraphic" in2="blur" operator="over"></feComposite>
          </filter>
          
          <marker id="insert-arrow-light" viewBox="0 0 10 10" refX="9" refY="5" 
                markerWidth="10" markerHeight="10" orient="auto">
            <path d="M 0 0 L 10 5 L 0 10 z" fill="#10b981" />
          </marker>
          
          <marker id="insert-arrow-dark" viewBox="0 0 10 10" refX="9" refY="5" 
                markerWidth="10" markerHeight="10" orient="auto">
            <path d="M 0 0 L 10 5 L 0 10 z" fill="#34d399" />
          </marker>
        </defs>
        <g :transform="`scale(${scale}) translate(${offsetX}, ${offsetY})`">
          <!-- Highlight columns involved in joins first (as background) -->
          <g v-for="(table, tableIndex) in parsedTables" :key="`join-highlight-${tableIndex}`">
            <g v-for="(column, colIndex) in table.columns" :key="`join-col-${colIndex}`">
              <rect v-if="isColumnPartOfJoin(tableIndex, colIndex)"
                 :x="table.x" 
                 :y="table.y + 40 + colIndex * 30" 
                 width="300" height="30" 
                 class="fill-indigo-100 dark:fill-indigo-900/30 stroke-indigo-300 dark:stroke-indigo-600 stroke-2" 
                 rx="4" ry="4" />
            </g>
          </g>
          
          <!-- Relationships will be rendered here (behind tables) -->
          <g v-for="(relation, index) in parsedRelations" :key="`rel-${index}`" 
             class="relation-line">
            <!-- Linha de guia pontilhada para tornar o caminho mais visível -->
            <path :d="relation.path" fill="none" stroke-width="1.5" 
                  stroke-dasharray="4,4"
                  class="stroke-indigo-300 dark:stroke-indigo-600 opacity-60" />
            
            <!-- Visual connector line showing the relationship -->
            <path :d="relation.path" fill="none" stroke-width="5" 
                  :class="[
                    isDarkMode ? 'dark-join-path filter-glow-dark' : 'light-join-path filter-glow-light'
                  ]"
                  :marker-end="isDarkMode ? 'url(#join-arrow-dark)' : 'url(#join-arrow-light)'" />
            
            <!-- Connection points at source column -->
            <circle :cx="relation.startX" :cy="relation.startY" r="6" 
                   class="fill-indigo-500 dark:fill-indigo-400 stroke-white dark:stroke-gray-800 stroke-2" />
            
            <!-- Connection points at target column -->
            <circle :cx="relation.endX" :cy="relation.endY" r="6"
                   class="fill-indigo-500 dark:fill-indigo-400 stroke-white dark:stroke-gray-800 stroke-2" />
            
            <!-- Highlight PK/FK relationship with badges -->
            <rect v-if="relation.isPkFkRelation" 
                  :x="relation.labelX - 35" :y="relation.labelY - 15" 
                  width="70" height="30" rx="15" ry="15" 
                  class="fill-purple-500/90 dark:fill-purple-600/90 stroke-white dark:stroke-gray-800 stroke-1" />
            <rect v-else
                  :x="relation.labelX - 30" :y="relation.labelY - 12" 
                  width="60" height="24" rx="12" ry="12" 
                  class="fill-indigo-500/90 dark:fill-indigo-600/90" />
            
            <!-- Join type text -->
            <text :x="relation.labelX" :y="relation.labelY + 5" text-anchor="middle" 
                  :class="relation.isPkFkRelation ? 'text-xs font-bold fill-white' : 'text-xs font-medium fill-white'">
              {{ relation.isPkFkRelation ? 'PK → FK' : relation.joinType }}
            </text>
            
            <!-- Clear relationship text below the line -->
            <text :x="relation.descriptionX" :y="relation.descriptionY" text-anchor="middle" 
                  class="text-sm font-bold fill-indigo-600 dark:fill-indigo-400 filter-glow">
              {{ relation.from.column }} = {{ relation.to.column }}
            </text>
          </g>
          
          <!-- Tables will be rendered here -->
          <g v-for="(table, index) in parsedTables" :key="`table-${index}`" 
             :transform="`translate(${table.x}, ${table.y})`"
             class="table-node">
            <!-- Table header -->
            <rect :width="table.width" height="40" rx="5" ry="5" 
                  :class="{'fill-indigo-500 dark:fill-indigo-700': queryType !== 'INSERT',
                           'fill-emerald-600 dark:fill-emerald-700': queryType === 'INSERT'}" />
            <text :x="table.width / 2" y="25" text-anchor="middle" 
                  class="font-bold text-white">{{ table.name }}</text>
            
            <!-- Table columns -->
            <rect x="0" y="40" :width="table.width" :height="table.columns.length * 30" 
                  fill="white" class="dark:fill-gray-700 stroke-gray-300 dark:stroke-gray-600" />
            <g v-for="(column, colIndex) in table.columns" :key="`${index}-${colIndex}`"
               :transform="`translate(0, ${40 + colIndex * 30})`">
              <!-- Highlight if column is involved in a join relationship -->
              <rect v-if="isColumnPartOfJoin(index, colIndex)" 
                    x="0" y="0" :width="table.width" height="30" 
                    class="fill-indigo-100/50 dark:fill-indigo-900/50 rounded-sm" />
              
              <!-- Column connection point indicator if part of join -->
              <circle v-if="isColumnPartOfJoin(index, colIndex)"
                      :cx="table.width < 250 ? table.width : table.width - 10" 
                      :cy="15"
                      r="4"
                      class="fill-indigo-400 dark:fill-indigo-300 stroke-white dark:stroke-gray-800 stroke-1" />
              
              <!-- Join indicator badge next to column -->
              <rect v-if="isColumnPartOfJoin(index, colIndex)"
                    :x="table.width - 40" 
                    :y="3"
                    width="26" height="22" 
                    rx="3" ry="3"
                    class="fill-indigo-500/80 dark:fill-indigo-600/80" />
              
              <text v-if="isColumnPartOfJoin(index, colIndex)"
                    :x="table.width - 27" 
                    :y="18"
                    class="text-xs font-bold fill-white">
                JOIN
              </text>
              
              <!-- Column name cell -->
              <text x="10" y="20" 
                    :class="[
                      'text-sm', 
                      isColumnPartOfJoin(index, colIndex) ? 'font-bold fill-indigo-700 dark:fill-indigo-300' : 'fill-gray-700 dark:fill-gray-300',
                      column.isPrimaryKey ? 'underline decoration-2 underline-offset-4' : ''
                    ]">
                {{ column.name }}
                <tspan v-if="column.isPrimaryKey" class="text-xs fill-purple-600 dark:fill-purple-400 font-bold"> PK</tspan>
                <tspan v-if="column.isForeignKey" class="text-xs fill-blue-600 dark:fill-blue-400 font-bold"> FK</tspan>
              </text>
              
              <!-- Type cell -->
              <text :x="Math.min(table.width * 0.5, 150)" y="20" 
                    :class="[
                      'text-sm', 
                      isColumnPartOfJoin(index, colIndex) ? 'font-medium fill-indigo-700 dark:fill-indigo-300' : 'fill-gray-500 dark:fill-gray-400'
                    ]">
                {{ column.type }}
              </text>
              
              <!-- Value cell (if present) -->
              <text v-if="column.value" :x="Math.min(table.width * 0.7, 220)" y="20" 
                    :class="[
                      'text-sm', 
                      isColumnPartOfJoin(index, colIndex) ? 'font-medium fill-indigo-700 dark:fill-indigo-300' : 'fill-gray-500 dark:fill-gray-400'
                    ]">
                {{ column.value }}
              </text>
            </g>
          </g>
          
          <!-- Result Box for SELECT queries with improved visibility -->
          <g v-if="resultBox && queryType === 'SELECT'" :transform="`translate(${resultBox.x}, ${resultBox.y})`">
            <!-- Enhanced outer border/shadow for result box -->
            <rect 
              x="-4" y="-4" 
              :width="resultBox.width + 8" 
              :height="getResultBoxHeight(resultBox.columns.length) + 8" 
              rx="8" ry="8"
              class="fill-none stroke-emerald-500 dark:stroke-emerald-400 stroke-2 filter-glow-light dark:filter-glow-dark opacity-70"
            />
            
            <!-- Connection lines from tables to results -->
            <g v-for="(column, colIndex) in resultBox.columns" :key="`result-conn-${colIndex}`">
              <path 
                v-if="column.sourceTable && column.table"
                :d="getResultConnectionPath(column, colIndex)"
                class="stroke-emerald-400/60 dark:stroke-emerald-400/50 result-connection"
                stroke-width="1.5"
                fill="none"
              />
              <!-- Círculo destacando o ponto de conexão na tabela -->
              <circle 
                v-if="column.sourceTable && column.table"
                :cx="getResultConnectionStartPoint(column).x - resultBox.x" 
                :cy="getResultConnectionStartPoint(column).y - resultBox.y" 
                r="3" 
                class="fill-emerald-400/80 dark:fill-emerald-400/70 stroke-white dark:stroke-gray-800 stroke-1"
              />
              <!-- Círculo destacando o ponto de conexão no resultado -->
              <circle 
                v-if="column.sourceTable && column.table"
                cx="-10" 
                :cy="getResultRowY(colIndex)" 
                r="3" 
                class="fill-emerald-400/80 dark:fill-emerald-400/70 stroke-white dark:stroke-gray-800 stroke-1"
              />
            </g>
            
            <!-- Result Header with simplified styling -->
            <rect 
              :width="resultBox.width" height="40" rx="6" ry="6" 
              class="fill-emerald-600 dark:fill-emerald-700 shadow-lg"
            />
            <text :x="resultBox.width / 2" y="25" text-anchor="middle" class="font-bold text-white text-base">
              Query Result
            </text>
            
            <!-- Result count badge (simplified) -->
            <rect 
              x="10" y="-15" 
              :width="resultBox.columns.length > 0 ? 55 : 80" 
              height="24" 
              rx="12" ry="12"
              class="fill-white dark:fill-gray-800 shadow-sm"
            />
            <text x="35" y="1" text-anchor="middle" class="text-xs font-medium fill-emerald-600 dark:fill-emerald-500">
              {{ resultBox.columns.length > 0 ? `${resultBox.columns.length} cols` : 'No columns' }}
            </text>
            
            <!-- Fundo das colunas do resultado -->
            <rect 
              x="0" y="40" 
              :width="resultBox.width" 
              :height="resultBox.columns.length * 30" 
              class="fill-white dark:fill-gray-700 stroke-gray-300 dark:stroke-gray-600"
            />
            
            <!-- Cabeçalho das colunas -->
            <rect 
              x="0" y="40" 
              :width="resultBox.width" 
              height="24" 
              class="fill-gray-100 dark:fill-gray-600/50 stroke-gray-200 dark:stroke-gray-600"
            />
            <text x="10" y="57" class="text-xs font-medium fill-gray-600 dark:fill-gray-300">
              COLUMN
            </text>
            <text :x="resultBox.width - 50" y="57" class="text-xs font-medium fill-gray-600 dark:fill-gray-300">
              SOURCE
            </text>
            
            <!-- Linhas do resultado -->
            <g v-for="(column, colIndex) in resultBox.columns" :key="`result-col-${colIndex}`">
              <!-- Fundo da linha -->
              <rect 
                x="0" 
                :y="64 + colIndex * 30" 
                :width="resultBox.width" 
                height="30" 
                class="fill-white dark:fill-gray-700 stroke-gray-200 dark:stroke-gray-600 stroke-1"
              />
              
              <!-- Nome da coluna com info de alias quando presente -->
              <text 
                x="10" 
                :y="64 + colIndex * 30 + 20" 
                class="text-sm font-medium fill-gray-800 dark:fill-gray-200"
              >
                {{ column.name }}
                <tspan v-if="column.originalName" class="text-xs fill-gray-500 dark:fill-gray-400">
                  (alias de {{ column.originalName }})
                </tspan>
              </text>
              
              <!-- Info da tabela de origem -->
              <text 
                :x="resultBox.width - 10" 
                :y="64 + colIndex * 30 + 20" 
                text-anchor="end"
                class="text-xs fill-gray-500 dark:fill-gray-400"
              >
                {{ column.table || '-' }}
              </text>
            </g>
          </g>
          
          <!-- Insert Visualization for INSERT queries -->
          <g v-if="queryType === 'INSERT'" 
             :transform="`translate(${insertVisualization.x}, ${insertVisualization.y})`">
            <!-- Arrow pointing to the table -->
            <path :d="insertVisualization.path" stroke-width="2" fill="none"
                  class="stroke-emerald-500 dark:stroke-emerald-400"
                  :marker-end="isDarkMode ? 'url(#insert-arrow-dark)' : 'url(#insert-arrow-light)'" />
            <!-- Insert label -->
            <text :x="insertVisualization.labelX" :y="insertVisualization.labelY" 
                  text-anchor="middle" class="text-lg font-bold fill-emerald-600 dark:fill-emerald-500">
              INSERT
            </text>
          </g>
        </g>
      </svg>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, computed, onUnmounted } from 'vue';
import { parse } from '../utils/sqlParser';
import Button from './Button.vue';

const props = defineProps<{
  query: string;
}>();

const svgElement = ref<SVGElement | null>(null);
const diagramContainer = ref<HTMLElement | null>(null);
const width = ref(800);
const height = ref(600);
const scale = ref(1);
const offsetX = ref(0);
const offsetY = ref(0);
const isDragging = ref(false);
const startX = ref(0);
const startY = ref(0);
const queryType = ref<'SELECT' | 'INSERT' | 'UPDATE' | 'DELETE' | 'OTHER'>('SELECT');

// Detect dark mode
const isDarkMode = ref(false);

// Store join relationships for highlighting
const joinRelationships = ref<Array<{tableIndex: number, columnIndex: number}>>([]);

// Function to check if a column is part of a join relationship
function isColumnPartOfJoin(tableIndex: number, columnIndex: number): boolean {
  return joinRelationships.value.some(
    rel => rel.tableIndex === tableIndex && rel.columnIndex === columnIndex
  );
}

// Check if dark mode is enabled
const updateDarkMode = () => {
  isDarkMode.value = document.documentElement.classList.contains('dark') || 
                     document.body.classList.contains('dark') ||
                     window.matchMedia('(prefers-color-scheme: dark)').matches;
};

// Table structure for visualization
interface Column {
  name: string;
  type: string;
  isPrimaryKey: boolean;
  isForeignKey: boolean;
  isSelected?: boolean;
  table?: string;
  value?: string; // Value for INSERT queries
  sourceTable?: string;  // Tabela de origem
  sourceColumn?: string; // Coluna de origem
  originalName?: string; // Nome original quando há alias
}

interface Table {
  name: string;
  columns: Column[];
  x: number;
  y: number;
  width: number; // Dynamic width
  height: number; // Dynamic height
}

interface ResultBox {
  columns: Column[];
  x: number;
  y: number;
  width: number;
}

interface InsertVisualization {
  x: number;
  y: number;
  path: string;
  labelX: number;
  labelY: number;
}

interface Relation {
  from: {
    table: string;
    column: string;
  };
  to: {
    table: string;
    column: string;
  };
  path: string;
  arrowPoints: string;
  joinType: string;
  labelX: number;
  labelY: number;
  descriptionX: number;
  descriptionY: number;
  startX: number;
  startY: number;
  endX: number;
  endY: number;
  isPkFkRelation: boolean;
}

const parsedTables = ref<Table[]>([]);
const parsedRelations = ref<Relation[]>([]);
const resultBox = ref<ResultBox | null>(null);
const insertVisualization = ref<InsertVisualization>({
  x: 0,
  y: 0,
  path: '',
  labelX: 0,
  labelY: 0
});

// Format value for display (truncate if too long, format strings)
const formatValue = (value: string): string => {
  // If it's a string (starts with quote), format it and truncate if needed
  if (value.startsWith("'") || value.startsWith('"')) {
    const cleaned = value.substring(1, value.length - 1);
    return cleaned.length > 15 ? `'${cleaned.substring(0, 15)}...'` : value;
  }
  
  // For other values (numbers, etc.)
  return value.length > 15 ? `${value.substring(0, 15)}...` : value;
};

const parseQuery = (sql: string) => {
  try {
    // Reset current data
    parsedTables.value = [];
    parsedRelations.value = [];
    resultBox.value = null;
    joinRelationships.value = [];
    
    // Parse SQL using our custom parser
    const result = parse(sql);
    queryType.value = result.queryType;
    
    if (result.error) {
      console.error('SQL parsing error:', result.error);
      return;
    }
    
    // Transform parsed data into diagram format
    const tables: Table[] = [];
    const relations: Relation[] = [];
    
    // Process tables
    result.tables.forEach((sqlTable, i) => {
      // Convert columns to the format needed for visualization
      const columns: Column[] = sqlTable.columns.map(col => {
        // Detect column type based on name patterns
        let type = 'varchar';
        
        // Try to infer type from column name
        const name = col.name.toLowerCase();
        if (name === 'id' || name.endsWith('_id') || name.includes('id_')) {
          type = 'int';
        } else if (name.includes('date') || name.includes('time')) {
          type = 'datetime';
        } else if (name.includes('price') || name.includes('amount') || name.includes('total')) {
          type = 'decimal';
        } else if (name === 'active' || name === 'enabled' || name === 'status') {
          type = 'boolean';
        }
        
        // Detect primary and foreign keys
        const isPrimaryKey = 
          col.isPrimaryKey || 
          col.name.toLowerCase() === 'id' || 
          sqlTable.name.toLowerCase() + '_id' === col.name.toLowerCase();
        
        const isForeignKey = 
          col.isForeignKey || 
          (col.name.toLowerCase().endsWith('_id') && col.name.toLowerCase() !== 'id');
        
        return {
          name: col.alias || col.name,
          type: type,
          isPrimaryKey: isPrimaryKey,
          isForeignKey: isForeignKey,
          isSelected: !!col.isSelected,
          value: col.value,
          sourceTable: sqlTable.name,
          sourceColumn: col.name,
          originalName: col.alias ? col.name : undefined
        };
      });
      
      // If no columns were found, add placeholder columns
      if (columns.length === 0) {
        columns.push(
          { name: 'id', type: 'int', isPrimaryKey: true, isForeignKey: false },
          { name: 'name', type: 'varchar', isPrimaryKey: false, isForeignKey: false }
        );
      }
      
      // Calculate optimal width based on content
      let minWidth = 300; // Default minimum width
      const tableNameWidth = Math.max(sqlTable.name.length * 10, minWidth);
      
      // Find the widest column name
      const maxColNameWidth = columns.reduce((max, col) => {
        // Calculate width needed for column name + badges
        let columnWidth = col.name.length * 8;
        if (col.isPrimaryKey) columnWidth += 30; // Space for PK badge
        if (col.isForeignKey) columnWidth += 30; // Space for FK badge
        return Math.max(max, columnWidth);
      }, 0);
      
      // Find the widest value
      const maxValueWidth = columns.reduce((max, col) => {
        if (col.value) {
          return Math.max(max, col.value.length * 8);
        }
        return max;
      }, 0);
      
      // Calculate width based on column content and type display
      const optimalWidth = Math.max(
        tableNameWidth, 
        maxColNameWidth + 100, // Space for column name + type
        maxColNameWidth + maxValueWidth + 140 // Space for name, type and value
      );
      
      // Use optimal width but ensure it's not too narrow or too wide
      const tableWidth = Math.max(minWidth, Math.min(600, optimalWidth));
      
      // Calculate height based on number of columns
      const tableHeight = 40 + columns.length * 30; // Header (40px) + rows
      
      tables.push({
        name: sqlTable.name,
        columns,
        x: 50 + (i % 2) * (tableWidth + 100), // More horizontal spacing
        y: 50 + Math.floor(i / 2) * (tableHeight + 80), // More vertical spacing
        width: tableWidth,
        height: tableHeight
      });
    });
    
    // Process joins/relationships for SELECT queries
    if (result.queryType === 'SELECT') {
      result.joins.forEach(join => {
        if (join.condition) {
          const { leftTable, leftColumn, rightTable, rightColumn } = join.condition;
          
          // Find the source and target tables - usando toLowerCase para melhorar a correspondência
          const sourceTableIndex = tables.findIndex(t => 
            t.name.toLowerCase() === leftTable.toLowerCase() || 
            t.name.toLowerCase() === rightTable.toLowerCase()
          );
          
          const targetTableIndex = tables.findIndex(t => 
            (t.name.toLowerCase() === rightTable.toLowerCase() && sourceTableIndex !== -1 && 
             tables[sourceTableIndex].name.toLowerCase() !== rightTable.toLowerCase()) || 
            (t.name.toLowerCase() === leftTable.toLowerCase() && sourceTableIndex !== -1 && 
             tables[sourceTableIndex].name.toLowerCase() !== leftTable.toLowerCase())
          );
          
          if (sourceTableIndex !== -1 && targetTableIndex !== -1) {
            const sourceTable = tables[sourceTableIndex];
            const targetTable = tables[targetTableIndex];
            
            // Find column indices for highlighting
            const sourceColumnName = tables[sourceTableIndex].name.toLowerCase() === leftTable.toLowerCase() 
              ? leftColumn.toLowerCase() 
              : rightColumn.toLowerCase();
            
            const targetColumnName = tables[targetTableIndex].name.toLowerCase() === rightTable.toLowerCase() 
              ? rightColumn.toLowerCase() 
              : leftColumn.toLowerCase();
            
            // Busca mais flexível da coluna por nome (tolerante a _id)
            const sourceColumnIndex = sourceTable.columns.findIndex(c => 
              c.name.toLowerCase() === sourceColumnName ||
              c.name.toLowerCase() === sourceColumnName.replace('_id', '') + '_id'
            );
            
            const targetColumnIndex = targetTable.columns.findIndex(c => 
              c.name.toLowerCase() === targetColumnName ||
              c.name.toLowerCase() === targetColumnName.replace('_id', '') + '_id'
            );
            
            console.log(`Joining ${sourceTable.name}.${sourceColumnName} to ${targetTable.name}.${targetColumnName}`, 
                       {sourceColumnIndex, targetColumnIndex});
            
            // Add join columns to the highlighting array
            if (sourceColumnIndex !== -1) {
              joinRelationships.value.push({
                tableIndex: sourceTableIndex,
                columnIndex: sourceColumnIndex
              });
            }
            
            if (targetColumnIndex !== -1) {
              joinRelationships.value.push({
                tableIndex: targetTableIndex,
                columnIndex: targetColumnIndex
              });
            }
            
            // Check if this is a primary key to foreign key relationship
            let isPkFkRelation = false;
            let sourcePk = sourceTable.columns[sourceColumnIndex]?.isPrimaryKey;
            let targetFk = targetTable.columns[targetColumnIndex]?.isForeignKey;
            let targetPk = targetTable.columns[targetColumnIndex]?.isPrimaryKey;
            let sourceFk = sourceTable.columns[sourceColumnIndex]?.isForeignKey;
            
            if ((sourcePk && targetFk) || (targetPk && sourceFk)) {
              isPkFkRelation = true;
            }
            
            // If column has "_id" in it, it's likely a foreign key to a primary key
            if (!isPkFkRelation && (
              sourceColumnName.includes('_id') || 
              targetColumnName.includes('_id'))
            ) {
              isPkFkRelation = true;
            }
            
            // Calcular os pontos de conexão exatos - alinhados com as colunas específicas
            const sourceY = sourceTable.y + 40 + (sourceColumnIndex !== -1 
              ? sourceColumnIndex * 30 + 15  // Middle of the column
              : sourceTable.height / 2); // Middle of the table
            
            const targetY = targetTable.y + 40 + (targetColumnIndex !== -1 
              ? targetColumnIndex * 30 + 15  // Middle of the column
              : targetTable.height / 2); // Middle of the table
            
            // Ajustar os pontos iniciais/finais para ficarem nas bordas das tabelas
            const startX = sourceTable.x + sourceTable.width; // Right edge of source table
            const endX = targetTable.x; // Left edge of target table
            
            // Criar um caminho de conexão curvo mais pronunciado para melhor visibilidade
            const pathCurveStrength = 150; // Higher = more curved
            const path = `M ${startX} ${sourceY} C ${startX + pathCurveStrength} ${sourceY}, ${endX - pathCurveStrength} ${targetY}, ${endX} ${targetY}`;
            
            // Calcular o ponto médio para o rótulo do tipo de join
            const midX = (startX + endX) / 2;
            const midY = (sourceY + targetY) / 2 - 35; // Mais acima da linha
            
            // Calcular o ponto para a descrição da relação
            const descX = midX;
            const descY = (sourceY + targetY) / 2 + 40; // Mais abaixo da linha
            
            relations.push({
              from: { 
                table: tables[sourceTableIndex].name.toLowerCase() === leftTable.toLowerCase() ? leftTable : rightTable, 
                column: tables[sourceTableIndex].name.toLowerCase() === leftTable.toLowerCase() ? leftColumn : rightColumn
              },
              to: { 
                table: tables[targetTableIndex].name.toLowerCase() === rightTable.toLowerCase() ? rightTable : leftTable, 
                column: tables[targetTableIndex].name.toLowerCase() === rightTable.toLowerCase() ? rightColumn : leftColumn
              },
              path,
              arrowPoints: "", // Not used with marker-end
              joinType: join.type || 'JOIN',
              labelX: midX,
              labelY: midY,
              descriptionX: descX,
              descriptionY: descY,
              startX, // Starting point X
              startY: sourceY, // Starting point Y
              endX, // End point X
              endY: targetY, // End point Y
              isPkFkRelation // Flag to indicate PK-FK relationship
            });
          }
        }
      });
      
      // Create result box with all selected columns
      const selectedColumns = result.selectedColumns.map(col => {
        // Separar o nome original e o alias para referência
        const hasAlias = col.alias && col.alias !== col.name;
        
        return {
          name: col.alias || col.name, // Nome com alias quando disponível
          originalName: hasAlias ? col.name : undefined, // Nome original se tiver alias
          type: 'unknown',
          isPrimaryKey: false,
          isForeignKey: false,
          table: col.table
        };
      });
      
      if (selectedColumns.length > 0) {
        // Calculate optimal width for result box
        const maxColNameWidth = selectedColumns.reduce((max, col) => {
          // Account for column name + table name if present
          let width = col.name.length * 8;
          if (col.table) width += col.table.length * 6 + 20; // Space for table name in parentheses
          return Math.max(max, width);
        }, 0);
        
        // Calculate width but not less than 250px or more than 400px
        const resultBoxWidth = Math.max(250, Math.min(400, maxColNameWidth + 60));
        
        // If we have multiple tables, position the result box at the bottom center
        let resultBoxX = 50;
        let resultBoxY = 50;
        
        if (tables.length > 0) {
          // Find the maximum Y position of all tables to place result box below
          const maxY = Math.max(...tables.map(t => t.y + t.height));
          // Find the center X position of all tables
          const totalX = tables.reduce((sum, t) => sum + t.x + t.width/2, 0);
          const avgX = totalX / tables.length - resultBoxWidth/2; // Center aligned with tables
          
          resultBoxX = Math.max(50, avgX);
          resultBoxY = maxY + 80; // More space below the lowest table
        }
        
        // Mapear colunas selecionadas para seus originais
        const enhancedSelectedColumns = selectedColumns.map(col => {
          // Se a coluna tem uma tabela de origem, mantemos essa informação
          const sourceInfo = {
            ...col,
            sourceTable: col.table,
            sourceColumn: col.name
          };
          return sourceInfo;
        });
        
        resultBox.value = {
          columns: enhancedSelectedColumns,
          x: resultBoxX,
          y: resultBoxY,
          width: resultBoxWidth
        };
      }
    } 
    // Process INSERT visualization
    else if (result.queryType === 'INSERT' && tables.length > 0) {
      const table = tables[0]; // In INSERT, we only have one table
      
      // Create INSERT arrow/visualization with dynamic position
      const arrowStartX = table.x - 100;
      const arrowStartY = table.y + table.height / 2; // Center of the table
      const arrowEndX = table.x;
      const arrowEndY = table.y + table.height / 2;
      
      insertVisualization.value = {
        x: arrowStartX,
        y: arrowStartY,
        path: `M ${arrowStartX} ${arrowStartY} L ${arrowEndX} ${arrowEndY}`,
        labelX: arrowStartX - 40,
        labelY: arrowStartY
      };
    }
    
    parsedTables.value = tables;
    parsedRelations.value = relations;
    
    // Update diagram size after all elements are added
    updateDiagramSize();
  } catch (error) {
    console.error('Error parsing SQL query:', error);
  }
};

// Zoom handlers
const zoomIn = () => {
  scale.value = Math.min(scale.value + 0.1, 2);
};

const zoomOut = () => {
  scale.value = Math.max(scale.value - 0.1, 0.5);
};

const resetZoom = () => {
  scale.value = 1;
  offsetX.value = 0;
  offsetY.value = 0;
};

const handleZoom = (event: WheelEvent) => {
  if (event.deltaY < 0) {
    zoomIn();
  } else {
    zoomOut();
  }
};

// Panning handlers
const startDrag = (event: MouseEvent) => {
  if (event.button !== 0) return; // Only left mouse button
  isDragging.value = true;
  startX.value = event.clientX;
  startY.value = event.clientY;
};

const drag = (event: MouseEvent) => {
  if (!isDragging.value) return;
  const dx = event.clientX - startX.value;
  const dy = event.clientY - startY.value;
  offsetX.value += dx / scale.value;
  offsetY.value += dy / scale.value;
  startX.value = event.clientX;
  startY.value = event.clientY;
};

const endDrag = () => {
  isDragging.value = false;
};

// Add this function to adjust SVG size based on content
const updateDiagramSize = () => {
  if (!parsedTables.value.length) return;

  // Find the furthest right and bottom elements
  const maxX = Math.max(
    ...parsedTables.value.map(t => t.x + t.width),
    resultBox.value ? resultBox.value.x + resultBox.value.width : 0
  );
  
  const maxY = Math.max(
    ...parsedTables.value.map(t => t.y + t.height),
    resultBox.value ? resultBox.value.y + getResultBoxHeight(resultBox.value.columns.length) : 0
  );
  
  // Set SVG size with padding
  width.value = Math.max(800, maxX + 100);
  height.value = Math.max(600, maxY + 100);
};

// Atualizar função getResultRowY para calcular corretamente a posição Y
function getResultRowY(colIndex: number): number {
  // Posição Y do centro da linha (64px é o início das linhas após o cabeçalho, 15px é metade da altura da linha)
  return 64 + colIndex * 30 + 15;
}

// Função para obter o caminho de conexão entre tabela e resultado
function getResultConnectionPath(column: Column, colIndex: number) {
  // Obter as coordenadas iniciais (na tabela)
  const startPoint = getResultConnectionStartPoint(column);
  
  // Obter as coordenadas do ponto final (no resultado)
  const endPointX = -10; // 10px à esquerda do início do resultado
  const endPointY = getResultRowY(colIndex);
  
  // Ajustar os pontos de controle para suavizar a curva
  const controlPointX1 = startPoint.x + (endPointX - startPoint.x) * 0.5;
  const controlPointY1 = startPoint.y;
  const controlPointX2 = endPointX - 20;
  const controlPointY2 = endPointY;
  
  // Retornar o caminho de conexão como uma curva Bezier
  return `M ${startPoint.x} ${startPoint.y} C ${controlPointX1} ${controlPointY1}, ${controlPointX2} ${controlPointY2}, ${endPointX} ${endPointY}`;
}

// Função auxiliar para obter o ponto de partida da conexão
function getResultConnectionStartPoint(column: Column): {x: number, y: number} {
  const sourceTableIndex = parsedTables.value.findIndex(
    t => t.name.toLowerCase() === (column.table || '').toLowerCase()
  );
  
  if (sourceTableIndex === -1) {
    return { x: 0, y: 0 };
  }
  
  const sourceTable = parsedTables.value[sourceTableIndex];
  
  // Encontrar a coluna de origem na tabela fonte
  const sourceColumnName = column.sourceColumn || column.originalName || column.name;
  const sourceColumnIndex = sourceTable.columns.findIndex(
    c => c.name.toLowerCase() === sourceColumnName.toLowerCase()
  );
  
  // Ponto de partida: lado direito da tabela de origem
  const startX = sourceTable.x + sourceTable.width;
  // Se encontramos a coluna exata, use sua posição, caso contrário, use o meio da tabela
  const startY = sourceTable.y + 40 + (sourceColumnIndex !== -1 
    ? sourceColumnIndex * 30 + 15  // Meio da coluna específica
    : (sourceTable.columns.length / 2) * 30); // Meio da tabela
  
  return { x: startX, y: startY };
}

// Função para calcular a altura do box de resultados com base no número de colunas
function getResultBoxHeight(columnCount: number): number {
  const headerHeight = 40;    // Altura do cabeçalho principal
  const columnHeaderHeight = 24; // Altura do cabeçalho das colunas
  const rowHeight = 30;       // Altura de cada linha de coluna
  
  // Se não tiver colunas, mostra apenas o cabeçalho
  if (columnCount <= 0) {
    return headerHeight + columnHeaderHeight;
  }
  
  return headerHeight + columnHeaderHeight + (columnCount * rowHeight);
}

onMounted(() => {
  if (diagramContainer.value) {
    width.value = diagramContainer.value.clientWidth;
    height.value = Math.max(500, diagramContainer.value.clientHeight);
    
    // Add mouse event listeners for panning
    svgElement.value?.addEventListener('mousedown', startDrag);
    window.addEventListener('mousemove', drag);
    window.addEventListener('mouseup', endDrag);
  }
  
  // Parse initial query
  if (props.query) {
    parseQuery(props.query);
    updateDiagramSize(); // Update SVG size after parsing
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
  svgElement.value?.removeEventListener('mousedown', startDrag);
  window.removeEventListener('mousemove', drag);
  window.removeEventListener('mouseup', endDrag);
  
  // Remove dark mode listeners
  window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', updateDarkMode);
});

// Update the watch for query changes to also update size
watch(() => props.query, (newQuery) => {
  if (newQuery) {
    parseQuery(newQuery);
    updateDiagramSize(); // Update SVG size after parsing
  }
});
</script>

<style scoped>
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
}

.table-node {
  cursor: move;
}

.relation-line {
  pointer-events: none;
}

.relation-line text {
  pointer-events: none;
  user-select: none;
}

/* Path styling for joins */
.light-join-path {
  stroke: #4f46e5 !important; /* indigo-600 - mais escuro */
  stroke-opacity: 0.9;
  stroke-dasharray: none;
}

.dark-join-path {
  stroke: #818cf8 !important; /* indigo-400 */
  stroke-opacity: 1;
  stroke-dasharray: none;
}

/* Glow effects */
.filter-glow-light {
  filter: url(#glow-light);
}

.filter-glow-dark {
  filter: url(#glow-dark);
}

.filter-glow {
  filter: url(#text-glow);
}

/* Ensure SVG elements are visible with proper z-index */
svg {
  overflow: visible;
}

svg marker {
  overflow: visible;
}

/* Dashed line styling */
.stroke-dasharray-2 {
  stroke-dasharray: 2 2;
}

/* Result connection styling */
.result-connection {
  stroke-dasharray: 4 2;
  pointer-events: none;
}

/* Dark gray alternate background for result rows */
.fill-gray-750 {
  fill: #2d3748;
}
</style> 