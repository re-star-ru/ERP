<template>
  <q-page padding>
    <div class="q-pa-md q-gutter-sm">
      <q-btn outline>Загрузить товары из 1с</q-btn>
    </div>
    <ProductCardEditor />
  </q-page>
</template>

<script>
import ProductCardEditor from 'components/ProductCardEditor'
export default {
  components: {
    ProductCardEditor,
  },
}
</script>
