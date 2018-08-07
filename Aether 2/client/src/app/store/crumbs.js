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
function makeCurrentBoardCrumb(context) {
    return createCrumb('board', context.state.currentBoard.name, 'board/' + context.state.currentBoard.fingerprint, context.state.currentBoard.fingerprint);
}
function makeCurrentThreadCrumb(context) {
    return createCrumb('thread', context.state.currentThread.name, 'board/' + context.state.currentBoard.fingerprint + '/thread/' + context.state.currentThread.fingerprint, context.state.currentThread.fingerprint);
}
function makeCurrentUserCrumb(context) {
    return createCrumb('user', '@' +
        globalMethods.GetUserName(context.state.currentUserEntity), 'user/' + context.state.currentUserEntity.fingerprint, context.state.currentUserEntity.fingerprint);
}
var crumbActions = {
    updateBreadcrumbs: function (context) {
        console.log('update crumbs hits');
        var updatedCrumbs = [];
        console.log("context.state.route.name is:");
        console.log(context.state.route.name);
        if (context.state.route.name === 'Board') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
        }
        else if (context.state.route.name === 'Thread') {
            updatedCrumbs.push(makeCurrentBoardCrumb(context));
            updatedCrumbs.push(makeCurrentThreadCrumb(context));
        }
        else if (context.state.route.name === 'Global') {
            updatedCrumbs.push({
                EntityType: '',
                VisibleName: 'Communities',
                Link: 'globalscope',
                Fingerprint: ''
            });
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
        else if (context.state.route.name === 'User') {
            updatedCrumbs.push(makeCurrentUserCrumb(context));
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