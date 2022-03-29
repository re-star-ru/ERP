<template>
  <q-list>
    <q-chip
      @click="unselect"
      :style="{ visibility: selected ? 'visible' : 'hidden' }"
      clickable
      icon="ion-close"
    >
      <div v-if="selected">{{ selectedNode().name }}</div>
    </q-chip>
    <q-tree
      ref="productMenu"
      accordion
      color="primary"
      no-connectors
      :duration="100"
      :nodes="simple"
      node-key="guid"
      label-key="name"
      :selected.sync="selected"
    />
  </q-list>
</template>

<script>
export default {
  data() {
    return {
      selected: '',
      simple: [],
      productTypes: [],
    }
  },
  methods: {
    async getProductTypes() {
      try {
        const resp = await this.$axios.get('/catalog/product-types')
        this.simple = this.createSimple(resp.data)
        await this.$nextTick()
        if (this.$route.query.t) {
          this.selected = this.$route.query.t
        }
      } catch (e) {
        console.log(e)
      }
    },
    selectedNode() {
      return this.$refs.productMenu.getNodeByKey(this.selected)
    },
    unselect() {
      this.selected = ''
      this.updateRoute()
    },
    updateRoute(node) {
      if (this.selected) {
        this.$router.push({
          name: 'catalog',
          query: { ...this.$route.query, t: node.guid },
        })
        return
      }
      this.selected = ''
      this.$router.push({
        name: 'catalog',
        query: { ...this.$route.query, t: undefined },
      })
    },
    createSimple(tree) {
      const simple = []

      for (let key in tree) {
        if (key !== 'children' && tree[key] === null) {
          continue
        }
        let el = {
          ...tree[key],
          name: tree[key].name,
          handler: this.nodeHandle,
          selectable: !tree[key]['isgroup'],
          children: this.createSimple(tree[key].children),
        }
        simple.push(el)
      }

      return simple
    },
    nodeHandle(node) {
      if (!node.isgroup) {
        this.updateRoute(node)
      }

      if (this.$refs.productMenu.isExpanded(node.guid)) {
        this.$refs.productMenu.setExpanded(node.guid, false)
        return
      }
      this.$refs.productMenu.setExpanded(node.guid, true)
    },
  },
  // watch: {
  //   selected(val) {
  //     this.updateRoute(val)
  //   },
  // },
  mounted() {
    this.getProductTypes()
  },
}

// function createSimple(tree) {
//   const simple = []
//
//   for (let key in tree) {
//     if (key !== 'children' && tree[key] === null) {
//       continue
//     }
//     let el = {
//       ...tree[key],
//       name: tree[key].name,
//       handler: (node) => {
//         console.log('handler', node)
//         this.$refs.productMenu.setExpanded(node.guid, true)
//       },
//       selectable: !tree[key]['isgroup'],
//       children: createSimple(tree[key].children),
//     }
//     simple.push(el)
//   }
//
//   return simple
// }
</script>

<style lang="scss">
.catalog-menu-link {
  color: white;
  background: #f2c037;
}

.q-tree__node--selected {
  background-color: $deep-purple-2;
}

.q-tree__children {
  padding-left: 10px;
}
</style>
