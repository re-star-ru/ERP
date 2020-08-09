const routes = [
  {
    path: '/',
    component: () => import('layouts/MainLayout.vue'),
    children: [
      { path: '', component: () => import('pages/Index.vue') },
      { path: '/offers', redirect: '/offers/1' },
      {
        path: '/offers/:page',
        name: 'page',
        component: () => import('pages/Offers')
      },
      { path: '/amount', component: () => import('pages/Amount') },
      {
        name: 'catalog',
        path: '/catalog',
        component: () => import('pages/Catalog')
      },
      { path: '/test', component: () => import('pages/Test') },
      { path: '/license', component: () => import('pages/License') }
    ]
  }
]

// Always leave this as last one
routes.push({
  path: '*',
  component: () => import('pages/Error404.vue')
})

export default routes
