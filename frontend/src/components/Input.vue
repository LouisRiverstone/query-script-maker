<template>
    <div class="flex flex-col w-full">
        <label :for="id" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">{{ label }}</label>
        <input :type="props.type" :id="props.id" :placeholder="props.placeholder" :value="value" @input="emitText"
            class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-indigo-500 focus:border-indigo-500 block w-full p-2.5 dark:bg-gray-700 dark:border-gray-600 dark:placeholder-gray-400 dark:text-white dark:focus:ring-indigo-500 dark:focus:border-indigo-500" />
    </div>
</template>

<script lang="ts" setup>
import { defineModel } from 'vue';

const value = defineModel()

const props = defineProps<{
    id: string,
    label: string,
    type: string,
    placeholder?: string
}>()

const emitText = (e: Event) => {
    value.value = (e.target as HTMLInputElement).value;
}

</script>