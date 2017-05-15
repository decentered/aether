function SubjectsLiteController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices, $timeout) {
    $scope.$watch('requestedId', function(){
        if ($rootScope.requestedId !== undefined) {
            gateReaderServices.getParentPost(parentIdArrived, $rootScope.requestedId)
        }
    })

    function parentIdArrived(data) {

        $rootScope.thereIsASelectedSubject = true
        $rootScope.parentId = data.PostFingerprint
        $scope.selectedTopic = data

        $scope.isCurrentlySubscribed = false
        // Check if this topic is a topic that is currently subscribed
        for (var i=0;i<$rootScope.userProfile.userDetails.selectedTopics.length;i++) {
            (function(){
                if ($rootScope.userProfile.userDetails.selectedTopics[i] ===
                    $scope.selectedTopic.PostFingerprint) {
                    $scope.isCurrentlySubscribed = true
                }
            })()
        }


        if ($rootScope.parentId) {
            gateReaderServices.getAllSpecificDepthPosts(subjectsArrived, $rootScope.parentId, 1)
        }

        function subjectsArrived(data) {
            // Remove the first element, because it is the
            // subject itself.
            data.splice(0,1)
            //data = data.sort($rootScope.sortByScore)
            $scope.subjects = data
            $timeout(function(){
                document.getElementById('post-'+$rootScope.requestedId).scrollIntoViewIfNeeded()
            },10)
        }


    }

    $scope.toggleSubscriptionState = function() {
        console.log('scope selectedtopic postfingerprint:', $scope.selectedTopic.PostFingerprint)
        if($rootScope.userProfile.userDetails.selectedTopics
            .indexOf($scope.selectedTopic.PostFingerprint) > -1 ) {
            // If it exists, remove
            var index = $rootScope.userProfile.userDetails.selectedTopics.indexOf($scope.selectedTopic.PostFingerprint)
            $rootScope.userProfile.userDetails.selectedTopics.splice(index, 1)
            $scope.isCurrentlySubscribed = false
        }
        else
        {
            $rootScope.userProfile.userDetails.selectedTopics.push($scope.selectedTopic.PostFingerprint)
            $scope.isCurrentlySubscribed = true
        }
    }

    $scope.returnToTopics = function() {
        frameViewStateBroadcast.receiveState('', 'topicsFeedLite', '')
    }

    // Learned hard way: if you capture an event on a higher scope, it won't go
    // deeper. this means any $scope.on(..) that has already been captured at
    // frame controllers won't even get triggered here.

    // TODO: Here, check if the controller can access 2ndframeid. If it can, it's
    // at 3rd frame. In that case, write a ng-class directive that can look for
    // item id and activate on match so the one lights up.

    // I also need to determine here if this is in 2nd or 3rd scope, or write a
    // different controller for lite version.
}
SubjectsLiteController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices', '$timeout']