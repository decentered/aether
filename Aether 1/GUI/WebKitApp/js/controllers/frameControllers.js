// First, second and third frame controllers handle the resize events.
// They also handle global events that result in change of pages, buttons etc.

function FirstFrameController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices) {

    $scope.$on('initialised', function() {
        document.getElementById('cloak').style.display = 'none';
    })

    // Check if there are any topics and their counts available. This is the global app state check.
    $rootScope.countsAvailable = false
    gateReaderServices.getUppermostTopics(function(data) {
        if (data.length) {
            // if there are topics, but not their counts.
            for (var i=0;i<data.length;i++) {
                if (data[i].ReplyCount) {
                    $rootScope.countsAvailable = true
                    return
                }
            }

        }
        $rootScope.countsAvailable = false
    })


    // Methods below are available to all scopes. They're inside the first frame
    // controller, because it seems I cannot place rootScope
    // objects outside one! They're still available everywhere though.

    // Yes, I know rootScope is evil. Let me know if you have a better solution.

    $rootScope.log = function(variable) {
        console.log(variable)
    }
    $rootScope.alert = function(text) {
        alert(text)
    }



    // Make markdown available globally
    $rootScope.mdConverter = Markdown.getSanitizingConverter()

    // The cache objects. These objects store the transient values for
    // thread and post creation. All the pages that open a creation interface
    // use this as the workspace, albeit after linking they can continue
    // to use their own names.

    $rootScope.postTextCache = ''
    $rootScope.subjectHeaderCache = ''


    // This entire thing below: I'm starting to get a distinct feeling that
    // I am poorly recreating ui-router. There's gotta be a better way of
    // doing this. Research if I can move to ui-router and take all that
    // cruft off my plate.

    // These two methods are available globally: they allow basic frame change
    // methods to work. Frame change changes the partial views, and frame size
    // change changes the sizes of the second and third frames.

    // IMPORTANT: If you pass "" as the 3rd value of changeState, it will RETAIN
    // the value of requestedId in memory so you can still access it from the
    // new partial. If you pass undefined as the 3rd value, it will FLUSH the
    // existing value from rootScope, so no one has access to the value anymore.

    // Ergo, be careful in passing undefined, as in one page there are probably
    // more than one partial making use of that value. The only use case I can
    // see for flushing the value so far is the entry, where the value is
    // empty from the start, but if left as "", it is mistaken as "retain"
    // rather than explicit empty.

    // If you want to use undefined, be careful, most controllers actively
    // filter out undefined to keep the values within safe from outside
    // interference.

    $rootScope.changeState = function(secondFrame, thirdFrame, id) {
        frameViewStateBroadcast.receiveState(secondFrame, thirdFrame, id)
    }

    // This is a general listener on view state change that makes requestedId
    // variable available on the scope on all possible controllers.

    $rootScope.$on('frameViewStateChanged', function() {
        if (frameViewStateBroadcast.id !== "") {
            $rootScope.requestedId = frameViewStateBroadcast.id
        }
    })

    $rootScope.$watch('viewportHeight', function() {
        $rootScope.viewportStyle = {
            'height': $scope.viewportHeight + 'px'
        }
    })

    $rootScope.clickToUsername = function(userName) {
        $rootScope.requestedUsername = userName
        $rootScope.thirdFrameCSSStyle = {}
        $scope.changeState('unregisteredUserFeed', '', '')
    }

    // These three things are how you manipulate the frames in DOM rather
    // than jQuery. Strictly speaking I don't need this, but it's mostly here
    // as a reminder. (Because evaluation on DOM is graceful to undefined)

    // These are default states at the beginning of the application.

    $rootScope.firstFrameCSSStyle = {
    }

    $rootScope.secondFrameCSSStyle = {
        'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
        (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2)) + 'px'
    }

    $rootScope.thirdFrameCSSStyle = {}

    // Below are the first call to userProfile object to make it globally
    // available at application start. The watch below will fire at every change.

    // The watch might need an object equality check rather than default check
    // in the case this does not update, add true as the last parameter of the
    // watch.

    // The fact that I could see this and get this to work before I made a million
    // save requests everywhere points to that being a rare moment of brilliance.

    gateReaderServices.getUserProfile(
        function(userProfile) {
            $rootScope.userProfile = userProfile
        })

    // If any part of the userProfile doesn't exist on first load, create it here.
    // I moved JSON creation to the backend.

