// Services > Tooltips

// This service provides a tooltip binder - for any view that has tooltips, this needs to run.

export { }
var Tippy: any = require('../../../../node_modules/tippy.js')

module.exports = {
  Mount() {
    Tippy('[hasTooltip]', {
      animation: 'fade',
      delay: 50,
      duration: [200, 500], // [show, hide]
      performance: true,
      placement: 'top',
      theme: 'dark',
      // size: 'small',
      // inertia: true
    })
  },
  MountInfomark() {
    Tippy('[hasTooltip]', {
      animation: 'fade',
      delay: 50,
      duration: [200, 500], // [show, hide]
      performance: true,
      placement: 'bottom',
      theme: 'infomark',
      // size: 'small',
      // inertia: true
      arrow: true,
      arrowType: 'round',
      hideOnClick: false,
      interactive: true,
      allowTitleHTML: true,
      offset: '0,5',
    })
  }
}