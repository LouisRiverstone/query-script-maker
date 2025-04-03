<template>
    <div class="container mx-auto py-6 px-4 sm:px-6 lg:px-8">
        <div class="bg-white dark:bg-gray-900 shadow overflow-hidden rounded-lg">
            <!-- Loading state -->
            <div v-if="loading" class="flex justify-center items-center py-12">
                <Loader />
            </div>
            
            <!-- Error state -->
            <div v-else-if="hasError" class="py-8 px-4 text-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-red-500 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                </svg>
                <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">Something went wrong</h3>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ errorMessage }}</p>
                <div class="mt-6">
                    <button 
                        @click="retryLoading" 
                        class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                    >
                        Try Again
                    </button>
                </div>
            </div>
            
            <!-- No database connection -->
            <div v-else-if="!hasConnection" class="py-8 px-4 text-center">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-400 dark:text-gray-500 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14M12 5v14M4 20h16a2 2 0 002-2V6a2 2 0 00-2-2H4a2 2 0 00-2 2v12a2 2 0 002 2z" />
                </svg>
                <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">No Database Connected</h3>
                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Please add a database connection in the configuration page.</p>
                <div class="mt-6">
                    <router-link to="/config" class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
                        Go to Config
                    </router-link>
                </div>
            </div>
            
            <!-- Database browser -->
            <div class="h-[calc(100vh-188px)]" v-else>
                <!-- Header with database info -->
                <div class="bg-gradient-to-r from-indigo-600 to-indigo-800 dark:from-indigo-800 dark:to-indigo-900 px-4 py-5 sm:px-6 text-white shadow-lg">
                    <div class="flex flex-col sm:flex-row justify-between">
                        <div>
                            <h3 class="text-lg leading-6 font-medium flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
                                </svg>
                                <span class="truncate">Database Browser: {{ databaseConnection.Database }}</span>
                            </h3>
                            <p class="mt-1 max-w-2xl text-sm opacity-80 flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                                </svg>
                                <span class="truncate">{{ databaseConnection.Username }}@{{ databaseConnection.Host }}:{{ databaseConnection.Port }}</span>
                            </p>
                        </div>
                        <div class="mt-3 sm:mt-0 flex items-center flex-wrap gap-2">
                            <button 
                                @click="showDiagramModal" 
                                class="inline-flex items-center px-3 py-1.5 border border-transparent text-sm font-medium rounded-md shadow-sm text-indigo-600 bg-white hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                                </svg>
                                <span class="hidden sm:inline">Show Diagram</span>
                                <span class="sm:hidden">Diagram</span>
                            </button>
                            <button 
                                @click="refreshStructure" 
                                class="inline-flex items-center px-3 py-1.5 border border-transparent text-sm font-medium rounded-md shadow-sm text-indigo-600 bg-white hover:bg-gray-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                                </svg>
                                <span class="hidden sm:inline">Refresh</span>
                            </button>
                            <button 
                                @click="() => loadDatabaseStructure(true)" 
                                class="inline-flex items-center px-3 py-1.5 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-700 hover:bg-indigo-800 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
                            >
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                                </svg>
                                <span class="hidden sm:inline">Refresh Structures</span>
                                <span class="sm:hidden">Refresh All</span>
                            </button>
                        </div>
                    </div>
                </div>
                
                <!-- Browser layout with responsive drawer -->
                <div class="flex flex-col md:flex-row relative">
                    <!-- Mobile sidebar toggle -->
                    <button 
                        @click="toggleSidebar" 
                        class="md:hidden absolute top-2 right-2 z-10 bg-indigo-600 text-white p-2 rounded-full shadow-lg"
                    >
                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                            <path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
                        </svg>
                    </button>
                
                    <!-- Table list (sidebar) - responsive with transition -->
                    <div 
                        class="w-full md:w-64 border-r border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 transition-all duration-300 ease-in-out overflow-hidden flex flex-col"
                        :class="isSidebarOpen || !isMobile ? 'max-h-[calc(100vh-188px)]' : 'max-h-0 md:max-h-[calc(100vh-188px)]'"
                    >
                        <div class="p-4 border-b border-gray-200 dark:border-gray-700">
                            <h4 class="font-medium text-gray-700 dark:text-gray-300 flex items-center">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2 text-indigo-500" viewBox="0 0 20 20" fill="currentColor">
                                    <path fill-rule="evenodd" d="M3 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 4a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
                                </svg>
                                Tables
                                <span class="ml-2 bg-indigo-100 dark:bg-indigo-900 text-indigo-700 dark:text-indigo-300 text-xs rounded-full px-2 py-0.5">
                                    {{ filteredTables.length }}
                                </span>
                            </h4>
                            <div class="mt-2 relative">
                                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                                    <svg class="h-4 w-4 text-gray-400 dark:text-gray-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                                        <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
                                    </svg>
                                </div>
                                <input 
                                    type="text" 
                                    placeholder="Filter tables..." 
                                    v-model="tableFilter"
                                    class="w-full pl-10 pr-3 py-2 border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
                                />
                            </div>
                            <!-- Enhanced filtering options -->
                            <div class="mt-2 grid grid-cols-2 gap-2">
                                <select
                                    v-model="tableFilterType"
                                    class="block w-full pl-3 pr-10 py-1.5 text-xs border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                >
                                    <option value="all">All tables</option>
                                    <option value="system">System tables</option>
                                    <option value="data">Data tables</option>
                                </select>
                                <select
                                    v-model="tableSortOrder"
                                    class="block w-full pl-3 pr-10 py-1.5 text-xs border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                                >
                                    <option value="name">Sort by name</option>
                                    <option value="columns">Sort by columns</option>
                                    <option value="recent">Recently viewed</option>
                                </select>
                            </div>
                        </div>
                        <div class="overflow-y-auto flex-grow scrollbar-thin scrollbar-thumb-gray-300 dark:scrollbar-thumb-gray-600">
                            <ul class="divide-y divide-gray-200 dark:divide-gray-700">
                                <li 
                                    v-for="table in filteredTables" 
                                    :key="table.name"
                                    class="cursor-pointer hover:bg-gray-100 dark:hover:bg-gray-700 transition-colors duration-150"
                                    :class="selectedTable === table.name ? 'bg-indigo-50 dark:bg-indigo-900/30 border-l-4 border-indigo-500' : ''"
                                >
                                    <div 
                                        @click="toggleTableExpand(table.name)"
                                        class="px-4 py-3 flex items-center justify-between"
                                    >
                                        <div class="flex items-center">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 text-indigo-500 mr-2" viewBox="0 0 20 20" fill="currentColor">
                                                <path fill-rule="evenodd" d="M5 4a3 3 0 00-3 3v6a3 3 0 003 3h10a3 3 0 003-3V7a3 3 0 00-3-3H5zm-1 9v-1h5v2H5a1 1 0 01-1-1zm7 1h4a1 1 0 001-1v-1h-5v2zm0-4h5V8h-5v2zM9 8H4v2h5V8z" clip-rule="evenodd" />
                                            </svg>
                                            <span class="text-sm font-medium text-gray-700 dark:text-gray-300 truncate max-w-[160px]">
                                                {{ table.name }}
                                            </span>
                                        </div>
                                        <div class="flex items-center">
                                            <span class="text-xs bg-indigo-100 dark:bg-indigo-900 text-indigo-700 dark:text-indigo-300 rounded-full px-2 py-0.5 mr-2">
                                                {{ table.columns.length }}
                                            </span>
                                            <svg 
                                                xmlns="http://www.w3.org/2000/svg" 
                                                class="h-4 w-4 transform transition-transform duration-200" 
                                                :class="expandedTables.includes(table.name) ? 'rotate-180' : ''"
                                                fill="none" 
                                                viewBox="0 0 24 24" 
                                                stroke="currentColor"
                                            >
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                                            </svg>
                                        </div>
                                    </div>
                                    <div 
                                        @click="selectTable(table.name)"
                                        class="border-t border-gray-100 dark:border-gray-700 px-4 py-1 bg-gray-50 dark:bg-gray-800/50"
                                        :class="selectedTable === table.name ? 'bg-indigo-50 dark:bg-indigo-900/30' : ''"
                                    >
                                        <div class="flex justify-end">
                                            <button class="text-xs text-indigo-600 dark:text-indigo-400 hover:text-indigo-800 dark:hover:text-indigo-300 px-2 py-1">
                                                Browse
                                            </button>
                                        </div>
                                    </div>
                                    <!-- Collapsible columns section -->
                                    <div v-if="expandedTables.includes(table.name)" class="bg-gray-50 dark:bg-gray-800/30 px-4 pb-2 border-t border-gray-100 dark:border-gray-700">
                                        <input 
                                            type="text" 
                                            v-model="columnFilters[table.name]" 
                                            placeholder="Filter columns..." 
                                            class="w-full mt-1 py-1 px-2 text-xs border border-gray-300 dark:border-gray-600 dark:bg-gray-700 dark:text-white rounded focus:outline-none focus:ring-1 focus:ring-indigo-500"
                                        />
                                        <div class="mt-1 max-h-48 overflow-y-auto scrollbar-thin scrollbar-thumb-gray-300 dark:scrollbar-thumb-gray-600">
                                            <div 
                                                v-for="column in filteredColumns(table)" 
                                                :key="`${table.name}-${column.name}`"
                                                class="flex items-center py-1 text-xs border-b border-gray-100 dark:border-gray-700 last:border-0"
                                            >
                                                <div class="flex-1 flex items-center">
                                                    <span 
                                                        class="w-2 h-2 rounded-full mr-2"
                                                        :class="column.isPrimary ? 'bg-amber-500' : column.key ? 'bg-blue-500' : 'bg-gray-400'"
                                                    ></span>
                                                    <span class="font-medium text-gray-700 dark:text-gray-300">{{ column.name }}</span>
                                                </div>
                                                <div class="text-gray-500 dark:text-gray-400 truncate max-w-[100px]">
                                                    {{ column.type }}
                                                </div>
                                            </div>
                                        </div>
                                    </div>
                                </li>
                            </ul>
                        </div>
                    </div>
                    
                    <!-- Main content area -->
                    <div class="flex-1 overflow-x-auto">
                        <div v-if="!selectedTable" class="flex items-center justify-center h-[calc(100vh-300px)]">
                            <div class="text-center px-4">
                                <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-400 dark:text-gray-500 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 9l4-4 4 4m0 6l-4 4-4-4" />
                                </svg>
                                <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100">Select a table</h3>
                                <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">Choose a table from the list to view its data</p>
                                
                                <!-- Mobile-only table selection -->
                                <div class="md:hidden mt-4">
                                    <button
                                        @click="toggleSidebar"
                                        class="px-4 py-2 bg-indigo-600 text-white rounded-md shadow-sm hover:bg-indigo-700 transition-colors duration-200"
                                    >
                                        Open Table List
                                    </button>
                                </div>
                            </div>
                        </div>
                        
                        <div v-else>
                            <!-- Tabs - scrollable on mobile -->
                            <div class="px-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800 overflow-x-auto">
                                <div class="flex space-x-8 whitespace-nowrap">
                                    <button 
                                        @click="activeTab = 'data'" 
                                        class="py-4 border-b-2 font-medium text-sm flex items-center"
                                        :class="activeTab === 'data' ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400' : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'"
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
                                        </svg>
                                        Browse
                                    </button>
                                     <button 
                                        @click="activeTab = 'structure'" 
                                        class="py-4 border-b-2 font-medium text-sm flex items-center"
                                        :class="activeTab === 'structure' ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400' : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'"
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10" />
                                        </svg>
                                        Structure
                                    </button>
                                    <button 
                                        @click="activeTab = 'sql'" 
                                        class="py-4 border-b-2 font-medium text-sm flex items-center"
                                        :class="activeTab === 'sql' ? 'border-indigo-500 text-indigo-600 dark:text-indigo-400' : 'border-transparent text-gray-500 hover:text-gray-700 dark:text-gray-400 dark:hover:text-gray-300'"
                                    >
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 20l4-16m4 4l4 4-4 4M6 16l-4-4 4-4" />
                                        </svg>
                                        SQL
                                    </button>
                                </div>
                            </div>
                            
                            <div class="p-4 border-b border-gray-200 dark:border-gray-700 bg-gray-50 dark:bg-gray-800">
                                <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
                                    <h3 class="text-lg font-medium text-gray-900 dark:text-gray-100 flex items-center mb-3 sm:mb-0">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-indigo-500 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                                            <path fill-rule="evenodd" d="M5 4a3 3 0 00-3 3v6a3 3 0 003 3h10a3 3 0 003-3V7a3 3 0 00-3-3H5zm-1 9v-1h5v2H5a1 1 0 01-1-1zm7 1h4a1 1 0 001-1v-1h-5v2zm0-4h5V8h-5v2zM9 8H4v2h5V8z" clip-rule="evenodd" />
                                        </svg>
                                        <span class="truncate max-w-[200px] sm:max-w-full">{{ selectedTable }}</span>
                                    </h3>
                                    <div class="flex flex-wrap gap-2">
                                        <button 
                                            @click="exportTableData" 
                                            class="px-3 py-1.5 text-xs font-medium rounded-md bg-green-600 text-white hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-colors duration-200 flex items-center"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                                            </svg>
                                            <span class="hidden sm:inline">Export</span>
                                        </button>
                                        <button 
                                            @click="refreshTableData()" 
                                            class="px-3 py-1.5 text-xs font-medium rounded-md bg-indigo-600 text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200 flex items-center"
                                        >
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                                            </svg>
                                            <span class="hidden sm:inline">Refresh</span>
                                        </button>
                                        <select 
                                            v-model="rowLimit" 
                                            @change="onRowLimitChange"
                                            class="px-3 py-1.5 text-xs font-medium rounded-md bg-white dark:bg-gray-700 border border-gray-300 dark:border-gray-600 text-gray-700 dark:text-gray-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200"
                                        >
                                            <option value="10">10 rows</option>
                                            <option value="25">25 rows</option>
                                            <option value="50">50 rows</option>
                                            <option value="100">100 rows</option>
                                        </select>
                                    </div>
                                </div>
                            </div>

                            <div class="p-4">
                                <!-- Data Tab -->
                                <div v-if="activeTab === 'data'">
                                    <div v-if="tableDataLoading" class="flex justify-center items-center py-12">
                                        <Loader />
                                    </div>
                                    <div v-else-if="tableData.length === 0" class="text-center py-12">
                                        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 mx-auto text-gray-400 dark:text-gray-500 mb-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
                                        </svg>
                                        <p class="text-gray-500 dark:text-gray-400">No data found in this table</p>
                                    </div>
                                    <div v-else>
                                        <div class="border border-gray-200 dark:border-gray-700 rounded-md shadow-sm overflow-x-auto">
                                            <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                                                <thead class="bg-gray-50 dark:bg-gray-800">
                                                    <tr>
                                                        <th v-for="column in tableColumns" :key="column"
                                                            class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider whitespace-nowrap">
                                                            {{ column }}
                                                        </th>
                                                    </tr>
                                                </thead>
                                                <tbody class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
                                                    <tr v-for="(row, rowIndex) in tableData" :key="rowIndex" class="hover:bg-gray-50 dark:hover:bg-gray-800/70 transition-colors duration-150">
                                                        <td v-for="column in tableColumns" :key="`${rowIndex}-${column}`"
                                                            class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400 max-w-[200px] truncate">
                                                            <span v-if="row[column] === null" class="text-gray-400 dark:text-gray-500 italic">NULL</span>
                                                            <span v-else-if="typeof row[column] === 'string' && row[column].length > 100" class="cursor-pointer hover:underline" @click="showFullCellContent(row[column])">
                                                                {{ row[column].substring(0, 100) }}...
                                                            </span>
                                                            <span v-else>{{ row[column] }}</span>
                                                        </td>
                                                    </tr>
                                                </tbody>
                                            </table>
                                        </div>
                                        
                                        <!-- Enhanced pagination -->
                                        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mt-4 gap-3">
                                            <div class="text-sm text-gray-700 dark:text-gray-300 flex items-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1 text-indigo-500 flex-shrink-0" viewBox="0 0 20 20" fill="currentColor">
                                                    <path d="M7 3a1 1 0 000 2h6a1 1 0 100-2H7zM4 7a1 1 0 011-1h10a1 1 0 110 2H5a1 1 0 01-1-1zM2 11a2 2 0 012-2h12a2 2 0 012 2v4a2 2 0 01-2 2H4a2 2 0 01-2-2v-4z" />
                                                </svg>
                                                <span class="text-xs sm:text-sm">
                                                    Showing <span class="font-medium">{{ tableData.length }}</span> rows
                                                    <span v-if="tableDataTotalRows[`${databaseConnection.Database}_${selectedTable}`]" class="ml-1">
                                                        of <span class="font-medium">{{ tableDataTotalRows[`${databaseConnection.Database}_${selectedTable}`] }}</span> total
                                                    </span>
                                                    <span v-if="currentPage > 0" class="ml-1">
                                                        - Page <span class="font-medium">{{ currentPage + 1 }}</span>
                                                    </span>
                                                </span>
                                            </div>
                                            <div class="flex space-x-2">
                                                <button
                                                    @click="loadPreviousPage"
                                                    :disabled="currentPage === 0"
                                                    class="px-3 py-1 text-sm font-medium rounded-md transition-colors duration-200 flex items-center"
                                                    :class="currentPage === 0 
                                                        ? 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-not-allowed' 
                                                        : 'bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300 hover:bg-indigo-200'"
                                                >
                                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
                                                    </svg>
                                                    Previous
                                                </button>
                                                <button
                                                    @click="loadNextPage"
                                                    :disabled="!hasMorePages"
                                                    class="px-3 py-1 text-sm font-medium rounded-md transition-colors duration-200 flex items-center"
                                                    :class="!hasMorePages 
                                                        ? 'bg-gray-200 dark:bg-gray-700 text-gray-500 dark:text-gray-400 cursor-not-allowed' 
                                                        : 'bg-indigo-100 dark:bg-indigo-900/50 text-indigo-700 dark:text-indigo-300 hover:bg-indigo-200'"
                                                >
                                                    Next
                                                    <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                                                    </svg>
                                                </button>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                                
                                <!-- SQL Tab optimized for mobile -->
                                <div v-else-if="activeTab === 'sql'">
                                    <div class="bg-gray-100 dark:bg-gray-800 rounded-lg p-4 mb-4">
                                        <h4 class="text-sm font-medium text-gray-700 dark:text-gray-300 mb-2 flex items-center">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-2 text-indigo-500" viewBox="0 0 20 20" fill="currentColor">
                                                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM9.555 7.168A1 1 0 008 8v4a1 1 0 001.555.832l3-2a1 1 0 000-1.664l-3-2z" clip-rule="evenodd" />
                                            </svg>
                                            SQL Query Editor
                                        </h4>
                                        <div class="mb-2">
                                            <Editor 
                                                v-model="sqlQuery" 
                                                placeholder="Enter your SQL query here..." 
                                                :height="isMobile ? '120px' : '180px'"
                                                class="dark:bg-gray-900 border border-gray-300 dark:border-gray-700 rounded-md"
                                            />
                                        </div>
                                        <div class="flex flex-col sm:flex-row sm:justify-between gap-2">
                                            <div class="text-xs text-gray-500 dark:text-gray-400 flex items-center">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                                </svg>
                                                <span class="text-xs">Use semicolons to separate queries</span>
                                            </div>
                                            <Button 
                                                type="button" 
                                                class="bg-indigo-600 hover:bg-indigo-700 flex items-center justify-center" 
                                                @click="executeCustomQuery"
                                                :disabled="sqlQueryLoading"
                                            >
                                                <svg v-if="!sqlQueryLoading" xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
                                                </svg>
                                                <span v-if="!sqlQueryLoading">Execute Query</span>
                                                <Loader v-else />
                                            </Button>
                                        </div>
                                    </div>

                                    <div v-if="sqlQueryResults.length > 0" class="mt-6">
                                        <h4 class="text-lg font-medium text-gray-900 dark:text-gray-100 mb-4 flex items-center">
                                            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                                            </svg>
                                            Query Results
                                        </h4>
                                        <div v-for="(result, resultIndex) in sqlQueryResults" :key="resultIndex" class="mb-6">
                                            <div v-if="Array.isArray(result) && result.length > 0 && !result[0].error" class="border border-gray-200 dark:border-gray-700 rounded-md shadow-sm overflow-x-auto">
                                                <div class="bg-gray-50 dark:bg-gray-800 px-4 py-2 border-b border-gray-200 dark:border-gray-700 flex justify-between items-center">
                                                    <span class="text-xs font-medium text-gray-500 dark:text-gray-400">Result set #{{ resultIndex + 1 }}</span>
                                                    <span class="text-xs text-gray-500 dark:text-gray-400">{{ result.length }} rows</span>
                                                </div>
                                                <table class="min-w-full divide-y divide-gray-200 dark:divide-gray-700">
                                                    <thead class="bg-gray-50 dark:bg-gray-800">
                                                        <tr>
                                                            <th v-for="column in Object.keys(result[0])" :key="column"
                                                                class="px-6 py-3 text-left text-xs font-medium text-gray-500 dark:text-gray-400 uppercase tracking-wider">
                                                                {{ column }}
                                                            </th>
                                                        </tr>
                                                    </thead>
                                                    <tbody class="bg-white dark:bg-gray-900 divide-y divide-gray-200 dark:divide-gray-700">
                                                        <tr v-for="(row, rowIndex) in result" :key="rowIndex" class="hover:bg-gray-50 dark:hover:bg-gray-800/70">
                                                            <td v-for="column in Object.keys(result[0])" :key="`${rowIndex}-${column}`"
                                                                class="px-6 py-4 whitespace-nowrap text-sm text-gray-500 dark:text-gray-400">
                                                                <span v-if="row[column] === null" class="text-gray-400 dark:text-gray-500 italic">NULL</span>
                                                                <span v-else>{{ row[column] }}</span>
                                                            </td>
                                                        </tr>
                                                    </tbody>
                                                </table>
                                            </div>
                                            <div v-else-if="result[0] && result[0].error" class="text-center py-6 bg-red-50 dark:bg-red-900/20 border border-red-200 dark:border-red-800 rounded-md text-red-700 dark:text-red-300">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto mb-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                                </svg>
                                                <p>{{ result[0].error }}</p>
                                            </div>
                                            <div v-else class="text-center py-6 bg-green-50 dark:bg-green-900/20 border border-green-200 dark:border-green-800 rounded-md">
                                                <svg xmlns="http://www.w3.org/2000/svg" class="h-6 w-6 mx-auto mb-2 text-green-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
                                                </svg>
                                                <p class="text-green-700 dark:text-green-300">
                                                    Query executed successfully ({{ Array.isArray(result) ? result.length : 0 }} rows affected)
                                                </p>
                                            </div>
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                </div>
            </div>
        </div>
    </div>
    
    <!-- Database diagram modal -->
    <DatabaseDiagramModal 
        :isOpen="isDiagramModalOpen" 
        :databaseStructure="databaseStructureForDiagram" 
        @close="isDiagramModalOpen = false"
        @refresh="handleDiagramRefresh"
    />
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch, onBeforeMount, nextTick, onBeforeUnmount } from 'vue';
import { GetDatabaseConnection, GetDatabaseStructure, TestQueryInDatabase, GetLatestDatabaseStructure } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import Loader from '../components/Loader.vue';
import Button from '../components/Button.vue';
import Editor from '../components/Editor.vue';
import DatabaseDiagramModal from '../components/DatabaseDiagramModal.vue';

