import Vue from 'vue'
import Axios from 'axios'

let baseURL = 'https://api.restar26.site/v1'
if (process.env.DEV) {
  // baseURL = 'http://192.168.0.64:3000/v1'
  // baseURL = 'http://10.51.0.128:3000/v1'
  baseURL = 'http://10.51.1.8:3000/v1'
  // baseURL = 'http://192.168.1.128:3000/v1'
}

Vue.prototype.$axios = Axios
Vue.prototype.$axios.defaults.baseURL = baseURL

let email = localStorage.getItem('email')
let token = localStorage.getItem('token')

if (email && token) {
  Vue.prototype.$axios.defaults.auth = {
    username: email,
    password: token
  }
}
