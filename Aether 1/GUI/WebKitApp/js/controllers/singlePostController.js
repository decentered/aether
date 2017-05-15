function SinglePostController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices, $timeout) {



    $scope.postButtonDisabled = true
    $scope.postText = $rootScope.postTextCache
    $scope.postCompleted = false // This is useful, because the cache events happen on destroy.

    $scope.$on('$destroy', function() {
        if ($scope.postCompleted) {
            $scope.postText = ''
            $rootScope.postTextCache = ''
        }
        else
        {
            var elm = angular.element(document.getElementsByClassName('subject-body-entry'))[0]
            $rootScope.postTextCache = elm.innerHTML // Cache the child HTML nodes (html input contented.) for future.
        }
    })


    gateReaderServices.getSinglePost(postArrived, $rootScope.requestedId)
    function postArrived(p) {
        p.Body = $rootScope.mdConverter.makeHtml(p.Body)
        $scope.post = p
        gateReaderServices.getParentSubjectOfGivenPost(parentSubjectArrived, p.PostFingerprint)
        function parentSubjectArrived(parentSubject) {
            if (parentSubject.length) { // Assign scope.subject only if there is a subject available.
                $scope.subject = parentSubject[0]
            }
        }
        gateReaderServices.getSinglePost(parentPostArrived, $scope.post.ParentPostFingerprint)
        function parentPostArrived(parent) {
            if (parent.Body) {
                parent.Body = $rootScope.mdConverter.makeHtml(parent.Body)
                $scope.parentPost = parent
            }
        }
    }

    $scope.exportButtonClick = function() {
        gateReaderServices.exportSinglePost(answerArrived, $scope.post)
        function answerArrived(answer) {
            console.log(answer)
        }
    }

    $scope.parentIsVisible = false

    $scope.toggleParentVisibility = function() {
        if ($scope.parentIsVisible) {
            $scope.parentIsVisible = false
        }
        else
        {
            $scope.parentIsVisible = true
        }
    }

    $scope.upvotePost = function(postFingerprint) {
        gateReaderServices.votePost(answerArrived, postFingerprint, 1)
        function answerArrived(answer) {
            if (answer === true) {
                var post = $scope.post
                if (!post.Upvoted) {
                    // If post is not already upvoted
                    post.UpvoteCount += 1
                    post.Upvoted = true
                    if (post.Downvoted) {
                        post.DownvoteCount -= 1
                        post.Downvoted = false
                    }
                }
                else {
                    // If it is already upvoted, remove the upvote.
                    gateReaderServices.votePost(answerArrived, postFingerprint, 2)
                    function answerArrived(answer) {
                        post.UpvoteCount -= 1
                        post.Upvoted = false
                    }
                }
            }
        }
    }

    $scope.downvotePost = function(postFingerprint) {
        gateReaderServices.votePost(answerArrived, postFingerprint, -1)
        function answerArrived(answer) {
            if (answer === true) {
                var post = $scope.post
                if (!post.Downvoted) {
                    // If this is a new downvote
                    post.DownvoteCount += 1
                    post.Downvoted = true
                    if (post.Upvoted) {
                        post.UpvoteCount -= 1
                        post.Upvoted = false
                    }
                }
                else {
                    // if already voted
                    gateReaderServices.votePost(answerArrived, postFingerprint, 2)
                    function answerArrived(answer) {
                        post.DownvoteCount -= 1
                        post.Downvoted = false
                    }
                }
            }
        }
    }

    $scope.savePost = function(postFingerprint) {
        gateReaderServices.savePost(answerArrived, postFingerprint)
        function answerArrived(answer) {
            if (answer === true) {
                var post = $scope.post
                post.Saved === true ? post.Saved = false : post.Saved = true
            }
        }
    }

    $scope.replyButtonClick = function() {
        $scope.replyPaneOpen = $scope.replyPaneOpen ? false : true
        $scope.target = angular.element(document.getElementsByClassName('scrolling-target'))[0]
        $timeout(function() {
            $scope.target.scrollIntoViewIfNeeded()
        }, 10)
        $timeout(function(){
            var elm = angular.element(document.getElementsByClassName('subject-body-entry'))[0]
            elm.focus()
        },100)

    }

    $scope.submitButtonClick = function() {

        // check the username.
        if ($rootScope.userProfile.userDetails.username === '') {
            $rootScope.userProfile.userDetails.username = 'no name given'
        }

        $scope.replyPaneOpen = false
        var content = angular.element(document.getElementsByClassName('subject-body-entry'))[0].innerText.trim()
        // IMPORTANT: LANGUAGE SANITY CHECK, the lang needs to exist in selected user langs.

        gateReaderServices.createPost(answerArrived, '', content, $scope.post.PostFingerprint,
            $rootScope.userProfile.userDetails.username, $scope.post.Language)

        function answerArrived(createdPostFingerprint) {
            console.log('The user has replied to one of his / her replies. The fingerprint of the new reply is: ', createdPostFingerprint)
            $scope.postCompleted = true
            $scope.goToThread(createdPostFingerprint)


        }
    }

    $scope.$watch('postText', function() {
        window.requestAnimationFrame(function(){
            if ($scope.target &&
                -document.getElementById('single-post').getBoundingClientRect().top + document.height + 100 >
                    document.getElementById('single-post').offsetHeight) {
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

    $scope.goToThread = function(lightupTarget) {
        if ($scope.subject) {
            $rootScope.secondFrameCSSStyle = {}
            $rootScope.thirdFrameCSSStyle = {
                'display':'block'
            }
            $rootScope.changeState('postsFeed', 'subjectsFeedLite', $scope.subject.PostFingerprint)
            $timeout(function(){
                if (lightupTarget) {
                    var targetElement = document.getElementById('post-'+lightupTarget)
                }
                else
                {
                    var targetElement = document.getElementById('post-'+$scope.post.PostFingerprint)
                }
                targetElement.classList.add('post-color-highlight')
                targetElement.scrollIntoViewIfNeeded()
            }, 10)
        }
    }

//    var subjectBodyEntry = angular.element(document.getElementsByClassName('subject-body-entry')[0])
//    var metaBox = angular.element(document.getElementsByClassName('meta-box'))
//    $scope.subjectBodyEntryStyle = {}
//
//    subjectBodyEntry.bind('paste keydown', function() {
//        $timeout(function() {
//            metaBox[0].scrollIntoViewIfNeeded()
//            $scope.bodyLetterCount = subjectBodyEntry.text().trim().length
//            $scope.bodyWordCount = subjectBodyEntry.text().trim().split(/\s+/).length - 1
//        }, 10)
//    })

}
SinglePostController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
    'gateReaderServices', '$timeout']
