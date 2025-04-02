<template>
    <div class="flex flex-col gap-3 w-full">
        <div v-for="(title, i) in headers" :key="i" class="flex flex-row gap-3">
            <Input type="text" v-model="fields[i]" :label="title" :id="`${title}[${i}]`" />
        </div>
    </div>
</template>

<script setup lang="ts">
import { computed, onMounted, ref } from 'vue';
import Input from './Input.vue';

const props = defineProps<{
    headers: string[]
}>()

const emit = defineEmits(['update:modelValue']);

const fields = ref<string[]>([]);

const variables = computed(() => {
    return fields.value.reduce((acc, field, i) => {
        console.log(field)
        acc.push({ Field: props.headers[i], Value: field, Position: i });
        
        return acc;
    }, [] as { Field: string, Value: string, Position: number }[]);
});

const parseAsVariableName = (str: string) => {
    return str.replace(/[^a-zA-Z0-9]/g, '_');
}

onMounted(() => {
    const tableVars = props.headers.map((header) => {
        return parseAsVariableName(header);
    });

    fields.value = [...tableVars];
});

defineExpose({
    variables
});
</script>