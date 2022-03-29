<template>
  <q-page class="flex justify-center">
    <div class="q-pa-md">
      <q-table
        :grid="$q.screen.lt.sm"
        title="Остатки"
        no-data-label="Остатков не найдено"
        :data="products"
        :columns="columns"
        :pagination.sync="pagination"
        row-key="index"
        :loading="loading"
        :visibleColumns="visibleColumns"
      >
        <template v-slot:loading>
          <q-inner-loading showing color="primary" />
        </template>
      </q-table>
    </div>
  </q-page>
</template>

<script>
export default {
  data() {
    return {
      text: '',
      pagination: {
        page: 1,
        rowsPerPage: 30
      },
      visibleColumns: ['sku', 'name', 'spec', 'description', 'amount'],
      columns: [
        {
          name: 'index',
          label: '№',
          field: 'index'
        },
        {
          name: 'sku',
          label: 'Артикул',
          field: 'sku'
        },
        {
          name: 'name',
          label: 'Наименование',
          field: 'name'
        },
        {
          name: 'spec',
          label: 'Характеристика',
          field: 'spec'
        },
        {
          name: 'description',
          label: 'Описание',
          field: 'description'
        },
        {
          name: 'amount',
          label: 'Остатки',
          field: 'amount',
          format: val => `${val} шт`
        }
      ]
    }
  },
  methods: {
    searchProducts() {}
  },
  computed: {
    products() {
      return this.$store.state.amount.data
    },
    loading() {
      return this.$store.state.amount.loading
    }
  }
}
</script>
