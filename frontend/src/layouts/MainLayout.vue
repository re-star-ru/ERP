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

        <q-toolbar-title @click="$router.push('/')">
          <q-img src="~assets/logo.png" style="height: 50px; max-width: 100px;" />
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
            <q-icon v-else name="ion-close" class="cursor-pointer" @click="clearProducts" />
          </template>
        </q-input>
      </q-toolbar>
    </q-header>

    <q-drawer v-model="leftDrawerOpen" show-if-above bordered content-class="bg-grey-1">
      <q-list>
        <q-item v-if="this.$store.getters.isLogged" clickable v-ripple to="/profile">
          <q-item-section avatar>
            <q-icon name="ion-person" />
          </q-item-section>
          <q-item-section>
            <q-item-label>{{ $store.state.auth.email }}</q-item-label>
            <q-item-label caption>
              {{
              $store.state.auth.aclGroup
              }}
            </q-item-label>
          </q-item-section>
        </q-item>
        <q-item v-if="!this.$store.getters.isLogged" @click="openLoginDialog" clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="ion-log-out" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Вход</q-item-label>
            <q-item-label caption>login</q-item-label>
          </q-item-section>
        </q-item>

        <q-item
          v-if="!this.$store.getters.isLogged"
          @click="openRegistrationDialog"
          clickable
          v-ripple
        >
          <q-item-section avatar>
            <q-icon name="ion-person-add" />
          </q-item-section>
          <q-item-section>
            <q-item-label>Регистрация</q-item-label>
            <q-item-label caption>registration</q-item-label>
          </q-item-section>
        </q-item>

        <q-separator color="primary" />

        <EssentialLink v-for="link in activeLinks" :key="link.title" v-bind="link" />

        <q-separator v-if="this.$store.getters.isLogged" color="primary" />
        <q-item v-if="this.$store.getters.isLogged" @click="logout" clickable v-ripple>
          <q-item-section avatar>
            <q-icon name="ion-log-out" />
          </q-item-section>
          <q-item-section>Выход</q-item-section>
        </q-item>

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
            <q-btn type="a" href="tel:+79887575225" outline>+7988-75-75-225</q-btn>
            <q-btn type="a" href="tel:+79283222555" outline>+7928-3-222-555</q-btn>
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
            <a href="https://webmaster.yandex.ru/siteinfo/?site=https://restar26.site">
              <img
                width="88"
                height="31"
                alt
                border="0"
                src="https://yandex.ru/cycounter?https://restar26.site&theme=light&lang=ru"
              />
            </a>
          </q-item-section>
        </q-item>
      </q-list>
    </q-drawer>

    <q-page-container>
      <router-view />
    </q-page-container>

    <q-dialog
      v-model="loginDialog"
      @keyup.enter="checkLogin"
      persistent
      transition-show="scale"
      transition-hide="scale"
    >
      <q-card style="width: 500px; max-width: 80vw;" class="bg-primary text-white">
        <q-card-section>
          <div class="text-h6">Вход</div>
        </q-card-section>
        <q-card-section class="bg-white text-teal">
          <q-input
            filled
            v-model="email"
            label="Email"
            hint="Введите ваш email"
            lazy-rules
            :rules="[val => (val && val.length > 0) || 'Неправильный email']"
          />
          <q-input
            type="password"
            filled
            v-model="password"
            label="Пароль"
            hint="Введите пароль"
            lazy-rules
            :rules="[val => (val && val.length > 0) || 'Неправильный пароль']"
          />
        </q-card-section>
        <q-card-actions align="right" class="bg-white text-primary">
          <q-btn flat label="OK" @click="checkLogin" />
          <q-btn flat label="Отмена" @click="closeLoginDialog" />
        </q-card-actions>
      </q-card>
    </q-dialog>

    <q-dialog
      v-model="registrationDialog"
      @keyup.enter="registration"
      persistent
      transition-show="scale"
      transition-hide="scale"
    >
      <q-card style="width: 500px; max-width: 80vw;" class="bg-primary text-white">
        <q-card-section>
          <div class="text-h6">Регистрация</div>
        </q-card-section>
        <q-card-section class="bg-white text-teal">
          <q-input
            filled
            v-model="email"
            label="Email"
            hint="Введите ваш email"
            lazy-rules
            :rules="[val => (val && val.length > 0) || 'Неправильный email']"
          />
          <q-input
            type="password"
            filled
            v-model="password"
            label="Пароль"
            hint="Введите пароль"
            lazy-rules
            :rules="[val => (val && val.length > 0) || 'Неправильный пароль']"
          />

          <div>
            <q-checkbox
              v-model="acceptLicense"
              class="text-primary"
              label="Я принимаю лицензионное"
            />
            <q-btn flat unelevated dense type="a" href="/license" target="_blank">Соглашение</q-btn>
          </div>
        </q-card-section>

        <q-card-actions align="around" class="bg-white text-primary">
          <q-btn flat label="OK" @click="registration" />
          <q-btn flat label="Отмена" @click="closeRegistrationDialog" />
        </q-card-actions>
      </q-card>
    </q-dialog>
  </q-layout>
