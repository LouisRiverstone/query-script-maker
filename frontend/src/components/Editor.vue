<template>
    <div class="flex flex-col gap-3 w-full" :class="{ 'editor-dark-mode': isDarkMode }">
        <slot></slot>
        <div class="flex flex-col lg:flex-row gap-6">
            <div class="flex flex-col w-full gap-3">
                <div class="flex justify-between items-center">
                    <Divider v-if="showBindedSql">Input</Divider>
                    <button 
                        @click="toggleDarkMode" 
                        title="Toggle dark mode (debug only)"
                        class="ml-auto text-xs text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 px-2 py-0.5 rounded"
                    >
                        {{ isDarkMode ? '☀️' : '🌙' }}
                    </button>
                </div>
                <small class="text-gray-600 dark:text-gray-400 text-xs italic mb-1 px-1">To use variables, use: {{ variableText }}</small>
                <code-mirror v-model="value" :lang="lang" :extensions="extensions" :linter="null" basic wrap tab class="w-full" />
            </div>
            <div v-if="showBindedSql" class="flex flex-col w-full gap-3">
                <Divider>Output</Divider>
                <code-mirror v-model="linesBinded" :lang="lang" :extensions="outputExtensions" :linter="null" basic wrap tab class="w-full" />
            </div>
        </div>
    </div>
</template>

<script lang="ts" setup>
import { ref, watch, onMounted, computed, onUnmounted } from 'vue';
import CodeMirror from 'vue-codemirror6';
import { sql, MySQL } from "@codemirror/lang-sql";
import { oneDarkTheme } from '@codemirror/theme-one-dark';
import { MakeBindedSQL } from '../../wailsjs/go/main/App';
import { main } from '../../wailsjs/go/models';
import { computedAsync } from '@vueuse/core';
import { autocompletion, CompletionContext, CompletionResult } from '@codemirror/autocomplete';
import { EditorView, Decoration, ViewPlugin, ViewUpdate, DecorationSet } from '@codemirror/view';

import Divider from './Divider.vue';

const variableText = ref("{{ variable }}");
const isDarkMode = ref(false);
let darkModeObserver: MutationObserver;

// Check if dark mode is enabled
const updateDarkMode = () => {
    // Check for .dark class on html or body
    const hasDarkClass = 
        document.documentElement.classList.contains('dark') || 
        document.body.classList.contains('dark');
    
    // Check for prefers-color-scheme media query if available
    const prefersDark = window.matchMedia && 
        window.matchMedia('(prefers-color-scheme: dark)').matches;
    
    // Check for data-theme="dark" attribute
    const hasDataThemeDark = 
        document.documentElement.getAttribute('data-theme') === 'dark' || 
        document.body.getAttribute('data-theme') === 'dark';
    
    // Check for any custom dark mode implementation (like .tw-dark)
    const hasCustomDarkClass = 
        document.documentElement.classList.contains('tw-dark') || 
        document.body.classList.contains('tw-dark');
    
    // Determine final dark mode state based on checks
    // Priority: .dark class > data-theme > custom class > media query
    const isDark = hasDarkClass || hasDataThemeDark || hasCustomDarkClass || prefersDark;
    
    // Log the determination factors (for debugging)
    console.log("Dark mode detection:", { 
        hasDarkClass, 
        prefersDark, 
        hasDataThemeDark, 
        hasCustomDarkClass,
        result: isDark
    });
    
    isDarkMode.value = isDark;
};

onMounted(() => {
    // Initial check
    updateDarkMode();
    
    // Watch for changes in dark mode
    darkModeObserver = new MutationObserver(updateDarkMode);
    
    // Observe changes to the class attribute of both html and body elements
    darkModeObserver.observe(document.documentElement, { 
        attributes: true, 
        attributeFilter: ['class'] 
    });
    
    darkModeObserver.observe(document.body, { 
        attributes: true, 
        attributeFilter: ['class'] 
    });
    
    // Listen for system preference changes
    window.matchMedia('(prefers-color-scheme: dark)').addEventListener('change', updateDarkMode);
});

onUnmounted(() => {
    // Clean up observers and event listeners
    if (darkModeObserver) {
        darkModeObserver.disconnect();
    }
    
    window.matchMedia('(prefers-color-scheme: dark)').removeEventListener('change', updateDarkMode);
});

const props = defineProps<{
    variables?: Array<main.Variable>
    data?: { [key: string]: any }[]
    modelValue: string,
    showBindedSql?: boolean
    minify?: boolean
}>()

const emit = defineEmits(['update:modelValue']);

const value = ref(props.modelValue);

const lang = sql({
    dialect: MySQL
});

