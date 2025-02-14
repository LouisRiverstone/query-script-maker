<template>
    <div v-if="value" class="fixed inset-0 flex items-center justify-center bg-transparent bg-opacity-50 shadow-2xl">
        <div class="text-black dark:text-white bg-white dark:bg-gray-700 rounded-lg shadow-lg p-6 w-1/3">
            <h2 class="text-black dark:text-white text-xl font-semibold mb-4">{{ props.title }}</h2>
            <p class="mb-6 text-black dark:text-white p">{{ props.message }}</p>
            <div class="flex justify-center">
                <button @click="confirm" class="bg-indigo-500 text-black dark:text-white px-4 py-2 rounded">OK</button>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>

const props = defineProps<{
    title: string,
    message: string
    modelValue: boolean
}>()

const emit = defineEmits(['cancel', 'confirm'])

const value = defineModel()

const confirm = () => {
    emit('confirm')
    value.value = false
}
</script>