<template>
    <div class="overflow-x-auto">
        <table class="min-w-full divide-y-2 divide-gray-200 bg-white text-sm">
            <thead class="ltr:text-left rtl:text-right">
                <tr>
                    <th v-for="(title, i) in header" :key="i" class="whitespace-nowrap px-4 py-2 font-medium text-gray-900">{{ title }}</th>
                </tr>
            </thead>

            <tbody class="divide-y divide-gray-200">
                <tr v-for="(row, i) in rows" :key="i">
                    <td v-for="(column,j) in row" :key="j" class="whitespace-nowrap px-4 py-2 font-medium text-gray-900">{{ column }}</td>
                </tr>
            </tbody>
        </table>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

const props = defineProps<{
    data: { [key: string]: any }[]
}>()

const header = computed(() => {
    return Object.keys(props.data[0]);
})

const rows = computed(() => {
    return props.data.map((row) => {
        return Object.values(row);
    })
})
</script>