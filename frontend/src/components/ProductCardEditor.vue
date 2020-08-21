<template>
  <div>
    {{ product }}
    <q-btn outline @click="newProduct = !newProduct"
      >Добавить новый товар</q-btn
    >
    <q-form
      v-if="newProduct || id"
      @submit="onSubmit"
      @reset="onReset"
      class="q-gutter-md"
    >
      <q-input label="Наименование" v-model="product.name"></q-input>
      <q-input label="Артикул" v-model="product.sku"></q-input>
      <q-btn outline>Добавить изображения</q-btn>

      <div>
        <q-btn label="Submit" type="submit" color="primary" />
        <q-btn
          label="Reset"
          type="reset"
          color="primary"
          flat
          class="q-ml-sm"
        />
      </div>
    </q-form>
  </div>
</template>

<script>
export default {
  // name: 'ComponentName',
  props: ['id'],
  data() {
    return {
      newProduct: false,

      product: {
        name: '',
        sku: '',
      },
    }
  },

  methods: {
    onReset() {
      this.newProduct.name = ''
      this.newProduct.sku = ''
    },

    async onSubmit() {
      const res = await this.$axios.post('/products', this.newProduct)
      if (res.status == 201) {
        this.$q.notify({
          color: 'green-4',
          textColor: 'white',
          icon: 'ion-cloud-done',
          message: 'Товар добавлен',
        })
        this.onReset()
        return
      }
      this.$q.notify({
        color: 'red-4',
        textColor: 'white',
        icon: 'ion-error',
        message: 'Товар не добавлен',
      })
    },

    async getProductInfo() {
      console.log('get product with id:', this.id)
      const res = this.$axios.get(`products/${this.id}`)
      this.product = res.data
    },
  },
  created() {
    this.getProductInfo()
  },
}
</script>
