<template>
    <main class="flex flex-col h-screen overflow-y-auto bg-gray-200 dark:bg-gray-800 dark:text-white text-gray-900">
        <Header />
        <transition name="fade">
            <router-view />
        </transition>
        <Footer />
        <Alert v-model="showVersionModal" title="New version available" message="A new version is available. Access the repository to download it." />
    </main>
</template>

<script lang="ts" setup>
import Header from './components/Header.vue';
import Footer from './components/Footer.vue';

import { CheckHasUpdate } from '../wailsjs/go/main/App'
import { onMounted, ref } from 'vue';
import Alert from './components/Alert.vue';

const showVersionModal = ref<boolean>(false)

const checkForUpdate = async () => {
  showVersionModal.value = await CheckHasUpdate()
}

onMounted(() => {
  checkForUpdate()
});

</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>