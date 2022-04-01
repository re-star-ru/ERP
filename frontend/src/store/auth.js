import axios from 'axios'
import { LocalStorage } from 'quasar'

export const state = () => {
  return {
    accessToken: LocalStorage.getItem('accessToken') || '',
    email: LocalStorage.getItem('email') || '',
    aclGroup: LocalStorage.getItem('aclGroup') || '',
    user: {}
  }
}

export const getters = {
  isLogged: (state) => !!state.accessToken,
  aclGroup: (state) => state.aclGroup
}

export const actions = {
  async login ({ commit }, { credentials }) {
    try {
      const res = await axios.post('auth/sign-in', credentials)
      console.log(res.data)

      commit('saveLocalAuth', {
        email: credentials.email,
        accessToken: res.data
      })
    } catch (e) {
      if (e.code === 5) {
        const err = new Error('Неправильные символы')
        console.log(err)
        throw err
      }
      console.dir(e)
      commit('clearLocalAuth')
      e.message = 'Неправильный логин или пароль'
      throw e
    }
  },

  async registration ({ commit }, { credentials }) {
    const res = await axios.post('auth/sign-up', credentials)
    commit('saveLocalAuth', {
      email: credentials.email,
      accessToken: res.data
    })
  },
  logout ({ commit }) {
    commit('clearLocalAuth')
  }
}

export const mutations = {
  clearLocalAuth (state) {
    LocalStorage.remove('email')
    LocalStorage.remove('accessToken')
    LocalStorage.remove('aclGroup')

    state.email = ''
    state.accessToken = ''
    state.aclGroup = ''
  },

  saveLocalAuth (state, data) {
    state.email = data.email
    state.accessToken = data.accessToken
    state.aclGroup = 'test'
    try {
      LocalStorage.set('email', data.email)
      LocalStorage.set('accessToken', data.accessToken)
      LocalStorage.set('aclGroup', 'test')
    } catch (error) {
      console.dir(error)
    }

    axios.defaults.headers.common.Authorization = `Bearer ${data.accessToken}`
  }
}