// SQL keywords for autocomplete
const sqlKeywords = [
  // Common SQL Commands
  { label: "SELECT", type: "keyword", info: "Retrieves data from a database" },
  { label: "FROM", type: "keyword", info: "Specifies which table to select or delete data from" },
  { label: "WHERE", type: "keyword", info: "Filters records based on a condition" },
  { label: "INSERT INTO", type: "keyword", info: "Inserts new data into a database" },
  { label: "UPDATE", type: "keyword", info: "Modifies existing database data" },
  { label: "DELETE", type: "keyword", info: "Deletes data from a database" },
  { label: "CREATE", type: "keyword", info: "Creates a new database object (table, view, etc.)" },
  { label: "ALTER", type: "keyword", info: "Modifies an existing database object" },
  { label: "DROP", type: "keyword", info: "Deletes an existing database object" },
  { label: "TRUNCATE", type: "keyword", info: "Removes all records from a table, but not the table itself" },
  { label: "RENAME", type: "keyword", info: "Renames a database object" },
  { label: "SHOW", type: "keyword", info: "Shows information about databases, tables, columns, or status" },
  { label: "DESCRIBE", type: "keyword", info: "Shows the structure of a table" },
  { label: "EXPLAIN", type: "keyword", info: "Shows the execution plan for a query" },
  { label: "USE", type: "keyword", info: "Selects a database" },
  
  // Join Types
  { label: "JOIN", type: "keyword", info: "Combines rows from two or more tables" },
  { label: "LEFT JOIN", type: "keyword", info: "Returns all records from the left table, and matched records from the right table" },
  { label: "RIGHT JOIN", type: "keyword", info: "Returns all records from the right table, and matched records from the left table" },
  { label: "INNER JOIN", type: "keyword", info: "Returns records that have matching values in both tables" },
  { label: "CROSS JOIN", type: "keyword", info: "Returns the Cartesian product of two tables" },
  { label: "NATURAL JOIN", type: "keyword", info: "Joins tables by automatically finding matching column names" },
  { label: "FULL JOIN", type: "keyword", info: "Returns all records when there is a match in either left or right table" },
  { label: "SELF JOIN", type: "keyword", info: "Joins a table to itself" },
  
  // Clauses
  { label: "GROUP BY", type: "keyword", info: "Groups rows that have the same values into summary rows" },
  { label: "HAVING", type: "keyword", info: "Filters records after GROUP BY is applied" },
  { label: "ORDER BY", type: "keyword", info: "Sorts the result set in ascending or descending order" },
  { label: "LIMIT", type: "keyword", info: "Limits the number of records returned" },
  { label: "OFFSET", type: "keyword", info: "Specifies where to start selecting records" },
  { label: "UNION", type: "keyword", info: "Combines the result sets of two or more SELECT statements" },
  { label: "UNION ALL", type: "keyword", info: "Combines the result sets of two or more SELECT statements (allows duplicates)" },
  { label: "INTERSECT", type: "keyword", info: "Returns the records that both queries have" },
  { label: "EXCEPT", type: "keyword", info: "Returns the records from the first query that are not in the second query" },
  { label: "WITH", type: "keyword", info: "Specifies temporary named result sets (Common Table Expressions)" },
  
  // Operators
  { label: "AND", type: "operator", info: "Logical AND operator" },
  { label: "OR", type: "operator", info: "Logical OR operator" },
  { label: "NOT", type: "operator", info: "Logical NOT operator" },
  { label: "IN", type: "operator", info: "Specifies multiple possible values for a column" },
  { label: "NOT IN", type: "operator", info: "Specifies values that are not in a list of values" },
  { label: "BETWEEN", type: "operator", info: "Selects values within a given range" },
  { label: "NOT BETWEEN", type: "operator", info: "Selects values outside a given range" },
  { label: "LIKE", type: "operator", info: "Searches for a specified pattern in a column" },
  { label: "NOT LIKE", type: "operator", info: "Searches for values that don't match a pattern" },
  { label: "REGEXP", type: "operator", info: "Matches a string against a regular expression pattern" },
  { label: "IS NULL", type: "operator", info: "Tests for NULL values" },
  { label: "IS NOT NULL", type: "operator", info: "Tests for non-NULL values" },
  { label: "EXISTS", type: "operator", info: "Tests for the existence of records in a subquery" },
  { label: "ANY", type: "operator", info: "Compares a value to any value in a list" },
  { label: "ALL", type: "operator", info: "Compares a value to all values in a list" },
  { label: "SOME", type: "operator", info: "Same as ANY" },
  
  // Aggregate Functions
  { label: "COUNT()", type: "function", info: "Returns the number of rows" },
  { label: "COUNT(DISTINCT)", type: "function", info: "Returns the number of distinct values" },
  { label: "SUM()", type: "function", info: "Returns the sum of values" },
  { label: "AVG()", type: "function", info: "Returns the average value" },
  { label: "MIN()", type: "function", info: "Returns the minimum value" },
  { label: "MAX()", type: "function", info: "Returns the maximum value" },
  { label: "GROUP_CONCAT()", type: "function", info: "Returns a concatenated string of values" },
  { label: "STD()", type: "function", info: "Returns the standard deviation" },
  { label: "STDDEV()", type: "function", info: "Returns the standard deviation" },
  { label: "STDDEV_POP()", type: "function", info: "Returns the population standard deviation" },
  { label: "STDDEV_SAMP()", type: "function", info: "Returns the sample standard deviation" },
  { label: "VAR_POP()", type: "function", info: "Returns the population variance" },
  { label: "VAR_SAMP()", type: "function", info: "Returns the sample variance" },
  { label: "VARIANCE()", type: "function", info: "Returns the variance" },
  { label: "BIT_AND()", type: "function", info: "Returns the bitwise AND of all bits in a group" },
  { label: "BIT_OR()", type: "function", info: "Returns the bitwise OR of all bits in a group" },
  { label: "BIT_XOR()", type: "function", info: "Returns the bitwise XOR of all bits in a group" },
  { label: "JSON_ARRAYAGG()", type: "function", info: "Returns a JSON array containing values from a group" },
  { label: "JSON_OBJECTAGG()", type: "function", info: "Returns a JSON object containing key-value pairs from a group" },
  
  // String Functions
  { label: "CONCAT()", type: "function", info: "Adds two or more strings together" },
  { label: "CONCAT_WS()", type: "function", info: "Adds two or more strings together with a separator" },
  { label: "SUBSTRING()", type: "function", info: "Extracts a string of characters from a string" },
  { label: "SUBSTR()", type: "function", info: "Extracts a substring from a string (alias for SUBSTRING)" },
  { label: "TRIM()", type: "function", info: "Removes leading and trailing spaces from a string" },
  { label: "LTRIM()", type: "function", info: "Removes leading spaces from a string" },
  { label: "RTRIM()", type: "function", info: "Removes trailing spaces from a string" },
  { label: "UPPER()", type: "function", info: "Converts a string to upper case" },
  { label: "LOWER()", type: "function", info: "Converts a string to lower case" },
  { label: "LCASE()", type: "function", info: "Converts a string to lower case (alias for LOWER)" },
  { label: "UCASE()", type: "function", info: "Converts a string to upper case (alias for UPPER)" },
  { label: "LENGTH()", type: "function", info: "Returns the length of a string" },
  { label: "CHAR_LENGTH()", type: "function", info: "Returns the number of characters in a string" },
  { label: "CHARACTER_LENGTH()", type: "function", info: "Returns the number of characters in a string" },
  { label: "REPLACE()", type: "function", info: "Replaces all occurrences of a substring within a string" },
  { label: "REVERSE()", type: "function", info: "Reverses a string" },
  { label: "REPEAT()", type: "function", info: "Repeats a string a specified number of times" },
  { label: "INSERT()", type: "function", info: "Inserts a substring at the specified position" },
  { label: "LEFT()", type: "function", info: "Returns the leftmost characters from a string" },
  { label: "RIGHT()", type: "function", info: "Returns the rightmost characters from a string" },
  { label: "LOCATE()", type: "function", info: "Returns the position of a substring in a string" },
  { label: "POSITION()", type: "function", info: "Returns the position of a substring in a string" },
  { label: "INSTR()", type: "function", info: "Returns the position of a substring in a string" },
  { label: "LPAD()", type: "function", info: "Left-pads a string with another string" },
  { label: "RPAD()", type: "function", info: "Right-pads a string with another string" },
  { label: "SPACE()", type: "function", info: "Returns a string of the specified number of spaces" },
  { label: "ELT()", type: "function", info: "Returns the string at the specified position" },
  { label: "FIELD()", type: "function", info: "Returns the index position of a value in a list" },
  { label: "FIND_IN_SET()", type: "function", info: "Returns the position of a string in a comma-separated list" },
  { label: "FORMAT()", type: "function", info: "Formats a number to a format like '#,###,###.##'" },
  { label: "HEX()", type: "function", info: "Converts a value to hexadecimal" },
  { label: "UNHEX()", type: "function", info: "Converts a hexadecimal value to a string" },
  { label: "BIN()", type: "function", info: "Returns a binary representation of a number" },
  { label: "OCT()", type: "function", info: "Returns an octal representation of a number" },
  { label: "ASCII()", type: "function", info: "Returns the ASCII value of the leftmost character" },
  { label: "ORD()", type: "function", info: "Returns the character code for the leftmost character" },
  { label: "SOUNDEX()", type: "function", info: "Returns a soundex string for a given string" },
  
  // Date Functions
  { label: "NOW()", type: "function", info: "Returns the current date and time" },
  { label: "SYSDATE()", type: "function", info: "Returns the current date and time" },
  { label: "CURDATE()", type: "function", info: "Returns the current date" },
  { label: "CURRENT_DATE()", type: "function", info: "Returns the current date" },
  { label: "CURTIME()", type: "function", info: "Returns the current time" },
  { label: "CURRENT_TIME()", type: "function", info: "Returns the current time" },
  { label: "CURRENT_TIMESTAMP()", type: "function", info: "Returns the current date and time" },
  { label: "LOCALTIME()", type: "function", info: "Returns the current date and time" },
  { label: "LOCALTIMESTAMP()", type: "function", info: "Returns the current date and time" },
  { label: "UTC_DATE()", type: "function", info: "Returns the current UTC date" },
  { label: "UTC_TIME()", type: "function", info: "Returns the current UTC time" },
  { label: "UTC_TIMESTAMP()", type: "function", info: "Returns the current UTC date and time" },
  { label: "YEAR()", type: "function", info: "Returns the year part of a date" },
  { label: "MONTH()", type: "function", info: "Returns the month part of a date" },
  { label: "DAY()", type: "function", info: "Returns the day part of a date" },
  { label: "HOUR()", type: "function", info: "Returns the hour part of a time" },
  { label: "MINUTE()", type: "function", info: "Returns the minute part of a time" },
  { label: "SECOND()", type: "function", info: "Returns the second part of a time" },
  { label: "MICROSECOND()", type: "function", info: "Returns the microsecond part of a time" },
  { label: "DATE_FORMAT()", type: "function", info: "Formats a date as specified" },
  { label: "TIME_FORMAT()", type: "function", info: "Formats a time as specified" },
  { label: "DATE_ADD()", type: "function", info: "Adds a time/date interval to a date" },
  { label: "ADDDATE()", type: "function", info: "Adds a time/date interval to a date" },
  { label: "DATE_SUB()", type: "function", info: "Subtracts a time/date interval from a date" },
  { label: "SUBDATE()", type: "function", info: "Subtracts a time/date interval from a date" },
  { label: "ADDTIME()", type: "function", info: "Adds a time interval to a time" },
  { label: "SUBTIME()", type: "function", info: "Subtracts a time interval from a time" },
  { label: "DATEDIFF()", type: "function", info: "Returns the difference in days between two dates" },
  { label: "TIMEDIFF()", type: "function", info: "Returns the difference between two times" },
  { label: "DATE()", type: "function", info: "Extracts the date part of a date or datetime" },
  { label: "TIME()", type: "function", info: "Extracts the time part of a time or datetime" },
  { label: "TIMESTAMP()", type: "function", info: "Returns a datetime value" },
  { label: "CONVERT_TZ()", type: "function", info: "Converts a datetime from one timezone to another" },
  { label: "EXTRACT()", type: "function", info: "Extracts a part of a date" },
  { label: "MAKEDATE()", type: "function", info: "Creates a date from a year and day value" },
  { label: "MAKETIME()", type: "function", info: "Creates a time from hour, minute, and second values" },
  { label: "TO_DAYS()", type: "function", info: "Converts a date to the number of days since year 0" },
  { label: "FROM_DAYS()", type: "function", info: "Converts a day number to a date" },
  { label: "TO_SECONDS()", type: "function", info: "Converts a date to the number of seconds since year 0" },
  { label: "FROM_UNIXTIME()", type: "function", info: "Converts a Unix timestamp to a date" },
  { label: "UNIX_TIMESTAMP()", type: "function", info: "Returns the Unix timestamp for a date" },
  { label: "SEC_TO_TIME()", type: "function", info: "Converts seconds to a time" },
  { label: "TIME_TO_SEC()", type: "function", info: "Converts a time to seconds" },
  { label: "DAYOFWEEK()", type: "function", info: "Returns the day of the week for a date" },
  { label: "WEEKDAY()", type: "function", info: "Returns the weekday index for a date" },
  { label: "DAYOFMONTH()", type: "function", info: "Returns the day of the month for a date" },
  { label: "DAYOFYEAR()", type: "function", info: "Returns the day of the year for a date" },
  { label: "DAYNAME()", type: "function", info: "Returns the name of the day for a date" },
  { label: "MONTHNAME()", type: "function", info: "Returns the name of the month for a date" },
  { label: "QUARTER()", type: "function", info: "Returns the quarter for a date" },
  { label: "WEEK()", type: "function", info: "Returns the week number for a date" },
  { label: "WEEKOFYEAR()", type: "function", info: "Returns the week number for a date" },
  { label: "YEARWEEK()", type: "function", info: "Returns year and week for a date" },
  { label: "LAST_DAY()", type: "function", info: "Returns the last day of the month for a date" },
  
  // Control Flow Functions
  { label: "CASE", type: "keyword", info: "Evaluates a list of conditions and returns a value when the first condition is met" },
  { label: "WHEN", type: "keyword", info: "Used in a CASE statement to specify a condition" },
  { label: "THEN", type: "keyword", info: "Used in a CASE statement to specify a result" },
  { label: "ELSE", type: "keyword", info: "Used in a CASE statement to specify a value to return if all conditions are false" },
  { label: "END", type: "keyword", info: "Used in a CASE statement to end the list of conditions" },
  { label: "IF()", type: "function", info: "Returns one value if a condition is TRUE, or another value if a condition is FALSE" },
  { label: "IFNULL()", type: "function", info: "Returns a specified value if the expression is NULL" },
  { label: "NULLIF()", type: "function", info: "Returns NULL if two expressions are equal" },
  { label: "COALESCE()", type: "function", info: "Returns the first non-NULL value in a list" },
  
  // Mathematical Functions
  { label: "ABS()", type: "function", info: "Returns the absolute value of a number" },
  { label: "ACOS()", type: "function", info: "Returns the arc cosine of a number" },
  { label: "ASIN()", type: "function", info: "Returns the arc sine of a number" },
  { label: "ATAN()", type: "function", info: "Returns the arc tangent of a number" },
  { label: "ATAN2()", type: "function", info: "Returns the arc tangent of two numbers" },
  { label: "CEIL()", type: "function", info: "Returns the smallest integer value greater than or equal to a number" },
  { label: "CEILING()", type: "function", info: "Returns the smallest integer value greater than or equal to a number" },
  { label: "COS()", type: "function", info: "Returns the cosine of a number" },
  { label: "COT()", type: "function", info: "Returns the cotangent of a number" },
  { label: "DEGREES()", type: "function", info: "Converts a value in radians to degrees" },
  { label: "EXP()", type: "function", info: "Returns e raised to the power of a number" },
  { label: "FLOOR()", type: "function", info: "Returns the largest integer value less than or equal to a number" },
  { label: "LN()", type: "function", info: "Returns the natural logarithm of a number" },
  { label: "LOG()", type: "function", info: "Returns the natural logarithm of a number, or the logarithm of a number to a specified base" },
  { label: "LOG10()", type: "function", info: "Returns the base-10 logarithm of a number" },
  { label: "LOG2()", type: "function", info: "Returns the base-2 logarithm of a number" },
  { label: "MOD()", type: "function", info: "Returns the remainder of a number divided by another number" },
  { label: "PI()", type: "function", info: "Returns the value of PI" },
  { label: "POW()", type: "function", info: "Returns the value of a number raised to the power of another number" },
  { label: "POWER()", type: "function", info: "Returns the value of a number raised to the power of another number" },
  { label: "RADIANS()", type: "function", info: "Converts a value in degrees to radians" },
  { label: "RAND()", type: "function", info: "Returns a random floating-point value between 0 and 1" },
  { label: "ROUND()", type: "function", info: "Rounds a number to a specified number of decimal places" },
  { label: "SIGN()", type: "function", info: "Returns the sign of a number" },
  { label: "SIN()", type: "function", info: "Returns the sine of a number" },
  { label: "SQRT()", type: "function", info: "Returns the square root of a number" },
  { label: "TAN()", type: "function", info: "Returns the tangent of a number" },
  { label: "TRUNCATE()", type: "function", info: "Truncates a number to a specified number of decimal places" },
  
  // Window Functions
  { label: "ROW_NUMBER()", type: "function", info: "Returns the row number of the current row" },
  { label: "RANK()", type: "function", info: "Returns the rank of the current row" },
  { label: "DENSE_RANK()", type: "function", info: "Returns the dense rank of the current row" },
  { label: "NTILE()", type: "function", info: "Returns the ntile group number of the current row" },
  { label: "LAG()", type: "function", info: "Returns the value of the expression evaluated at the row previous to the current row" },
  { label: "LEAD()", type: "function", info: "Returns the value of the expression evaluated at the row following the current row" },
  { label: "FIRST_VALUE()", type: "function", info: "Returns the value of the expression evaluated at the first row" },
  { label: "LAST_VALUE()", type: "function", info: "Returns the value of the expression evaluated at the last row" },
  { label: "NTH_VALUE()", type: "function", info: "Returns the value of the expression evaluated at the nth row" },
  { label: "PERCENT_RANK()", type: "function", info: "Returns the percent rank of the current row" },
  { label: "CUME_DIST()", type: "function", info: "Returns the cumulative distribution of the current row" },
  
  // JSON Functions (MySQL 5.7+)
  { label: "JSON_ARRAY()", type: "function", info: "Creates a JSON array" },
  { label: "JSON_OBJECT()", type: "function", info: "Creates a JSON object" },
  { label: "JSON_QUOTE()", type: "function", info: "Quotes a string as a JSON value" },
  { label: "JSON_CONTAINS()", type: "function", info: "Returns whether a JSON document contains a specific value" },
  { label: "JSON_CONTAINS_PATH()", type: "function", info: "Returns whether a JSON document contains a specific path" },
  { label: "JSON_EXTRACT()", type: "function", info: "Extracts a value from a JSON document" },
  { label: "JSON_KEYS()", type: "function", info: "Returns the keys from a JSON object" },
  { label: "JSON_SEARCH()", type: "function", info: "Searches a JSON document for a value" },
  { label: "JSON_ARRAY_APPEND()", type: "function", info: "Appends a value to a JSON array" },
  { label: "JSON_ARRAY_INSERT()", type: "function", info: "Inserts a value into a JSON array" },
  { label: "JSON_INSERT()", type: "function", info: "Inserts values into a JSON document" },
  { label: "JSON_REPLACE()", type: "function", info: "Replaces values in a JSON document" },
  { label: "JSON_REMOVE()", type: "function", info: "Removes values from a JSON document" },
  { label: "JSON_SET()", type: "function", info: "Sets values in a JSON document" },
  { label: "JSON_MERGE()", type: "function", info: "Merges JSON documents" },
  { label: "JSON_MERGE_PATCH()", type: "function", info: "Merges JSON documents using JSON Merge Patch" },
  { label: "JSON_MERGE_PRESERVE()", type: "function", info: "Merges JSON documents preserving duplicate keys" },
  { label: "JSON_TYPE()", type: "function", info: "Returns the type of a JSON value" },
  { label: "JSON_VALID()", type: "function", info: "Returns whether a value is valid JSON" },
  { label: "JSON_DEPTH()", type: "function", info: "Returns the maximum depth of a JSON document" },
  { label: "JSON_LENGTH()", type: "function", info: "Returns the length of a JSON document" },
  { label: "JSON_PRETTY()", type: "function", info: "Formats a JSON document for readability" },
  { label: "JSON_STORAGE_SIZE()", type: "function", info: "Returns the storage size of a JSON document" },
  { label: "JSON_TABLE()", type: "function", info: "Returns a relational table from JSON data" },
  { label: "JSON_UNQUOTE()", type: "function", info: "Unquotes a JSON value" },
  
  // MySQL 8.0+ Functions
  { label: "GROUPING()", type: "function", info: "Indicates whether a specified column expression in a GROUP BY clause is aggregated" },
  { label: "LATERAL", type: "keyword", info: "Used with derived tables to refer to preceding tables in the FROM clause" },
  { label: "OVER()", type: "function", info: "Defines a window for a window function" },
  { label: "PARTITION BY", type: "keyword", info: "Divides the result set into partitions" },
  { label: "WITH ROLLUP", type: "keyword", info: "Adds extra rows to the result set of GROUP BY to represent subtotals" },
  { label: "RECURSIVE", type: "keyword", info: "Used in a CTE to define a recursive query" },
  
  // Deprecated functions (included for compatibility)
  { label: "DATABASE()", type: "function", info: "Returns the name of the current database" },
  { label: "SCHEMA()", type: "function", info: "Returns the name of the current database (alias for DATABASE)" },
  { label: "USER()", type: "function", info: "Returns the current MySQL user name and host name" },
  { label: "VERSION()", type: "function", info: "Returns the current version of the MySQL server" },
  { label: "CURRENT_USER()", type: "function", info: "Returns the user name and host name for the MySQL account" },
  { label: "LAST_INSERT_ID()", type: "function", info: "Returns the AUTO_INCREMENT value generated by the last INSERT statement" },
  { label: "PASSWORD()", type: "function", info: "Calculates and returns a password string" },
  { label: "BENCHMARK()", type: "function", info: "Executes an expression repeatedly" },
  { label: "CONVERT()", type: "function", info: "Converts a value to a different data type" },
  { label: "CAST()", type: "function", info: "Converts a value to a different data type" },
];

