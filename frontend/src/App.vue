<template>
    <main class="flex flex-col h-screen overflow-y-auto bg-gray-200 dark:bg-gray-800 dark:text-white text-gray-900">
        <Header />
        <transition name="fade">
            <router-view />
        </transition>
        <Footer />
    </main>
</template>

<script lang="ts" setup>
import Header from './components/Header.vue';
import Footer from './components/Footer.vue';

import { CheckHasUpdate } from '../wailsjs/go/main/App'
import { onMounted } from 'vue';


const checkForUpdate = async () => {
  const hasUpdate = await CheckHasUpdate()

  alert(hasUpdate)

  if (!hasUpdate) {
    return;
  }


  const update = confirm("There is a new version available. Do you want to update?");

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