function RepliesController($scope, $rootScope, frameViewStateBroadcast, gateReaderServices) {
    var repliesArray = []
    gateReaderServices.getReplies(dataArrived)
    function dataArrived(data) {
        for (var i=0;i<data.length;i++) {
            (function(i){
                if ($rootScope.userProfile.userDetails.unreadReplies.indexOf(data[i].PostFingerprint) === -1
                    && $rootScope.userProfile.userDetails.readReplies.indexOf(data[i].PostFingerprint) === -1) {
                    // If the coming reply does not exist both in unread and read, it's new.
                    // Add the fingerprint to the unreads.
                    $rootScope.userProfile.userDetails.unreadReplies.push(data[i].PostFingerprint)
                    // Flag as Unread to the frontend.
                    data[i].Unread = true
                }
                else if ($rootScope.userProfile.userDetails.readReplies.indexOf(data[i].PostFingerprint) > -1) {
                    // If it exists in the read
                    data[i].Unread = false
                }
                else if ($rootScope.userProfile.userDetails.unreadReplies.indexOf(data[i].PostFingerprint) > -1) {
                    // If it exists in the unread
                    data[i].Unread = true
                }

                // And add to the result array.
                repliesArray.push(data[i])
            })(i)
        }
    }

    $scope.replies = repliesArray

    $scope.clickToReply = function(postFingerprint) {
        if ($rootScope.userProfile.userDetails.unreadReplies.indexOf(postFingerprint) > -1) {
            // If it exists in the unread
            var index = $rootScope.userProfile.userDetails.unreadReplies.indexOf(postFingerprint)
            // remove from the unreads,
            $rootScope.userProfile.userDetails.unreadReplies.splice(index, 1)
            // and add to the reads.
            $rootScope.userProfile.userDetails.readReplies.push(postFingerprint)
        }
        $rootScope.changeState('singleReply', '', postFingerprint)
    }

    $scope.clickToMarkAllAsRead = function() {
        gateReaderServices.markAllRepliesAsRead(answerArrived)
        function answerArrived(answer) {
            if (answer === true) {
                $rootScope.changeState('homeFeed', '', '')
                $rootScope.userProfile.userDetails.unreadReplies = []
                $rootScope.userProfile.userDetails.readReplies = []
                $rootScope.totalReplyCount = 0
            }
        }
    }
}
RepliesController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast', 'gateReaderServices']