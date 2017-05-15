function UnregisteredUserProfileController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {

    // This is the pagination level for the profile page. Profile shows only 10 results at a time for performance reasons.
    $scope.profileViewOffset = 10
    console.log('profileviewoffset FIRSTLOAD: ', $scope.profileViewOffset)

    gateReaderServices.getUnregisteredUserPosts(dataArrived, $rootScope.userProfile.userDetails.username, $scope.profileViewOffset)
    function dataArrived(data) {
        $scope.feed = data
    }

    $scope.showMore = function() {
        $scope.profileViewOffset += 10
        console.log('profileviewoffset: ', $scope.profileViewOffset)
    }

    $scope.$watch('profileViewOffset', function(){
        if ($scope.profileViewOffset != 10) {
            gateReaderServices.getUnregisteredUserPosts(offsetDataArrived, $rootScope.userProfile.userDetails.username, $scope.profileViewOffset)
            function offsetDataArrived(data) {
                $scope.feed = data
            }
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
UnregisteredUserProfileController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']