//    if ($rootScope.userProfile.userDetails.selectedTopics === undefined) {
//        $rootScope.userProfile.userDetails.selectedTopics = []
//    }
//
//    if ($rootScope.userProfile.userDetails === undefined) {
//        $rootScope.userProfile.userDetails = {
//            Username: '',
//            userLanguages: ['English'],
//            StartAtBoot: true,
//            maxInboundCount: 3,
//            maxOutboundCount: 10,
//            cooldown: 5
//        }
//    }
//
//    if ($rootScope.userProfile.userDetails.unreadReplies === undefined) {
//        $rootScope.userProfile.userDetails.unreadReplies = []
//    }
//
//    if ($rootScope.userProfile.userDetails.readReplies === undefined) {
//        $rootScope.userProfile.userDetails.readReplies = []
//    }

    // Set theme

    if ($rootScope.userProfile.userDetails.theme === undefined)
    {
        $rootScope.userProfile.userDetails.theme = false
    }

    $rootScope.setTheme = function(value) { // 0: Light, 1: Dark
        var darkSheet = document.getElementsByTagName('link')[2]
        var lightSheet = document.getElementsByTagName('link')[1]
        if (value === 1) {
            // We went dark
            darkSheet.disabled = false
            lightSheet.disabled = true
        }
        else if (value === 0) {
            // Let there be light
            darkSheet.disabled = true
            lightSheet.disabled = false
        }
    }

    if ($rootScope.userProfile.userDetails.theme === false) {
        // Light.
        $rootScope.setTheme(0)
    }
    else if($rootScope.userProfile.userDetails.theme === true) {
        // Dark. change.
        $rootScope.setTheme(1)
    }

    $rootScope.appIsPaused = false

    $scope.$watch('userProfile', function() {
        // Some guards about what can get into the profile. This is the place to protect and sanitise input from user.
        $scope.userProfile.userDetails.username = $scope.userProfile.userDetails.username.trim()
        var maxInbound = parseInt($scope.userProfile.machineDetails.maxInboundCount, 10) // base 10
        var maxOutbound = parseInt($scope.userProfile.machineDetails.maxOutboundCount, 10)
        var cooldown = parseInt($scope.userProfile.machineDetails.cooldown, 10)
        if (maxInbound < 1 || isNaN(maxInbound)) {
            $scope.userProfile.machineDetails.maxInboundCount = 3
        }
        else
        {
            $scope.userProfile.machineDetails.maxInboundCount = maxInbound
        }

        if (maxOutbound < 1 || isNaN(maxOutbound)) {
            $scope.userProfile.machineDetails.maxOutboundCount = 10
        }
        else
        {
            $scope.userProfile.machineDetails.maxOutboundCount = maxOutbound
        }

        if (cooldown < 1 || isNaN(cooldown)) {
            $scope.userProfile.machineDetails.cooldown = 5
        }
        else
        {
            $scope.userProfile.machineDetails.cooldown = cooldown
        }
        gateReaderServices.saveUserProfile(function() {}, $rootScope.userProfile)
        //console.log('autocommit fired, new profile: ',$rootScope.userProfile)
    }, true)


    // This is a helper function to check if the pane queried is currently
    // active.

    $rootScope.isActive = function(paneName) {
        return frameViewStateBroadcast.secondFrame == paneName ||
        frameViewStateBroadcast.thirdFrame == paneName ?
        true : false
    }

    $rootScope.reloadSecondFrameScrollPosition = function() {
        var scrollLocation = $rootScope.secondFrameSelectedTemplate.lastScrollLocation
            document.getElementById('second-frame').scrollTop = scrollLocation
    }

    $rootScope.scrollThirdFrameToTop = function() {
        document.getElementById('scroll-target-third-frame').scrollIntoView()
    }

//    $rootScope.sortByTime = function(a,b) {
//
//    }

    // This algorithm below is for subjects. Posts have their own much simpler idea, they just get sorted by the upvotes.
//    $rootScope.sortByScore = function(a, b) {
//        var now = Math.round(new Date().getTime() / 1000 / 60 / 60) // unix timestamp in hour
//        var creationDateA = Math.round(a.CreationDate / 60 / 60) // because this already comes in secs.
//        var timedeltaA = now - creationDateA
//        var scoreA = (a.UpvoteCount - a.DownvoteCount - 1) / Math.pow((timedeltaA + 2), 1.5)
//        var creationDateB = Math.round(b.CreationDate / 60 / 60) // because this already comes in secs.
//        var timedeltaB = now - creationDateB
//        var scoreB = (b.UpvoteCount - b.DownvoteCount - 1) / Math.pow((timedeltaB + 2), 1.5)
//        return scoreB - scoreA
            // is this useless? I am doing this sorting on the backend.
