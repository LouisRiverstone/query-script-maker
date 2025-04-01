<template>
    <div v-if="isVisible" 
        class="fixed inset-0 flex items-center justify-center z-50"
        @click.self="cancel"
        v-bind="$attrs">
        <div class="fixed inset-0 bg-black bg-opacity-40 backdrop-blur-sm transition-opacity duration-300" 
            :class="isVisible ? 'opacity-100' : 'opacity-0'"></div>
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow-xl p-6 w-full max-w-md mx-4 transform transition-all duration-300"
            :class="isVisible ? 'scale-100 opacity-100' : 'scale-95 opacity-0'">
            <h2 class="text-xl font-semibold mb-4 text-gray-900 dark:text-white">{{ title }}</h2>
            <p class="mb-6 text-gray-700 dark:text-gray-300">{{ message }}</p>
            <div class="flex justify-end gap-3">
                <button @click="cancel" 
                    class="px-4 py-2 rounded-md bg-gray-200 text-gray-800 hover:bg-gray-300 dark:bg-gray-700 dark:text-gray-200 dark:hover:bg-gray-600 transition-colors duration-200">
                    Cancel
                </button>
                <button @click="confirm" 
                    class="px-4 py-2 rounded-md bg-indigo-600 text-white hover:bg-indigo-700 dark:bg-indigo-500 dark:hover:bg-indigo-600 transition-colors duration-200">
                    Confirm
                </button>
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { computed } from 'vue';

const props = defineProps<{
    title: string,
    message: string,
    modelValue: boolean
}>();

const emit = defineEmits(['update:modelValue', 'cancel', 'confirm']);

const isVisible = computed({
    get: () => props.modelValue,
    set: (value) => emit('update:modelValue', value)
});

const cancel = () => {
    isVisible.value = false;
    emit('cancel');
};

const confirm = () => {
    isVisible.value = false;
    emit('confirm');
};
</script>