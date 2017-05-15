function CreateSubjectController($scope, $rootScope, $timeout, frameViewStateBroadcast,
    gateReaderServices, refreshService) {

    $scope.bodyLetterCount = 0
    $scope.subjectLetterCount = 0
    $scope.bodyWordCount = 0
    $scope.postLanguage = ''
    $scope.subjectLanguage = ''
    $scope.langInUserLangs = undefined // We don't know yet.
    $scope.writeInProgress = true
    $scope.postButtonDisabled = true
    $scope.subjectHeaderText = $rootScope.subjectHeaderCache
    $scope.subjectBodyText = $rootScope.postTextCache
    $scope.possLangs = []
    $scope.possLangs.push('Auto')
    $scope.target = angular.element(document.getElementsByClassName('scrolling-target'))[0]
    $scope.postCompleted = false // This is useful, because I need to know this on destroy.

    $timeout(function(){
        var headerElm = angular.element(document.getElementsByClassName('subject-name-entry'))[0]
        headerElm.focus()
    },100)

    $scope.$on('$destroy', function() {

        if ($scope.postCompleted) {
            $rootScope.postTextCache = '' // Remove the cache after successful post.
            $rootScope.subjectHeaderCache = ''
            $scope.subjectHeaderText = ''
            $scope.subjectBodyText = ''
        }
        else
        {
            var elm = angular.element(document.getElementsByClassName('subject-body-entry'))[0]
            var headerElm = angular.element(document.getElementsByClassName('subject-name-entry'))[0]
            $rootScope.postTextCache = elm.innerHTML // Cache the child HTML nodes (html input contented.) for future.
            $rootScope.subjectHeaderCache = headerElm.innerText
        }
    })


    $scope.cancel = function() {
        $rootScope.changeState('subjectsFeed', '', '')
    }

    $scope.setLang = function(lang) {
        console.log('manually selected a lang:', lang)
        $scope.language = lang
        //$scope.postButtonDisabled = false
        $scope.langManuallySelected = true
    }

    $scope.postSubject = function() {

        // check the username.
        if ($rootScope.userProfile.userDetails.username === '') {
            $rootScope.userProfile.userDetails.username = 'no name given'
        }
        $scope.writeInProgress = false
        // These two keep the newlines.
        var subjectHeader = angular.element(document.getElementsByClassName('subject-name-entry'))[0].innerText.trim()
        var subjectBody = angular.element(document.getElementsByClassName('subject-body-entry'))[0].innerText.trim()
        var found = false
        // The routine below handles the language detection implementation.
        if (!$scope.langManuallySelected) {

            // Figure out the most sound answer,
            if (subjectBody.length > subjectHeader.length) {
                $scope.language = detectLanguage(subjectBody)
            }
            else {
                $scope.language = detectLanguage(subjectHeader)
            }
            // Produce the model for dropdown to bind to, (possibleLanguages)
            $scope.possibleLanguages = []
            if ($scope.language != 'Tg_unknown_language') {
                $scope.possibleLanguages.push({name:$scope.language, langGroup:'Detected'})
            }
            for (var i=0; i<$rootScope.userProfile.userDetails.userLanguages.length; i++) {
                var langItem = {
                    name: $rootScope.userProfile.userDetails.userLanguages[i],
                    langGroup: 'Your Language(s)'
                }
                $scope.possibleLanguages.push(langItem)
            }
            $scope.selectedLanguage = $scope.possibleLanguages[0]
            // Do checking on the model if I need to show the dropdown.
            found = false
            for (var i=0; i<$rootScope.userProfile.userDetails.userLanguages.length; i++) {
                if ($scope.language === $rootScope.userProfile.userDetails.userLanguages[i]) {
                    found = true
                    break
                }
            }
            if (!found) {
                $scope.langManuallySelected = true
                return false // So the function exits here if it matches.
            }
        }


        if (!found && $scope.langManuallySelected) {
            if ($rootScope.userProfile.userDetails.userLanguages.indexOf($scope.language) === -1) {
                $rootScope.userProfile.userDetails.userLanguages.push($scope.language)
            }
        }

        gateReaderServices.createPost(subjectPostCreated,subjectHeader, '', $scope.requestedId,
            $rootScope.userProfile.userDetails.username, $scope.language)

        function subjectPostCreated(createdSubjectPostFingerprint) {
            console.log('this is the id of the created subject', createdSubjectPostFingerprint)
            gateReaderServices.createPost(bodyPostCreated, '', subjectBody,
                createdSubjectPostFingerprint, $rootScope.userProfile.userDetails.username,
                $scope.language)

            function bodyPostCreated(createdPostFingerprint) {
                console.log('this is the id of the created post', createdPostFingerprint)
                $scope.postCompleted = true
                $rootScope.changeState('subjectsFeed', '', '')

                $timeout(function(){
                    console.log('lock into the newly created subject triggered')
                    var singleColumnTarget = document.querySelector('.single-column #post-'+createdSubjectPostFingerprint)
                    var multiColumnTarget = document.querySelector('.multi-column #post-'+createdSubjectPostFingerprint)
                    singleColumnTarget.classList.add('post-color-highlight')
                    multiColumnTarget.classList.add('post-color-highlight')
                    singleColumnTarget.scrollIntoViewIfNeeded()
                    multiColumnTarget.scrollIntoViewIfNeeded()
                }, 10)
            }
        }
    }

    // Below are watches that fire on every keypress. This is quite low performance,
    // find a way around this. TODO.

    $scope.$watch('subjectHeaderText + subjectBodyText', function() {
        $scope.trimmedSubjectBodyText = $scope.subjectBodyText.trim()
        $scope.trimmedSubjectHeaderText = $scope.subjectHeaderText.trim()
        if ($scope.trimmedSubjectHeaderText.length > 10 &&
            $scope.trimmedSubjectBodyText.length > 5 &&
            $scope.trimmedSubjectBodyText.length < 60000) {
            $scope.postButtonDisabled = false
        }
        else {
            $scope.postButtonDisabled = true
        }
    })

    $scope.$watch('subjectBodyText', function() {
        // This will fire on every keypress on body
        // If this is an edit after user clicks to post and gets stopped,
        $scope.trimmedSubjectBodyText = $scope.subjectBodyText.trim()
        $scope.bodyLetterCount = $scope.trimmedSubjectBodyText.length
        $scope.bodyWordCount = $scope.trimmedSubjectBodyText.split(/\s+/).length - 1

        if (!$scope.writeInProgress) {
            $scope.langManuallySelected = false
        }
        $scope.writeInProgress = true
        //$scope.postButtonDisabled = false
        $scope.selectedLanguage = $scope.possLangs[0]
    })

    $scope.$watch('subjectHeaderText', function() {
        $scope.trimmedSubjectHeaderText = $scope.subjectHeaderText.trim()
        enforceHeaderMaxLength($scope.trimmedSubjectHeaderText, 225)
        $scope.subjectLetterCount = $scope.trimmedSubjectHeaderText.length

        if (!$scope.writeInProgress) {
            $scope.langManuallySelected = false
        }
        $scope.writeInProgress = true
        //$scope.postButtonDisabled = false
        $scope.selectedLanguage = $scope.possLangs[0]
    })

    $scope.$watch('subjectBodyText', function() {
        window.requestAnimationFrame(function(){
            if ($scope.target &&
                -document.getElementById('create-subject-feed').getBoundingClientRect().top + document.height + 100 >
                    document.getElementById('create-subject-feed').offsetHeight) {
                $scope.target.scrollIntoViewIfNeeded()
            }
        })
    })

    // Below are watches that fire only on certain intervals.

    $scope.$watch(function() { return Math.floor($scope.subjectBodyText.length / 200) }, function() {
        // This will fire every 200 words written, which is about when an update
        // is needed.
        $scope.timeItTakesToRead = calculateHowLongItTakesToRead($scope.bodyWordCount)
    })

    $scope.$watch(function() { return Math.floor($scope.subjectBodyText.length / 20) }, function() {
        // This will fire every 20 letters written, which is about when an update
        // is needed to guessed lang.
        $scope.possLangs = ['Auto']
        var possLang = detectLanguage(
                angular.element(document.getElementsByClassName('subject-body-entry')).text())
        if (possLang != 'Tg_unknown_language') {
            $scope.possLangs.push(possLang)
        }
        // For each of the user's languages,
        for (var i=0; i<$rootScope.userProfile.userDetails.userLanguages.length; i++) {
            // If the guessed language (possLangs[1], bc [0] is 'Auto') is not it,
            if ($rootScope.userProfile.userDetails.userLanguages[i] != $scope.possLangs[1]) {
                // Add it to the user's languages.
                $scope.possLangs.push($rootScope.userProfile.userDetails.userLanguages[i])
            }
        }
        $scope.selectedLanguage = $scope.possLangs[0]
    })

    // Algorithmic functions.

    function calculateHowLongItTakesToRead(wordCount) {
    // An average American reads 250 words per minute. Taking 200 to accommodate
    // everyone. There could be a few niceties to be put here...
        var minutes = Math.floor(wordCount / 200)
        var t = ''
        if (minutes === 0) {
            t = 'less than a min'
        }
        else if (minutes === 1) {
            t = 'about 1 min'
        }
        else {
            t = 'about ' + minutes + ' min'
        }
        return t
    }

    // FAQ: Why is this function duplicated?
    // Answer: For this function to work, it needs to have access to the immediate
    // scope, which is not possible if this function is generalised and carried
    // to parent scopes. Angular does not allow (as of in time of this writing)
    // access to child scopes.

    function enforceHeaderMaxLength(subjectHeaderText, maxLength) {
        $scope.headerLetterCount = subjectHeaderText.length
        if (subjectHeaderText.length > maxLength) {
            $scope.subjectHeaderText = subjectHeaderText.slice(0, -(subjectHeaderText.length - maxLength))
            // [0] of any jQuery or angular.element element is the element itself.
            setCaretToEndOfContenteditable(
                angular.element(document.getElementsByClassName('subject-name-entry'))[0])
        }

        function setCaretToEndOfContenteditable(target) {
            var range = document.createRange() // A range is an invisible selection
            range.selectNodeContents(target) // select node contents invisibly
            range.collapse(false) // collapse the selection to the end
            var selection = window.getSelection() // create a visible selection obj.
            selection.removeAllRanges() // remove the range within even if exists
            selection.addRange(range) // and make the range the selection's range.
        }

    }

}

CreateSubjectController.$inject = ['$scope', '$rootScope', '$timeout', 'frameViewStateBroadcast',
    'gateReaderServices', 'refreshService']