// Create a completion source for variables
const variableCompletions = computed(() => {
    if (!props.variables) return [];
    
    return props.variables.map(variable => ({
        label: `{{ ${variable.Value} }}`,
        type: "variable",
        detail: `Field from .xlsx`,
        info: `Field from .xlsx: ${variable.Field}`,
        apply: `{{ ${variable.Value} }}`,
        boost: 99, // Give variables higher priority
    }));
});

// Autocomplete function for variables
const completeVariables = (context: CompletionContext): CompletionResult | null => {
    // Get the word at cursor
    const word = context.matchBefore(/\{\{\s*\w*\s*\}?\}?/);
    
    // If we're not typing a variable, don't show completions
    if (!word || (word.from === word.to && !context.explicit && !word.text.startsWith("{{"))) {
        return null;
    }
    
    return {
        from: word.from,
        options: variableCompletions.value,
        validFor: /\{\{\s*\w*\s*\}?\}?/,
    };
};

// Autocomplete function for SQL
const completeSql = (context: CompletionContext): CompletionResult | null => {
    // Skip if we're inside a variable template
    const variablePattern = /\{\{\s*\w*\s*\}?\}?/;
    const beforeCursor = context.state.doc.sliceString(
        Math.max(0, context.pos - 20), 
        context.pos
    );
    
    if (variablePattern.test(beforeCursor)) {
        return null;
    }
    
    // Get the word at cursor
    const word = context.matchBefore(/\w+/);
    
    if (!word) {
        return null;
    }
    
    // Filter SQL keywords based on input
    const filteredKeywords = sqlKeywords.filter(keyword => 
        keyword.label.toLowerCase().startsWith(word.text.toLowerCase())
    );
    
    if (filteredKeywords.length === 0 && !context.explicit) {
        return null;
    }
    
    return {
        from: word.from,
        options: filteredKeywords,
        validFor: /\w+/,
    };
};

