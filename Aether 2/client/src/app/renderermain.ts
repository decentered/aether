/*
This is the main entry point to the client app. See app.vue for the start logic, and globally-applicable css.
*/

/*----------  Electron and our non-GUI services  ----------*/
/*
  Main thing to realise is that there are two processes on the electron side, the main and the renderer. Main is basically a node process, and renderer is basically very similar to script that linked via the <script> tag of a frame that the node has surfaced, which means it has much fewer privileges (though still has access to Node APIs).

  As a result, we start the frontend binary from the main side, but we establish the client gRPC server on the renderer side when it is time to connect to that server, since the data frontend provides needs to be delivered to the renderer, not the main process.
*/

// Electron IPC setup before doing anything else
var ipc = require('../../node_modules/electron-better-ipc')
const clapiserver = require('./services/clapiserver/clapiserver')
const feapiconsumer = require('./services/feapiconsumer/feapiconsumer')
const clientAPIServerPort: number = clapiserver.StartClientAPIServer()
console.log('renderer client api server port: ', clientAPIServerPort)
ipc.callMain('SetClientAPIServerPort', clientAPIServerPort).then(function(feDaemonStarted: boolean) {
  if (!feDaemonStarted) {
    // It's an Electron refresh, not a cold start.
    feapiconsumer.Initialise()
  }
})

/*----------  Vue + its plugins  ----------*/


var Vue = require('../../node_modules/vue/dist/vue.js')
var VueRouter = require('../../node_modules/vue-router').default
Vue.use(VueRouter)

// Register icons for our own use.
var Icon = require('../../node_modules/vue-awesome')
Vue.component('icon', Icon)

// Register the click-outside component
var ClickOutside = require('../../node_modules/v-click-outside')
Vue.use(ClickOutside)


/*----------  Third party dependencies  ----------*/

var Mousetrap = require('../../node_modules/mousetrap')
// var Spinner = require('../../node_modules/vue-simple-spinner')

/*----------  Components  ----------*/

// Global component declarations - do it here once.
Vue.component('a-app', require('./components/a-app.vue').default)
Vue.component('a-header', require('./components/a-header.vue').default)
Vue.component('a-header-icon', require('./components/a-header-icon.vue').default)
Vue.component('a-sidebar', require('./components/a-sidebar.vue').default)
Vue.component('a-boardheader', require('./components/a-boardheader.vue').default)
Vue.component('a-tabs', require('./components/a-tabs.vue').default)
Vue.component('a-thread-entity', require('./components/a-thread-entity.vue').default)
Vue.component('a-vote-action', require('./components/a-vote-action.vue').default)
Vue.component('a-thread-header-entity', require('./components/a-thread-header-entity.vue').default)
Vue.component('a-post', require('./components/a-post.vue').default)
Vue.component('a-side-header', require('./components/a-side-header.vue').default)
Vue.component('a-breadcrumbs', require('./components/a-breadcrumbs.vue').default)
Vue.component('a-username', require('./components/a-username.vue').default)
Vue.component('a-timestamp', require('./components/a-timestamp.vue').default)
Vue.component('a-globalscopeheader', require('./components/a-globalscopeheader.vue').default)
Vue.component('a-board-entity', require('./components/a-board-entity.vue').default)
Vue.component('a-hashimage', require('./components/a-hashimage.vue').default)
Vue.component('a-no-content', require('./components/a-no-content.vue').default)
Vue.component('a-markdown', require('./components/a-markdown.vue').default)
Vue.component('a-avatar-block', require('./components/a-avatar-block.vue').default)
Vue.component('a-composer', require('./components/a-composer.vue').default)
Vue.component('a-ballot', require('./components/a-ballot.vue').default)
Vue.component('a-progress-bar', require('./components/a-progress-bar.vue').default)
Vue.component('a-inflight-info', require('./components/a-inflight-info.vue').default)
Vue.component('a-info-marker', require('./components/a-info-marker.vue').default)
Vue.component('a-spinner', require('./components/a-spinner.vue').default)
Vue.component('a-notfound', require('./components/a-notfound.vue').default)

