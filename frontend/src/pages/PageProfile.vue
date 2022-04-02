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

    <q-btn unelevated color="primary">Сменить пароль</q-btn>

    <div>
      <h4>Корзина пользователя {{ info.id }}</h4>
    </div>

    <CartTable />
  </q-page>
</template>

<script>
import CartTable from 'components/CartTable'

export default {
  components: {
    CartTable
  },
  // name: 'PageName',
  data: () => {
    return {
      info: {
        email: ''
      }
    }
  },
  methods: {
    async getUserInfo () {
      console.log('users info getting')
      const res = await this.$axios.get('auth/whoami')
      this.info = res.data
      console.log(res)
    }
  },
  created () {
    this.getUserInfo()
  }
}
</script>
