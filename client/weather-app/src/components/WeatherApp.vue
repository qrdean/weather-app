<template>
    <main>
        <label for="longitude">Longitude </label>
        <input v-model="longitude">
        <label for="latitude"> Latitude </label>
        <input v-model="latitude">
        <label for="units"> Unit </label>
        <select v-model="units" @change="changeUnits($event)">
            <option v-for="unit in unitSelections" :value="unit.value" :key="unit.id">{{ unit.name}}</option>
        </select>
        <!-- <input v-model="units"> -->
        <div class="" v-for="day in forecast.days" :key="day.date">
            <div class="date_value">
                {{ day.day }}
                {{ day.date }}
            </div>
            <div class="temperature_value">High: {{ day.temperatureHigh }}</div>
            <div class="temperature_value">Low: {{ day.temperatureLow }}</div>
        </div>
    </main>
</template>

<script>
import WeatherForecast from '../services/Forecast'
import _ from 'lodash'

export default {
    name: 'WeatherApp',

    watch: {
        longitude: function() {
            this.debounceForecast()
        },
        latitude: function() {
            this.debounceForecast()
        },
    },
    created: function() {
        this.debounceForecast = _.debounce(this.setForecast, 1000)
        this.forecast.longitude = this.longitude
        this.forecast.latitude = this.latitude
        this.forecast.units = this.units
        this.forecast.update()
    },

    methods: {
        setForecast: function() {
            console.log(this.forecast)
            if (this.forecast) {
                this.forecast.longitude = this.longitude
                this.forecast.latitude = this.latitude
                this.forecast.units = this.units
                this.forecast.update()
            }
        },
        changeUnits (event) {
            this.units = event.target.options[event.target.options.selectedIndex].value
            this.setForecast()
        }
    },

    data() {
        return {
            forecast: new WeatherForecast(),
            longitude: -99.771335,
            latitude: 30.489772,
            units: 'metric',
            unitSelections: [
                { name: 'Metric', id: 1 , value: 'metric'},
                { name: 'Imperial', id: 2, value: 'imperial'},
                { name: 'Standard', id: 3, value: 'standard'}
            ]
        }
    }
}
</script>

<style scoped>
.date_value {
    font-size: 1.5em;
    color: rgba(3, 65, 24, 0.75)
}
.temperature_value {
    font-size: 1.5em;
    color: rgba(8, 128, 240, 0.75)
}

</style>