/*----------  Third party components  ----------*/

Vue.component('vue-simple-spinner', require('../../node_modules/vue-simple-spinner'))


/*----------  Places  ----------*/

const Home = require('./components/locations/home.vue').default
const Popular = require('./components/locations/popular.vue').default

/*----------  Global scope (whole network, i.e. list of boards)  ----------*/
const GlobalScope = require('./components/locations/globalscope.vue').default
const GlobalRoot = require('./components/locations/globalscope/globalroot.vue').default
const GlobalSubbed = require('./components/locations/globalscope/subbedroot.vue').default

/*----------  Board scope (board entity + list of threads)  ----------*/
const NewBoard = require('./components/locations/globalscope/newboard.vue').default
const BoardScope = require('./components/locations/boardscope.vue').default
const BoardRoot = require('./components/locations/boardscope/boardroot.vue').default
const BoardInfo = require('./components/locations/boardscope/boardinfo.vue').default
const ModActivity = require('./components/locations/boardscope/modactivity.vue').default
const Elections = require('./components/locations/boardscope/elections.vue').default

/*----------  Thread scope (thread entity + list of posts)  ----------*/
const NewThread = require('./components/locations/boardscope/newthread.vue').default
const ThreadScope = require('./components/locations/threadscope.vue').default

/*----------  Settings scope  ----------*/
const SettingsScope = require('./components/locations/settingsscope.vue').default
const SettingsRoot = require('./components/locations/settingsscope/settingsroot.vue').default
const AdvancedSettings = require('./components/locations/settingsscope/advancedsettings.vue').default
const About = require('./components/locations/settingsscope/about.vue').default
const Membership = require('./components/locations/settingsscope/membership.vue').default
const Changelog = require('./components/locations/settingsscope/changelog.vue').default
const AdminsQuickstart = require('./components/locations/settingsscope/adminsquickstart.vue').default
const Intro = require('./components/locations/settingsscope/intro.vue').default
const NewUser = require('./components/locations/settingsscope/newuser.vue').default

/*----------  User scope  ----------*/
const UserScope = require('./components/locations/userscope.vue').default
const UserRoot = require('./components/locations/userscope/userroot.vue').default
const UserBoards = require('./components/locations/userscope/userboards.vue').default
const UserThreads = require('./components/locations/userscope/userthreads.vue').default
const UserPosts = require('./components/locations/userscope/userposts.vue').default

/*----------  Status scope  ----------*/

const Status = require('./components/locations/status.vue').default

/*----------  Routes  ----------*/

