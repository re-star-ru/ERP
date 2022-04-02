import { defineStore } from 'pinia'
import { api } from 'boot/axios'

export const useAmountStore = defineStore('counter', {
  state: () => ({
    data: [],
    loading: false
  }),

  getters: {
    data: (state) => state.data
  },

  actions: {
    async searchProducts (text) {
      try {
        this.toggleLoading()
        const res = await api.get('/search/' + text)
        this.setProductsData(res.data)
        this.toggleLoading()
      } catch (e) {
        this.toggleLoading()
        throw new StoreError(e)
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

class StoreError extends Error {
  constructor (e) {
    super('Store error')
    this.error = 'ошибка сервера'
    this.message = 'неизвестная ошибка'

    console.dir(e)

    // const j = e.toJSON()
    // console.dir(j)

    if (e.response || e.request) {
      this.error = e.response.data.error
      this.message = e.response.data.message
    } else {
      this.error = 'ошибка клиента'
      this.message = e.message

      console.log('Error', e.message)
    }
  }
}
