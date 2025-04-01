<template>
    <div class="overflow-x-auto w-full bg-white dark:bg-gray-800 rounded-lg shadow-sm border border-gray-200 dark:border-gray-700">
        <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700 text-sm">
            <thead class="bg-gray-100 dark:bg-gray-900">
                <tr>
                    <th 
                        v-for="(title, i) in header" 
                        :key="i" 
                        class="px-4 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider"
                    >
                        {{ title }}
                    </th>
                </tr>
            </thead>

            <tbody class="divide-y divide-gray-200 dark:divide-gray-700 bg-white dark:bg-gray-800">
                <tr 
                    v-for="(row, i) in rows" 
                    :key="i"
                    :class="[i % 2 === 0 ? 'bg-white dark:bg-gray-800' : 'bg-gray-50 dark:bg-gray-900/50', 'hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-150']"
                >
                    <td 
                        v-for="(column, j) in row" 
                        :key="j" 
                        class="px-4 py-3 whitespace-nowrap text-gray-900 dark:text-gray-300"
                    >
                        {{ column }}
                    </td>
                </tr>
                
                <!-- Empty state -->
                <tr v-if="rows.length === 0">
                    <td :colspan="header.length" class="px-4 py-8 text-center text-gray-500 dark:text-gray-400">
                        No data available
                    </td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

const props = defineProps<{
    data: { [key: string]: any }[]
}>();

const header = computed(() => {
    if (!props.data.length) return [];
    return Object.keys(props.data[0]);
});

const rows = computed(() => {
    return props.data.map((row) => {
        return Object.values(row);
    });
});
</script>