import { LocalStorage } from 'quasar'
import { boot } from 'quasar/wrappers'
import axios from 'axios'

export const api = axios.create({
  baseURL: process.env.API,
  headers: { Authorization: `Bearer ${LocalStorage.getItem('accessToken')}` },
  timeout: 5000
})

export default boot(({ app }) => {
  app.config.globalProperties.$axios = axios
  app.config.globalProperties.$api = api
})