// Add custom styling for autocomplete items
const myTheme = EditorView.baseTheme({
    ".cm-tooltip.cm-tooltip-autocomplete": {
        border: "1px solid #ddd",
        background: "white",
        fontSize: "90%",
        borderRadius: "4px"
    },
    ".cm-tooltip.cm-tooltip-autocomplete.cm-tooltip-autocomplete-dark": {
        backgroundColor: "#1f2937", // dark:bg-gray-800
        border: "1px solid #374151", // dark:border-gray-700
        color: "white"
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-variable": {
        "&:after": { content: "'{{}}'" },
        color: "#0284c7", // text-sky-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-variable": {
        color: "#0284c7" // text-sky-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-variable": {
        "&:after": { content: "'{{}}'" },
        color: "#38bdf8", // dark:text-sky-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-variable": {
        color: "#38bdf8" // dark:text-sky-400
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-keyword": {
        "&:after": { content: "'k'" },
        color: "#7c3aed", // text-purple-600
        fontWeight: "bold" 
    },
    ".cm-completionLabel.cm-completionLabel-keyword": {
        color: "#7c3aed", // text-purple-600 
        fontWeight: "bold"
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-keyword": {
        "&:after": { content: "'k'" },
        color: "#a78bfa", // dark:text-purple-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-keyword": {
        color: "#a78bfa" // dark:text-purple-400
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-function": {
        "&:after": { content: "'f'" },
        color: "#0891b2", // text-cyan-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-function": {
        color: "#0891b2" // text-cyan-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-function": {
        "&:after": { content: "'f'" },
        color: "#22d3ee", // dark:text-cyan-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-function": {
        color: "#22d3ee" // dark:text-cyan-400
    },
    ".cm-tooltip-autocomplete .cm-completionIcon.cm-completionIcon-operator": {
        "&:after": { content: "'op'" },
        color: "#ea580c", // text-orange-600
        fontWeight: "bold"
    },
    ".cm-completionLabel.cm-completionLabel-operator": {
        color: "#ea580c" // text-orange-600
    },
    ".cm-tooltip-autocomplete-dark .cm-completionIcon.cm-completionIcon-operator": {
        "&:after": { content: "'op'" },
        color: "#fb923c", // dark:text-orange-400
    },
    ".cm-tooltip-autocomplete-dark .cm-completionLabel.cm-completionLabel-operator": {
        color: "#fb923c" // dark:text-orange-400
    },
    ".cm-variable-template": {
        color: "#0284c7", // text-sky-600
        background: "rgba(2, 132, 199, 0.1)", // bg-sky-50/10
        borderRadius: "3px",
        padding: "1px 0"
    },
    "&dark .cm-variable-template": {
        color: "#38bdf8", // dark:text-sky-400
        background: "rgba(56, 189, 248, 0.1)" // dark:bg-sky-400/10
    },
    
    // Selected item style
    ".cm-tooltip-autocomplete ul li[aria-selected=true]": {
        backgroundColor: "#e0f2fe", // bg-sky-100
        color: "#0c4a6e" // text-sky-900
    },
    ".cm-tooltip-autocomplete-dark ul li[aria-selected=true]": {
        backgroundColor: "#075985", // dark:bg-sky-800
        color: "#e0f7ff" // dark:text-sky-50
    }
});

