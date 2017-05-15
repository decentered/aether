function CreateLiteController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices, refreshService, $timeout, $window, $document) {

    $scope.postButtonDisabled = true
    $scope.postText = $rootScope.postTextCache
    $scope.target = angular.element(document.getElementsByClassName('scrolling-target'))[0]
    $scope.postCompleted = false
    $timeout(function(){
        var elm = angular.element(document.getElementsByClassName('create-box'))[0]
        elm.focus()
    },100)


     $scope.$on('$destroy', function() {
         if ($scope.postCompleted) {
            $scope.postText = ''
            $rootScope.postTextCache = '' // Posted successfully. Remove the post from cache.
         }
         else
         {
            var elm = angular.element(document.getElementsByClassName('create-box'))[0]
            $rootScope.postTextCache = elm.innerHTML // Cache the child HTML nodes (html input contented.) for future.
         }
    })

    var cachedSecondFrameCSSStyle, cachedThirdFrameCSSStyle

    $scope.closeReplyPane = function() {
        frameViewStateBroadcast.receiveState("", "subjectsFeedLite", $rootScope.requestedId)

        $rootScope.secondFrameCSSStyle = {
            'width': '850px'
        }
        $rootScope.thirdFrameCSSStyle = {
            'display':'block',
            'width': '243px'
        }
    }

    $scope.postReply = function(postText) {

        // check the username.
        if ($rootScope.userProfile.userDetails.username === '') {
            $rootScope.userProfile.userDetails.username = 'no name given'
        }

        //var content = $scope.postText.trim()
        // Why? Because this one keeps newlines intact. The other does not. Both are useful, but in the final
        // post you need the exact data.
        var content = angular.element(document.getElementsByClassName('create-box'))[0].innerText.trim()
        gateReaderServices.createPost(answerArrived, '', content, $rootScope.targetPost,
            $rootScope.userProfile.userDetails.username, $scope.targetPostObject.Language)
        function answerArrived(createdPostFingerprint) {
            console.log("this is the id of the currently created post: ",createdPostFingerprint)
            frameViewStateBroadcast.receiveState("", "subjectsFeedLite", $rootScope.requestedId)
            $rootScope.secondFrameCSSStyle = {
                'width': '850px'
            }
            $rootScope.thirdFrameCSSStyle = {
                'display':'block',
                'width': '243px'
            }
            $scope.postCompleted = true
            refreshService()

            $timeout(function(){
                var targetElement = document.getElementById('post-'+createdPostFingerprint)
                targetElement.classList.add('post-color-highlight')
                targetElement.scrollIntoViewIfNeeded()
            }, 10)


        }
    }

    $scope.togglePaneSize = function () {

        if ($rootScope.createPaneIsFullscreen === undefined ||
            $rootScope.createPaneIsFullscreen === false) {
            $rootScope.createPaneIsFullscreen = true
            // First, cache the styles,
            cachedSecondFrameCSSStyle = $rootScope.secondFrameCSSStyle
            cachedThirdFrameCSSStyle = $rootScope.thirdFrameCSSStyle
            // then assign new ones.
            $rootScope.secondFrameCSSStyle = {
                'display':'none'
            }
            $rootScope.thirdFrameCSSStyle = {
                'display':'block',
                'width':'100%'
            }
        }
        else {
            $rootScope.createPaneIsFullscreen = false
            // Returning to normal, call the old data from caches,
            $rootScope.secondFrameCSSStyle = cachedSecondFrameCSSStyle
            $rootScope.thirdFrameCSSStyle = cachedThirdFrameCSSStyle
            // and invalidate the caches.
            cachedSecondFrameCSSStyle = {}
            cachedThirdFrameCSSStyle = {}
        }
    }

    $scope.$watch('targetPost', function() {
        gateReaderServices.getSinglePost(subjectArrived, $rootScope.targetPost)
    })

    $scope.$watch('postText', function() {
        window.requestAnimationFrame(function(){
            if ($scope.target &&
                -document.getElementById('create-feed-lite').getBoundingClientRect().top + document.height + 100 >
                    document.getElementById('create-feed-lite').offsetHeight) {
                $scope.target.scrollIntoViewIfNeeded()
            }
        })
        $scope.trimmedPostText = $scope.postText.trim()
        if ($scope.trimmedPostText.length > 5 && $scope.trimmedPostText.length < 60000) {
            $scope.postButtonDisabled = false
        }
        else {
            $scope.postButtonDisabled = true
        }
    })

    $scope.$watch('createPaneIsFullscreen', function() {
        $timeout(function() {
            $scope.target.scrollIntoViewIfNeeded()
        }, 10)
        $timeout(function(){
            var elm = angular.element(document.getElementsByClassName('create-box'))[0]
            elm.focus()
        },100)
    })

    function subjectArrived(data) {
        $scope.targetPostObject = data
        if (data.Body === "") {
            $scope.parentIsSubject = true
        }
        else {
            $scope.parentIsSubject = false
            data.Body = $rootScope.mdConverter.makeHtml(data.Body)
        }
    }
}

CreateLiteController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
    'gateReaderServices', 'refreshService', '$timeout']