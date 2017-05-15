function UnregisteredUserFeed($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices) {
    $scope.profileViewOffset = 10
    gateReaderServices.getUnregisteredUserPosts(dataArrived, $rootScope.requestedUsername, 10)
    function dataArrived(data) {
        var leftFeed = []
        var rightFeed = []
        for (var i =0 ; i<data.length; i++) {
            if (i%2===0) {
                (function(i){
                    data[i].Body = $rootScope.mdConverter.makeHtml(data[i].Body)
                    leftFeed.push(data[i])
                })(i)

            }
            else {
                (function(i){
                    data[i].Body = $rootScope.mdConverter.makeHtml(data[i].Body)
                    rightFeed.push(data[i])
                })(i)
            }
        }
        $scope.leftFeed = leftFeed
        $scope.rightFeed = rightFeed
    }

    $scope.showMore = function() {
        $scope.profileViewOffset += 10
        console.log('profileviewoffset: ', $scope.profileViewOffset)
    }

    $scope.$watch('profileViewOffset', function(){
        if ($scope.profileViewOffset != 10) {
            gateReaderServices.getUnregisteredUserPosts(dataArrived, $rootScope.userProfile.userDetails.username, $scope.profileViewOffset)
        }

    })

    $scope.decideNgInclude = function(postSubject) {
        if (postSubject === '') {
            // It is a post
            return 'contentBlocks/userPostItem.html'
        }
        else {
            // It is a subject
            return 'contentBlocks/userSubjectItem.html'
        }
    }

    $scope.clickToPostItem = function(postFingerprint) {
        $rootScope.changeState('singleReply', '', postFingerprint)

    }

    $scope.clickToSubjectItem = function(postFingerprint) {
        $rootScope.changeState('postsFeed', 'subjectsFeedLite', postFingerprint)
        $rootScope.secondFrameCSSStyle = {}
        $rootScope.thirdFrameCSSStyle = {
            'display':'block'
        }
    }

}
UnregisteredUserFeed.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']