//    }

    gateReaderServices.countReplies(countArrived)
    function countArrived(count) {
        $rootScope.totalReplyCount = count
    }

    gateReaderServices.getOperatingSystem(function(osString){
        $rootScope.PLATFORM = osString
    })


    gateReaderServices.getAppVersion(function(reply) {
        $rootScope.appVersion = reply.slice(0,1) + '.' + reply.slice(1,2) + '.' + reply.slice(2,3)
    })

    $scope.homeButtonClick = function() {
        $rootScope.changeState('homeFeed', 'topicGroupsFeedLite', '')
        $rootScope.secondFrameCSSStyle = {
            'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
            (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2)) + 'px'
        }
        $rootScope.thirdFrameCSSStyle = {}
    }

    $scope.threadsButtonClick = function() {
        $rootScope.changeState('subjectsFeed', 'topicsFeedLite', $rootScope.parentId)
        $rootScope.thereIsASelectedSubject = false

        $rootScope.secondFrameCSSStyle = {}
        $rootScope.thirdFrameCSSStyle = {
            'display':'block'
        }
    }

    $scope.topicsButtonClick = function() {
        $rootScope.changeState('findOrCreateTopic', 'topicsFeedLite', undefined)
        $rootScope.thereIsASelectedSubject = false

        $rootScope.secondFrameCSSStyle = {}
        $rootScope.thirdFrameCSSStyle = {
            //'display':'block'
        }
    }

    $scope.savedItemsButtonClick = function() {
        $rootScope.changeState('savedItemsFeed', 'savedItemsFeedLite', '')
        $rootScope.secondFrameCSSStyle = {
            'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
            (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2)) + 'px'
        }
        $rootScope.thirdFrameCSSStyle = {}
    }

    $scope.settingsButtonClick = function() {
        $rootScope.changeState('settingsFeed', 'settingsSelectorFeedLite', undefined)
        $rootScope.thirdFrameCSSStyle = {
            'display':'block',
            'width': '150px'
        }
        $rootScope.secondFrameCSSStyle = {
            'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
            (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2))
            - (angular.element(document.getElementById('third-frame')).css('width').slice(0, -2)) + 'px'
        }

    }

    $scope.repliesButtonClick = function() {
        console.log('repliesbuttonclicked')
        $rootScope.changeState('repliesFeed', 'topicGroupsFeedLite', undefined)
        $rootScope.thirdFrameCSSStyle = {
        }
        $rootScope.secondFrameCSSStyle = {
            'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
            (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2))
            - (angular.element(document.getElementById('third-frame')).css('width').slice(0, -2)) + 'px'
        }

    }

    $scope.profileButtonClick = function() {
        $rootScope.changeState('unregisteredUserProfile', 'topicGroupsFeedLite', undefined)
        $rootScope.thirdFrameCSSStyle = {
        }
        $rootScope.secondFrameCSSStyle = {
            'width': (angular.element(document.getElementsByTagName('body')).css('width').slice(0, -2)) -
            (angular.element(document.getElementById('first-frame')).css('width').slice(0, -2))
            - (angular.element(document.getElementById('third-frame')).css('width').slice(0, -2)) + 'px'
        }

    }


}

FirstFrameController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
'gateReaderServices']



