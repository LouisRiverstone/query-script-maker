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