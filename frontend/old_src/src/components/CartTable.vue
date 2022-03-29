<template>
  <div>
    <q-markup-table separator="cell" flat bordered>
      <thead>
        <tr>
          <th class="text-left">Наименование</th>
          <th class="text-right">Количество</th>
          <th class="text-right">Цена</th>
          <th class="text-right">Сумма</th>
        </tr>
      </thead>
      <tbody>
        <tr class="q-ma-md" v-for="(v, i) in cart.addedProducts" :key="i">
          <td>{{v.productID}}</td>
          <td>{{v.count}}</td>
          <td>{{v.price||1}}</td>
          <td>{{v.count||1*v.price||1}}</td>
        </tr>
      </tbody>
    </q-markup-table>
    <h4>Итого к оплате: {{total}} ₽</h4>
    <q-btn outline>Перейти к оплате</q-btn>
  </div>
</template>

<script>
export default {
  data() {
    return {
      cart: '',
      total: 0,
    }
  },
  methods: {
    calculate() {
      if (this.cart.addedProducts != undefined) {
        console.log('us obj')
        console.log(Object.values(this.cart.addedProducts))
        this.total = Object.values(this.cart.addedProducts).reduce(
          (acc, sum) => {
            return acc + (sum.price || 1) * (sum.count || 1)
          },
          0
        )
      }
    },
    async getUserCart() {
      console.log('getting users cart')
      const res = await this.$axios.get('cart')
      this.cart = res.data
      this.calculate()
      console.log(res)
    },
  },
  created() {
    this.getUserCart()
  },
}
</script>
