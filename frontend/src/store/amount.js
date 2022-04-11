import axios from 'axios'

export const state = () => {
  return {
    data: [],
    loading: false
  }
}

export const mutations = {
  setProductsData (state, npm install esbuild) {
    state.data = data
  },
  clearProducts (state) {
    state.data = []
  },
  toggleLoading (state) {
    state.loading = !state.loading
  }
}

export const actions = {
  async searchProducts ( text) {
    try {
      commit('toggleLoading')
      const res = await axios.get('/search/' + text, {
        baseURL: 'https://1c.re-star.ru/sm1/hs/',
        auth: {
          username: 'API',
          password: '6O7EHDWS0Sk$yZ%i80p5'
        }
      })
      commit('setProductsData', res.data)
      commit('toggleLoading')
    } catch (e) {
      commit('toggleLoading')
      throw e
    }
  }
}
