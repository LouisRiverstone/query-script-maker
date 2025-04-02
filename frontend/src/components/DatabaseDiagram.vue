<template>
  <div class="database-diagram-container">
    <div class="diagram-header">
      <div class="diagram-controls">
        <Button type="button" @click="zoomIn" class="mr-1">+</Button>
        <Button type="button" @click="zoomOut" class="mr-1">-</Button>
        <Button type="button" @click="fitView">Fit</Button>
      </div>
    </div>
    <div ref="diagramContainer" class="diagram-content bg-white dark:bg-gray-800 rounded-lg p-4 h-[70vh]">
      <VueFlow v-if="nodes.length > 0"
        :default-zoom="1"
        :min-zoom="0.2"
        :max-zoom="4"
        :nodes="nodes"
        :edges="edges"
        class="vue-flow-wrapper h-full"
        fit-view-on-init
      >
        <!-- Custom Node Types -->
        <template #node-table="nodeProps">
          <div class="bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 rounded-md shadow-sm p-2">
            <div class="table-header bg-indigo-100 dark:bg-indigo-900 p-2 mb-2 rounded font-bold text-center">
              {{ nodeProps.data.label }}
            </div>
            <div class="table-columns">
              <div v-for="(column, index) in nodeProps.data.columns" :key="index" 
                class="column-item flex items-center p-1 border-b border-gray-200 dark:border-gray-700 last:border-0"
                :class="{'bg-yellow-50 dark:bg-yellow-900': column.isPrimary}">
                <div class="flex items-center">
                  <span v-if="column.isPrimary" class="text-yellow-600 dark:text-yellow-400 mr-1">🔑</span>
                  <span v-else-if="column.isForeign" class="text-blue-600 dark:text-blue-400 mr-1">🔗</span>
                  <span v-else class="mr-1 w-3.5"></span>
                </div>
                <div class="flex-1 text-left whitespace-nowrap pr-2">{{ column.name }}</div>
                <div class="text-gray-500 dark:text-gray-400 text-xs whitespace-nowrap">{{ column.type }}</div>
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
            <div class="legend-item flex items-center">
              <div class="w-3 h-3 bg-blue-500 dark:bg-blue-400 rounded-sm mr-2"></div>
              <span class="text-xs text-gray-600 dark:text-gray-300">Foreign Key</span>
            </div>
          </div>
        </Panel>
      </VueFlow>
      <div v-else class="empty-diagram flex items-center justify-center h-full">
        <p class="text-gray-500 dark:text-gray-400 text-center">No database structure found</p>
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
  MarkerType
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
}>();

const diagramContainer = ref<HTMLElement | null>(null);
const isDarkMode = ref(false);

// Elements for Vue Flow
const nodes = ref<any[]>([]);
const edges = ref<any[]>([]);

// Get Vue Flow instance
const { zoomIn, zoomOut, fitView } = useVueFlow({
  nodes: [],
  edges: []
});

// Function to check if dark mode is enabled
const updateDarkMode = () => {
  isDarkMode.value = document.documentElement.classList.contains('dark') || 
                   document.body.classList.contains('dark') ||
                   window.matchMedia('(prefers-color-scheme: dark)').matches;
};

// Node color function for minimap
const getNodeColor = (node: Node) => {
  return '#4f46e5'; // Default indigo color
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
  margin-bottom: 8px;
}

.diagram-content {
  flex-grow: 1;
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
</style> 