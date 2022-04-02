import Vue from 'vue'
import Axios from 'axios'
import { LocalStorage } from 'quasar'

let baseURL = 'https://api.re-star.ru/v1'
if (process.env.DEV) {
  baseURL = 'http://localhost:3000'
}

Axios.defaults.baseURL = baseURL
Axios.defaults.timeout = 3000

let accessToken = LocalStorage.getItem('accessToken')

if (accessToken) {
  Axios.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
}

Vue.prototype.$axios = Axios
