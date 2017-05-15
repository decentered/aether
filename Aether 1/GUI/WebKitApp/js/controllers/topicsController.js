function TopicsController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices, $timeout) {



    $scope.topics = []
    $scope.setSubjectsToSingleColumn = function() {
        $rootScope.userProfile.userDetails.subjectsSingleColumnLayout = true

    }
    $scope.setSubjectsToMultiColumn = function() {
        $rootScope.userProfile.userDetails.subjectsSingleColumnLayout = false
    }


    $scope.$on('refreshPage', function() {
        gateReaderServices.getUppermostTopics(topicsArrived)
    })

    gateReaderServices.getNewborn(function(answer) {$scope.newborn = answer})

    gateReaderServices.getUppermostTopics(topicsArrived)
    function topicsArrived(posts)
    {
        // Before anything else: remove all 'null's from the list. This is weird, IDK how this can happen..
        // oh I know, if you subscribe to a NULL board.
        if($rootScope.userProfile.userDetails.selectedTopics
                .indexOf(null) > -1 ) {
            console.log('there are nulls in the selected topics. Oops.')
            for (var i=0;i<$rootScope.userProfile.userDetails.selectedTopics.length;i++) {
                (function(i){
                    if ($rootScope.userProfile.userDetails.selectedTopics[i] === null) {
                        $rootScope.userProfile.userDetails.selectedTopics.splice(i,1)
                    }
                })(i)
            }
        }
        // First: Sort topics alphabetically before posting them to DOM.
        posts.sort(function(a,b) {
            var aTopic = a.Subject.toUpperCase()
            var bTopic = b.Subject.toUpperCase()
            return (aTopic<bTopic) ? -1 : (aTopic>bTopic) ? 1 : 0
        })

        // Then: Sort topics by popularity before posting them to the DOM.
        posts.sort(compareByReplyCount)
        function compareByReplyCount(b, a) {
                 return a.ReplyCount - b.ReplyCount
            }


        for (var i=0; i<posts.length; i++) {
            // check if posts[i] exists in userProfile.selectedTopics.
            if($rootScope.userProfile.userDetails.selectedTopics
                .indexOf(posts[i].PostFingerprint) > -1 ) {
                // If exists in userProfile
                $scope.topics.push(posts[i])
            }
        }

        //console.log('these are topics: ', $scope.topics)

        $scope.$watch('requestedId + topics', function(){
            if ($rootScope.requestedId === undefined && $scope.topics[0] != undefined)
            {
                $rootScope.requestedId = $scope.topics[0].PostFingerprint
            }
        })

    }

    $scope.topicClick = function(PostFingerprint) {
        $rootScope.thereIsASelectedSubject = false
        $scope.changeState('subjectsFeed', '', PostFingerprint)
    }

    $scope.returnToSubjects = function () {
        frameViewStateBroadcast.receiveState('', 'subjectsFeedLite', '')
    }

    $scope.toggleFindOrCreateTopic = function() {
        frameViewStateBroadcast.secondFrame == 'findOrCreateTopic' ?
        $scope.changeState('subjectsFeed', '', '') :
        $scope.changeState('findOrCreateTopic', '', '')
        $rootScope.thereIsASelectedSubject = false
    }

    $scope.isSelected = function(fingerprint) {
        if ($rootScope.thereIsASelectedSubject && $rootScope.parentId === fingerprint) {
            $timeout(function(){
                document.getElementById('post-'+$rootScope.parentId).scrollIntoViewIfNeeded()
            },100)
            return true
        }
        else if (fingerprint == $scope.requestedId &&
        frameViewStateBroadcast.secondFrame != 'findOrCreateTopic') {
            return true
        }
        else {
            return false
        }
    }




}
TopicsController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices', '$timeout']