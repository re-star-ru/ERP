import { defineStore } from 'pinia'
import { api } from 'boot/axios'

export const useAmountStore = defineStore('counter', {
  state: () => ({
    data: [],
    loading: false
  }),

  // getters: {
  //   doubleCount (state) {
  //     return state.counter * 2
  //   }
  // },

  actions: {
    async searchProducts (text) {
      try {
        this.toggleLoading()
        const res = await api.get('/search/' + text)
        this.setProductsData(res.data)
        this.toggleLoading()
      } catch (e) {
        this.toggleLoading()
        throw e
      }
    },

    increment () {
      this.counter++
    },

    clearProducts () {
      this.data = []
    },

    toggleLoading () {
      this.loading = !this.loading
    },

    setProductsData (data) {
      this.data = data
    }
  }
})
