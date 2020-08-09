import axios from 'axios'

let baseURL = 'https://api.restar26.site/'
if (process.env.DEV) {
  baseURL = 'http://10.51.0.128:3000/'
  // baseURL = 'http://192.168.0.64:3000/'
}

export const state = () => {
  return {
    token: localStorage.getItem('token') || '',
    email: localStorage.getItem('email') || '',
    aclGroup: localStorage.getItem('aclGroup') || '',
    user: {}
  }
}

export const getters = {
  isLogged: state => !!state.token,
  aclGroup: state => state.aclGroup
}

export const actions = {
  async login({ commit }, { credentials }) {
    // commit('auth request')
    try {
      console.log(credentials)
      const res = await axios.get(`${baseURL}auth/login`, {
        auth: {
          username: credentials.email,
          password: credentials.password
        }
      })
      console.log(res.data)
      commit('saveLocalAuth', {
        email: credentials.email,
        token: res.data.token,
        aclGroup: res.data.aclGroup
      })
    } catch (e) {
      if (e.code === 5) {
        let err = new Error('Неправильные символы')
        console.log(err)
        throw err
      }
      console.dir(e)
      // commit('auth_error')
      // localStorage.removeItem('token')
      commit('clearLocalAuth')
      e.message = 'Неправильный логин или пароль'
      throw e
    }
  },

  async registration({ commit }, { credentials }) {
    console.log('registrations')
    console.log(credentials)

    try {
      const res = await axios.post(`${baseURL}auth/registration`, {
        email: credentials.email,
        password: credentials.password
      })
      commit('saveLocalAuth', {
        email: credentials.email,
        token: res.data.token,
        aclGroup: res.data.aclGroup
      })
      console.log(res)
    } catch (e) {
      throw e
    }
  },
  logout({ commit }) {
    commit('clearLocalAuth')
  }
}

export const mutations = {
  clearLocalAuth(state) {
    localStorage.removeItem('email')
    localStorage.removeItem('token')
    localStorage.removeItem('aclGroup')

    delete axios.defaults.auth
    state.email = ''
    state.token = ''
    state.aclGroup = ''
  },
  saveLocalAuth(state, data) {
    state.email = data.email
    state.token = data.token
    state.aclGroup = data.aclGroup
    localStorage.setItem('email', data.email)
    localStorage.setItem('token', data.token)
    localStorage.setItem('aclGroup', data.aclGroup)
    axios.defaults.headers.common['Authorization'] =
      'Basic ' + btoa(data.email + ':' + data.token)
    console.log('saved')
  }
}
