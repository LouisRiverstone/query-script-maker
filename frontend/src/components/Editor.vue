<template>
    <div class="flex flex-col gap-3 w-full">
        <slot></slot>
        <div class="flex flex-col lg:flex-row gap-6">
            <div class="flex flex-col w-full gap-3">
                <Divider v-if="showBindedSql">Entrada</Divider>
                <code-mirror v-model="value" :lang="lang" :extensions="[oneDarkTheme]" :linter="null" basic wrap tab class="w-full" />
            </div>
            <div v-if="showBindedSql" class="flex flex-col w-full gap-3">
                <Divider>Saída</Divider>
                <code-mirror v-model="linesBinded" :lang="lang" :extensions="[oneDarkTheme]" :linter="null" basic wrap tab class="w-full" />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted } from 'vue';
import CodeMirror from 'vue-codemirror6';
import { sql, MySQL } from "@codemirror/lang-sql";
import { oneDarkTheme } from '@codemirror/theme-one-dark';
import { MakeBindedSQL } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import { computedAsync } from '@vueuse/core';

import Button from './Button.vue';
import Divider from './Divider.vue';

const props = defineProps<{
    variables?: Array<main.Variable>
    data?: { [key: string]: any }[]
    modelValue: string,
    showBindedSql?: boolean
}>()

const emit = defineEmits(['update:modelValue']);

const value = ref(props.modelValue);

const lang = sql({
    dialect: MySQL
});

const linesBinded = computedAsync(async () => {
    if(!props.showBindedSql) {
        return "";
    }

    return await MakeBindedSQL(value.value, props.data!, props.variables!) ?? "";
}, "");

const getBindedSQL = async (): Promise<string> => {
    if(!props.showBindedSql) {
        return "";
    }

    return await MakeBindedSQL(value.value, props.data!, props.variables!);
}

onMounted(() => {
    getDatabaseConnection();
});

watch(() => props.modelValue, (newValue) => {
    value.value = newValue;
});

watch(() => value.value, (val) => {
    emit('update:modelValue', val);
});

defineExpose({
    getBindedSQL
});
</script>