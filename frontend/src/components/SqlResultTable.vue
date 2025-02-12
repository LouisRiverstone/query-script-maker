<template>
    <div v-if="value" class="fixed inset-0 flex items-center justify-center bg-transparent bg-opacity-50 shadow-2xl z-50">
        <div class=" text-black dark:text-white bg-white dark:bg-gray-700 rounded-lg shadow-lg p-6 h-1/2 w-1/2">
            <div class="flex flex-col h-full">
                <h2 class=" text-black dark:text-white text-xl font-semibold mb-4">Sql Result</h2>
                <div class="flex flex-col flex-grow overflow-y-auto">
                    <Table :data="props.data"></Table>
                </div>
                <div class="flex justify-end">
                    <button @click="close"
                        class="bg-indigo-500 text-black dark:text-white px-4 py-2 rounded">Close</button>
                </div>
            </div>
        </div>
    </div>
</template>


<script setup lang="ts">
import Table from './Table.vue';

const props = defineProps<{
    data: { [key: string]: any }[]
}>()

const emit = defineEmits(['close'])

const value = defineModel()

const close = () => {
    emit('close')
    value.value = false
}
</script>