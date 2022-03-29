<template>
  <div>
    <q-btn @click="showRaw = !showRaw">show raw data</q-btn>
    <pre v-if="showRaw">{{ product }}</pre>

    <q-form
      v-if="newProduct || id"
      @submit="onSubmit"
      @reset="onReset"
      class="q-gutter-md"
    >
      <q-input label="Наименование" v-model="product.name"></q-input>
      <q-input label="Артикул" v-model="product.sku"></q-input>

      <q-field label="Дата создания" dense stack-label>
        <template v-slot:prepend>
          <q-icon name="ion-calendar" />
        </template>

        <template v-slot:control>
          <div>
            {{ product.createdAt }}
          </div>
        </template>
      </q-field>
      <q-field label="Дата изменения" dense stack-label>
        <template v-slot:prepend>
          <q-icon name="ion-calendar" />
        </template>

        <template v-slot:control>
          <div>
            {{ product.lastModified }}
          </div>
        </template>
      </q-field>
      <q-btn outline>Добавить изображения</q-btn>
      <div>
        <q-btn label="Сохранить" type="submit" color="primary" />
        <q-btn
          label="Сбросить"
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
const productType = {
  id: String,
  name: String,
  guid: String,
  sku: String,
  description: String,
  manufacturer: String,
  typeGUID: String,
  typeName: String,
  characteristics: Array,
  properties: Array,
  creator: {
    id: String,
    email: String,
    name: String
  },
  createdAt: String,
  lastModified: String
}

export default {
  // name: 'ComponentName',
  props: ['id'],
  data () {
    return {
      showRaw: false,
      newProduct: false,

      product: productType
    }
  },

  methods: {
    onReset () {
      this.product.name = ''
      this.product.sku = ''
    },

    async onSubmit () {
      const res = await this.$axios.put(
        `products/${this.product.id}`,
        this.product
      )
      console.log(res)
      if (res.status >= 200 || res.status < 300) {
        this.$q.notify({
          color: 'green-4',
          textColor: 'white',
          icon: 'ion-cloud-done',
          message: 'Товар успешно обновлен'
        })
        return
      }
      this.$q.notify({
        color: 'red-4',
        textColor: 'white',
        icon: 'ion-bug',
        message: 'Товар не добавлен'
      })
    },

    async getProductInfo () {
      console.log('get product with id:', this.id)
      try {
        const res = await this.$axios.get(`products/${this.id}`)
        this.product = res.data
      } catch (e) {
        console.log(e)
      }
    }
  },
  created () {
    this.getProductInfo()
  }
}
</script>
