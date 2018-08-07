"use strict";
/*
This is the main entry point to the client app. See app.vue for the start logic, and globally-applicable css.
*/
/*----------  Electron and our non-GUI services  ----------*/
/*
  Main thing to realise is that there are two processes on the electron side, the main and the renderer. Main is basically a node process, and renderer is basically very similar to script that linked via the <script> tag of a frame that the node has surfaced, which means it has much fewer privileges (though still has access to Node APIs).

  As a result, we start the frontend binary from the main side, but we establish the client gRPC server on the renderer side when it is time to connect to that server, since the data frontend provides needs to be delivered to the renderer, not the main process.
*/
// Electron IPC setup before doing anything else
var ipc = require('../../node_modules/electron-better-ipc');
var clapiserver = require('./services/clapiserver/clapiserver');
var feapiconsumer = require('./services/feapiconsumer/feapiconsumer');
var clientAPIServerPort = clapiserver.StartClientAPIServer();
console.log('renderer client api server port: ', clientAPIServerPort);
ipc.callMain('SetClientAPIServerPort', clientAPIServerPort).then(function (feDaemonStarted) {
    if (!feDaemonStarted) {
        // It's an Electron refresh, not a cold start.
        feapiconsumer.Initialise();
    }
});
/*----------  Vue + its plugins  ----------*/
var Vue = require('../../node_modules/vue/dist/vue.js');
var VueRouter = require('../../node_modules/vue-router').default;
Vue.use(VueRouter);
// Register icons for our own use.
var Icon = require('../../node_modules/vue-awesome');
Vue.component('icon', Icon);
// Register the click-outside component
var ClickOutside = require('../../node_modules/v-click-outside');
Vue.use(ClickOutside);
/*----------  Third party dependencies  ----------*/
var Mousetrap = require('../../node_modules/mousetrap');
/*----------  Components  ----------*/
// Global component declarations - do it here once.
Vue.component('a-app', require('./components/a-app.vue').default);
Vue.component('a-header', require('./components/a-header.vue').default);
Vue.component('a-header-icon', require('./components/a-header-icon.vue').default);
Vue.component('a-sidebar', require('./components/a-sidebar.vue').default);
Vue.component('a-boardheader', require('./components/a-boardheader.vue').default);
Vue.component('a-tabs', require('./components/a-tabs.vue').default);
Vue.component('a-thread-entity', require('./components/a-thread-entity.vue').default);
Vue.component('a-vote-action', require('./components/a-vote-action.vue').default);
Vue.component('a-thread-header-entity', require('./components/a-thread-header-entity.vue').default);
Vue.component('a-post', require('./components/a-post.vue').default);
Vue.component('a-side-header', require('./components/a-side-header.vue').default);
Vue.component('a-breadcrumbs', require('./components/a-breadcrumbs.vue').default);
Vue.component('a-username', require('./components/a-username.vue').default);
Vue.component('a-timestamp', require('./components/a-timestamp.vue').default);
Vue.component('a-globalscopeheader', require('./components/a-globalscopeheader.vue').default);
Vue.component('a-board-entity', require('./components/a-board-entity.vue').default);
Vue.component('a-hashimage', require('./components/a-hashimage.vue').default);
Vue.component('a-no-content', require('./components/a-no-content.vue').default);
Vue.component('a-markdown', require('./components/a-markdown.vue').default);
Vue.component('a-avatar-block', require('./components/a-avatar-block.vue').default);
// const tippy = require('../../node_modules/tippy.js/dist/tippy.js')
// tippy('[tippytitle]')
/*----------  Places  ----------*/
var Home = require('./components/locations/home.vue').default;
var Popular = require('./components/locations/popular.vue').default;
var GlobalScope = require('./components/locations/globalscope.vue').default;
var GlobalRoot = require('./components/locations/globalscope/globalroot.vue').default;
var GlobalSubbed = require('./components/locations/globalscope/subbedroot.vue').default;
var BoardScope = require('./components/locations/boardscope.vue').default;
var BoardRoot = require('./components/locations/boardscope/boardroot.vue').default;
var ThreadScope = require('./components/locations/threadscope.vue').default;
var SettingsScope = require('./components/locations/settingsscope.vue').default;
var SettingsRoot = require('./components/locations/settingsscope/settingsroot.vue').default;
var AdvancedSettings = require('./components/locations/settingsscope/advancedsettings.vue').default;
var About = require('./components/locations/settingsscope/about.vue').default;
var Membership = require('./components/locations/settingsscope/membership.vue').default;
var Changelog = require('./components/locations/settingsscope/changelog.vue').default;
var Intro = require('./components/locations/settingsscope/intro.vue').default;
var UserScope = require('./components/locations/userscope.vue').default;
var UserRoot = require('./components/locations/userscope/userroot.vue').default;
/*----------  Routes  ----------*/
var routes = [
    { path: '/', component: Home, name: 'Home', },
    { path: '/popular', component: Popular, name: 'Popular', },
    {
        path: '/globalscope', component: GlobalScope,
        children: [
            { path: '', component: GlobalRoot, name: 'Global', },
            { path: '/globalscope/subbed', component: GlobalSubbed, name: 'Global>Subbed', },
        ]
    },
    {
        path: '/board/:boardfp', component: BoardScope,
        children: [
            { path: '', component: BoardRoot, name: 'Board', },
            { path: '/board/:boardfp/thread/:threadfp', component: ThreadScope, name: 'Thread', },
            { path: '*', redirect: '/board/:boardfp' }
        ]
    },
    {
        path: '/settings', component: SettingsScope,
        children: [
            { path: '/intro', component: Intro, name: 'Intro', },
            { path: '', component: SettingsRoot, name: 'Settings', },
            { path: '/settings/advanced', component: AdvancedSettings, name: 'Settings>Advanced', },
            { path: '/about', component: About, name: 'About', },
            { path: '/membership', component: Membership, name: 'Membership', },
            { path: '/changelog', component: Changelog, name: 'Changelog', },
        ]
    },
    {
        path: '/user/:userfp', component: UserScope,
        children: [
            { path: '', component: UserRoot, name: 'User' },
            { path: '*', redirect: '/user/:userfp' }
        ]
    },
    { path: '*', redirect: '/' }
];
// { path: '/user/:userfp/posts', component: UserPosts, name: 'User>Posts', },
// { path: '/user/:userfp/threads', component: UserThreads, name: 'User>Threads', },
/*----------  Plumbing  ----------*/
var router = new VueRouter({
    routes: routes,
});
var Store = require('./store').default;
new Vue({
    el: '#app',
    template: '<a-app></a-app>',
    router: router,
    store: Store,
});
var Sync = require('../../node_modules/vuex-router-sync').sync;
Sync(Store, router);
/*
^ It adds a route module into the store, which contains the state representing the current route:
store.state.route.path   // current path (string)
store.state.route.params // current params (object)
store.state.route.query  // current query (object)
*/
// Disable events that are meaningless in this context.
// Drag start is being able to click and drag a link inside the app to outside of it. Since the app is a local one, that link will just be a local file, and it won't be useful to anybody.
document.addEventListener('dragstart', function (event) { event.preventDefault(); });
// Dragover is the event that gets fired when a dragged item is on a droppable target, every few hundred milliseconds. We have no drop targets.
document.addEventListener('dragover', function (event) { event.preventDefault(); });
// Cancelling drop prevents anything from being dropped into the container. This can be a mild security risk, if someone can convince you (or somehow automate dropping inside the app container), it can make the container ping a web address. This also assumes the container has the dropped remote address whitelisted, though, so it's a long shot. Still, defence in depth is preferable.
document.addEventListener('drop', function (event) { event.preventDefault(); });
/*----------  Some basic keyboard shortcuts  ----------*/
// Backspace: go back (except when a text field or otherwise editable object is selected)
// Some people just want to watch the world burn. ;)
Mousetrap.bind('backspace', function (event) {
    if (event.keyCode === 8 && event.target.nodeName.toLowerCase() !== 'textarea' && event.target.nodeName.toLowerCase() !== 'input' && event.target.contentEditable !== 'true') {
        history.back();
    }
});
// Ctrl + Backspace: go forward (except when a text field or otherwise editable object is selected)
Mousetrap.bind('ctrl+backspace', function (event) {
    if (event.keyCode === 8 && event.target.nodeName.toLowerCase() !== 'textarea' && event.target.nodeName.toLowerCase() !== 'input' && event.target.contentEditable !== 'true') {
        history.forward();
    }
});
/*
  Why? This is an controlled environment, not a browser. All text field entries are saved to cache and cannot be lost. If for some reason you manage to back out of your compose screen, you can always forward or click again and it will be back into the exact same spot, with your inserted text still there.

  The problem with backspace is caused by its propensity to cause data loss in a browser environment. When data loss is not possible, its potential harm is neutralised.
*/ 
//# sourceMappingURL=renderermain.js.map