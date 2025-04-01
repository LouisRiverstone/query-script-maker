<template>
    <header class="bg-white dark:bg-gray-900 shadow-sm border-b border-gray-200 dark:border-gray-800">
        <div class="mx-auto flex h-16 max-w-screen-xl items-center justify-between px-4 sm:px-6 lg:px-8">
            <router-link to="/" class="flex items-center text-indigo-600 dark:text-indigo-400 transition-transform hover:scale-105">
                <span class="sr-only">Home</span>
                <svg xmlns="http://www.w3.org/2000/svg" version="1.1" width="32" height="32" viewBox="0 0 132 132">
                    <svg viewBox="0 0 163.01300048828125 139" width="110.00000000000001" height="110.00000000000001" x="10.999999999999993" y="10.999999999999993">
                        <svg width="163" height="134" viewBox="0 0 163 134" xmlns="http://www.w3.org/2000/svg">
                            <g clip-path="url(#clip0)" fill="currentColor">
                                <path d="M79.703 44.9371H140.585C144.773 44.9371 148.192 41.5127 148.192 37.3167C148.192 33.1207 144.773 29.6963 140.585 29.6963H79.703C75.5145 29.6963 72.0962 33.1207 72.0962 37.3167C72.0962 41.5127 75.5145 44.9371 79.703 44.9371Z"></path>
                                <path d="M96.5243 15.2408H125.763C129.951 15.2408 133.37 11.8164 133.37 7.6204C133.37 3.42444 129.938 0 125.763 0H96.5243C92.3358 0 88.9175 3.42444 88.9175 7.6204C88.9175 11.8164 92.3358 15.2408 96.5243 15.2408Z"></path>
                                <path d="M37.2506 15.2408H66.489C70.6775 15.2408 74.0959 11.8164 74.0959 7.6204C74.0959 3.42444 70.664 0 66.489 0H37.2506C33.0622 0 29.6438 3.42444 29.6438 7.6204C29.6438 11.8164 33.0622 15.2408 37.2506 15.2408Z"></path>
                                <path d="M48.8432 44.5182H55.6123C59.8008 44.5182 63.2192 41.0937 63.2192 36.8977C63.2192 32.7018 59.8008 29.2773 55.6123 29.2773H48.8432C44.6547 29.2773 41.2363 32.7018 41.2363 36.8977C41.2363 41.0937 44.6547 44.5182 48.8432 44.5182Z"></path>
                                <path d="M155.393 59.3799H94.5248C90.3363 59.3799 86.918 62.8043 86.918 67.0003C86.918 71.1962 90.3363 74.6207 94.5248 74.6207H155.407C159.595 74.6207 163.013 71.1962 163.013 67.0003C163 62.8043 159.582 59.3799 155.393 59.3799Z"></path>
                                <path d="M78.0277 66.5804C78.0277 62.3844 74.6093 58.96 70.4208 58.96H52.059C47.8705 58.96 44.4521 62.3844 44.4521 66.5804C44.4521 70.7763 47.8705 74.2008 52.059 74.2008H70.4208C74.6093 74.2008 78.0277 70.7763 78.0277 66.5804Z"></path>
                                <path d="M33.5755 66.5804C33.5755 62.3844 30.1572 58.96 25.9687 58.96H7.60685C3.41835 58.96 0 62.3844 0 66.5804C0 70.7763 3.41835 74.2008 7.60685 74.2008H25.9687C30.1572 74.2008 33.5755 70.7763 33.5755 66.5804Z"></path>
                                <path d="M140.571 89.0762H79.703C75.5145 89.0762 72.0962 92.5006 72.0962 96.6966C72.0962 100.893 75.5145 104.317 79.703 104.317H140.585C144.773 104.317 148.192 100.893 148.192 96.6966C148.192 92.5006 144.76 89.0762 140.571 89.0762Z"></path>
                                <path d="M125.763 118.759H96.5243C92.3358 118.759 88.9175 122.183 88.9175 126.379C88.9175 130.575 92.3358 134 96.5243 134H125.763C129.951 134 133.37 130.575 133.37 126.379C133.37 122.183 129.938 118.759 125.763 118.759Z"></path>
                                <path d="M66.489 118.759H37.2506C33.0622 118.759 29.6438 122.183 29.6438 126.379C29.6438 130.575 33.0622 134 37.2506 134H66.489C70.6775 134 74.0959 130.575 74.0959 126.379C74.0959 122.183 70.664 118.759 66.489 118.759Z"></path>
                                <path d="M48.8432 104.737H55.6123C59.8008 104.737 63.2192 101.312 63.2192 97.1165C63.2192 92.9205 59.8008 89.4961 55.6123 89.4961H48.8432C44.6547 89.4961 41.2363 92.9205 41.2363 97.1165C41.2363 101.312 44.6547 104.737 48.8432 104.737Z"></path>
                            </g>
                        </svg>
                    </svg>
                </svg>
            </router-link>

            <div class="flex flex-1 items-center justify-end">
                <nav aria-label="Global" class="hidden md:block">
                    <ul class="flex items-center gap-8 text-sm">
                        <li>
                            <router-link to="/"
                                class="relative flex items-center px-2 py-1 font-medium transition-colors"
                                :class="[$route.path === '/' ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-700 hover:text-indigo-600 dark:text-gray-300 dark:hover:text-indigo-400']"
                                active-class="text-indigo-600 dark:text-indigo-400">
                                <span>Script Maker</span>
                                <span v-if="$route.path === '/'" class="absolute inset-x-0 -bottom-px h-0.5 bg-indigo-600 dark:bg-indigo-400"></span>
                            </router-link>
                        </li>

                        <li>
                            <router-link to="/config"
                                class="relative flex items-center px-2 py-1 font-medium transition-colors"
                                :class="[$route.path === '/config' ? 'text-indigo-600 dark:text-indigo-400' : 'text-gray-700 hover:text-indigo-600 dark:text-gray-300 dark:hover:text-indigo-400']"
                                active-class="text-indigo-600 dark:text-indigo-400">
                                <span>Config</span>
                                <span v-if="$route.path === '/config'" class="absolute inset-x-0 -bottom-px h-0.5 bg-indigo-600 dark:bg-indigo-400"></span>
                            </router-link>
                        </li>

                        <li>
                            <a href="/" 
                                @click.prevent="openRepository"
                                class="relative flex items-center px-2 py-1 font-medium text-gray-700 hover:text-indigo-600 dark:text-gray-300 dark:hover:text-indigo-400 transition-colors">
                                Github
                            </a>
                        </li>
                    </ul>
                </nav>

                <!-- Mobile menu button -->
                <button type="button" @click="mobileMenuOpen = !mobileMenuOpen" class="rounded p-2 text-gray-600 md:hidden">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path v-if="mobileMenuOpen" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                        <path v-else stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
                    </svg>
                </button>
            </div>
        </div>

        <!-- Mobile menu -->
        <div v-if="mobileMenuOpen" class="md:hidden">
            <nav class="flex flex-col space-y-2 py-4 px-6 border-t border-gray-200 dark:border-gray-800">
                <router-link to="/"
                    class="py-2 text-gray-700 hover:text-indigo-600 dark:text-gray-300 dark:hover:text-indigo-400"
                    @click="mobileMenuOpen = false">
                    Script Maker
                </router-link>
                <router-link to="/config"
                    class="py-2 text-gray-700 hover:text-indigo-600 dark:text-gray-300 dark:hover:text-indigo-400"
                    @click="mobileMenuOpen = false">
                    Config
                </router-link>
                <a href="/" 
                    @click.prevent="openRepositoryAndCloseMenu"
                    class="py-2 text-gray-700 hover:text-indigo-600 dark:text-gray-300 dark:hover:text-indigo-400">
                    Github
                </a>
            </nav>
        </div>
    </header>
</template>

<script lang="ts" setup>
import { BrowserOpenURL } from '../../wailsjs/runtime';
import { ref } from 'vue';

const mobileMenuOpen = ref(false);

const openRepository = () => {
  BrowserOpenURL("https://github.com/LouisRiverstone/query-script-maker");
};

const openRepositoryAndCloseMenu = () => {
  openRepository();
  mobileMenuOpen.value = false;
};
</script>