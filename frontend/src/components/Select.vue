<template>
<div class="flex flex-col w-full">
  <label :for="id" class="block mb-2 text-sm font-medium text-gray-900 dark:text-white">{{ label }}</label>
  <select
    v-model="value"
    :disabled="disabled"
    :name="id"
    :id="id"
    class="bg-gray-50 border border-gray-300 text-gray-900 text-sm rounded-lg focus:ring-indigo-500 focus:border-indigo-500 block w-full p-2.5 dark:bg-gray-600 dark:border-gray-500 dark:placeholder-gray-400 dark:text-white dark:focus:ring-indigo-500 dark:focus:border-indigo-500"
  >
    <option value="">Please select</option>
    <option v-for="(option,i) in props.options" :key="i" :value="option.value">{{ option.label }}</option>
  </select>
</div>
</template>

<script setup lang="ts">
const props = defineProps<{
    options: {
        value: string
        label: string
    }[]
    id: string
    disabled?: boolean
    label: string
}>()

const value = defineModel()
</script>