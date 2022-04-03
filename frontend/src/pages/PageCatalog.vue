<template>
  <q-page class="row content-start q-pa-sm">
    <div class="col-12">
      <q-input
        style="font-size: 1.8em;"
        standout="bg-primary black"
        type="search"
        v-model="searchString"
        debounce="300"
        @input="getCatalog"
        clearable
        clear-icon="ion-close"
      >
        <template v-slot:append>
          <q-icon v-if="searchString === ''" name="ion-search" />
        </template>
        <template v-slot:after>
          <q-btn class="full-height" color="primary" push label="Искать" @click="getCatalog"></q-btn>
        </template>
      </q-input>
    </div>

    <div class="col-12 row">
      <div class="q-pa-xs col-xs-12 col-sm-8 col-md-10 row justify-center">
        <SkuCard v-for="v in finded" :key="v.id" v-bind="v"/>
      </div>
    </div>

  </q-page>
</template>

<script>
import SkuCard from 'components/SkuCard'

export default {
  components: {
    SkuCard
  },
  data () {
    return {
      perPage: 48,
      // maxGroupsCount: 0,
      finded: [],
      // catalogData: [],
      // productTypes: [],
      searchString: ''
    }
  },
  methods: {
    async getCatalog () {
      console.log('get catalog')
      try {
        console.log('req query', this.$route.query)
        const resp = await this.$api.get('/catalog', {
          params: {
            ...this.$route.query,
            q: this.searchString
          }
        })

        this.finded = resp.data
        // this.catalogData = resp.data.groups
        // this.maxGroupsCount = resp.data.count
        // this.perPage = resp.data.perPage

        console.log(resp.data)
      } catch (e) {
        console.dir(e)
        this.$q.notify({
          type: 'warning',
          message: 'Ничего не найдено',
          timeout: 200
        })
        // this.catalogData = []
        // this.maxGroupsCount = 0
      }
    },

    toPage (page) {
      return { query: { ...this.$route.query, p: page } }
    }
    // changeProductType (val) {
    //   console.log('change product type', val)
    //   this.getCatalog({ ...this.query, t: val.guid })
    // }
  },
  computed: {
    // maxPages () {
    //   return Math.ceil(this.maxGroupsCount / this.perPage)
    // },
    currentPage () {
      const p = Number(this.$route.query.p)
      if (Number.isNaN(p)) {
        return 1
      }
      return p
    }
  },
  watch: {
    $route (to) {
      console.log(to)
      this.getCatalog()
    }
  },
  mounted () {
    this.getCatalog()
  }
}
</script>