// Definição de tipos para melhorar compatibilidade
interface TableDataCacheType {
  [key: string]: Record<number, any[]>;
}

interface TableDataTotalRowsType {
  [key: string]: number;
}

interface StructureCacheType {
  [key: string]: any;
}

const prefetching = ref<boolean>(true);
const loading = ref<boolean>(true);
const componentReady = ref<boolean>(false);
const hasConnection = ref<boolean>(false);
const hasError = ref<boolean>(false);
const errorMessage = ref<string>('An unexpected error occurred. Please try again.');
const databaseConnection = ref<main.DatabaseConnection>({
    Username: '',
    Password: '',
    Host: '',
    Port: 3306,
    Database: ''
});

const dbStructure = ref<any>(null);
const tableFilter = ref<string>('');
const selectedTable = ref<string>('');
const activeTab = ref<string>('data');

// Table data browsing
const tableData = ref<any[]>([]);
const tableColumns = ref<string[]>([]);
const tableDataLoading = ref<boolean>(false);
const rowLimit = ref<string>('25');
const currentPage = ref<number>(0);
const tableDataCache = ref<TableDataCacheType>({}); // Cache para dados de tabela por página
const tableDataTotalRows = ref<TableDataTotalRowsType>({}); // Armazenar o total de linhas por tabela

