/*
This is the data store for the client.

This data store does not hold any persistent data, nor does it cache it. The point of this is to hold the instance data. The frontend is the actual caching and compile logic that regenerates the data to be used as needed.
*/

var Vue = require('../../../node_modules/vue/dist/vue.js')
var Vuex = require('../../../node_modules/vuex').default
Vue.use(Vuex)

var fe = require('../services/feapiconsumer/feapiconsumer')

let dataLoaders = require('./dataloaders').default
let contentRelations = require('./contentrelations')
let crumbs = require('./crumbs')



const dataLoaderPlugin = function(store: any) {
  store.watch(
    // When the returned result changes,
    function(state: any) {
      return state.route.params
    },
    // Run this callback
    function(newValue: any, oldValue: any) {
      fe.SendInflightsPruneRequest(function(resp: any) {
        console.log(resp)
      })
      // First, check if we should refresh.
      if (oldValue === newValue && !store.state.frontendHasUpdates) {
        // if the values are the same, and frontend has no updates, bail.
        return
      }
      let routeParams: any = newValue
      if (store.state.route.name === "Board" || store.state.route.name === "Board>ThreadsNewList" || store.state.route.name === "Board>ModActivity" || store.state.route.name === "Board>Elections") {
        store.dispatch('loadBoardScopeData', routeParams.boardfp)
        store.dispatch('setLastSeenForBoard', { fp: routeParams.boardfp })
        return
      }

      if (store.state.route.name === "Board>NewThread") {
        store.dispatch('loadBoardScopeData', routeParams.boardfp)
        return
      }
      if (store.state.route.name === "Board>BoardInfo") {
        store.dispatch('loadBoardScopeData', routeParams.boardfp)
        return
      }

      if (store.state.route.name === "Thread") {
        store.dispatch('loadThreadScopeData', {
          boardfp: routeParams.boardfp,
          threadfp: routeParams.threadfp,
        })
        return
      }

      if (store.state.route.name === "Global" || store.state.route.name === "Global>Subbed") {
        store.dispatch('loadGlobalScopeData')
        return
      }

      if (store.state.route.name === "User" || store.state.route.name === 'User>Boards' || store.state.route.name === 'User>Threads' || store.state.route.name === 'User>Posts') {
        store.dispatch('loadUserScopeData', {
          fp: routeParams.userfp,
          userreq: true,
          boardsreq: false,
          threadsreq: false,
          postsreq: false,
        })
        return
      }

      // If none of the special cases, just trigger an update for breadcrumbs.
      store.dispatch('updateBreadcrumbs')
    }
  )
}

