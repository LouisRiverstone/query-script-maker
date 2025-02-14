<template>
<footer class="bg-gray-100 dark:bg-gray-900">
  <div class="relative mx-auto max-w-screen-xl px-1 py-4 sm:px-1 lg:px-2 lg:pt-6 flex flex-col items-center">
    <p class="text-center text-sm text-gray-500 dark:text-gray-400">
      <a href="/" @click="openRepository">
        SQL Script Maker - An Open Source Project by Luiz Neves
      </a>
    </p>
    <small class="text-gray-500 dark:text-gray-400">Version: {{ version }}</small>
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