// SQL query tab
const sqlQuery = ref<string>('');
const sqlQueryResults = ref<any[]>([]);
const sqlQueryLoading = ref<boolean>(false);

// Cache da estrutura
const structureCache = ref<StructureCacheType>({});
const lastRefreshTime = ref<number>(0);
const CACHE_TTL = 5 * 60 * 1000; // 5 minutes in milliseconds

const tableFilterType = ref<string>('all');
const tableSortOrder = ref<string>('name');
const expandedTables = ref<string[]>([]);
const columnFilters = ref<{[key: string]: string}>({});
const recentlyViewedTables = ref<string[]>([]);

// Toggle table expansion
const toggleTableExpand = (tableName: string) => {
    if (expandedTables.value.includes(tableName)) {
        expandedTables.value = expandedTables.value.filter(t => t !== tableName);
    } else {
        expandedTables.value.push(tableName);
        // Initialize column filter if not exists
        if (!columnFilters.value[tableName]) {
            columnFilters.value[tableName] = '';
        }
    }
};

// Filtered columns based on search
const filteredColumns = (table: any) => {
    const filter = columnFilters.value[table.name] || '';
    if (!filter) return table.columns;
    
    return table.columns.filter((column: any) => 
        column.name.toLowerCase().includes(filter.toLowerCase()) ||
        column.type.toLowerCase().includes(filter.toLowerCase())
    );
};

// Optimized and enhanced table filtering based on multiple criteria
const filteredTables = computed(() => {
    if (!dbStructure.value || !dbStructure.value.tables) return [];
    
    // Step 1: Apply name filter
    let result = dbStructure.value.tables;
    if (tableFilter.value) {
        result = result.filter((table: any) => 
            table.name.toLowerCase().includes(tableFilter.value.toLowerCase())
        );
    }
    
    // Step 2: Apply type filter
    if (tableFilterType.value !== 'all') {
        if (tableFilterType.value === 'system') {
            // System tables typically start with these prefixes
            const systemPrefixes = ['mysql', 'information_schema', 'performance_schema', 'sys', 'pg_'];
            result = result.filter((table: any) => 
                systemPrefixes.some(prefix => table.name.toLowerCase().startsWith(prefix))
            );
        } else if (tableFilterType.value === 'data') {
            // Non-system tables
            const systemPrefixes = ['mysql', 'information_schema', 'performance_schema', 'sys', 'pg_'];
            result = result.filter((table: any) => 
                !systemPrefixes.some(prefix => table.name.toLowerCase().startsWith(prefix))
            );
        }
    }
    
    // Step 3: Apply sorting
    result = [...result].sort((a: any, b: any) => {
        if (tableSortOrder.value === 'columns') {
            return b.columns.length - a.columns.length;
        } else if (tableSortOrder.value === 'recent') {
            const aIndex = recentlyViewedTables.value.indexOf(a.name);
            const bIndex = recentlyViewedTables.value.indexOf(b.name);
            
            // If both tables are in recently viewed, sort by recency
            if (aIndex !== -1 && bIndex !== -1) {
                return aIndex - bIndex;
            }
            // If only one table is in recently viewed, prioritize it
            if (aIndex !== -1) return -1;
            if (bIndex !== -1) return 1;
            
            // Otherwise, fall back to name sorting
            return a.name.localeCompare(b.name);
        } else {
            // Default to name sorting
            return a.name.localeCompare(b.name);
        }
    });
    
    return result;
});

// Get columns for the selected table
const selectedTableColumns = computed(() => {
    if (!selectedTable.value || !dbStructure.value || !dbStructure.value.tables) return [];
    
    const table = dbStructure.value.tables.find((t: any) => t.name === selectedTable.value);
    return table ? table.columns : [];
});

// Retry loading after an error
const retryLoading = async () => {
    hasError.value = false;
    loading.value = true;
    
    try {
        await loadDatabaseConnection();
        if (hasConnection.value) {
            await loadDatabaseStructure();
        }
    } catch (error) {
        console.error('Failed to retry loading:', error);
        hasError.value = true;
        errorMessage.value = error instanceof Error 
            ? `Error: ${error.message}` 
            : 'Failed to connect to the database. Please check your connection.';
        dbStructure.value = { tables: [] };
    } finally {
        loading.value = false;
    }
};

// Enhanced loading mechanism
onBeforeMount(() => {
    prefetching.value = true;
    loading.value = true;
});

// Initialize the page with improved transitions
onMounted(async () => {
    // Add a small delay for smoother transitions
    setTimeout(() => {
        prefetching.value = false;
    }, 50);
    
    // Add a safety timeout to prevent endless loading
    const timeout = setTimeout(() => {
        if (loading.value) {
            loading.value = false;
            hasError.value = true;
            errorMessage.value = 'Loading timed out. The database might be unavailable.';
            dbStructure.value = { tables: [] };
            console.error('Loading timed out');
        }
    }, 5 * 60 * 1000); // 5 minutes timeout
    
    try {
        await loadDatabaseConnection();
        if (hasConnection.value) {
            await loadDatabaseStructure();
        }
    } catch (error) {
        console.error('Failed to initialize database browser:', error);
        hasError.value = true;
        errorMessage.value = error instanceof Error 
            ? `Error: ${error.message}` 
            : 'Failed to load the database browser.';
        dbStructure.value = { tables: [] };
    } finally {
        loading.value = false;
        clearTimeout(timeout);
        
        // Mark component as ready after a slight delay for smoother transition
        setTimeout(() => {
            componentReady.value = true;
        }, 100);
    }
});

// Load database connection
const loadDatabaseConnection = async () => {
    try {
        const connection = await GetDatabaseConnection();
        if (connection && connection.ID) {
            databaseConnection.value = connection;
            hasConnection.value = true;
        }
    } catch (error) {
        console.error('Failed to load database connection:', error);
        hasConnection.value = false;
    }
};

// Optimize query execution with cancellation support
const executeQueryWithTimeout = async (query: string, timeout: number = 5 * 60 * 1000): Promise<any> => {
    // Use AbortController for better cancellation control
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), timeout);
    
    try {
        // Execute query with abort signal
        const queryPromise = TestQueryInDatabase(databaseConnection.value, query, false);
        
        // Race between the query and the timeout
        const result = await Promise.race<any>([
            queryPromise,
            new Promise<never>((_, reject) => 
                setTimeout(() => reject(new Error('Query timed out')), timeout)
            )
        ]);
        
        clearTimeout(timeoutId);
        return result;
    } catch (error) {
        console.error('Query execution error:', error);
        throw error;
    } finally {
        clearTimeout(timeoutId);
    }
};

// Load database structure with better performance
const loadDatabaseStructure = async (forceRefresh = false) => {
    try {
        // Check if we have a valid cached structure
        const cacheKey = `${databaseConnection.value.Host}_${databaseConnection.value.Port}_${databaseConnection.value.Database}`;
        const now = Date.now();
        
        // Use the cache if it exists, refresh is not forced, and it hasn't expired
        if (!forceRefresh && 
            structureCache.value[cacheKey] && 
            (now - lastRefreshTime.value) < CACHE_TTL) {
            console.log('Using cached database structure');
            dbStructure.value = structureCache.value[cacheKey];
            return;
        }
        
        console.log('Fetching fresh database structure');
        
        // First try to load the most recent structure from the server's local cache
        const cachedStructure = await GetLatestDatabaseStructure()
            .catch(err => {
                console.warn('Failed to load cached structure:', err);
                return null;
            });
            
        if (cachedStructure && cachedStructure.length > 0) {
            try {
                dbStructure.value = JSON.parse(cachedStructure);
                structureCache.value[cacheKey] = dbStructure.value;
                lastRefreshTime.value = now;
                
                // If refresh wasn't forced, use the cached structure and return
                if (!forceRefresh) {
                    return;
                }
                // If refresh was forced, continue and fetch updated structure
            } catch (parseError) {
                console.error('Failed to parse cached structure:', parseError);
            }
        }
        
        // Load structure with timeout
        const structure = await Promise.race<string>([
            GetDatabaseStructure(databaseConnection.value),
            new Promise<string>((_, reject) => 
                setTimeout(() => reject(new Error('Database structure loading timed out')), 5 * 60 * 1000)
            )
        ]);
        
        if (!structure) {
            console.warn('No database structure returned');
            dbStructure.value = { tables: [] };
            return;
        }
        
        try {
            dbStructure.value = JSON.parse(structure);
            structureCache.value[cacheKey] = dbStructure.value;
            lastRefreshTime.value = now;
        } catch (parseError) {
            console.error('Failed to parse database structure:', parseError);
            dbStructure.value = { tables: [] };
        }
    } catch (error) {
        console.error('Failed to load database structure:', error);
        dbStructure.value = { tables: [] };
    }
};

// Refresh database structure
const refreshStructure = async () => {
    loading.value = true;
    try {
        await loadDatabaseStructure(true); // Forçar atualização
        // Reset selected table if it no longer exists
        if (selectedTable.value && dbStructure.value && dbStructure.value.tables && 
            !dbStructure.value.tables.some((t: any) => t.name === selectedTable.value)) {
            selectedTable.value = '';
        }
    } catch (error) {
        console.error('Failed to refresh database structure:', error);
        dbStructure.value = { tables: [] };
    } finally {
        loading.value = false;
    }
};

// Select a table
const selectTable = async (tableName: string) => {
    // Close sidebar on mobile when selecting a table
    if (isMobile.value) {
        isSidebarOpen.value = false;
    }
    
    // Limpar cache de dados de tabela ao mudar de tabela
    if (selectedTable.value !== tableName) {
        tableData.value = [];
        tableColumns.value = [];
        currentPage.value = 0;
        
        // Track recently viewed tables
        recentlyViewedTables.value = recentlyViewedTables.value.filter(t => t !== tableName);
        recentlyViewedTables.value.unshift(tableName);
        // Keep only the 10 most recent tables
        if (recentlyViewedTables.value.length > 10) {
            recentlyViewedTables.value = recentlyViewedTables.value.slice(0, 10);
        }
    }
    
    selectedTable.value = tableName;
    activeTab.value = 'data';
    
    // Pre-populate SQL query with SELECT statement
    sqlQuery.value = `SELECT * FROM \`${tableName}\` LIMIT 100;`;
    
    // Load data if data tab is active or when tab becomes active
    if (activeTab.value === 'data') {
        await refreshTableData();
    }
};

// Refresh table data with optimized query execution
const refreshTableData = async (forceRefresh = false) => {
    if (!selectedTable.value) return;
    
    const table = selectedTable.value;
    const page = currentPage.value;
    const limit = parseInt(rowLimit.value);
    const cacheKey = `${databaseConnection.value.Database}_${table}`;
    
    // Use cached data if available and refresh is not forced
    if (!forceRefresh && 
        tableDataCache.value[cacheKey] && 
        tableDataCache.value[cacheKey][page]) {
        console.log(`Using cached data for ${table} page ${page}`);
        tableData.value = tableDataCache.value[cacheKey][page];
        if (tableData.value.length > 0) {
            tableColumns.value = Object.keys(tableData.value[0]);
        }
        return;
    }
    
    tableDataLoading.value = true;
    
    try {
        // Get total row count in parallel (but only if we don't already have it)
        if (!tableDataTotalRows.value[cacheKey]) {
            const countQuery = `SELECT COUNT(*) AS total FROM \`${table}\``;
            executeQueryWithTimeout(countQuery)
                .then(countResult => {
                    if (countResult && countResult.length > 0 && countResult[0].total !== undefined) {
                        tableDataTotalRows.value[cacheKey] = parseInt(countResult[0].total);
                    }
                })
                .catch(err => console.warn('Failed to get row count:', err));
        }
        
        // Build optimized query
        let query = `SELECT * FROM \`${table}\` LIMIT ${limit} OFFSET ${page * limit}`;
        
        // Optimization for large tables
        if (dbStructure.value && dbStructure.value.tables) {
            const tableInfo = dbStructure.value.tables.find((t: any) => t.name === table);
            if (tableInfo && tableInfo.columns && tableInfo.columns.length > 15) {
                // For tables with many columns, select only the first 15
                // and add ID/primary key if available
                const primaryKey = tableInfo.columns.find((c: any) => c.isPrimary)?.name;
                const selectedColumns = tableInfo.columns.slice(0, 15).map((c: any) => `\`${c.name}\``);
                
                if (primaryKey && !selectedColumns.includes(`\`${primaryKey}\``)) {
                    selectedColumns.unshift(`\`${primaryKey}\``);
                }
                
                query = `SELECT ${selectedColumns.join(', ')} FROM \`${table}\` LIMIT ${limit} OFFSET ${page * limit}`;
            }
        }
        
        // Execute query with timeout handling
        const result = await executeQueryWithTimeout(query);
        
        if (result && Array.isArray(result) && result.length > 0) {
            tableData.value = result;
            tableColumns.value = Object.keys(result[0]);
            
            // Store in cache
            if (!tableDataCache.value[cacheKey]) {
                tableDataCache.value[cacheKey] = {};
            }
            tableDataCache.value[cacheKey][page] = result;
        } else {
            tableData.value = [];
            tableColumns.value = [];
        }
    } catch (error) {
        console.error('Failed to load table data:', error);
        tableData.value = [];
        tableColumns.value = [];
    } finally {
        tableDataLoading.value = false;
    }
};

// Computed to check if there are more pages
const hasMorePages = computed(() => {
    // If no table is selected, there are no more pages
    if (!selectedTable.value) {
        return false;
    }
    
    const cacheKey = `${databaseConnection.value.Database}_${selectedTable.value}`;
    const limit = parseInt(rowLimit.value);
    
    // If we know the total number of rows in the table
    if (tableDataTotalRows.value[cacheKey]) {
        const totalRows = tableDataTotalRows.value[cacheKey];
        const currentRows = (currentPage.value + 1) * limit;
        return currentRows < totalRows;
    }
    
    // Otherwise, check if the current page has the maximum number of rows
    return tableData.value.length >= limit;
});

// Load previous page
const loadPreviousPage = async () => {
    if (currentPage.value > 0) {
        currentPage.value -= 1;
        await refreshTableData();
    }
};

// Load next page
const loadNextPage = async () => {
    // Verificar se tem mais páginas antes de carregar a próxima
    const canLoadMore = hasMorePages.value;
    
    if (canLoadMore) {
        currentPage.value += 1;
        await refreshTableData();
    }
};

