"use strict";
// Store > Data Loaders
Object.defineProperty(exports, "__esModule", { value: true });
var fe = require('../services/feapiconsumer/feapiconsumer');
var dataLoaders = {
    loadBoardScopeData: function (context, boardfp) {
        context.dispatch('setCurrentBoardFp', boardfp);
    },
    loadThreadScopeData: function (context, _a) {
        var boardfp = _a.boardfp, threadfp = _a.threadfp;
        context.dispatch('setCurrentThreadFp', { boardfp: boardfp, threadfp: threadfp });
    },
    loadGlobalScopeData: function (context) {
        fe.GetAllBoards(function (resp) {
            console.log('received the all boards payload from fe');
            context.commit('SET_ALL_BOARDS', resp);
            context.dispatch('updateBreadcrumbs');
        });
    },
    loadUserScopeData: function (context, _a) {
        var fp = _a.fp, userreq = _a.userreq, boardsreq = _a.boardsreq, threadsreq = _a.threadsreq, postsreq = _a.postsreq;
        console.log('the fingerprint we are asking for');
        console.log(fp);
        fe.GetUserAndGraph(fp, userreq, boardsreq, threadsreq, postsreq, function (resp) {
            console.log('Received user scope data');
            // We need to set in the query values we asked, so that the mutation will know what to override, and what not to.
            resp.userentityrequested = userreq;
            resp.userboardsrequested = boardsreq;
            resp.userthreadsrequested = threadsreq;
            resp.userpostsrequested = postsreq;
            context.commit('SET_USER_SCOPE_DATA', resp);
            context.dispatch('updateBreadcrumbs');
        });
    }
};
exports.default = dataLoaders;
//# sourceMappingURL=dataloaders.js.map