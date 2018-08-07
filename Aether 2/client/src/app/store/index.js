"use strict";
/*
This is the data store for the client.

This data store does not hold any persistent data, nor does it cache it. The point of this is to hold the instance data. The frontend is the actual caching and compile logic that regenerates the data to be used as needed.
*/
var __assign = (this && this.__assign) || Object.assign || function(t) {
    for (var s, i = 1, n = arguments.length; i < n; i++) {
        s = arguments[i];
        for (var p in s) if (Object.prototype.hasOwnProperty.call(s, p))
            t[p] = s[p];
    }
    return t;
};
Object.defineProperty(exports, "__esModule", { value: true });
var Vue = require('../../../node_modules/vue/dist/vue.js');
var Vuex = require('../../../node_modules/vuex').default;
Vue.use(Vuex);
var fe = require('../services/feapiconsumer/feapiconsumer');
var dataLoaders = require('./dataloaders').default;
var contentRelations = require('./contentrelations');
var crumbs = require('./crumbs');
var dataLoaderPlugin = function (store) {
    store.watch(
    // When the returned result changes,
    function (state) {
        return state.route.params;
    }, 
    // Run this callback
    function (newValue, oldValue) {
        // First, check if we should refresh.
        if (oldValue === newValue && !store.state.frontendHasUpdates) {
            // if the values are the same, and frontend has no updates, bail.
            return;
        }
        var routeParams = newValue;
        if (store.state.route.name === "Board") {
            store.dispatch('loadBoardScopeData', routeParams.boardfp);
            store.dispatch('setLastSeenForBoard', { fp: routeParams.boardfp });
            return;
        }
        if (store.state.route.name === "Thread") {
            store.dispatch('loadThreadScopeData', {
                boardfp: routeParams.boardfp,
                threadfp: routeParams.threadfp,
            });
            return;
        }
        if (store.state.route.name === "Global" || store.state.route.name === "Global>Subbed") {
            store.dispatch('loadGlobalScopeData');
            return;
        }
        if (store.state.route.name === "User") {
            store.dispatch('loadUserScopeData', {
                fp: routeParams.userfp,
                userreq: true,
                boardsreq: false,
                threadsreq: false,
                postsreq: false,
            });
            return;
        }
        // If none of the special cases, just trigger an update for breadcrumbs.
        store.dispatch('updateBreadcrumbs');
    });
};
var actions = __assign({ 
    // within any of those, context.state is how you access state above.
    setSidebarState: function (context, sidebarOpen) {
        context.commit('SET_SIDEBAR_STATE', sidebarOpen);
    },
    setAmbientBoards: function (context, ambientBoards) {
        context.commit('SET_AMBIENT_BOARDS', ambientBoards);
    },
    setCurrentBoardFp: function (context, fp) {
        // This is a boardscope load. Retrieve board + its threads.
        if (context.state.currentBoardFp !== fp) {
            context.dispatch('setCurrentBoardAndThreads', fp);
            context.commit('SET_CURRENT_BOARD_FP', fp);
            // Same scope, but fe has updated.
            return;
        }
        // current board fp is the same as what we asked for, but FE has updates.
        if (context.state.frontendHasUpdates) {
            context.dispatch('setCurrentBoardAndThreads', fp);
            context.commit('SET_CURRENT_BOARD_FP', fp);
        }
    },
    setCurrentThreadFp: function (context, _a) {
        var boardfp = _a.boardfp, threadfp = _a.threadfp;
        if (context.state.currentThreadFp !== threadfp) {
            context.dispatch('setCurrentThreadAndPosts', { boardfp: boardfp, threadfp: threadfp });
            context.commit('SET_CURRENT_THREAD_FP', threadfp);
            // Same scope, but fe has updated.
            return;
        }
        // current thread fp is the same as what we asked for, but FE has updates.
        if (context.state.frontendHasUpdates) {
            context.dispatch('setCurrentThreadAndPosts', { boardfp: boardfp, threadfp: threadfp });
            context.commit('SET_CURRENT_THREAD_FP', threadfp);
        }
    },
    setCurrentBoardAndThreads: function (context, fp) {
        fe.GetBoardAndThreads(fp, function (resp) {
            context.commit('SET_CURRENT_BOARD', resp.board);
            context.commit('SET_CURRENT_BOARDS_THREADS', resp.threadsList);
            context.dispatch('updateBreadcrumbs');
            // ^ This has to be here because otherwise the BC compute process runs before the data is ready, resulting in empty breadcrumbs.
        });
    },
    setCurrentThreadAndPosts: function (context, _a) {
        var boardfp = _a.boardfp, threadfp = _a.threadfp;
        fe.GetThreadAndPosts(boardfp, threadfp, function (resp) {
            context.commit('SET_CURRENT_BOARD', resp.board);
            context.commit('SET_CURRENT_THREAD', resp.thread);
            context.commit('SET_CURRENT_THREADS_POSTS', resp.postsList);
            context.dispatch('updateBreadcrumbs');
        });
    } }, dataLoaders, crumbs.crumbActions, contentRelations.actions);
var mutations = __assign({ SET_SIDEBAR_STATE: function (state, sidebarOpen) {
        state.sidebarOpen = sidebarOpen;
    },
    SET_AMBIENT_BOARDS: function (state, ambientBoards) {
        state.ambientBoards = ambientBoards;
    },
    SET_CURRENT_BOARD_FP: function (state, fp) {
        state.currentBoardFp = fp;
    },
    SET_CURRENT_THREAD_FP: function (state, fp) {
        state.currentThreadFp = fp;
    },
    SET_CURRENT_BOARD: function (state, board) {
        state.currentBoard = board;
    },
    SET_CURRENT_THREAD: function (state, thread) {
        state.currentThread = thread;
    },
    SET_CURRENT_BOARDS_THREADS: function (state, threads) {
        state.currentBoardsThreads = threads;
    },
    SET_CURRENT_THREADS_POSTS: function (state, posts) {
        state.currentThreadsPosts = posts;
    },
    SET_ALL_BOARDS: function (state, boards) {
        state.allBoards = boards;
    },
    SET_USER_SCOPE_DATA: function (state, resp) {
        if (resp.UserEntityRequested) {
            state.currentUserEntity = resp.user;
        }
        if (resp.BoardsRequested) {
            state.currentUserBoards = resp.Boards;
        }
        if (resp.ThreadsRequested) {
            state.currentUserThreads = resp.Threads;
        }
        if (resp.PostsRequested) {
            state.currentUserPosts = resp.Posts;
        }
    } }, crumbs.crumbMutations, contentRelations.mutations);
var st = new Vuex.Store({
    state: {
        sidebarOpen: true,
        breadcrumbs: [],
        frontendPort: 0,
        ambientBoards: {},
        currentBoardFp: "",
        currentBoard: {},
        currentBoardsThreads: [],
        currentThreadFp: "",
        currentThread: {},
        currentThreadsPosts: [],
        currentUserEntity: {},
        currentUserBoards: [],
        currentUserThreads: [],
        currentUserPosts: [],
        allBoards: [],
        frontendHasUpdates: true,
        route: {},
    },
    actions: actions,
    mutations: mutations,
    plugins: [dataLoaderPlugin],
});
exports.default = st;
/*

Reminder:

changeTestData(context: any) {
  // console.log("yo yo yo ")
  context.commit('editTestData')
}

is the same as:

changeTestData({commit}) {
  // console.log("yo yo yo ")
  commit('editTestData')
}
*/
//# sourceMappingURL=index.js.map