let actions = {
  /*----------  Refreshers  ----------*/
  /*
    These are smaller, less encompassing versions of the loaders, and they're meant to be used after the principal payload is brought in.
  */
  refreshCurrentBoardAndThreads(context: any, boardfp: string) {
    fe.GetBoardAndThreads(boardfp, function(resp: any) {
      actions.pruneInflights()
      context.commit('SET_CURRENT_BOARD', resp.board)
      context.commit('SET_CURRENT_BOARDS_THREADS', resp.threadsList)
      context.commit('SET_CURRENT_BOARD_LOAD_COMPLETE', true)
    })
  },
  pruneInflights() {
    fe.SendInflightsPruneRequest(function() { })
  },
  refreshCurrentThreadAndPosts(context: any, { boardfp, threadfp }: { boardfp: string, threadfp: string }) {
    fe.GetThreadAndPosts(boardfp, threadfp, function(resp: any) {
      actions.pruneInflights()
      context.commit('SET_CURRENT_BOARD', resp.board)
      context.commit('SET_CURRENT_THREAD', resp.thread)
      context.commit('SET_CURRENT_THREADS_POSTS', resp.postsList)
      // context.commit('SET_CURRENT_THREAD_LOAD_COMPLETE', true)
    })
  },
  /*----------  Refreshers END  ----------*/

  // within any of those, context.state is how you access state above.
  setSidebarState(context: any, sidebarOpen: boolean) {
    context.commit('SET_SIDEBAR_STATE', sidebarOpen)
  },
  setAmbientBoards(context: any, ambientBoards: any) {
    context.commit('SET_AMBIENT_BOARDS', ambientBoards)
  },
  setAmbientStatus(context: any, ambientStatus: any) {
    context.commit('SET_AMBIENT_STATUS', ambientStatus)
  },
  setAmbientLocalUserEntity(context: any, ambientLocalUserEntityPayload: any) {
    context.commit('SET_AMBIENT_LOCAL_USER_ENTITY', ambientLocalUserEntityPayload)
  },
  setCurrentBoardFp(context: any, fp: string) {
    // This is a boardscope load. Retrieve board + its threads.
    // if (context.state.currentBoardFp !== fp) {
    //   context.dispatch('setCurrentBoardAndThreads', fp)
    //   context.commit('SET_CURRENT_BOARD_FP', fp)
    //   // Same scope, but fe has updated.
    //   return
    // }
    context.dispatch('setCurrentBoardAndThreads', fp)
    context.commit('SET_CURRENT_BOARD_FP', fp)
    // // current board fp is the same as what we asked for, but FE has updates.
    // if (context.state.frontendHasUpdates) {
    //   context.dispatch('setCurrentBoardAndThreads', fp)
    //   context.commit('SET_CURRENT_BOARD_FP', fp)
    // }
  },
  setCurrentThreadFp(context: any, { boardfp, threadfp }: { boardfp: string, threadfp: string }) {
    context.dispatch('setCurrentThreadAndPosts',
      { boardfp: boardfp, threadfp: threadfp })
    context.commit('SET_CURRENT_THREAD_FP', threadfp)
    // Same scope, but fe has updated.
    return
    // current thread fp is the same as what we asked for, but FE has updates.
    // if (context.state.frontendHasUpdates) {
    //   context.dispatch('setCurrentThreadAndPosts',
    //     { boardfp: boardfp, threadfp: threadfp })
    //   context.commit('SET_CURRENT_THREAD_FP', threadfp)
    // }
  },
  setCurrentBoardAndThreads(context: any, boardfp: string) {
    if (context.state.currentBoardFp === boardfp) {
      fe.GetBoardAndThreads(boardfp, function(resp: any) {
        context.commit('SET_CURRENT_BOARD', resp.board)
        context.commit('SET_CURRENT_BOARDS_THREADS', resp.threadsList)
        context.commit('SET_CURRENT_BOARD_LOAD_COMPLETE', true)
        context.dispatch('updateBreadcrumbs')
      })
      return
      // If we're already here, update board but without false/true current board load complete flash.
    }
    context.commit('SET_CURRENT_BOARD_LOAD_COMPLETE', false)
    fe.GetBoardAndThreads(boardfp, function(resp: any) {
      context.commit('SET_CURRENT_BOARD', resp.board)
      context.commit('SET_CURRENT_BOARDS_THREADS', resp.threadsList)
      context.commit('SET_CURRENT_BOARD_LOAD_COMPLETE', true)
      context.dispatch('updateBreadcrumbs')
      // ^ This has to be here because otherwise the BC compute process runs before the data is ready, resulting in empty breadcrumbs.
    })
  },
  setCurrentThreadAndPosts(context: any, { boardfp, threadfp }: { boardfp: string, threadfp: string }) {
    if (context.state.currentThreadFp === threadfp) {
      context.dispatch('updateBreadcrumbs')
      return
    }
    context.commit('SET_CURRENT_THREAD_LOAD_COMPLETE', false)
    fe.GetThreadAndPosts(boardfp, threadfp, function(resp: any) {
      context.commit('SET_CURRENT_BOARD', resp.board)
      context.commit('SET_CURRENT_THREAD', resp.thread)
      context.commit('SET_CURRENT_THREADS_POSTS', resp.postsList)
      context.commit('SET_CURRENT_THREAD_LOAD_COMPLETE', true)
      context.dispatch('updateBreadcrumbs')
    })
  },
  ...dataLoaders,
  ...crumbs.crumbActions,
  ...contentRelations.actions,
}

