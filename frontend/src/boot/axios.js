import { boot } from 'quasar/wrappers'
import axios from 'axios'

export const api = axios.create({
  baseURL: process.env.API,
  timeout: 10000
})

export default boot(({ app }) => {
  app.config.globalProperties.$axios = axios
  app.config.globalProperties.$api = api
})
