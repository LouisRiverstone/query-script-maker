<template>
    <div v-if="value"
        class="fixed inset-0 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm z-50 transition-opacity duration-300 ease-in-out">
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl p-6 w-4/5 max-w-6xl h-4/5 flex flex-col border border-gray-200 dark:border-gray-700 transform transition-all">
            <div class="flex items-center justify-between mb-4 border-b border-gray-200 dark:border-gray-700 pb-3">
                <h2 class="text-xl font-bold text-gray-800 dark:text-white flex items-center gap-2">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-500" viewBox="0 0 20 20" fill="currentColor">
                        <path fill-rule="evenodd" d="M10 2a8 8 0 100 16 8 8 0 000-16zm0 14a6 6 0 100-12 6 6 0 000 12z" clip-rule="evenodd" />
                        <path fill-rule="evenodd" d="M10 4a1 1 0 100 2 1 1 0 000-2zM8.5 8.5a1.5 1.5 0 113 0v4a1.5 1.5 0 11-3 0v-4z" clip-rule="evenodd" />
                    </svg>
                    SQL Results <span class="text-sm text-gray-500 dark:text-gray-400 font-normal ml-2">({{ totalRows }} rows)</span>
                </h2>
                <button @click="close" class="text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-200 transition-colors">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                </button>
            </div>

            <div class="flex-grow overflow-y-auto mb-4 scrollbar-thin scrollbar-thumb-gray-300 dark:scrollbar-thumb-gray-600 scrollbar-track-transparent">
                <div v-if="props.data.length === 0" class="flex flex-col items-center justify-center h-full text-gray-500 dark:text-gray-400">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mb-4 text-gray-300 dark:text-gray-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                    </svg>
                    <p class="text-lg font-medium">No results found</p>
                    <p class="text-sm">Your query returned no data</p>
                </div>
                
                <div v-else>
                    <div v-for="(row, i) in props.data" :key="i" class="mb-6">
                        <div class="flex items-center mb-2 gap-2">
                            <div class="bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300 px-2 py-1 rounded-md text-xs font-medium">
                                Row {{ i + 1 }}
                            </div>
                            <div v-if="showCollapseButtons" @click="toggleRowCollapse(i)" 
                                class="text-xs text-gray-500 dark:text-gray-400 hover:text-indigo-500 dark:hover:text-indigo-400 cursor-pointer font-medium flex items-center gap-1">
                                {{ rowCollapsed[i] ? 'Expand' : 'Collapse' }}
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" :class="rowCollapsed[i] ? 'transform rotate-180' : ''" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                                </svg>
                            </div>
                        </div>
                        <div v-show="!rowCollapsed[i]">
                            <Table :data="row" />
                        </div>
                    </div>
                </div>
            </div>

            <div class="flex justify-between items-center pt-3 border-t border-gray-200 dark:border-gray-700">
                <div class="text-sm text-gray-500 dark:text-gray-400">
                    Showing {{ props.data.length }} result set{{ props.data.length !== 1 ? 's' : '' }}
                </div>
                <div class="flex gap-3">
                    <button v-if="props.data.length > 1" @click="toggleAllRows" 
                        class="text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-300 text-sm font-medium">
                        {{ allRowsCollapsed ? 'Expand All' : 'Collapse All' }}
                    </button>
                    <button @click="close"
                        class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-md text-sm font-medium transition-colors duration-150 flex items-center gap-2">
                        <span>Close</span>
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
            </div>
        </div>
    </div>
</template>


<script setup lang="ts">
import Table from './Table.vue';
import { ref, computed, onMounted } from 'vue';

const props = defineProps<{
    data: { [key: string]: any }[][]
}>();

const emit = defineEmits(['close']);
const value = defineModel();

// Row collapse state
const rowCollapsed = ref<Record<number, boolean>>({});
const allRowsCollapsed = computed(() => {
    return Object.values(rowCollapsed.value).every(val => val);
});

// Total number of rows across all result sets
const totalRows = computed(() => {
    return props.data.reduce((sum, resultSet) => sum + resultSet.length, 0);
});

// Only show collapse buttons if we have multiple result sets
const showCollapseButtons = computed(() => props.data.length > 1);

// Initialize collapse state
onMounted(() => {
    props.data.forEach((_, index) => {
        rowCollapsed.value[index] = false;
    });
});

const toggleRowCollapse = (index: number) => {
    rowCollapsed.value[index] = !rowCollapsed.value[index];
};

const toggleAllRows = () => {
    const newState = !allRowsCollapsed.value;
    props.data.forEach((_, index) => {
        rowCollapsed.value[index] = newState;
    });
};

const close = () => {
    emit('close');
    value.value = false;
};
</script>

<style scoped>
/* Custom scrollbar styles */
.scrollbar-thin::-webkit-scrollbar {
    width: 6px;
}
.scrollbar-thin::-webkit-scrollbar-track {
    background: transparent;
}
.scrollbar-thin::-webkit-scrollbar-thumb {
    background-color: #d1d5db;
    border-radius: 3px;
}
.dark .scrollbar-thin::-webkit-scrollbar-thumb {
    background-color: #4b5563;
}
</style>