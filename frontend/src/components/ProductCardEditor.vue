<template>
  <div>
    <q-btn>Добавить новый товар</q-btn>

    <q-form @reset="onReset">
      <q-input label="Наименование" v-model="newProduct.name"></q-input>
      <q-input label="Артикул" v-model="newProduct.sku"></q-input>
      <q-btn outline>Добавить изображения</q-btn>

      <div>
        <q-btn label="Submit" type="submit" color="primary" />
        <q-btn label="Reset" type="reset" color="primary" flat class="q-ml-sm" />
      </div>
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
