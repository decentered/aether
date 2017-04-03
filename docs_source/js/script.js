"use strict"


// insert fastclick

window.addEventListener('load', function() {
    FastClick.attach(document.body)
}, false)

// insert the old browser warning
var $buoop = {
  vs:{i:10,f:15,o:17,s:6,n:9},
  reminder: 0,
}
$buoop.ol = window.onload;
window.onload=function(){
 try {if ($buoop.ol) $buoop.ol();}catch (e) {}
 var e = document.createElement("script");
 e.setAttribute("type", "text/javascript");
 e.setAttribute("src", "http://browser-update.org/update.js");
 document.body.appendChild(e);
}




document.addEventListener('DOMContentLoaded', function() {

  // // Mobile hidden menu
  // var button = document.getElementsByClassName('mobile-menu-container')[0]
  // var centralArea = document.getElementsByClassName('central-area')[0]
  // var sidebar = document.getElementsByClassName('sidebar')[0]
  // var layoutInner = document.getElementsByClassName('layout')[0]
  // var mobileMenuOpen = false

  // function openSidebar() {
  //   centralArea.style.position = 'absolute'
  //   centralArea.style.left = '210px'
  //   sidebar.style.display = 'block'
  //   layoutInner.style.paddingLeft = '210px'
  //   button.style.left = '222px'
  //   mobileMenuOpen = true
  // }

  // function closeSidebar() {
  //   centralArea.style.position = null
  //   centralArea.style.left = null
  //   sidebar.style.display = null
  //   layoutInner.style.paddingLeft = null
  //   button.style.left = null
  //   mobileMenuOpen = false
  // }

  // button.addEventListener('click', function(){
  //   if (!mobileMenuOpen) {
  //     openSidebar()
  //   }
  //   else
  //   {
  //     closeSidebar()
  //   }
  // })


  // Shorten email placeholders for mobile devices.

  var firstEmailEntry = document.getElementsByClassName('email-signup-input')[0]
  var secondEmailEntry = document.getElementsByClassName('email-signup-input')[1]

  if (window.innerWidth < 768) {
    firstEmailEntry.placeholder = 'Enter your email for updates.'
    secondEmailEntry.placeholder = 'Enter your email for updates.'
  }

  var wH = window.innerHeight
  var entryView = document.getElementsByClassName('header')[0]
  var previews = document.getElementsByClassName('preview')
  startFlippingImages()
  for(var i=0;i<previews.length;i++) {
    previews[i].style.height = wH+'px'
  }
  var previewBox = document.getElementsByClassName('previews-box')[0]

  if (entryView && previewBox) {
    entryView.style.height = wH+'px'
    previewBox.style.height = wH+'px'
    // if (window.innerWidth > 768) {
    //   entryView.classList.add('darker')
    // }

    var moveButton = document.getElementsByClassName('scroll-bottom button')[0]
    var moveButtonTwo = document.getElementsByClassName('scroll-bottom button')[1]

    // Socials height implementation.
    var socialBox = document.getElementsByClassName('social-box')[0]
    socialBox.style.height = wH / 6.7+'px'

    moveButton.addEventListener('click', function(event) {
      scrollTo(wH, 800, (function (t) { return t<.5 ? 2*t*t : -1+(4-2*t)*t }))
    })

    moveButtonTwo.addEventListener('click', function(event) {
      scrollTo(wH*2, 800, (function (t) { return t<.5 ? 2*t*t : -1+(4-2*t)*t }))
    })


  }

  // spinning logo implementation
  var elem1 = document.body // Chrome default
  var elem2 = document.documentElement // Firefox default

  var logo = document.getElementsByClassName('aether-logo')[0]

  // Do not add the spinner if it is in mobile or ipad. (because scroll doesn't fire until after dragging ends)
  if (window.innerWidth > 1024) {
    window.addEventListener('scroll', function(){ requestAnimationFrame(logoSpinner) })
  }

  function logoSpinner() {
    var location = elem1.scrollTop > elem2.scrollTop ? elem1.scrollTop : elem2.scrollTop
    logo.style.webkitTransform = 'rotate('+location/100+'deg)'
    logo.style.transform = 'rotate('+location/100+'deg)'
  }


})

function startFlippingImages() {
  var firstImage = document.getElementById('first-screenshot')
  var secondImage = document.getElementById('second-screenshot')
  var thirdImage = document.getElementById('third-screenshot')
  var fourthImage = document.getElementById('fourth-screenshot')
  var fifthImage = document.getElementById('fifth-screenshot')
  var sixthImage = document.getElementById('sixth-screenshot')

  var images = [firstImage, secondImage, thirdImage,
                fourthImage, fifthImage, sixthImage]

  function flipper(){

      var elementToHide = (n-1>=0) ? images[n-1] : images[5]
      var elementToShow = images[n]
      handleVisibility(elementToHide, elementToShow)
      if (n>=5) {
        n=0
      } else {
        n++
      }
    }

  if (firstImage) {
    var n = 0
    flipper()
    setInterval(flipper ,2000)
  }

}

function handleVisibility(toHide, toShow) {
  // make it go to the back of the stack
  toHide.style.zIndex = '9'
  // Make it disappear 1 second later (from below the stack)
  setTimeout(function(){
    toHide.style.opacity = '0'
    toHide.style.visibility = 'hidden'
  },1000)

  toShow.style.opacity = '1'
  toShow.style.visibility = 'visible'
  // bring it to the front of the stack
  toShow.style.zIndex = '10'

}

function scrollTo(Y, duration, easingFunction, callback) {
    var browserChecked = false
    var start = Date.now()
    var elem1 = document.body // Chrome default
    var elem2 = document.documentElement // Firefox default
    //var elem = document.documentElement.scrollTop ? document.documentElement : document.body
    var from = elem1.scrollTop > elem2.scrollTop ? elem1.scrollTop : elem2.scrollTop
    if(from === Y) {
        if(callback) callback();
        return; // Prevent scrolling to the Y point if already there
    }
    function min(a,b) {
      return a<b?a:b;
    }
    function scroll(timestamp) {
        var currentTime = Date.now(),
            time = min(1, ((currentTime - start) / duration)),
            easedT = easingFunction(time);
        elem1.scrollTop = (easedT * (Y - from)) + from;
        elem2.scrollTop = (easedT * (Y - from)) + from;
        if(time < 1) requestAnimationFrame(scroll);
        else
            if(callback) callback();
    }
    requestAnimationFrame(scroll)
}
