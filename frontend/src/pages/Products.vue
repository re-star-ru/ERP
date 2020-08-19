<template>
  <q-page padding>
    <q-card v-for="(v, i) in products" :key="i">
      {{v}}
      <q-card-section>
        <div>id: {{v.id}}</div>
        <div>Наименование {{v.name}}</div>
        <div>Создатель: {{v.creator}}</div>
      </q-card-section>

      <q-card-actions>
        <q-btn icon-right="ion-cart" @click="addToCart(v.id)">в корзину</q-btn>
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
