<script>

export default {
    name: 'PageNavigator',
    props: {
        total_pages: {
            type: Number,
            required: true
        },
        page: {
            type: Number,
            required: true
        },
        max_pages: {
            type: Number,
            default: 10
        }
    },
    watch: {
        page() {
            this.recalculatePages();
        },
        total_pages() {
            this.recalculatePages();
        },
    },
    data() {
        const showed_pages = [];
        const half_max_pages = Math.floor(this.max_pages / 2);
        const start = Math.max(1, this.page - half_max_pages);
        const end = Math.min(this.total_pages, this.page + half_max_pages);
        for (let i = start; i <= end; i++) {
            showed_pages.push(i);
        }

        return {
            showed_pages,
            showStart: start > 1,
            startIsFar: start > 2,
            showEnd: end < this.total_pages,
            endIsFar: end < this.total_pages - 1,
            inputPage: this.page
        }
    },
    methods: {
        recalculatePages() {
            const showed_pages = [];
            // Only show 10 pages max (first, last, and 8 pages near current)
            const start = Math.max(1, this.page - 4);
            const end = Math.min(this.total_pages, this.page + 4);
            for (let i = start; i <= end; i++) {
                showed_pages.push(i);
            }

            this.showed_pages = showed_pages;
            this.showStart = start > 1;
            this.startIsFar = start > 2;
            this.showEnd = end < this.total_pages;
            this.endIsFar = end < this.total_pages - 1;
        },
        changePage(new_page) {
            this.$emit('change-page', new_page);
        },
        submitCustomPage() {
            if (this.inputPage < 1) return;
            if (this.inputPage > this.total_pages) return;
            this.changePage(this.inputPage);
        }
    }
}

</script>

<template>
    <nav class="d-flex justify-content-center">
        <ul class="pagination justify-content-center mb-0">
            <li class="page-item" :class="{ disabled: page === 1 }">
                <a class="page-link" href="#" aria-label="Previous" v-on:click.prevent="changePage(page - 1)">
                    <span aria-hidden="true">&laquo;</span>
                </a>
            </li>

            <li class="page-item" :class="{ disabled: page === 1 }" v-if="showStart">
                <a class="page-link" href="#" v-on:click.prevent="changePage(1)">1</a>
            </li>

            <li class="page-item disabled" v-if="startIsFar">
                <a class="page-link" href="#">...</a>
            </li>

            <li class="page-item" :class="{ active: i === page }" v-for="i in showed_pages" :key="i">
                <a class="page-link" href="#" v-on:click.prevent="changePage(i)">{{ i }}</a>
            </li>

            <li class="page-item disabled" v-if="endIsFar">
                <a class="page-link" href="#">...</a>
            </li>

            <li class="page-item" :class="{ disabled: page === total_pages }" v-if="showEnd">
                <a class="page-link" href="#" v-on:click.prevent="changePage(total_pages)">{{ total_pages }}</a>
            </li>

            <li class="page-item" :class="{ disabled: page === total_pages }">
                <a class="page-link" href="#" aria-label="Next" v-on:click.prevent="changePage(page + 1)">
                    <span aria-hidden="true">&raquo;</span>
                </a>
            </li>
        </ul>

        <div class="input-group input-group-sm ms-1" style="width: 6rem;">
            <input type="number" class="form-control form-control-sm"
                v-model.number="inputPage" min="1" :max="total_pages">
            <button class="btn btn-sm btn-success" @click.prevent="submitCustomPage">
                <i class="fa fa-arrow-right"></i>
            </button>
        </div>
    </nav>
</template>

<style scoped>
.page-link {
    color: white;
}

.active>.page-link, .page-link.active {
    background-color: #198754 !important;
    border-color: #198754 !important;
}

.page-link:hover {
    color: white;
    background-color: #157347 !important;
    border-color: #157347 !important;
}
</style>
