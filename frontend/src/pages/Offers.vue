<template>
  <q-page class="q-pa-md q-gutter-sm flex flex-center">
    <h5>Горячие заказы</h5>

    <div class="q-pa-md row items-center justify-around q-gutter-md">
      <q-card
        v-if="offersDownloaded"
        class="my-card"
        v-for="n in getCardPerPage"
        :key="n"
      >
        <q-img
          class="cursor-pointer"
          @click="selectOffer(getOfferIndex(n))"
          :src="offers[getOfferIndex(n)].mainImage || imageNotFound"
        >
          <div class="absolute-bottom text-h6">{{ offers[n - 1].name }}</div>
        </q-img>

        <q-card-section>
          <!--          idx = {{ getOfferIndex(n) }}-->
          Артикул
          {{ offers[getOfferIndex(n)].sku }}
          <br />
          Характеристика {{ offers[getOfferIndex(n)].spec }}
          <br />
          Остаток {{ offers[getOfferIndex(n)].amount }} шт.
          <br />
          <!--          Цена {{ offers[getOfferIndex(n)].price }}-->
          <br />
          Цена договорная
          <br />

          Описание {{ offers[getOfferIndex(n)].description }}
        </q-card-section>
      </q-card>
    </div>
    <q-pagination
      :to-fn="goToPage"
      class="self-center"
      v-model="currentPage"
      :max="maxPages"
      :input="true"
    ></q-pagination>
    <q-dialog full-width v-model="carousel" @hide="onOfferHide">
      <q-responsive :ratio="1">
        <q-carousel
          control-type="push"
          swipeable
          animated
          arrows
          :fullscreen.sync="fullscreen"
          v-model="slide"
          control-color="primary"
          navigation
          class="bg-white shadow-1 rounded-borders col"
          infinite
        >
          <q-carousel-slide
            v-for="(v, i) in currentOffer.images"
            :key="i"
            :name="i"
            :img-src="v"
          >
          </q-carousel-slide>

          <template v-slot:control>
            <q-carousel-control
              class="space"
              position="bottom-right"
              :offset="[18, 18]"
            >
              <q-btn
                push
                round
                dense
                color="white"
                text-color="primary"
                :icon="fullscreen ? 'ion-close' : 'ion-expand'"
                @click="fullscreen = !fullscreen"
              />
            </q-carousel-control>

            <q-carousel-control position="top-right" :offset="[18, 18]">
              <q-btn-group>
                <q-btn
                  v-if="enforceAclGroup"
                  icon="ion-camera"
                  @click="openFile"
                  color="accent"
                >
                </q-btn>
                <q-btn
                  color="secondary"
                  icon="ion-open"
                  @click="openImg(currentOffer.images[slide])"
                >
                </q-btn>

                <q-btn color="primary" icon="ion-menu" v-if="enforceAclGroup">
                  <q-menu
                    transition-show="flip-right"
                    transition-hide="flip-left"
                  >
                    <q-list style="min-width: 100px">
                      <q-item>
                        <q-btn
                          text-color="black"
                          color="amber"
                          label="Сделать главной"
                        />
                      </q-item>
                      <q-item>
                        <q-btn
                          label="удалить"
                          icon="ion-trash"
                          @click="deleteImage"
                          color="negative"
                        >
                        </q-btn>
                      </q-item>
                    </q-list>
                  </q-menu>
                </q-btn>
              </q-btn-group>
            </q-carousel-control>
          </template>
        </q-carousel>
      </q-responsive>
    </q-dialog>
    <input
      style="display: none"
      @change="updateImageDisplay"
      ref="imgCapture"
      type="file"
      accept="image/*"
    />
  </q-page>
</template>

<script>
import { openURL } from 'quasar'