function SecondFrameController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices, $anchorScroll) {

    // This is the list of partials this scope can load. This is also one of
    // the key reasons the first / second / third frame controllers cannot be
    // removed.

    // REGIONS describe the states application is in. Region change is
    // required for first frame buttons' awareness, but could be useful in other
    // places.

    $rootScope.secondFrameTemplates = [{
        name: 'homeFeed',
        url: 'partials/homeFeed.html',
        region: 'Home',
        lastScrollLocation:0
    },{
        name:'onboarding',
        url: 'partials/onboarding.html',
        region: 'NoRegion',
        lastScrollLocation:0
    },{
        name: 'postsFeed',
        url: 'partials/postsFeed.html',
        region: 'Details',
        lastScrollLocation:0 // This is always zero. I do not retain scroll location on this.
    },{
        name: 'subjectsFeed',
        url: 'partials/subjectsFeed.html',
        region: 'Details',
        lastScrollLocation:0
    },{
        name: 'savedItemsFeed',
        url: 'partials/savedItemsFeed.html',
        region: 'SavedItems',
        lastScrollLocation:0
    },{
        name: 'createSubject',
        url: 'partials/createSubject.html',
        region: 'Details',
        lastScrollLocation:0
    },{
        name: 'findOrCreateTopic',
        url: 'partials/findOrCreateTopic.html',
        region: 'Topics',
        lastScrollLocation:0
    },{
        name: 'settingsFeed',
        url: 'partials/settingsBase.html',
        region: 'Settings',
        lastScrollLocation:0
    },{
        name: 'singlePost',
        url: 'partials/singlePost.html',
        region: 'NoRegion',
        lastScrollLocation:0
    },{
        name: 'repliesFeed',
        url: 'partials/repliesFeed.html',
        region: 'Replies',
        lastScrollLocation:0
    },{
        name: 'singleReply',
        url: 'partials/singlePost.html',
        region: 'NoRegion',
        lastScrollLocation:0
    },{
        name: 'unregisteredUserFeed',
        url: 'partials/unregisteredUserItemsFeed.html',
        region: 'NoRegion', // This is for now..
        lastScrollLocation:0
    },{
        name: 'unregisteredUserProfile',
        url: 'partials/unregisteredUserProfile.html',
        region: 'Profile',
        lastScrollLocation:0
    }]

    gateReaderServices.getOnboardingComplete(onboardingCompleteArrived)
    function onboardingCompleteArrived(onboardingComplete) {
        $rootScope.onboardingComplete = onboardingComplete // normal state = no !
        console.log('onboardingComplete status:', onboardingComplete)
        if (onboardingComplete === true) { //normal state = true
            // If not newborn
            $rootScope.secondFrameSelectedTemplate = $rootScope.secondFrameTemplates[0] // 0 is home
            $rootScope.currentApplicationRegion = 'Home'
        }
        else
        {
            $rootScope.secondFrameSelectedTemplate = $rootScope.secondFrameTemplates[1] //1 is onboarding
            $rootScope.currentApplicationRegion = 'NoRegion'
        }
    }


    // This part below fires on frame state change and greps for the incoming
    // frame name in templates. If I built this logic, this is smart, but I have
    // no recollection of writing it. Oh god.

    $scope.$on('frameViewStateChanged', function() {
        // Save the scroll state.
        if ($rootScope.secondFrameSelectedTemplate.name != 'postsFeed') {
            // Because it does not make sense to retain the scroll location from one view to another.
            $rootScope.secondFrameSelectedTemplate.lastScrollLocation =
            -document.getElementById('second-frame-contents').getBoundingClientRect().top
        }
        var searchResult
        for (var i = 0; i<$rootScope.secondFrameTemplates.length; i++) {
            if ($rootScope.secondFrameTemplates[i].name === frameViewStateBroadcast.secondFrame) {
                searchResult = $rootScope.secondFrameTemplates[i]
            }
        }
        // Change only if 2nd variable exists i.e. it isn't ""
        if (frameViewStateBroadcast.secondFrame !== "") {
            $rootScope.secondFrameSelectedTemplate = searchResult
            $rootScope.currentApplicationRegion = searchResult.region

        }
        $rootScope.reloadSecondFrameScrollPosition()
    })
}

SecondFrameController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
'gateReaderServices', '$anchorScroll']

function ThirdFrameController($scope, $rootScope, frameViewStateBroadcast) {
    $scope.templates = [{
        name: 'subjectsFeedLite',
        url: 'partials/subjectsFeedLite.html'
    }, {
        name: 'topicsFeedLite',
        url: 'partials/topicsFeedLite.html'
    }, {
        name: 'createLite',
        url: 'partials/createLite.html'
    }, {
        name: 'settingsSelectorFeedLite',
        url: 'partials/settingsSelectorFeedLite.html'
    }  ]

    $scope.selectedTemplate = $scope.templates[0]

    $scope.$on('frameViewStateChanged', function() {
        var searchResult
        for (var i = 0; i<$scope.templates.length; i++) {
            if ($scope.templates[i].name === frameViewStateBroadcast.thirdFrame) {
                searchResult = $scope.templates[i]
            }
        }
        // Change only if 2nd variable exists i.e. it isn't ""
        if (frameViewStateBroadcast.thirdFrame !== "") {
            $scope.selectedTemplate = searchResult
        }
    })
}

ThirdFrameController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast']