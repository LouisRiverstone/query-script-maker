<template>
  <div v-if="isOpen" class="fixed inset-0 overflow-hidden z-50 transition-opacity duration-300">
    <div class="flex items-center justify-center min-h-screen p-2">
      <!-- Backdrop -->
      <div class="fixed inset-0 bg-gray-900 bg-opacity-75 dark:bg-black dark:bg-opacity-80 backdrop-blur-sm transition-opacity duration-300" @click="close"></div>
      
      <!-- Modal container -->
      <div class="relative bg-white dark:bg-gray-850 border border-gray-200 dark:border-gray-700 rounded-xl w-full max-w-[95vw] max-h-[92vh] flex flex-col h-[92vh] shadow-2xl transition-all duration-300 ease-out transform">
        <!-- Header -->
        <div class="flex justify-between items-center px-4 py-3 border-b border-gray-200 dark:border-gray-700 flex-shrink-0 bg-white dark:bg-gray-800 rounded-t-xl">
          <div class="flex items-center">
            <svg class="w-5 h-5 text-indigo-600 dark:text-indigo-400 mr-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4m0 5c0 2.21-3.582 4-8 4s-8-1.79-8-4" />
            </svg>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Database Structure Diagram</h3>
          </div>
          <div class="flex items-center gap-2">
            <Button 
              @on-click="refreshDiagram" 
              :loading="isRefreshing" 
              variant="secondary" 
              size="sm"
              class="flex items-center"
            >
              <template #icon-left>
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
              </template>
              Refresh
            </Button>
            
            <!-- Theme toggle button -->
            <Button 
              @on-click="toggleTheme" 
              variant="ghost" 
              size="sm"
              title="Toggle theme"
              class="text-gray-500 dark:text-gray-400 hover:text-indigo-600 dark:hover:text-indigo-400"
            >
              <!-- Sun icon for dark mode -->
              <svg v-if="isDarkMode" xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
              </svg>
              
              <!-- Moon icon for light mode -->
              <svg v-else xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
              </svg>
            </Button>
            
            <!-- Close button -->
            <Button 
              @on-click="close" 
              variant="ghost" 
              size="sm"
              class="text-gray-500 dark:text-gray-400 hover:text-red-600 dark:hover:text-red-400"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </Button>
          </div>
        </div>
        
        <!-- Content -->
        <div class="flex-grow overflow-hidden bg-white dark:bg-gray-900">
          <div v-if="isRefreshing" class="flex items-center justify-center h-full">
            <div class="text-center p-8 bg-white dark:bg-gray-800 rounded-lg shadow-lg">
              <div class="inline-block">
                <svg class="animate-spin h-12 w-12 text-indigo-600 dark:text-indigo-400" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
              <p class="mt-4 text-gray-600 dark:text-gray-300 font-medium">Refreshing database structure...</p>
              <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">This may take a moment</p>
            </div>
          </div>
          <DatabaseDiagram v-else :databaseStructure="databaseStructure" :isDarkMode="isDarkMode" class="h-full" />
        </div>
        
        <!-- Footer (optional) -->
        <div class="px-4 py-2 border-t border-gray-200 dark:border-gray-700 flex-shrink-0 bg-white dark:bg-gray-850 text-xs text-gray-500 dark:text-gray-400 flex justify-between items-center rounded-b-xl">
          <span>Use the mouse wheel to zoom, drag to pan</span>
          <span v-if="tableCount > 0">{{ tableCount }} tables displayed</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref, computed, onMounted, watch } from 'vue';
import DatabaseDiagram from './DatabaseDiagram.vue';
import Button from './Button.vue';
import { GetLatestDatabaseStructure, GetDatabaseConnection, GetDatabaseStructure } from '../../wailsjs/go/main/App';

const props = defineProps<{
  isOpen: boolean;
  databaseStructure: string;
}>();

const emit = defineEmits(['close', 'refresh']);

const isRefreshing = ref(false);
const isDarkMode = ref(false);
const dbStructure = computed(() => {
  try {
    return JSON.parse(props.databaseStructure || '{}');
  } catch (e) {
    return {};
  }
});

const tableCount = computed(() => {
  return dbStructure.value?.tables?.length || 0;
});

const close = () => {
  emit('close');
};

const refreshDiagram = async () => {
  try {
    isRefreshing.value = true;
    
    // Try to fetch structure directly from database first
    try {
      const connection = await GetDatabaseConnection();
      if (connection && connection.ID) {
        // If we have a valid connection, fetch fresh structure
        const freshStructure = await GetDatabaseStructure(connection);
        emit('refresh', freshStructure);
        return;
      }
    } catch (error) {
      console.log('Could not refresh from database directly, using cached structure');
    }
    
    // Fall back to cached structure in SQLite
    const newStructure = await GetLatestDatabaseStructure();
    emit('refresh', newStructure);
  } catch (error) {
    console.error('Error refreshing database structure:', error);
  } finally {
    isRefreshing.value = false;
  }
};

const toggleTheme = () => {
  if (isDarkMode.value) {
    document.documentElement.classList.remove('dark');
    localStorage.theme = 'light';
  } else {
    document.documentElement.classList.add('dark');
    localStorage.theme = 'dark';
  }
  isDarkMode.value = !isDarkMode.value;
};

// Check system preference for dark mode
const updateDarkMode = () => {
  isDarkMode.value = document.documentElement.classList.contains('dark') || 
                    document.body.classList.contains('dark') ||
                    window.matchMedia('(prefers-color-scheme: dark)').matches ||
                    localStorage.theme === 'dark';
};

onMounted(() => {
  updateDarkMode();
  
  // Apply dark mode if saved in localStorage or from system preference
  if (localStorage.theme === 'dark' || (!('theme' in localStorage) && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
    document.documentElement.classList.add('dark');
    isDarkMode.value = true;
  } else {
    document.documentElement.classList.remove('dark');
    isDarkMode.value = false;
  }
});

// Pass isDarkMode to the DatabaseDiagram component when it changes
watch(() => isDarkMode.value, (newValue) => {
  // Any additional logic needed when dark mode changes
});
</script>

<style>

/* Smooth transitions */
.transition-all {
  transition-property: all;
}

/* Optional custom scrollbar for browsers that support it */
.overflow-auto::-webkit-scrollbar,
.overflow-y-auto::-webkit-scrollbar {
  width: 8px;
}

.overflow-auto::-webkit-scrollbar-track,
.overflow-y-auto::-webkit-scrollbar-track {
  background: rgba(0, 0, 0, 0.05);
  border-radius: 8px;
}

.dark .overflow-auto::-webkit-scrollbar-track,
.dark .overflow-y-auto::-webkit-scrollbar-track {
  background: rgba(255, 255, 255, 0.05);
}

.overflow-auto::-webkit-scrollbar-thumb,
.overflow-y-auto::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.2);
  border-radius: 8px;
}

.dark .overflow-auto::-webkit-scrollbar-thumb,
.dark .overflow-y-auto::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.2);
}

.overflow-auto::-webkit-scrollbar-thumb:hover,
.overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: rgba(0, 0, 0, 0.3);
}

.dark .overflow-auto::-webkit-scrollbar-thumb:hover,
.dark .overflow-y-auto::-webkit-scrollbar-thumb:hover {
  background: rgba(255, 255, 255, 0.3);
}
</style> 