export default {
  name: 'PageName',
  data: () => {
    return {
      cardsPerPage: 30,
      lastPageCards: 1,
      currentPage: 1,
      maxPages: 1,
      offersDownloaded: false,
      carousel: false,
      imageNotFound: 'back.png',
      slide: 0,
      fullscreen: false,
      currentOffer: {
        guid: String,
        name: String,
        spec: String,
        amount: Number,
        price: Number,
        mainImage: String,
        images: []
      },
      offers: [
        {
          guid: '',
          name: '',
          spec: '',
          amount: 0,
          price: 0,
          mainImage: '',
          images: ['']
        }
      ]
    }
  },
  methods: {
    getOfferIndex(i) {
      let currentPageIndex = this.currentPage - 1
      let currentIndex = i - 1
      return currentPageIndex * this.cardsPerPage + currentIndex
    },
    openFile() {
      this.$refs.imgCapture.click()
    },
    goToPage(page) {
      return { name: 'page', params: { page: page } }
    },
    async updateOffer() {
      console.log('update')
      try {
        let res = await this.$axios.get(`offer/${this.currentOffer.guid}`)
        this.currentOffer = res.data
        this.slide = 0

        for (let i = 0; i < this.offers.length; i++) {
          if (this.offers[i].guid === res.data.guid) {
            this.offers[i] = res.data
          }
        }
      } catch (e) {
        console.log(e)
      }
    },
    async deleteImage() {
      let params = {
        link: this.currentOffer.images[this.slide],
        el: this.slide
      }
      try {
        let res = await this.$axios.delete(`/image/${this.currentOffer.guid}`, {
          data: params
        })
        await this.updateOffer()
      } catch (e) {
        console.dir(e)
      }

      console.log(params)
    },
    selectOffer(i) {
      this.currentOffer = this.offers[i]
      console.log(i)
      this.carousel = !this.carousel
    },
    openImg(img) {
      openURL(img)
    },
    onOfferHide() {
      this.currentOffer = {
        guid: String,
        name: String,
        spec: String,
        amount: Number,
        sku: String,
        description: String,
        price: Number,
        mainImage: String,
        images: []
      }
    },
    async updateImageDisplay() {
      let files = this.$refs.imgCapture.files
      const form = new FormData()

      form.append('file', files[0])

      try {
        this.$q.loading.show()
        await this.$axios.post(`/image/${this.currentOffer.guid}`, form)
        this.$q.loading.hide()
        this.$q.notify({
          message: 'Загружено успешно',
          color: 'green'
        })
        this.updateOffer()
      } catch (e) {
        this.$q.loading.hide()
        this.$q.notify({
          message: 'Ошибка',
          color: 'red'
        })
      }
      this.$q.loading.hide()
    },
    async downloadOffers() {
      try {
        this.$q.loading.show()
        let res = await this.$axios.get('/offers')
        this.$q.loading.hide()
        console.log(res)
        this.offers = res.data
        this.offersDownloaded = true
        this.setPagesConfig()
      } catch (e) {
        this.$q.loading.hide()
        console.log(e)
      }
    },
    setCurrentPage() {
      console.log(this.$route.params)
      let page = parseInt(this.$route.params.page)
      if (!Number.isNaN(page)) {
        this.currentPage = page
        return
      }
      this.currentPage = 1
    },
    setPagesConfig() {
      this.maxPages = Math.ceil(this.offers.length / this.cardsPerPage)
      this.lastPageCards = this.offers.length % this.cardsPerPage
    }
  },
  computed: {
    enforceAclGroup() {
      return this.$store.getters.aclGroup === 'manager'
    },
    getCardPerPage() {
      console.log('card per page', this.currentPage, this.maxPages)
      if (this.currentPage === this.maxPages) {
        return this.lastPageCards
      }
      return this.cardsPerPage
    }
  },
  mounted() {
    this.downloadOffers()
    this.setCurrentPage()
  },
  beforeRouteUpdate(to, from, next) {
    this.currentPage = to.params.page
    next()
  }
}
</script>

<style lang="scss" scoped>
#input {
  display: none;
}
.my-card {
  width: 100%;
  max-width: 300px;
}

.q-carousel__slide {
  background-size: contain;
  background-repeat: no-repeat;
}

label {
  cursor: pointer;
}
</style>