// Execute custom SQL query
const executeCustomQuery = async () => {
    if (!sqlQuery.value.trim()) return;
    
    sqlQueryLoading.value = true;
    try {
        // Check if it's a multi-query (separated by semicolons)
        const queries = sqlQuery.value
            .split(';')
            .map(q => q.trim())
            .filter(q => q.length > 0);
        
        const results = [];
        
        for (const query of queries) {
            try {
                // Add a safety timeout for each query
                const timeoutPromise = new Promise((_, reject) => 
                    setTimeout(() => reject(new Error('Query timed out')), 5 * 60 * 1000)
                );
                
                const queryPromise = TestQueryInDatabase(databaseConnection.value, query, false);
                
                // Race between the query and the timeout
                const result = await Promise.race([queryPromise, timeoutPromise]);
                results.push(result);
            } catch (error) {
                console.error(`Error executing query: ${query}`, error);
                // Add a message to show the user that the query failed
                results.push([{ error: error instanceof Error ? error.message : 'Query execution failed' }]);
            }
        }
        
        sqlQueryResults.value = results;
    } catch (error) {
        console.error('Failed to execute SQL query:', error);
        sqlQueryResults.value = [];
    } finally {
        sqlQueryLoading.value = false;
    }
};

// Export table data to CSV
const exportTableData = async () => {
    if (!selectedTable.value || tableData.value.length === 0) return;
    
    try {
        // Create CSV content
        const headers = tableColumns.value;
        const csvRows = [];
        
        // Add headers
        csvRows.push(headers.join(','));
        
        // Add data rows
        for (const row of tableData.value) {
            const values = headers.map(header => {
                const value = row[header];
                // Handle null values, quotes, and commas in values
                if (value === null) return 'NULL';
                if (typeof value === 'string') {
                    // Escape quotes and wrap in quotes if needed
                    if (value.includes(',') || value.includes('"') || value.includes('\n')) {
                        return `"${value.replace(/"/g, '""')}"`;
                    }
                    return value;
                }
                return value;
            });
            csvRows.push(values.join(','));
        }
        
        // Combine into a single CSV string
        const csvContent = csvRows.join('\n');
        
        // Create a blob and download
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.setAttribute('href', url);
        link.setAttribute('download', `${selectedTable.value}_export.csv`);
        link.style.visibility = 'hidden';
        document.body.appendChild(link);
        link.click();
        document.body.removeChild(link);
    } catch (error) {
        console.error('Failed to export table data:', error);
        alert('Failed to export data. Please try again.');
    }
};

// Handler for row limit change
const onRowLimitChange = () => {
    // Reset to first page when changing the number of rows
    currentPage.value = 0;
    refreshTableData();
};

// Observe mudanças na conexão de banco e recarregue a estrutura quando necessário
watch(() => databaseConnection.value, (newConn, oldConn) => {
    if (newConn.ID !== oldConn.ID || 
        newConn.Host !== oldConn.Host || 
        newConn.Port !== oldConn.Port || 
        newConn.Database !== oldConn.Database) {
        loadDatabaseStructure(true);
    }
}, { deep: true });

// Observar mudanças na guia ativa
watch(() => activeTab.value, async (newTab) => {
    if (newTab === 'data' && selectedTable.value && tableData.value.length === 0) {
        await refreshTableData();
    }
});

// Database diagram modal state
const isDiagramModalOpen = ref<boolean>(false);

// Safely prepare database structure for the diagram - use computed with cache to prevent unnecessary re-renders
const databaseStructureForDiagram = computed(() => {
    // Only stringify when the modal is open to avoid unnecessary computation
    if (!isDiagramModalOpen.value) {
        return "{}";
    }
    
    try {
        return JSON.stringify(dbStructure.value || { tables: [] });
    } catch (error) {
        console.error('Error preparing database structure for diagram:', error);
        return JSON.stringify({ tables: [] });
    }
});

// Show diagram modal
const showDiagramModal = () => {
    isDiagramModalOpen.value = true;
};

// Handle refresh from diagram modal with debounce
let refreshStructureTimer: number | null = null;
const handleDiagramRefresh = (newStructure: string) => {
    // Prevent multiple rapid updates
    if (refreshStructureTimer) {
        clearTimeout(refreshStructureTimer);
    }
    
    refreshStructureTimer = window.setTimeout(() => {
        try {
            dbStructure.value = JSON.parse(newStructure);
            refreshStructureTimer = null;
        } catch (error) {
            console.error('Failed to parse refreshed structure:', error);
        }
    }, 200);
};

// Add mobile responsiveness
const isMobile = ref<boolean>(false);
const isSidebarOpen = ref<boolean>(false);

const checkMobileState = () => {
    isMobile.value = window.innerWidth < 768;
    // If switching to desktop, ensure sidebar is visible
    if (!isMobile.value) {
        isSidebarOpen.value = true;
    }
};

const toggleSidebar = () => {
    isSidebarOpen.value = !isSidebarOpen.value;
};

// Handle responsive layout
onMounted(() => {
    checkMobileState();
    window.addEventListener('resize', checkMobileState);
});

// Clean up event listeners
onBeforeUnmount(() => {
    window.removeEventListener('resize', checkMobileState);
});

// Function to show full cell content in a modal
const showFullCellContent = (content: string) => {
    if (!content) return;
    
    // Simple alert for now, but could be replaced with a proper modal
    alert(content);
};
</script>

<style scoped>
/* Smooth fade-in transitions for the database browser components */
.container {
    transition: opacity 0.3s ease-in-out;
}

.fade-enter-active,
.fade-leave-active {
    transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
    opacity: 0;
}

/* Table-specific transitions */
.table-fade-enter-active,
.table-fade-leave-active {
    transition: all 0.2s ease-in-out;
}

.table-fade-enter-from,
.table-fade-leave-to {
    opacity: 0;
    transform: translateY(10px);
}

/* Remove flash of unstyled content */
[v-cloak] {
    display: none;
}
</style> 