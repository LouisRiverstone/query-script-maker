<template>
    <div class="flex flex-col flex-grow gap-10 w-full px-5 py-10">
        <div class="flex flex-col gap-5">
            <Divider>Local Query Database</Divider>
            <div v-if="!loadingDatabase" class="flex flex-col gap-5">
                <p>
                    Manage your database by importing or exporting data files. Click "Import" to upload a database file
                    from
                    your local system, or click "Export" to download the current database file.
                </p>
                <div class="flex flex-row gap-2">
                    <Button type="button" class="w-full" @click="importDatabaseFile">Import</Button>
                    <Button type="button" class="w-full" @click="exportDatabaseFile">Export</Button>
                </div>
            </div>
            <div v-else>
                <Loader />
            </div>
        </div>
        <div class="flex flex-col gap-5">
            <Divider>Query Editor</Divider>
            <div v-if="!loadingQueries" class="flex flex-col gap-5">
                <div class="flex flex-row gap-5 items-end">
                    <div class="flex flex-grow w-full flex-col">
                        <Select label="Queries" id="select-query" :options="mapQueriesForSelect(queries)"
                            v-model="selectedQueryId" />
                    </div>
                    <Button type="button" @click="resetQuery">NEW QUERY</Button>
                </div>
                <div v-if="query" class="flex flex-col gap-3">
                    <Input id="query-title" type="text" v-model="query.Title" label="Title" />
                    <Input id="query-title" type="text" v-model="query.Description" label="Description" />
                    <span>Sql Query:</span>
                    <Editor v-model="query.Query" :show-binded-sql="false" />
                </div>

                <div class="flex flex-row justify-end gap-2">
                    <Button :disabled="!!!query.ID" type="button" @click="() => showDeleteQueryModal = true">DELETE
                        QUERY</Button>
                    <Button type="button" @click="() => showSaveQueryModal = true">SAVE QUERY</Button>
                </div>
            </div>
            <div v-else>
                <Loader />
            </div>
        </div>
        <div class="flex flex-col gap-5">
            <Divider>Database Connection</Divider>
            <div v-if="!loadingDatabaseConnection" class="flex flex-col gap-5">
                <p>
                    Configure your database connection by filling in the fields below. Click "Save" to save the
                    connection.
                </p>
                <div v-if="databaseConnection" class="flex flex-col gap-3">
                    <div class="flex md:flex-row flex-col gap-5">
                        <Input id="database-host" type="text" v-model="databaseConnection.Host" label="Host" />
                        <Input id="database-port" type="number" v-model="databaseConnection.Port" label="Port" />
                    </div>
                    <div class="flex md:flex-row flex-col gap-5">
                        <Input id="database-database" type="text" v-model="databaseConnection.Database" label="Database" />
                    </div>
                    <div class="flex md:flex-row flex-col gap-5">
                        <Input id="database-username" type="text" v-model="databaseConnection.Username" label="Username" />
                        <Input id="database-password" type="password" v-model="databaseConnection.Password" label="Password" />
                    </div>
                </div>
                <div class="flex flex-row gap-2 justify-end">
                    <div>
                        <Button type="button" class="w-full" @click="() => showSaveDatabaseConnection = true">SAVE CONNECTION</Button>
                    </div>
                </div>
            </div>
            <div v-else>
                <Loader />
            </div>
        </div>
        <ConfirmationModal v-model="showDeleteQueryModal" title="Delete Query" message="Are you sure?"
            @cancel="() => showDeleteQueryModal = false" @confirm="deleteQuery" />
        <ConfirmationModal v-model="showSaveQueryModal" title="Save Query" message="Are you sure?"
            @cancel="() => showSaveQueryModal = false" @confirm="saveQuery" />
        <ConfirmationModal v-model="showSaveDatabaseConnection" title="Save Database Connection" message="Are you sure?"
            @cancel="() => showSaveDatabaseConnection = false" @confirm="saveDatabaseConnection" />
    </div>
</template>

<script setup lang="ts">
import { onMounted, ref, watch } from 'vue'