const routes = [
  { path: '/', component: Home, name: 'Home', },
  { path: '/popular', component: Popular, name: 'Popular', },
  {
    path: '/globalscope', component: GlobalScope,
    children: [
      { path: '', component: GlobalRoot, name: 'Global', },
      { path: '/globalscope/subbed', component: GlobalSubbed, name: 'Global>Subbed', },
      { path: '/globalscope/newboard', component: NewBoard, name: 'Global>NewBoard', },
    ]
  },
  {
    path: '/board/:boardfp', component: BoardScope,
    children: [
      { path: '', component: BoardRoot, name: 'Board', },
      { path: '/board/:boardfp/new', component: BoardRoot, name: 'Board>ThreadsNewList', },
      { path: '/board/:boardfp/info', component: BoardInfo, name: 'Board>BoardInfo', },
      { path: '/board/:boardfp/modactivity', component: ModActivity, name: 'Board>ModActivity', },
      { path: '/board/:boardfp/elections', component: Elections, name: 'Board>Elections', },
      { path: '/board/:boardfp/newthread', component: NewThread, name: 'Board>NewThread', },
    ]
  }, {
    path: '/board/:boardfp/thread/:threadfp', component: ThreadScope, name: 'Thread',
  },
  {
    path: '/settings', component: SettingsScope,
    children: [
      { path: '', component: SettingsRoot, name: 'Settings', },
      { path: '/settings/advanced', component: AdvancedSettings, name: 'Settings>Advanced', },
      /*This is a little weird, these things are in settings scope but they're not in a settings path. That's because they exist in a router link that is in the settings structure. If you move this outside and try to use it, it uses the router link outside settings, which is the main main-block router link, which means the settings frame box won't be rendered. So this is not an oversight. */
      { path: '/intro', component: Intro, name: 'Intro', },
      { path: '/about', component: About, name: 'About', },
      { path: '/membership', component: Membership, name: 'Membership', },
      { path: '/changelog', component: Changelog, name: 'Changelog', },
      { path: '/adminsquickstart', component: AdminsQuickstart, name: 'AdminsQuickstart', },
      { path: '/newuser', component: NewUser, name: 'NewUser', },
    ]
  },
  {
    path: '/user/:userfp', component: UserScope,
    children: [
      { path: '', component: UserRoot, name: 'User' },
      { path: '/user/:userfp/boards', component: UserBoards, name: 'User>Boards' },
      { path: '/user/:userfp/threads', component: UserThreads, name: 'User>Threads' },
      { path: '/user/:userfp/posts', component: UserPosts, name: 'User>Posts' },
      { path: '*', redirect: '/user/:userfp' }
    ]
  },
  { path: '/status', component: Status, name: 'Status', },

  { path: '*', redirect: '/' }
]

// { path: '/user/:userfp/posts', component: UserPosts, name: 'User>Posts', },
// { path: '/user/:userfp/threads', component: UserThreads, name: 'User>Threads', },

/*----------  Plumbing  ----------*/

const router = new VueRouter({
  scrollBehavior() { // Always return to top while navigating.
    return { x: 0, y: 0 }
  },
  routes: routes,
  // mode: 'history'
})

const Store = require('./store').default

new Vue({
  el: '#app',
  template: '<a-app></a-app>',
  router: router,
  store: Store,
})


let Sync = require('../../node_modules/vuex-router-sync').sync
Sync(Store, router)
/*
^ It adds a route module into the store, which contains the state representing the current route:
store.state.route.path   // current path (string)
store.state.route.params // current params (object)
store.state.route.query  // current query (object)
*/

// Disable events that are meaningless in this context.

// Drag start is being able to click and drag a link inside the app to outside of it. Since the app is a local one, that link will just be a local file, and it won't be useful to anybody.
document.addEventListener('dragstart', function(event: any) { event.preventDefault() })

// Dragover is the event that gets fired when a dragged item is on a droppable target, every few hundred milliseconds. We have no drop targets.
document.addEventListener('dragover', function(event: any) { event.preventDefault() })

// Cancelling drop prevents anything from being dropped into the container. This can be a mild security risk, if someone can convince you (or somehow automate dropping inside the app container), it can make the container ping a web address. This also assumes the container has the dropped remote address whitelisted, though, so it's a long shot. Still, defence in depth is preferable.
document.addEventListener('drop', function(event: any) { event.preventDefault() })


/*----------  Some basic keyboard shortcuts  ----------*/

Mousetrap.bind('<', function() {
  history.back()
  // if (event.target.nodeName.toLowerCase() !== 'textarea' && event.target.nodeName.toLowerCase() !== 'input' && event.target.contentEditable !== 'true') {
  //   history.back()
  // }
})

// Ctrl + Backspace: go forward (except when a text field or otherwise editable object is selected)
Mousetrap.bind('>', function() {
  history.forward()
  // if (event.target.nodeName.toLowerCase() !== 'textarea' && event.target.nodeName.toLowerCase() !== 'input' && event.target.contentEditable !== 'true') {
  //   history.forward()
  // }
})

