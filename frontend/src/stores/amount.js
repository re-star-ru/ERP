import { defineStore } from 'pinia'
import { api } from 'boot/axios'
import { Notify } from 'quasar'

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
        Notify(e)
        this.toggleLoading()
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