</template>

<script>
import EssentialLink from 'components/EssentialLink'
import { openURL } from 'quasar'

export default {
  name: 'MainLayout',

  components: {
    EssentialLink,
  },

  data() {
    return {
      acceptLicense: false,
      openLicense: false,
      email: '',
      password: '',
      loginDialog: false,
      registrationDialog: false,
      text: '',
      leftDrawerOpen: false,
      essentialLinks: [
        {
          title: 'Главная',
          caption: 'restar',
          icon: 'ion-home',
          link: '/',
        },
        {
          title: 'Продукты',
          caption: 'products',
          icon: 'ion-home',
          link: '/products',
        },
        {
          title: 'Каталог',
          caption: 'catalog',
          icon: 'ion-folder',
          link: '/catalog',
        },
        {
          title: 'Предложения',
          caption: 'offers',
          icon: 'ion-flame',
          link: '/offers',
        },
        {
          title: 'Остатки',
          caption: 'amount',
          icon: 'ion-search',
          link: '/amount',
        },
        {
          title: 'test',
          caption: 'test',
          icon: 'ion-test',
          link: '/test',
          onlyManager: true,
        },
      ],
    }
  },
  methods: {
    async searchProducts() {
      try {
        await this.$store.dispatch('searchProducts', this.text)
      } catch (e) {
        this.$q.notify({
          message: 'Ошибка',
          color: 'red',
        })
        console.log(e)
      }
    },
    clearProducts() {
      this.text = ''
      this.$store.commit('clearProducts')
    },
    openLoginDialog() {
      this.loginDialog = true
    },
    closeLoginDialog() {
      this.email = ''
      this.password = ''
      this.loginDialog = false
    },
    openRegistrationDialog() {
      console.log('open registration')
      this.registrationDialog = true
    },
    closeRegistrationDialog() {
      this.registrationDialog = false
    },
    async checkLogin() {
      console.log('checklogin')
      let credentials = {
        email: this.email,
        password: this.password,
      }

      try {
        await this.$store.dispatch('login', { credentials })
        this.$q.notify({
          message: ' Вы вошли',
          color: 'accent',
        })
        this.closeLoginDialog()
      } catch (e) {
        this.$q.notify({
          message: e.message,
          color: 'red',
        })
      }
    },

    logout() {
      console.log(this.$store.getters.isLogged)
      this.$store.dispatch('logout')
    },
    async registration() {
      let credentials = {
        email: this.email,
        password: this.password,
      }

      try {
        await this.$store.dispatch('registration', { credentials })
        this.$q.notify({
          message: ' Вы зарегистрировались и вошли',
          color: 'accent',
        })
        this.closeRegistrationDialog()
      } catch (e) {
        this.$q.notify({
          message: e.response.data,
          color: 'red',
        })
        console.dir(e)
      }
    },
  },
  computed: {
    activeLinks() {
      if (
        this.$store.getters.isLogged &&
        this.$store.getters.aclGroup === 'manager'
      ) {
        return this.essentialLinks
      }
      return this.essentialLinks.filter((link) => link.onlyManager !== true)
    },
  },
}
</script>
