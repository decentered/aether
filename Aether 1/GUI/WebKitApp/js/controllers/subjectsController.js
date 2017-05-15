function SubjectsController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {

    $scope.$watch('requestedId + timespan', function() {
        if ($rootScope.requestedId === undefined) {
            // Do nothing. Make the entire thing invisible, so it doesn't look ugly.
            document.getElementsByClassName('subject-global-actions-box')[0].style.opacity = 0
            document.getElementById('aether-footer-brand').style.opacity = 0
        }
        else
        {
            document.getElementsByClassName('subject-global-actions-box')[0].style.opacity = 1
            document.getElementById('aether-footer-brand').style.opacity = 1
            if ($scope.timespan === undefined || $scope.timespan === 0) {
                console.log("timespan undefined or zero")
                gateReaderServices.getSubjects(subjectsArrived,
                    $rootScope.requestedId, 180) // six months.

                $scope.timespanIsFiltered = false
            }
            else {
                console.log("timespan NOT undefined or zero")
                gateReaderServices.getSubjects(subjectsArrived,
                    $rootScope.requestedId, $scope.timespan)

                $scope.timespanIsFiltered = true
            }

            function subjectsArrived(data) {

                // here, sort data according to score.

                //data = data.sort($rootScope.sortByScore)

                $scope.singleCol = data

                var column1 = []
                var column2 = []
                var column3 = []
                for (var i=0;i<data.length;i++) {
                    if (i%3===0) {
                        column1.push(data[i])
                    }
                    else if(i%3===1) {
                        column2.push(data[i])
                    }
                    else {
                        column3.push(data[i])
                    }
                }
                $scope.subjectsCol1 = column1
                $scope.subjectsCol2 = column2
                $scope.subjectsCol3 = column3

                var length = column1.length + column2.length + column3.length
                // This allows me to hide the unsightly page in case a topic has
                // no available subjects in it to show a good sadface.
                console.log("subjects in scope ", $scope.subjects)
                if(length === 0) {
                    $rootScope.noSubjectsAvailable = true
                }
                else {
                    $rootScope.noSubjectsAvailable = false
                }

            }

            gateReaderServices.getSinglePost(topicArrived, $rootScope.requestedId)
            function topicArrived(data) {
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
            }
        }




    })

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







    // TODO: Here, check if the controller can access 2ndframeid. If it can, it's
    // at 3rd frame. In that case, write a ng-class directive that can look for
    // item id and activate on match so the one lights up.

    // I also need to determine here if this is in 2nd or 3rd scope, or write a
    // different controller for lite version.

}
SubjectsController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']
