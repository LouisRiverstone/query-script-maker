<template>
    <main class="bg-gray-200 dark:bg-gray-800 dark:text-white text-gray-900 flex flex-col gap-5 w-full p-5 h-screen overflow-y-auto">
        <Steps :steps="stepsHeaders" v-model="step" :disable-next-button="disableNextButton">
            <section v-show="step === 0">
                <div class="flex flex-col gap-3">
                    <div class="flex flex-row justify-center ">
                        <Button type="button" @click="importXLSX">Import .XLSX</Button>
                    </div>
                    <div v-if="content.length > 0" class="flex flex-row justify-center">
                        <Table :data="content" />
                    </div>
                </div>
            </section>
            <section v-show="step === 1">
                <div v-if="headers.length > 0" class="flex flex-row justify-center">
                    <VariablesCaster :headers="headers" ref="variblesCasterRef" />
                </div>
            </section>
            <section v-show="step === 2">
                <div v-if="hasEditor" class="flex flex-row justify-center">
                        <Button type="button" @click="createSqlFile">Save .SQL</Button>
                </div>
                <div class="flex flex-col gap-3">
                    <div v-if="hasVariablesAssigned" class="flex flex-row justify-center">
                        <Editor v-model="query" :variables="variables" :data="content" ref="editorRef" />
                    </div>
                </div>
            </section>
        </Steps>
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
import Steps from './components/Steps.vue';

const content = ref<{ [k: string]: string }[]>([]);
const variblesCasterRef = ref<typeof VariablesCaster | null>(null);
const editorRef = ref<typeof Editor | null>(null);
const query = ref<string>('SELECT * FROM users where email = {{email}} and pode = {{pode}}');

const stepsHeaders = computed(() => {
    return ['Import .XLSX', 'Assign Variables', 'Create SQL'];
})

const step = ref(0);

const headers = computed(() => {
    return Object.keys(content.value[0] ?? []);
})

const disableNextButton = computed(() => {
    if (step.value === 0) {
        return !headers.value.length > 0
    }

    if (step.value === 1) {
        return !hasVariablesAssigned.value;
    }
})

const hasVariablesAssigned = computed(() => {
    if (!variblesCasterRef.value) {
        return false;
    }

    return variblesCasterRef.value.variables.length > 0;
})

const hasEditor = computed(() => {
    return !!editorRef.value;
})

const variables = computed(() => {
    if (!variblesCasterRef.value) {
        return [];
    }

    return variblesCasterRef.value.variables;
})

const importXLSX = async () => {
    content.value = JSON.parse(await ReadFile());
}

const createSqlFile = async () => {
    if (!hasEditor.value) {
        return;
    }

    const sql = editorRef.value!.getBindedSQL();

    await CreateSQLFile(sql);
}
</script>

<style></style>