// Create a decoration for highlighting variables
const createVariableHighlighter = () => {
    const variableRegex = /\{\{\s*[a-zA-Z0-9_]+\s*\}\}/g;
    
    return ViewPlugin.fromClass(class {
        decorations: DecorationSet;
        
        constructor(view: EditorView) {
            this.decorations = this.buildDecorations(view);
        }
        
        update(update: ViewUpdate) {
            if (update.docChanged || update.viewportChanged) {
                this.decorations = this.buildDecorations(update.view);
            }
        }
        
        buildDecorations(view: EditorView) {
            const decorations = [];
            
            for (const { from, to } of view.visibleRanges) {
                const text = view.state.doc.sliceString(from, to);
                const matches = [...text.matchAll(variableRegex)];
                
                for (const match of matches) {
                    const start = from + match.index!;
                    const end = start + match[0].length;
                    
                    decorations.push(Decoration.mark({
                        class: "cm-variable-template"
                    }).range(start, end));
                }
            }
            
            return Decoration.set(decorations);
        }
    }, {
        decorations: v => v.decorations
    });
};

// Combine all extensions
const extensions = computed(() => {
    const exts = [
        myTheme,
        createVariableHighlighter(),
        autocompletion({
            override: [completeVariables, completeSql],
            icons: true,
            optionClass: option => option.type ? `cm-completion-${option.type}` : "",
            addToOptions: [
                {
                    render(completion, state) {
                        const dom = document.createElement("div");
                        dom.style.padding = "2px 0";
                        dom.style.display = "flex";
                        dom.style.alignItems = "center";
                        dom.style.opacity = "0.7";
                        dom.style.fontSize = "0.85em";
                        
                        // Show different info based on completion type
                        if (completion.type === "variable" && completion.detail) {
                            dom.textContent = completion.detail.length > 20 
                                ? completion.detail.substring(0, 20) + "..." 
                                : completion.detail;
                        } else if (completion.info) {
                            dom.textContent = typeof completion.info === 'string' 
                                ? completion.info 
                                : '';
                        }
                        
                        return dom;
                    },
                    position: 80
                }
            ]
        })
    ];
    
    // Add oneDarkTheme only if in dark mode
    if (isDarkMode.value) {
        exts.unshift(oneDarkTheme);
    }
    
    return exts;
});

