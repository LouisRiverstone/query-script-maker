<template>
  <div class="database-diagram-container">
    <div class="flex flex-row h-full">
      <!-- Table List Sidebar -->
      <div class="w-72 flex-shrink-0 bg-white dark:bg-gray-800 border-r border-gray-200 dark:border-gray-700 overflow-hidden flex flex-col">
        <h3 class="font-semibold text-gray-700 dark:text-gray-300 p-3 border-b border-gray-200 dark:border-gray-700">
          Tables ({{ filteredNodes.length }}/{{ nodes.length }})
        </h3>
        
        <!-- Search input -->
        <div class="px-2 pt-2 pb-1">
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search tables..."
              class="w-full px-3 py-2 bg-gray-100 dark:bg-gray-700 text-gray-700 dark:text-gray-200 rounded-md text-sm focus:outline-none focus:ring-2 focus:ring-indigo-500 dark:focus:ring-indigo-400 border border-gray-200 dark:border-gray-600"
            />
            <span v-if="searchQuery" @click="searchQuery = ''" class="absolute right-3 top-2.5 cursor-pointer text-gray-500 dark:text-gray-400 hover:text-gray-700 dark:hover:text-gray-300">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </span>
            <span v-else class="absolute right-3 top-2.5 text-gray-500 dark:text-gray-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </span>
          </div>
        </div>
        
        <div class="overflow-y-auto flex-1 p-2 bg-white dark:bg-gray-800">
          <div v-if="filteredNodes.length === 0" class="text-center py-4 text-gray-500 dark:text-gray-400 text-sm italic">
            No tables match your search
          </div>
          <div v-else class="space-y-1">
            <div v-for="node in filteredNodes" :key="node.id" class="table-list-item mb-3">
              <div 
                class="cursor-pointer p-2 text-sm rounded-md bg-gray-50 dark:bg-gray-800 hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors border border-gray-200 dark:border-gray-600"
                :class="{'bg-indigo-50 dark:bg-indigo-900/30 border-indigo-300 dark:border-indigo-600': selectedTable === node.id}"
              >
                <div class="font-medium text-gray-800 dark:text-gray-200 flex items-center justify-between" @click="toggleExpandTable(node.id)">
                  <div class="flex items-center">
                    <span class="text-indigo-600 dark:text-indigo-400 mr-1">
                      <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
                      </svg>
                    </span>
                    <span :class="{'bg-yellow-100 dark:bg-yellow-900/50 px-1 rounded': tableMatchesSearch(node)}">{{ node.data.label }}</span>
                    <span class="ml-1 text-xs text-gray-500 dark:text-gray-400">({{ node.data.columns.length }})</span>
                  </div>
                  
                  <!-- Expand/Collapse Icon -->
                  <button class="text-gray-500 dark:text-gray-400 focus:outline-none">
                    <svg v-if="expandedTables.includes(node.id)" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                    </svg>
                    <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                    </svg>
                  </button>
                </div>
                
                <!-- Column list (collapsible) -->
                <div v-if="expandedTables.includes(node.id)" class="mt-2 ml-2 space-y-0.5 max-h-60 overflow-y-auto border-l-2 border-gray-200 dark:border-gray-600 pl-2">
                  <div 
                    v-for="column in node.data.columns" 
                    :key="`${node.id}-${column.name}`"
                    class="text-xs p-1 rounded-md hover:bg-gray-200 dark:hover:bg-gray-600 cursor-pointer transition-colors flex flex-col bg-transparent dark:bg-transparent"
                    :class="{
                      'bg-indigo-50 dark:bg-indigo-900/30': selectedColumn === column.name && selectedTable === node.id
                    }"
                    @click.stop="focusNodeColumn(node.id, column.name)"
                  >
                    <div class="flex items-center">
                      <span v-if="column.isPrimary" class="text-yellow-600 dark:text-yellow-400 mr-1" title="Primary Key">🔑</span>
                      <span v-else-if="column.isForeign" class="text-blue-600 dark:text-blue-400 mr-1" title="Foreign Key">🔗</span>
                      <span v-else class="mr-1 w-3.5"></span>
                      <span :class="{'bg-yellow-100 dark:bg-yellow-900/50 px-1 rounded': columnMatchesSearch(column)}">{{ column.name }}</span>
                    </div>
                    <div class="flex items-center mt-0.5 ml-4 text-gray-500 dark:text-gray-400">
                      <span class="inline-block px-1 py-0.5 bg-gray-100 dark:bg-gray-700 rounded text-[10px]">{{ column.type }}</span>
                      <span v-if="column.nullable === 'YES'" class="ml-1 text-[10px]">nullable</span>
                      <span v-if="column.default" class="ml-1 text-[10px]" title="Default value">
                        default: {{ column.default.length > 10 ? column.default.substring(0, 10) + '...' : column.default }}
                      </span>
                    </div>
                    <div v-if="column.extra" class="ml-4 text-[10px] text-gray-500 dark:text-gray-400 mt-0.5">
                      {{ column.extra }}
                    </div>
                  </div>
                </div>

                <!-- Focus button -->
                <div v-if="selectedTable !== node.id" class="mt-2 flex justify-end">
                  <button 
                    @click.stop="focusNode(node.id)" 
                    class="text-xs px-2 py-1 bg-indigo-50 dark:bg-indigo-900/20 text-indigo-600 dark:text-indigo-400 rounded hover:bg-indigo-100 dark:hover:bg-indigo-900/40 transition-colors"
                  >
                    Focus
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Diagram Content -->
      <div class="flex-1 flex flex-col overflow-hidden">
        <div ref="diagramContainer" class="diagram-content flex-1 bg-white dark:bg-gray-800 rounded-lg p-2 overflow-hidden">
          <VueFlow v-if="nodes.length > 0"
            :default-zoom="1"
            :min-zoom="0.2"
            :max-zoom="4"
            :nodes="nodes"
            :edges="edges"
            class="vue-flow-wrapper h-full"
            fit-view-on-init
            @node-click="onNodeClick"
          >
            <!-- Custom Node Types -->
            <template #node-table="nodeProps">
              <div 
                class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm p-2"
                :class="{'ring-2 ring-indigo-500 dark:ring-indigo-400': selectedTable === nodeProps.id}"
              >
                <div class="table-header bg-indigo-100 dark:bg-indigo-900 p-2 mb-2 rounded font-bold text-center">
                  {{ nodeProps.data.label }}
                </div>
                <div class="table-columns">
                  <div 
                    v-for="(column, index) in nodeProps.data.columns" 
                    :key="index" 
                    class="column-item flex items-center p-1 border-b border-gray-200 dark:border-gray-700 last:border-0"
                    :class="{
                      'bg-yellow-50 dark:bg-yellow-900': column.isPrimary,
                      'bg-indigo-50 dark:bg-indigo-900/30': selectedColumn === column.name && selectedTable === nodeProps.id
                    }"
                  >
                    <div class="flex items-center">
                      <span v-if="column.isPrimary" class="text-yellow-600 dark:text-yellow-400 mr-1" title="Primary Key">🔑</span>
                      <span v-else-if="column.isForeign" class="text-blue-600 dark:text-blue-400 mr-1" title="Foreign Key">🔗</span>
                      <span v-else class="mr-1 w-3.5"></span>
                    </div>
                    <div class="flex-1 text-left whitespace-nowrap pr-2">
                      {{ column.name }}
                      <span v-if="column.nullable === 'YES'" class="text-gray-400 dark:text-gray-500 text-xs italic ml-1" title="Nullable">(null)</span>
                    </div>
                    <div class="text-gray-500 dark:text-gray-400 text-xs whitespace-nowrap flex items-center">
                      <span>{{ column.type }}</span>
                      <span v-if="column.default" class="ml-1 text-gray-400 dark:text-gray-500 italic" title="Default value">
                        = {{ column.default.length > 8 ? column.default.substring(0, 8) + '...' : column.default }}
                      </span>
                      <span v-if="column.extra" class="ml-1 px-1 py-0.5 bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 rounded text-[10px]" title="Extra attribute">
                        {{ column.extra }}
                      </span>
                    </div>
                  </div>
                </div>
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
            
            <!-- Add Minimap -->
            <MiniMap :node-color="getNodeColor" />
            
            <!-- Add panel with legend -->
            <Panel :position="'top-right'" class="legend-panel">
              <div class="legend bg-white dark:bg-gray-700 p-2 rounded shadow-md">
                <h4 class="text-sm font-semibold text-gray-700 dark:text-white mb-1">Legend</h4>
                <div class="legend-item flex items-center mb-1">
                  <div class="w-3 h-3 bg-yellow-500 dark:bg-yellow-400 rounded-sm mr-2"></div>
                  <span class="text-xs text-gray-600 dark:text-gray-300">Primary Key</span>
                </div>
                <div class="legend-item flex items-center mb-1">
                  <div class="w-3 h-3 bg-blue-500 dark:bg-blue-400 rounded-sm mr-2"></div>
                  <span class="text-xs text-gray-600 dark:text-gray-300">Foreign Key</span>
                </div>
                <div class="legend-item flex items-center mb-1">
                  <span class="text-xs text-gray-400 dark:text-gray-500 italic mr-2">(null)</span>
                  <span class="text-xs text-gray-600 dark:text-gray-300">Nullable Field</span>
                </div>
                <div class="legend-item flex items-center">
                  <span class="mr-2 px-1 py-0.5 bg-purple-100 dark:bg-purple-900/30 text-purple-600 dark:text-purple-400 rounded text-[10px]">extra</span>
                  <span class="text-xs text-gray-600 dark:text-gray-300">Extra Attribute</span>
                </div>
              </div>
            </Panel>
          </VueFlow>
          <div v-else class="empty-diagram flex items-center justify-center h-full">
            <p class="text-gray-500 dark:text-gray-400 text-center">No database structure found</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, watch, onUnmounted } from 'vue';
