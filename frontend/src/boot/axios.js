// import Vue from 'vue'
import { LocalStorage } from 'quasar'
import { boot } from 'quasar/wrappers'
import axios from 'axios'
// import Axios from 'axios'

const api = axios.create({
  baseURL: 'https://api.re-star.ru/v1',
  headers: { Authorization: `Bearer ${LocalStorage.getItem('accessToken')}` },
  timeout: 3000
})

if (process.env.DEV) {
  api.baseURL = 'http://localhost:3000'
}

export default boot(({ app }) => {
  app.config.globalProperties.$axios = axios
  app.config.globalProperties.$api = api
})
