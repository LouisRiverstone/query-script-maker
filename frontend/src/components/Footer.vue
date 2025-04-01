<template>
<footer class="bg-gray-100 dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800">
  <div class="mx-auto max-w-screen-xl px-4 py-4 lg:py-6 flex flex-col sm:flex-row items-center justify-between">
    <p class="text-center text-sm text-gray-600 dark:text-gray-400 hover:text-indigo-600 dark:hover:text-indigo-400 transition-colors duration-200">
      <a href="/" @click.prevent="openRepository" class="inline-flex items-center gap-1">
        <span>SQL Script Maker</span>
        <span class="text-xs opacity-75">- An Open Source Project by Luiz Neves</span>
      </a>
    </p>
    <div class="flex items-center mt-2 sm:mt-0">
      <span class="text-xs text-gray-500 dark:text-gray-500 bg-gray-200 dark:bg-gray-800 px-2 py-1 rounded-full">
        v{{ version }}
      </span>
    </div>
  </div>
</footer>
</template>

<script setup lang="ts">
import { BrowserOpenURL } from '../../wailsjs/runtime'
import { GetBuildParams } from '../../wailsjs/go/main/App'
import { onMounted, ref } from 'vue';

const version = ref<string>("0.0.0");

const getVersion = async () => {
  const data = await GetBuildParams()
  version.value = data.version
}

const openRepository = () => {
  BrowserOpenURL("https://github.com/LouisRiverstone/query-script-maker");
}

onMounted(() => {
  getVersion()
});
</script>