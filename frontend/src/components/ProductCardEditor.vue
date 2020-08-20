<template>
  <div>
    <q-btn>Добавить новый товар</q-btn>
    <q-form>
      <q-input label="Наименование" v-model="newProduct.name"></q-input>
      <q-input label="Артикул" v-model="newProduct.sku"></q-input>
      <q-btn outline>Добавить изображения</q-btn>
    </q-form>
  </div>
</template>

<script>
export default {
  // name: 'ComponentName',
  data() {
    return {
      newProduct: {
        name: '',
        sku: '',
      },
    }
  },
}
</script>
