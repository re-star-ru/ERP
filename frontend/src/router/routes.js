const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/PageIndex.vue') },
      { path: '/offers', redirect: '/offers/1' },
      {
        path: '/profile',
        name: 'profile',
        component: () => import('pages/PageProfile')
      },

      {
        path: '/products',
        name: 'products',
        component: () => import('pages/PageProducts')
      },
      {
        path: '/cart',
        name: 'cart',
        component: () => import('pages/PageCart')
      },
      {
        path: '/catalog-manager',
        name: 'catalog-manager',
        component: () => import('pages/PageCatalogManager'),
        children: [
          {
            path: ':id',
            component: () => import('components/ProductCardEditor'),
            name: 'editproductbyid',
            props: true
          }
        ]
      },
      {
        path: '/offers/:page',
        name: 'page',
        component: () => import('pages/PageOffers')
      },
      { path: '/amount', component: () => import('pages/PageAmount') },
      {
        name: 'catalog',
        path: '/catalog',
        component: () => import('pages/PageCatalog')
      },
      { path: '/test', component: () => import('pages/PageTest') },
      { path: '/license', component: () => import('pages/PageLicense') }
    ]
  }
]

// Always leave this as last one
routes.push({
  path: '/:catchAll(.*)"',
  component: () => import('pages/PageError404.vue')
})

export default routes
