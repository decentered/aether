"use strict";
/*
This is the main entry point to the client app. See app.vue for the start logic, and globally-applicable css.
*/
var Vue = require('../../node_modules/vue/dist/vue.js');
var VueRouter = require('../../node_modules/vue-router').default;
Vue.use(VueRouter);
// Register icons for our own use.
var Icon = require('../../node_modules/vue-awesome');
Vue.component('icon', Icon);
// Register the click-outside component
var ClickOutside = require('../../node_modules/v-click-outside');
Vue.use(ClickOutside);
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
// const tippy = require('../../node_modules/tippy.js/dist/tippy.js')
// tippy('[tippytitle]')
var Home = require('./components/locations/home.vue').default;
var Popular = require('./components/locations/popular.vue').default;
var BoardScope = require('./components/locations/boardscope.vue').default;
var BoardRoot = require('./components/locations/boardscope/boardroot.vue').default;
var ThreadScope = require('./components/locations/threadscope.vue').default;
// Define routes
var routes = [
    { path: '/', component: Home },
    { path: '/popular', component: Popular },
    {
        path: '/board/:board_id', component: BoardScope,
        children: [
            { path: '', component: BoardRoot },
            { path: '/board/:board_id/thread/:thread_id', component: ThreadScope },
            { path: '*', redirect: '/board/:board_id' }
        ]
    },
    { path: '*', redirect: '/' }
];
var router = new VueRouter({
    routes: routes,
});
var Store = require('./store').default;
new Vue({
    el: '#app',
    template: '<a-app></a-app>',
    router: router,
    store: Store
});
//# sourceMappingURL=main.js.map