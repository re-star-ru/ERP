<template>
  <q-page padding class="row">
    <q-card
      flat
      bordered
      square
      class="q-pa-sm col-12 col-md-6 col-lg-4"
      v-for="(v, i) in products"
      :key="i"
    >
      <q-card-section>
        <q-img src="https://s3.restar26.site/public/485f9ba5-f205-4eea-80e9-353a9760d4e1.jpg" />

        <q-card-section></q-card-section>
        <div>id: {{v.id}}</div>
        <div>Наименование {{v.name}}</div>
        <div>Создатель: {{v.creator}}</div>
      </q-card-section>

      <q-card-actions>
        <q-btn outline icon-right="ion-cart" @click="addToCart(v.id)">в корзину</q-btn>
      </q-card-actions>
    </q-card>
  </q-page>
</template>

<script>
export default {
  data: () => {
    return {
      products: [],
    }
  },
  methods: {
    async getAllProducts() {
      console.log('get all prducts')
      const res = await this.$axios.get('products')
      this.products = res.data.products
    },

    async addToCart(id) {
      console.log('add to cart')
      const res = await this.$axios.post('/cart', {
        id: id,
        count: 1,
      })
      console.log(res.status)
    },
  },

  created() {
    this.getAllProducts()
  },
}
</script>

