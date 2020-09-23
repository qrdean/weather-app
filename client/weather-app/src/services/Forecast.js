'use strict'

const ENDPOINT_DEV = 'http://localhost:8080'

class WeatherForecast {
    constructor() {
        this.days = []

        this.latitude = null,
        this.longitude = null,
        this.units = ''
    }

    update() {
        this.updateForecast(this.longitude, this.latitude, this.units)
    }

    populate(data) {
        this.days = []
        data.forecast.forEach(forecast_ => {
            let daily = new DailyForecast(forecast_.day, forecast_.max, forecast_.min, forecast_.dt)
            this.days.push(daily)
        })
    }

    async updateForecast(longitude, latitude, units) {
        let data = null
        
        if (longitude === null || latitude === null) {
            data = this.getError()
            return
        }
        try {
            data = await this.getForecast({longitude, latitude}, units)
        } catch (e) {
            data = this.getError()
        }

        this.populate(data)
    }

    async getForecast(coordinates, units) {
        console.log(coordinates)
        let endpoint = `${ENDPOINT_DEV}/weatherapi/v1/latitude/${coordinates.latitude}/longitude/${coordinates.longitude}/units/${units}`
        console.log(endpoint)
        let response = await fetch(endpoint)

        return await response.json()
    }

    getError() {
        return {
            day: null,
            max: null,
            min: null,
            dt: null,
        }
    }
}

class DailyForecast {
    constructor(day, min, max, date) {
        this.day = day,
        this.temperatureHigh = min,
        this.temperatureLow = max,
        this.date = date
    }

}

export default WeatherForecast