import Button from './Button.vue';
import { 
  VueFlow, 
  useVueFlow,
  Panel,
  Node,
  Edge,
  MarkerType,
  NodeMouseEvent,
  useZoomPanHelper
} from '@vue-flow/core';
import '@vue-flow/core/dist/style.css';
import '@vue-flow/core/dist/theme-default.css';

// Additional components
import { Background } from '@vue-flow/background';
import { Controls } from '@vue-flow/controls';
import { MiniMap } from '@vue-flow/minimap';

// Layout utilities
import dagre from 'dagre';

// Additional styles
import '@vue-flow/minimap/dist/style.css';
import '@vue-flow/controls/dist/style.css';

// Types for database structure
interface Column {
  name: string;
  type: string;
  nullable: string;
  key: string;
  default: string;
  extra: string;
  isPrimary: boolean;
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

const props = defineProps<{
  databaseStructure: string;
  isDarkMode?: boolean;
}>();

const diagramContainer = ref<HTMLElement | null>(null);
const localIsDarkMode = ref(false);

// Use the provided isDarkMode prop if available, otherwise detect it locally
const isDarkMode = computed(() => {
  return props.isDarkMode !== undefined ? props.isDarkMode : localIsDarkMode.value;
});

// Elements for Vue Flow
const nodes = ref<any[]>([]);
const edges = ref<any[]>([]);

// Selected table/column tracking
const selectedTable = ref<string | null>(null);
const selectedColumn = ref<string | null>(null);

// Get Vue Flow instance
const { zoomIn, zoomOut, fitView } = useVueFlow({
  nodes: [],
  edges: []
});

// Get zoom pan helper for custom focus
const { setCenter } = useZoomPanHelper();

// Function to focus on a specific node
const focusNode = (nodeId: string) => {
  const node = nodes.value.find(node => node.id === nodeId);
  if (node) {
    // Select the table
    selectedTable.value = nodeId;
    selectedColumn.value = null;
    
    // Ensure the table is expanded
    if (!expandedTables.value.includes(nodeId)) {
      expandedTables.value.push(nodeId);
    }
    
    // Center the view on the node with zoom
    const x = node.position.x + node.style.width / 2;
    const y = node.position.y + 100; // Add some offset for better centering
    
    setCenter(x, y, { zoom: 1.5, duration: 800 });
  }
};

// Function to focus on a specific column within a node
const focusNodeColumn = (nodeId: string, columnName: string) => {
  // Select the table and column
  selectedTable.value = nodeId;
  selectedColumn.value = columnName;
  
  // First focus the node
  focusNode(nodeId);
};

// Handle node click in the diagram
const onNodeClick = (event: NodeMouseEvent) => {
  const nodeId = event.node.id as string;
  focusNode(nodeId);
};

// Function to check if dark mode is enabled
const updateDarkMode = () => {
  localIsDarkMode.value = document.documentElement.classList.contains('dark') || 
                  document.body.classList.contains('dark') ||
                  window.matchMedia('(prefers-color-scheme: dark)').matches;
};

// Node color function for minimap
const getNodeColor = (node: Node) => {
  if (node.id === selectedTable.value) {
    return '#4f46e5'; // Highlight selected table
  }
  return '#6b7280'; // Default gray color
};

// Parse database structure and create diagram
const parseDbStructure = (structureStr: string) => {
  try {
    if (!structureStr) {
      nodes.value = [];
      edges.value = [];
      return;
    }

    const structure: DatabaseStructure = JSON.parse(structureStr);
    
    if (!structure.tables || !structure.tables.length) {
      nodes.value = [];
      edges.value = [];
      return;
    }
    
    const diagramNodes: Node[] = [];
    const diagramEdges: Edge[] = [];
    
    // Track foreign keys for edge creation
    const foreignKeyMap = new Map<string, ForeignKey[]>();
    
    // Process tables as nodes
    structure.tables.forEach((table) => {
      // Keep track of foreign key columns
      const foreignKeyColumns = new Set(table.foreignKeys.map(fk => fk.columnName));
      
      // Prepare columns with visual indicators
      const columns = table.columns.map(col => ({
        ...col,
        isForeign: foreignKeyColumns.has(col.name)
      }));
      
      // Create node for the table
      diagramNodes.push({
        id: table.name,
        type: 'table',
        data: { 
          label: table.name,
          columns
        },
        position: { x: 0, y: 0 } // Will be positioned by layout algorithm
      });
      
      // Store foreign keys for edge creation
      if (table.foreignKeys.length > 0) {
        foreignKeyMap.set(table.name, table.foreignKeys);
      }
    });
    
    // Process foreign keys as edges
    foreignKeyMap.forEach((foreignKeys, tableName) => {
      foreignKeys.forEach(fk => {
        const edgeId = `${tableName}-${fk.columnName}-to-${fk.referencedTable}-${fk.referencedColumn}`;
        
        diagramEdges.push({
          id: edgeId,
          source: tableName,
          target: fk.referencedTable,
          label: `${fk.columnName} → ${fk.referencedColumn}`,
          animated: false,
          style: { stroke: '#3b82f6' },
          labelBgStyle: { fill: '#3b82f6', color: '#ffffff', fillOpacity: 0.7 },
          labelStyle: { fill: '#ffffff', fontWeight: 500, fontSize: 10 },
          markerEnd: {
            type: MarkerType.ArrowClosed,
            color: '#3b82f6'
          }
        });
      });
    });
    
    // Apply layout
    const positionedNodes = getLayoutedElements(diagramNodes, diagramEdges);
    
    nodes.value = positionedNodes;
    edges.value = diagramEdges;
    
    // Reset selection
    selectedTable.value = null;
    selectedColumn.value = null;
  } catch (error) {
    console.error('Error parsing database structure:', error);
    nodes.value = [];
    edges.value = [];
  }
};

// Calculate layout using dagre
const getLayoutedElements = (nodes: Node[], edges: Edge[], direction = 'TB') => {
  const dagreGraph = new dagre.graphlib.Graph();
  dagreGraph.setDefaultEdgeLabel(() => ({}));
  
  const nodeWidth = 220;
  const nodeHeight = 180;
  
  // Set graph direction and spacing
  dagreGraph.setGraph({ rankdir: direction, ranksep: 150, nodesep: 100 });
  
  // Add nodes with dimensions
  nodes.forEach((node) => {
    // Calculate height based on number of columns (each being ~32px)
    const height = 60 + (node.data.columns.length * 32);
    
    dagreGraph.setNode(node.id, { width: nodeWidth, height });
  });
  
  // Add edges
  edges.forEach((edge) => {
    dagreGraph.setEdge(edge.source, edge.target);
  });
  
  // Calculate layout
  dagre.layout(dagreGraph);
  
  // Apply calculated positions to nodes
  return nodes.map((node) => {
    const nodeWithPosition = dagreGraph.node(node.id);
    
    return {
      ...node,
      position: {
        x: nodeWithPosition.x - nodeWidth / 2,
        y: nodeWithPosition.y - nodeWithPosition.height / 2
      },
      style: {
        width: nodeWidth
      }
    };
  });
};

// Watch for dark mode changes
onMounted(() => {
  updateDarkMode();
  window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', updateDarkMode);
});

onUnmounted(() => {
  // Clean up event listener
  window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', updateDarkMode);
});

// Watch for database structure changes
watch(() => props.databaseStructure, (newValue) => {
  parseDbStructure(newValue);
}, { immediate: true });

// Watch for dark mode changes from parent
watch(() => props.isDarkMode, (newValue) => {
  if (newValue !== undefined) {
    // If the parent is controlling dark mode, we don't need to update it locally
    // This prevents conflicts between parent and local dark mode detection
  }
});

// Filter nodes based on search query
const searchQuery = ref('');
const filteredNodes = computed(() => {
  if (!nodes.value || !searchQuery.value) {
    return nodes.value;
  }
  
  const query = searchQuery.value.toLowerCase();
  
  return nodes.value.filter(node => {
    // Check if table name matches
    if (node.data.label.toLowerCase().includes(query)) {
      return true;
    }
    
    // Check if any column name matches
    return node.data.columns.some((column: { name: string }) => 
      column.name.toLowerCase().includes(query)
    );
  });
});

// Expand/collapse table functionality
const expandedTables = ref<string[]>([]);
const toggleExpandTable = (tableId: string) => {
  if (expandedTables.value.includes(tableId)) {
    expandedTables.value = expandedTables.value.filter(id => id !== tableId);
  } else {
    expandedTables.value.push(tableId);
  }
};

// Helper to check if table name matches search
const tableMatchesSearch = (node: any) => {
  if (!searchQuery.value) return false;
  return node.data.label.toLowerCase().includes(searchQuery.value.toLowerCase());
};

// Helper to check if column matches search
const columnMatchesSearch = (column: { name: string }) => {
  if (!searchQuery.value) return false;
  return column.name.toLowerCase().includes(searchQuery.value.toLowerCase());
};

// Watch search query to auto-expand tables with matches
watch(() => searchQuery.value, (newQuery) => {
  if (!newQuery) {
    // If query is cleared, collapse all tables except the selected one
    expandedTables.value = selectedTable.value ? [selectedTable.value] : [];
    return;
  }
  
  // Find tables that match the search query (either table name or column names)
  const matchingTableIds = nodes.value.filter(node => {
    const query = newQuery.toLowerCase();
    
    // Check if table name matches
    if (node.data.label.toLowerCase().includes(query)) {
      return true;
    }
    
    // Check if any column name matches
    return node.data.columns.some((column: { name: string }) => 
      column.name.toLowerCase().includes(query)
    );
  }).map(node => node.id);
  
  // Auto expand tables with matches
  expandedTables.value = [...new Set([...expandedTables.value, ...matchingTableIds])];
});
</script>

<style scoped>
.database-diagram-container {
  height: 100%;
  display: flex;
  flex-direction: column;
}

.diagram-header {
  display: flex;
  justify-content: flex-end;
}

.diagram-content {
  position: relative;
}

.vue-flow-wrapper {
  width: 100%;
  height: 100%;
}

:deep(.vue-flow__minimap) {
  transform: scale(0.8);
  transform-origin: bottom right;
}

.table-list-item {
  transition: all 0.2s ease;
}

/* Define a custom dark mode bg-gray-750 color for consistent styling */
.dark .bg-gray-750 {
  background-color: #242937;
}

/* Fix any background issues for columns */
.dark .table-columns .column-item {
  background-color: transparent;
}
</style> 