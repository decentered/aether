function OnboardingController($scope, $rootScope, frameViewStateBroadcast,
    gateReaderServices, $timeout) {

    var possibleStates = ['welcome', 'license', 'shorthandLicense',
                            'handle', 'defaults', 'ethics', 'finish']

    $scope.currentState = 'welcome'

    $scope.currentLangShown = 1

    $scope.quitApp = function() {
        gateReaderServices.quitApp()
    }

    setInterval(function(){
        if ($scope.currentLangShown != 6) {
            $scope.currentLangShown++
        }
        else
        {
            $scope.currentLangShown = 1
        }
        $scope.$apply()
    }, 2000)

    $rootScope.currentStateNumber = 1 // This is used for bubble counting on onboarding.

    $scope.advanceState = function(callback) {
        for (var i=0;i<possibleStates.length;i++) { // If in the last state, don't touch.
                if ($scope.currentState === possibleStates[i]) {
                    if($scope.currentState === 'ethics') {
                        gateReaderServices.setOnboardingComplete()
                        $rootScope.onboardingComplete = !$rootScope.onboardingComplete
                        frameViewStateBroadcast.receiveState("findOrCreateTopic", "topicsFeedLite", "")
                        $rootScope.secondFrameCSSStyle = {}
                        $rootScope.thirdFrameCSSStyle = {
                        }
                    }
                    else {
                        $scope.currentState = possibleStates[i+1]
                        $rootScope.currentStateNumber++
                        if (callback != undefined) {
                            callback()
                        }
                        break
                    }

                }
            }
        }

        $scope.insertDefaultOptions = function() {
            startAtBoot = document.getElementById("openAtStartup").checked
            $rootScope.userProfile.machineDetails.startAtBoot = startAtBoot
            $rootScope.userProfile.machineDetails.lastConnectedBareIP = $scope.optionalBootstrapNodeIP
                $rootScope.userProfile.machineDetails.lastConnectedBarePort = $scope.optionalBootStrapNodePort
            if ($scope.optionalBootstrapNodeIP && $scope.optionalBootStrapNodePort) {
                gateReaderServices.connectToNodeWithIP(
                    $scope.optionalBootstrapNodeIP,
                    $scope.optionalBootStrapNodePort)
            }
            else
            {
                $rootScope.userProfile.machineDetails.lastConnectedBareIP = '151.236.11.192'
                $rootScope.userProfile.machineDetails.lastConnectedBarePort = 39994
                gateReaderServices.connectToNodeWithIP(
                    '151.236.11.192', 39994)
            }
        }

        // Defaults & Languages part

        $scope.EnglishSelected = true
        $rootScope.userProfile.userDetails.userLanguages = ['English']

        $scope.langClick = function(language) {
//            $scope[language+'Selected']
            console.log($scope[language+'Selected'])
            if ($scope[language+'Selected']) {
                if ($rootScope.userProfile.userDetails.userLanguages.length > 1) {
                    $scope[language+'Selected'] = false
                    var index = $rootScope.userProfile.userDetails.userLanguages.indexOf(language)
                    if (index > -1) {
                        $rootScope.userProfile.userDetails.userLanguages.splice(index, 1);
                    }
                }
            }
            else
            {
                $scope[language+'Selected'] = true
                $rootScope.userProfile.userDetails.userLanguages.push(language)
            }
        }

    }

OnboardingController.$inject = ['$scope', '$rootScope', 'frameViewStateBroadcast',
    'gateReaderServices', '$timeout']


