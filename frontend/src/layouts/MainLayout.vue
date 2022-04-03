<template>
  <q-layout view="lHh Lpr lFf">
    <q-header elevated>
      <q-toolbar>
        <q-btn
          flat
          dense
          round
          @click="leftDrawerOpen = !leftDrawerOpen"
          icon="ion-menu"
          aria-label="Menu"
        />

        <q-toolbar-title>
          <q-img
            @click="$router.push('/')"
            src="~assets/logo.png"
            style="height: 50px; max-width: 100px;"
          />
        </q-toolbar-title>

        <q-input
          v-if="$route.path === '/amount'"
          @change="searchProducts"
          dark
          dense
          standout
          v-model="text"
          input-class="text-right"
          class="q-ml-md text-h5"
          placeholder="Поиск"
        >
          <template v-slot:append>
            <q-icon v-if="text === ''" name="ion-search" />
            <q-icon
              v-else
              name="ion-close"
              class="cursor-pointer"
              @click="clearProducts"
            />
          </template>
        </q-input>

        <q-btn
          unelevated
          icon-right="ion-cart"
          v-if="$route.path === '/products'"
          to="cart"
          color="accent"
          >Корзина</q-btn
        >
      </q-toolbar>
    </q-header>

    <q-drawer
      v-model="leftDrawerOpen"
      show-if-above
      bordered
      content-class="bg-grey-1"
    >

      <q-list>
        <EssentialLink
          v-for="link in essentialLinks"
          :key="link.title"
          v-bind="link"
        />

        <q-separator color="primary" />

        <q-item>
          <q-item-section class="q-gutter-xs">
            <p class="text-subtitle1">Контакты</p>
            <p class="text-body1 text-bold">
              г.Пятигорск
              <q-btn
                type="a"
                unelevated
                outline
                icon="ion-map"
                href="https://yandex.ru/maps/-/C0CbRCPy"
              />
            </p>
            <p class="text-body1 text-bold">
              г.Пятигорск, пос.Горячеводск ул.Совхозная, 85
              <q-btn
                type="a"
                unelevated
                outline
                icon="ion-map"
                href="https://yandex.ru/maps/-/C0Cb7Jmi"
              />
            </p>
            <p class="text-body2">трасса Ростов-Баку, рынок "Бетта"</p>
            <q-btn type="a" href="tel:+79887575225" outline
              >+7988-75-75-225</q-btn
            >
            <q-btn type="a" href="tel:+79283222555" outline
              >+7928-3-222-555</q-btn
            >
            <div class="row justify-center q-gutter-xs">
              <q-btn
                color="primary"
                icon="ion-logo-instagram"
                type="a"
                href="https://www.instagram.com/restar_26/"
                label="restar_26"
              />
            </div>
          </q-item-section>
        </q-item>

        <q-item>
          <q-item-section>
            <a href="https://webmaster.yandex.ru/siteinfo/?site=https://re-star.ru">
              <img width="88" height="31" alt="" border="0" src="https://yandex.ru/cycounter?https://re-star.ru&theme=light&lang=ru"/>
            </a>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

  </q-layout>
</template>

<script>
import EssentialLink from 'components/EssentialLink'
import { useQuasar } from 'quasar'
import { useAmountStore } from 'src/stores/amount'
import { ref } from 'vue'

export default {
  name: 'MainLayout',

  components: {
    EssentialLink
  },

  setup (props) {
    const store = useAmountStore()
    const $q = useQuasar()
    const text = ref('')

    const searchProducts = async () => {
      if (text.value === '') {
        $q.notify({
          message: 'минимум 3 символа для поиска',
          color: 'red'
        })
        return
      }

      try {
        await store.searchProducts(text.value)
      } catch (e) {
        $q.notify({
          color: 'red',
          message: e.message,
          caption: e.error
        })
      }
    }

    return {
      text,
      searchProducts
    }
  },

  data () {
    return {
      acceptLicense: false,
      openLicense: false,
      email: '',
      password: '',
      loginDialog: false,
      registrationDialog: false,
      leftDrawerOpen: false,
      essentialLinks: [
        {
          title: 'Главная',
          caption: 'restar',
          icon: 'ion-home',
          link: '/'
        },
        // {
        //   title: 'Продукты',
        //   caption: 'products',
        //   icon: 'ion-home',
        //   link: '/products'
        // },
        // {
        //   title: 'Корзина',
        //   caption: 'cart',
        //   icon: 'ion-cart',
        //   link: '/cart'
        // },
        // {
        //   title: 'Менеджер каталога',
        //   caption: 'catalog-manager',
        //   icon: 'ion-folder',
        //   link: '/catalog-manager'
        // },
        {
          title: 'Каталог',
          caption: 'catalog',
          icon: 'ion-folder',
          link: '/catalog'
        },
        // {
        //   title: 'Предложения',
        //   caption: 'offers',
        //   icon: 'ion-flame',
        //   link: '/offers'
        // },
        {
          title: 'Остатки',
          caption: 'amount',
          icon: 'ion-search',
          link: '/amount'
        }
        // {
        //   title: 'test',
        //   caption: 'test',
        //   icon: 'ion-test',
        //   link: '/test'
        // }
      ]
    }
  },
  methods: {
    clearProducts () {
      this.text = ''
    }
  }
}
</script>