// Output editor extensions (simpler, no autocomplete)
const outputExtensions = computed(() => {
    const exts = [myTheme];
    
    // Add oneDarkTheme only if in dark mode
    if (isDarkMode.value) {
        exts.unshift(oneDarkTheme);
    }
    
    return exts;
});

const linesBinded = computedAsync(async () => {
    if(!props.showBindedSql) {
        return "";
    }

    return await MakeBindedSQL(value.value, props.data!, props.variables!, props.minify) ?? "";
}, "");

const getBindedSQL = async (): Promise<string> => {
    if(!props.showBindedSql) {
        return "";
    }

    return await MakeBindedSQL(value.value, props.data!, props.variables!, props.minify) ?? "";
}

watch(() => props.modelValue, (newValue) => {
    value.value = newValue;
});

watch(() => value.value, (val) => {
    emit('update:modelValue', val);
});

// Manual toggle for dark mode (for testing and debugging)
const toggleDarkMode = () => {
    isDarkMode.value = !isDarkMode.value;
    console.log("Manual dark mode toggle:", isDarkMode.value);
};

// Watch for dark mode changes
watch(() => isDarkMode.value, (newValue) => {
    // The computed extensions will automatically update
    console.log("Dark mode changed to:", newValue);
});

defineExpose({
    getBindedSQL,
    toggleDarkMode,
    isDarkMode
});
</script>

<style scoped>
.editor-dark-mode :deep(.cm-editor) {
    background-color: #1f2937; /* dark:bg-gray-800 */
    color: #f9fafb; /* dark:text-gray-100 */
}

/* Custom dark styles that can be applied directly */
.editor-dark-mode :deep(.cm-gutters) {
    background-color: #111827; /* dark:bg-gray-900 */
    color: #6b7280; /* dark:text-gray-500 */
    border-right-color: #374151; /* dark:border-gray-700 */
}

.editor-dark-mode :deep(.cm-activeLineGutter) {
    background-color: #374151; /* dark:bg-gray-700 */
    color: #9ca3af; /* dark:text-gray-400 */
}

.editor-dark-mode :deep(.cm-activeLine) {
    background-color: rgba(55, 65, 81, 0.3); /* dark:bg-gray-700/30 */
}

.editor-dark-mode :deep(.cm-selectionMatch) {
    background-color: rgba(37, 99, 235, 0.2); /* dark:bg-blue-600/20 */
}

.editor-dark-mode :deep(.cm-tooltip) {
    background: #1f2937; /* dark:bg-gray-800 */
    border: 1px solid #374151; /* dark:border-gray-700 */
    color: #f9fafb; /* dark:text-gray-100 */
}
</style>