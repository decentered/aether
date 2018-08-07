// Store > Data Loaders

// These actions are the high-level loaders that correspond roughly to page contexts we have.

export { }
var fe = require('../services/feapiconsumer/feapiconsumer')

let dataLoaders = {

  loadBoardScopeData(context: any, boardfp: string) {
    context.dispatch('setCurrentBoardFp', boardfp)
  },

  loadThreadScopeData(context: any, { boardfp, threadfp }: { boardfp: string, threadfp: string }) {
    context.dispatch('setCurrentThreadFp', { boardfp: boardfp, threadfp: threadfp })
  },

  loadGlobalScopeData(context: any) {
    context.commit('SET_ALL_BOARDS_LOAD_COMPLETE', false)
    fe.GetAllBoards(function(resp: any) {
      console.log('received the all boards payload from fe')
      context.commit('SET_ALL_BOARDS', resp)
      context.dispatch('updateBreadcrumbs')
      context.commit('SET_ALL_BOARDS_LOAD_COMPLETE', true)
    })
  },

  loadUserScopeData(context: any, { fp, userreq, boardsreq, threadsreq, postsreq }: { fp: string, userreq: boolean, boardsreq: boolean, threadsreq: boolean, postsreq: boolean }) {
    context.commit('SET_CURRENT_USER_LOAD_COMPLETE', false)
    fe.GetUserAndGraph(fp, userreq, boardsreq, threadsreq, postsreq, function(resp: any) {
      console.log('Received user scope data')
      // We need to set in the query values we asked, so that the mutation will know what to override, and what not to.
      // resp.userentityrequested = userreq
      // resp.userboardsrequested = boardsreq
      // resp.userthreadsrequested = threadsreq
      // resp.userpostsrequested = postsreq
      context.commit('SET_USER_SCOPE_DATA', resp)
      context.commit('SET_CURRENT_USER_LOAD_COMPLETE', true)
      context.dispatch('updateBreadcrumbs')
    })
  }
}

export default dataLoaders