<template>
<footer class="bg-gray-100 dark:bg-gray-900">
  <div class="relative mx-auto max-w-screen-xl px-1 py-4 sm:px-1 lg:px-2 lg:pt-6">
    <p class="text-center text-sm text-gray-500 dark:text-gray-400">
      <a href="/" @click="openRepository">
        SQL Script Maker - An Open Source Project by Luiz Neves
      </a>
    </p>
  </div>
</footer>
</template>

<script setup lang="ts">
import { BrowserOpenURL } from '../../wailsjs/runtime'

const openRepository = () => {
  BrowserOpenURL("https://github.com/LouisRiverstone/query-script-maker");
}
</script>