let mutations = {
  SET_SIDEBAR_STATE(state: any, sidebarOpen: boolean) {
    state.sidebarOpen = sidebarOpen
  },

  SET_AMBIENT_BOARDS(state: any, ambientBoards: any) {
    state.ambientBoards = ambientBoards
  },
  SET_AMBIENT_STATUS(state: any, ambientStatus: any) {
    state.ambientStatus = ambientStatus
  },
  SET_AMBIENT_LOCAL_USER_ENTITY(state: any, payload: any) {
    state.localUserArrived = true
    // ^ Always set to true since some message arrived.
    state.localUserExists = payload.localuserexists
    state.localUser = payload.localuserentity
  },
  SET_CURRENT_BOARD_FP(state: any, fp: string) {
    state.currentBoardFp = fp
  },
  SET_CURRENT_THREAD_FP(state: any, fp: string) {
    state.currentThreadFp = fp
  },
  SET_CURRENT_BOARD(state: any, board: any) {
    state.currentBoard = board
  },
  SET_CURRENT_THREAD(state: any, thread: any) {
    state.currentThread = thread
  },
  SET_CURRENT_BOARDS_THREADS(state: any, threads: any) {
    state.currentBoardsThreads = threads
  },
  SET_CURRENT_THREADS_POSTS(state: any, posts: any) {
    state.currentThreadsPosts = posts
  },
  SET_ALL_BOARDS(state: any, boards: any) {
    state.allBoards = boards
  },
  SET_USER_SCOPE_DATA(state: any, resp: any) {
    if (resp.userentityrequested) {
      state.currentUserEntity = resp.user
    }
    if (resp.boardsrequested) {
      state.currentUserBoards = resp.boards
    }
    if (resp.threadsrequested) {
      state.currentUserThreads = resp.threads
    }
    if (resp.postsrequested) {
      state.currentUserPosts = resp.posts
    }
  },
  ...crumbs.crumbMutations,
  ...contentRelations.mutations,
  /*----------  Loader mutations that mark a pull done  ----------*/
  /*
    These are important because when these are complete and there is no data, we know that we should show a 404. These only apply to singular entities, not lists, so effectively thread view, board view, user view.
  */
  SET_ALL_BOARDS_LOAD_COMPLETE(state: any, loadComplete: boolean) {
    state.allBoardsLoadComplete = loadComplete
  },
  SET_CURRENT_BOARD_LOAD_COMPLETE(state: any, loadComplete: boolean) {
    state.currentBoardLoadComplete = loadComplete
  },
  SET_CURRENT_THREAD_LOAD_COMPLETE(state: any, loadComplete: boolean) {
    state.currentThreadLoadComplete = loadComplete
  },
  SET_CURRENT_USER_LOAD_COMPLETE(state: any, loadComplete: boolean) {
    state.currentUserLoadComplete = loadComplete
  },
}

let st = new Vuex.Store({
  state: {
    /*----------  All boards main  ----------*/
    allBoards: [],
    allBoardsLoadComplete: false,

    /*----------  Current board main  ----------*/
    currentBoard: {},
    currentBoardFp: "",
    currentBoardLoadComplete: false,
    /*----------  Current board sub data  ----------*/
    currentBoardsThreads: [],

    /*----------  Current thread main  ----------*/
    currentThread: {}, // todo - insert 404 here
    currentThreadFp: "",
    currentThreadLoadComplete: false,
    /*----------  Current thread sub data  ----------*/
    currentThreadsPosts: [],
    currentUserBoards: [],

    /*----------  Current user main  ----------*/
    currentUserEntity: {},  // This is the last user entity loaded into the user scope, *not* the current user occupying the client.
    currentUserLoadComplete: false,
    /*----------  Current user sub data  ----------*/
    currentUserPosts: [],
    currentUserThreads: [],

    /*----------  Ambient data pushed in from frontend  ----------*/
    ambientBoards: {},
    ambientStatus: {
      backendambientstatus: {},
      frontendambientstatus: {},
      inflights: {
        boardsList: [],
        threadsList: [],
        postsList: [],
        votesList: [],
        keysList: [],
        truststatesList: []
      },
    },

    /*----------  Local user data  ----------*/

    localUser: {},
    localUserExists: false,
    localUserArrived: false,
    // ^ Did we ever get a payload from FE? Until this is true, you can hide unready parts.

    /*----------  Misc  ----------*/
    frontendHasUpdates: true,
    frontendPort: 0,
    route: {},
    sidebarOpen: true,
    breadcrumbs: [],
  },
  actions: actions,
  mutations: mutations,
  plugins: [dataLoaderPlugin],
})

export default st
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

