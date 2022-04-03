<template>
  <div class="q-pa-xs col-xs-12 col-sm-6 col-md-4 col-lg-3">
    <q-card class="sku-card" flat bordered>
      <q-img style="max-height: 300px;" src="~assets/logo.png"></q-img>
      <q-list>
        <q-item v-ripple clickable>
          <q-item-section>
            <div class="text-h6">{{ sku }}</div>
            <div class="text-h6">{{ name }}</div>
            <div class="text-h6">{{ char }}</div>

            <p>цена {{ price }}</p>
            <p>остаток {{ amount }}</p>
          </q-item-section>
        </q-item>
      </q-list>
    </q-card>
  </div>
</template>

<script setup>

defineProps({
  id: String,
  name: String,
  char: String,
  amount: Number,
  sku: String,
  price: Number
})

</script>
