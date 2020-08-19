<template>
  <q-page padding>
    <div>
      <h4>Профиль</h4>
    </div>

    <q-field label="Почта" stack-label>
      <template v-slot:control>
        <div class="self-center full-width no-outline" tabindex="0">{{ info.email }}</div>
      </template>
    </q-field>

    <div>
      <h4>Корзина пользователя {{ cart.userID }}</h4>
    </div>

    <div>
      <q-card v-for="(v, i) in cart.addedProducts" :key="i">
        <q-card-section>{{ v }}</q-card-section>
      </q-card>
    </div>
  </q-page>
</template>

<script>
export default {
  // name: 'PageName',
  data: () => {
    return {
      info: {
        email: '',
      },
      cart: '',
    }
  },
  methods: {
    async getUserInfo() {
      console.log('users info getting')
      const res = await this.$axios.get('auth/whoami')
      this.info = res.data
      console.log(res)
    },
    async getUserCart() {
      console.log('getting users cart')
      const res = await this.$axios.get('cart')
      this.cart = res.data
      console.log(res)
    },
  },
  created() {
    this.getUserInfo()
    this.getUserCart()
  },
}
</script>
