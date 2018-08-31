"use strict";
// Store > Crumbs
Object.defineProperty(exports, "__esModule", { value: true });
var globalMethods = require('../services/globals/methods');
function createCrumb(entityType, visibleName, link, fingerprint) {
    var c = {
        EntityType: entityType,
        VisibleName: visibleName,
        Link: link,
        Fingerprint: fingerprint
    };
    return c;
}
// createEndCrumb always creates a 'cap' crumb for the current entity based on the current path
function createEndCrumb(context, entityType, visibleName, fingerprint) {
    var c = {
        EntityType: entityType,
        VisibleName: visibleName,
        Link: context.state.route.fullPath.substr(1),
        Fingerprint: fingerprint
    };
    return c;
}
function makeCurrentBoardCrumb(context) {
    if (context.state.currentBoard.fingerprint.length === 0) {
        // Object not found
        return createCrumb('board', 'Not found', 'board/' + context.state.currentBoard.fingerprint, context.state.currentBoard.fingerprint);
    }
    return createCrumb('board', context.state.currentBoard.name, 'board/' + context.state.currentBoard.fingerprint, context.state.currentBoard.fingerprint);
}
function makeCurrentThreadCrumb(context) {
    if (context.state.currentThread.fingerprint.length === 0) {
        // Object not found
        return createCrumb('thread', 'Not found', 'board/' + context.state.currentBoard.fingerprint + '/thread/' + context.state.currentThread.fingerprint, context.state.currentThread.fingerprint);
    }
    return createCrumb('thread', context.state.currentThread.name, 'board/' + context.state.currentBoard.fingerprint + '/thread/' + context.state.currentThread.fingerprint, context.state.currentThread.fingerprint);
}
function makeCurrentUserCrumb(context) {
    return createCrumb('user', globalMethods.GetUserName(context.state.currentUserEntity), 'user/' + context.state.currentUserEntity.fingerprint, context.state.currentUserEntity.fingerprint);
}
function makeGlobalCrumb() {
    return {
        EntityType: '',
        VisibleName: 'Communities',
        Link: 'globalscope',
        Fingerprint: ''
    };
}
var crumbActions = {
    updateBreadcrumbs: function (context) {
        // console.log('update crumbs hits')
        var updatedCrumbs = [];
        // console.log("context.state.route.name is:")
        // console.log(context.state.route.name)
        if (context.state.route.name === 'Board') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
        }
        else if (context.state.route.name === 'Board>ThreadsNewList') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'threadsNewList', 'New threads', ''));
        }
        else if (context.state.route.name === 'Board>NewThread') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'newThread', 'New thread', ''));
        }
        else if (context.state.route.name === 'Board>ModActivity') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'modActivity', 'Mod Activity', ''));
        }
        else if (context.state.route.name === 'Board>Elections') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'elections', 'Elections', ''));
        }
        else if (context.state.route.name === 'Board>Reports') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'reports', 'Reports', ''));
        }
        else if (context.state.route.name === 'Board>BoardInfo') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'info', 'Info', ''));
        }
        else if (context.state.route.name === 'Thread') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
            updatedCrumbs.push(makeCurrentThreadCrumb(context));
        }
        else if (context.state.route.name === 'Global') {
            updatedCrumbs.push(makeGlobalCrumb());
        }
        else if (context.state.route.name === 'Global>Subbed') {
            updatedCrumbs.push(makeGlobalCrumb());
            updatedCrumbs.push(createEndCrumb(context, 'subscribed', 'Subscribed', ''));
        }
        else if (context.state.route.name === 'Global>NewBoard') {
            updatedCrumbs.push(makeGlobalCrumb());
            updatedCrumbs.push(createEndCrumb(context, 'newBoard', 'New community', ''));
        }
        else if (context.state.route.name === 'Intro') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "A Beginner's Guide to the Galaxy",
                Link: 'intro',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'Settings') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Settings",
                Link: 'settings',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'Settings>Advanced') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Settings",
                Link: 'settings',
                Fingerprint: ''
            });
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Advanced",
                Link: 'settings/advanced',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'About') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "About",
                Link: 'about',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'Membership') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Membership",
                Link: 'membership',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'Changelog') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Changelog",
                Link: 'changelog',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'AdminsQuickstart') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Admin's Quickstart",
                Link: 'adminsquickstart',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'SFWList') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Safe-for-work list",
                Link: 'sfwlist',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'Modship') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Mod mode",
                Link: 'modship',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'NewUser') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: "Create New User",
                Link: 'newuser',
                Fingerprint: ''
            });
        }
        else if (context.state.route.name === 'User') {
            updatedCrumbs.push(makeCurrentUserCrumb(context));
        }
        else if (context.state.route.name === 'User>Boards') {
            updatedCrumbs.push(makeCurrentUserCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'boards', 'Communities', ''));
        }
        else if (context.state.route.name === 'User>Threads') {
            updatedCrumbs.push(makeCurrentUserCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'threads', 'Threads', ''));
        }
        else if (context.state.route.name === 'User>Posts') {
            updatedCrumbs.push(makeCurrentUserCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'posts', 'Posts', ''));
        }
        else if (context.state.route.name === 'User>Notifications') {
            updatedCrumbs.push(makeCurrentUserCrumb(context));
            updatedCrumbs.push(createEndCrumb(context, 'notifications', 'Notifications', ''));
        }
        else if (context.state.route.name === 'Status') {
            updatedCrumbs.push(createEndCrumb(context, 'status', 'Status', ''));
        }
        else if (context.state.route.name === 'Popular') {
            updatedCrumbs.push(createEndCrumb(context, 'popular', 'Popular', ''));
        }
        context.state.breadcrumbs = updatedCrumbs;
    },
    setBreadcrumbs: function (context, breadcrumbs) {
        context.commit('SET_BREADCRUMBS', breadcrumbs);
    },
};
var crumbMutations = {
    SET_BREADCRUMBS: function (state, breadcrumbs) {
        state.breadcrumbs = breadcrumbs;
    },
};
module.exports = {
    crumbActions: crumbActions,
    crumbMutations: crumbMutations
};
//# sourceMappingURL=crumbs.js.map