import Divider from '../components/Divider.vue'
import Button from '../components/Button.vue';
import Select from '../components/Select.vue';
import Loader from '../components/Loader.vue';
import Editor from '../components/Editor.vue';
import Input from '../components/Input.vue';
import ConfirmationModal from '../components/ConfirmationModal.vue';

import { ImportDatabaseFile, ExportDatabaseFile, GetQueriesList, InsertQueryInDatabase, UpdateQuery, DeleteQuery, CreateOrUpdateDatabaseConnection , GetDatabaseConnection} from '../../wailsjs/go/main/App'
import { main } from '../../wailsjs/go/models';

const queries = ref<Array<main.Query>>([])
const query = ref<main.Query>({
    ID: 0,
    Title: "",
    Query: "",
    Description: "",
})

const databaseConnection = ref<main.DatabaseConnection>({
    ID: 0,
    Host: "",
    Port: 0,
    Username: "",
    Password: "",
    Database: ""
})

const selectedQueryId = ref<string>("")
const loadingQueries = ref<boolean>(false)
const loadingDatabase = ref<boolean>(false)
const showDeleteQueryModal = ref<boolean>(false)
const showSaveQueryModal = ref<boolean>(false)
const showSaveDatabaseConnection = ref<boolean>(false)
const loadingDatabaseConnection = ref<boolean>(false)

const mapQueriesForSelect = (queries: Array<main.Query>) => {
    return queries.map((query) => {
        return {
            label: query.Title,
            value: String(query.ID),
        }
    })
}

const saveQuery = async () => {
    try {
        loadingQueries.value = true

        if (query.value.ID !== 0) {
            await UpdateQuery(query.value.ID!, query.value)
        }

        if (query.value.ID === 0) {
            await InsertQueryInDatabase(query.value)
        }

        resetQuery()

        await getQueries()
    } catch (error) {
        console.error(error)
    } finally {
        loadingQueries.value = false
    }
}

const getQueries = async () => {
    queries.value = await GetQueriesList(false)
}

const deleteQuery = async () => {
    try {
        loadingQueries.value = true

        if (query.value.ID === 0) {
            return
        }

        await DeleteQuery(Number(query.value.ID))

        resetQuery()

        await getQueries();
    } catch (error) {
        console.error(error)
    } finally {
        loadingQueries.value = false
    }
}

const importDatabaseFile = async () => {
    try {
        loadingDatabase.value = true
        await ImportDatabaseFile()

        queries.value = await GetQueriesList(false)
    } catch (error) {
        console.error(error)
    } finally {
        loadingDatabase.value = false
    }
}

const exportDatabaseFile = async () => {
    try {
        loadingDatabase.value = true
        await ExportDatabaseFile()
    } catch (error) {
        console.error(error)
    } finally {
        loadingDatabase.value = false
    }
}

const resetQuery = () => {
    query.value = {
        ID: 0,
        Title: "",
        Query: "",
        Description: "",
    }

    selectedQueryId.value = ""
}

const getDatabaseConnection = async () => {
    try {
        loadingDatabaseConnection.value = true
        databaseConnection.value = await GetDatabaseConnection()
    } catch (error) {
        console.error(error)
    } finally {
        loadingDatabaseConnection.value = false
    }
}

const saveDatabaseConnection = async () => {
    try {
        loadingDatabaseConnection.value = true

        const dbConn: main.DatabaseConnection = {
            ID: databaseConnection.value.ID,
            Host: databaseConnection.value.Host,
            Port: Number(databaseConnection.value.Port),
            Username: databaseConnection.value.Username,
            Password: databaseConnection.value.Password,
            Database: databaseConnection.value.Database,
        }

        databaseConnection.value = await CreateOrUpdateDatabaseConnection(dbConn)
    } catch (error) {
        console.error(error)
    } finally {
        loadingDatabaseConnection.value = false
    }
}

const mount = async () => {
    try {
        loadingQueries.value = true
        await getQueries()
        await getDatabaseConnection()
    } catch (error) {
        console.error(error)
    } finally {
        loadingQueries.value = false
    }
}

onMounted(async () => {
    await mount()
})

watch(() => selectedQueryId.value, (value) => {
    query.value = queries.value.find((query) => query.ID === Number(value))!
})
</script>