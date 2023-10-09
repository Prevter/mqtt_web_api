<script setup>
import { ref } from 'vue';

const cookieModal = ref(null);

function acceptCookies() {
    cookieModal.value.classList.add('opacity-0', 'translate-y-4');
    cookieModal.value.classList.remove('opacity-100', 'translate-y-0');
    setTimeout(() => {
        cookieModal.value.classList.add('hidden');
    }, 300);
    localStorage.setItem('cookiesAccepted', true);
}

if (!localStorage.getItem('cookiesAccepted')) {
    setTimeout(() => {
        cookieModal.value.classList.remove('hidden');
        cookieModal.value.classList.add('opacity-100', 'translate-y-0');
    }, 1000);
}
</script>

<template>
    <div ref="cookieModal" class="cookie-modal fixed bottom-4 left-2 mb-2 max-w-sm p-6 dark:bg-slate-900 bg-white border dark:border-gray-950 border-gray-300 rounded shadow-md transform translate-y-4 opacity-0">
        <p class="mb-4">
            🍪 Цей веб-сайт використовує файли cookie для зберігання даних на вашому комп'ютері.
        </p>
        <button @click="acceptCookies"
            class="flex justify-end bg-green-600 text-white px-4 py-2 rounded-md hover:bg-green-700 transition-colors duration-300 ml-auto focus:outline-none">
            Добре
        </button>
    </div>
</template>

<style scoped>
.cookie-modal {
    z-index: 100;
    @apply transition-all duration-300;
}
</style>

