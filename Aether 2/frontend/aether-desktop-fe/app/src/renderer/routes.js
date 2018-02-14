export default [
  {
    path: '/aether-entity-view',
    name: 'Aether entity view',
    component: require('components/AetherEntityView')
  },
  {
    path: '/board-entity-view',
    name: 'Board entity view',
    component: require('components/BoardEntityView')
  },
  {
    path: '/thread-entity-view',
    name: 'Thread entity view',
    component: require('components/ThreadEntityView')
  },
  {
    path: '/post-entity-view',
    name: 'Post entity view',
    component: require('components/PostEntityView')
  },
  {
    path: '/user-entity-view',
    name: 'User entity view',
    component: require('components/UserEntityView')
  },
  {
    path: '/front-page-view',
    name: 'Front page view',
    component: require('components/FrontPageView')
  },
  {
    path: '/notifications-view',
    name: 'Notifications',
    component: require('components/NotificationsView')
  },
  { // This will eventually end up being split to a few other pages, likely.
    path: '/settings',
    name: 'Settings view',
    component: require('components/SettingsView')
  },
  {
    path: '*',
    redirect: '/thread-entity-view'
  }
]
