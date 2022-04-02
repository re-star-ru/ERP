<template>
  <div class="q-pa-xs col-xs-12 col-sm-6 col-md-4 col-lg-3">
    <q-card class="sku-card" flat bordered>
      <q-img style="max-height: 300px;" src="~assets/logo.png"></q-img>
      <q-list>
        <q-item
          class="cursor-pointer "
          @click="openPopupSkuCard"
          v-ripple
          clickable
        >
          <q-item-section>
            <div class="text-h6">{{ SkuData.mainproductname }}</div>

            <div class="text-subtitle2">
              {{ SkuData.mainproducttypename }} <br />
              артикул {{ SkuData.sku }}
            </div>
          </q-item-section>
        </q-item>
        <q-table
          v-if="SkuData.mainproductproperties.length !== 0"
          title="Характеристики"
          dense
          :data="SkuData.mainproductproperties"
          :columns="mainproductpropertiesChars"
          row-key="name"
          hide-header
          flat
          hide-bottom
          :pagination="pagination"
          binary-state-sort
        />
      </q-list>
    </q-card>
    <q-dialog maximized v-model="popupOpened">
      <q-layout view="lhh LpR lff" container class="bg-white">
        <q-header class="bg-primary">
          <q-toolbar>
            <q-toolbar-title>Аналоги</q-toolbar-title>
            <q-btn flat v-close-popup round dense icon="ion-close" />
          </q-toolbar>
        </q-header>

        <q-page-container>
          <q-page>
            <div class="row items-start justify-center">
              <q-card square class="col-12">
                <q-card-section>
                  <div class="text-h6">Артикул</div>
                  <div class="text-subtitle2">{{ SkuData.sku }}</div>
                </q-card-section>

                <q-card-section class="row no-padding">
                  <q-card-section class="col-xs-12 col-sm-6 col-lg-8">
                    <q-markup-table dense flat>
                      <thead>
                        <tr>
                          <th colspan="3">
                            <div class="text-left text-subtitle1">
                              Характеристики
                            </div>
                          </th>
                        </tr>
                      </thead>
                      <tbody>
                        <tr
                          v-for="(v, i) in SkuData.mainproductproperties"
                          :key="i"
                        >
                          <td class="text-left" style="width: 5px;">vo:</td>
                          <td class="text-left">{{ v.name }}</td>
                          <td class="text-right">{{ v.value }}</td>
                        </tr>
                      </tbody>
                    </q-markup-table>
                  </q-card-section>
                  <q-card-section class="col-xs-12 col-sm-6 col-lg-4">
                    <q-carousel
                      swipeable
                      animated
                      v-model="slide"
                      thumbnails
                      infinite
                      class="main"
                    >
                      <q-carousel-slide :name="1" img-src="~assets/logo.png" />
                    </q-carousel>
                  </q-card-section>
                </q-card-section>

                <q-card-section class="row">
                  <q-markup-table flat class="col-12 col-sm-6" dense>
                    <thead>
                      <tr>
                        <th colspan="3">
                          <div class="text-left text-subtitle1">
                            Оригинальные и справочные номера
                          </div>
                        </th>
                      </tr>
                    </thead>
                    <tbody>
                      <tr>
                        <td class="text-left" style="width: 5px;">
                          Производитель
                        </td>
                        <td class="text-left">Номера</td>
                      </tr>
                      <tr>
                        <td class="text-left">Bo</td>
                        <td class="text-left">
                          сы143, cs123, 1231, casdf234, 21340
                        </td>
                      </tr>
                      <tr>
                        <td class="text-left">Cargo</td>
                        <td class="text-left">123123,1231,123 12,31,23,1,23</td>
                      </tr>
                    </tbody>
                  </q-markup-table>

                  <div class="col-12 col-sm-6">
                    <div class="text-subtitle2">Применяется для авто</div>
                    <q-scroll-area visible style="height: 200px;">
                      <q-list dense>
                        <q-item>
                          DAEWOO Lacetti 1.4 16V [F14D3]
                        </q-item>
                        <q-item>
                          DAEWOO Lacetti 1.6 [F16D3]
                        </q-item>
                      </q-list>
                    </q-scroll-area>
                  </div>
                </q-card-section>
              </q-card>
              <div class="col-12">
                <h6 class="text-center">Доступные аналоги : Количество</h6>
              </div>

              <div
                v-for="(product, i) in SkuData.products"
                :key="i"
                class="q-pa-xs col-xs-12 col-sm-6 col-md-4 col-lg-3"
              >
                <q-card flat bordered>
                  <q-card-section>
                    <div class="text-h6">Артикул {{ product.sku }}</div>
                    <div class="text-h6">Наименование {{ product.name }}</div>
                    <div class="text-subtitle2">
                      Производитель {{ product.manufacturer }}
                    </div>
                    <div class="text-subtitle1">
                      {{ product.typename }}
                    </div>
                  </q-card-section>
                  <q-img src="~assets/no-photo.svg"></q-img>
                  <!--                  <q-img src="~assets/no-photo.svg"></q-img>-->
                  <q-card-section>
                    <div class="text-h6">
                      Описание {{ product.description }}
                    </div>
                    <!--                    <div class="text-subtitle2">cs150</div>-->
                  </q-card-section>
                </q-card>
              </div>
            </div>
          </q-page>
        </q-page-container>
      </q-layout>
    </q-dialog>
  </div>
</template>

<script>
export default {
  props: {
    SkuData: {
      GUID: String,
      Name: String,
      Spec: String,
      Amount: Number,
      SKU: String,
      Description: String,
      Price: Number,
      MainImage: String,
      mainproductproperties: Array
    }
  },
  data() {
    return {
      pagination: {
        rowsPerPage: 0,
        sortBy: 'name',
        descending: false
      },
      mainproductpropertiesChars: [
        {
          descending: true,
          name: 'name',
          field: row => row.name.split(';', 1),
          align: 'left'
        },
        {
          name: 'value',
          field: 'value'
        }
      ],
      expanded: false,
      slide: 1,
      popupOpened: false,
      lorem:
        'lorem lorem lorem loremlorem loremlorem loremlorem loremlorem lorem'
    }
  },
  methods: {
    openPopupSkuCard() {
      this.popupOpened = true
    }
  }
}
</script>

<style lang="scss" scoped>
.main {
  min-height: 50vh;
}
</style>
