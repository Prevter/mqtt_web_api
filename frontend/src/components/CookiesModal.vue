<script setup>
import { ref } from 'vue';

const cookieModal = ref(null);

function acceptCookies() {
    cookieModal.value.classList.add('opacity-0');
    setTimeout(() => {
        cookieModal.value.classList.add('d-none');
    }, 300);
    localStorage.setItem('cookiesAccepted', true);
}

if (!localStorage.getItem('cookiesAccepted')) {
    setTimeout(() => {
        cookieModal.value.classList.remove('opacity-0');
    }, 1000);
}
</script>

<template>
    <div ref="cookieModal" class="cookie-modal opacity-0">
        <span>🍪 Цей веб-сайт використовує файли cookie для зберігання даних на вашому комп'ютері.</span>
        <div class="mt-2 d-flex align-items-center justify-content-center g-2">
            <button @click="acceptCookies" class="btn btn-success mr-1">Добре</button>
        </div>
    </div>
</template>

<style scoped>
.cookie-modal {
    position: fixed;
    bottom: 4rem;
    left: 2rem;
    margin-bottom: 2rem;
    max-width: 20rem;
    padding: 1.5rem;
    background-color: #454545;
    border: 1px solid #707070;
    border-radius: 0.375rem;
    box-shadow: 0 0.25rem 0.5rem rgba(0, 0, 0, 0.05);
    transition: transform 0.3s ease, opacity 0.3s ease;
}
</style>

