<template>
    <Bar :data="data" :options="options" />
</template>

<script>
import { Bar } from 'vue-chartjs'
import { Chart as ChartJS, Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale } from 'chart.js'

ChartJS.register(Title, Tooltip, Legend, BarElement, CategoryScale, LinearScale)

export default {
    name: 'BarChart',
    props: {
        data: {
            type: Array,
            required: true
        },
        labels: {
            type: Array,
            required: true
        }
    },
    components: { Bar },
    data() {
        return {
            options: {
                responsive: true,
                // coloring:
                color: 'white',
                // legend:
                legend: {
                    display: true,
                    position: 'bottom',
                    labels: {
                        fontColor: 'white'
                    }
                },
                // scales:
                scales: {
                    xAxes: [{
                        color: 'white',
                        tickColor: 'white',
                    }],
                    yAxes: [{
                        color: 'white',
                        tickColor: 'white',
                    }]
                }
            }
        }
    },
}
</script>