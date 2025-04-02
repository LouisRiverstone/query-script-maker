<template>
  <div v-if="isOpen" class="fixed inset-0 overflow-y-auto z-50">
    <div class="flex items-center justify-center min-h-screen p-4">
      <div class="fixed inset-0 bg-black opacity-50"></div>
      <div class="relative bg-white dark:bg-gray-800 rounded-lg w-full max-w-7xl overflow-hidden max-h-[90vh] flex flex-col">
        <div class="flex justify-between items-center p-4 border-b border-gray-200 dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white">Database Structure Diagram</h3>
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
            <Button 
              @on-click="close" 
              variant="ghost" 
              size="sm"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </Button>
          </div>
        </div>
        <div class="p-4 flex-grow overflow-auto">
          <div v-if="isRefreshing" class="flex items-center justify-center h-full">
            <div class="text-center">
              <div class="inline-block">
                <svg class="animate-spin h-10 w-10 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
              </div>
              <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">Refreshing database structure...</p>
            </div>
          </div>
          <DatabaseDiagram v-else :databaseStructure="databaseStructure" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { defineProps, defineEmits, ref } from 'vue';
import DatabaseDiagram from './DatabaseDiagram.vue';
import Button from './Button.vue';
import { GetLatestDatabaseStructure, GetDatabaseConnection, GetDatabaseStructure } from '../../wailsjs/go/main/App';

const props = defineProps<{
  isOpen: boolean;
  databaseStructure: string;
}>();

const emit = defineEmits(['close', 'refresh']);

const isRefreshing = ref(false);

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
</script> 