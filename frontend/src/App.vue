<template>
    <main class="bg-gray-200 dark:bg-gray-700 dark:text-white text-gray-900 h-screen flex flex-col gap-5 w-full p-5">
        <div class="flex flex-col gap-3">
            <div class="flex flex-row justify-center">
                <Button type="button" @click="importXLSX">Import .XLSX</Button>
            </div>
            <div v-if="content.length > 0" class="flex flex-row justify-center">
                <Table :data="content"/>
            </div>
            <div v-if="headers.length > 0" class="flex flex-row justify-center">
                <VariablesCaster  :headers="headers" ref="variblesCasterRef" />
            </div>
            <div v-if="hasVariablesAssigned" class="flex flex-row justify-center">
                <Editor v-model="query"  :variables="variables" :data="content"  ref="editorRef"/>
            </div>
            <div v-if="hasEditor" class="flex flex-row justify-center">
                <Button type="button" @click="createSqlFile">Save .SQL</Button>
            </div>
        </div>
    </main>
</template>

<script lang="ts" setup>
import { computed, ref } from 'vue';
import Button from './components/Button.vue';
import { ReadFile } from '../wailsjs/go/main/App'
import { CreateSQLFile } from '../wailsjs/go/main/App'
import Table from './components/Table.vue';
import VariablesCaster from './components/VariablesCaster.vue';
import Editor from './components/Editor.vue';

const content = ref<{[k: string] : string}[]>([]);
const variblesCasterRef = ref<typeof VariablesCaster | null>(null);
const editorRef = ref<typeof Editor | null>(null);
const query = ref<string>('SELECT * FROM users where email = {{email}} and pode = {{pode}}');

const headers = computed(() => {
    return Object.keys(content.value[0] ?? []);
})

const hasVariablesAssigned = computed(() => {
    if(!variblesCasterRef.value){
        return false;
    }

    return variblesCasterRef.value.variables.length > 0;
})

const hasEditor = computed(() => {
    return !!editorRef.value;
})

const variables = computed(() => {
    if(!variblesCasterRef.value){
        return [];
    }

    return variblesCasterRef.value.variables;
})

const importXLSX = async () =>{
    content.value = JSON.parse(await ReadFile());
}

const createSqlFile = async () =>{
    if(!hasEditor.value){
        return;
    }

    const sql = editorRef.value.getBindedSQL();

    await CreateSQLFile(sql);
}
</script>

<style>

</style>