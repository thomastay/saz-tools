/* global sazTheme:writable, changeTheme, updateThemeSwitcher,
          ensureDarkOverrides, localStorage */

import $ from './jquery.js'

function initializeThemeSwitcher () {
  $('#theme-switcher').on('click', switchTheme)
}

function switchTheme () {
  const body = $(document.body)
  body.fadeOut(200, function () {
    $('#theme,#dark-overrides').remove()
    switch (sazTheme) {
      case 'dark': sazTheme = 'light'; break
      case 'light': sazTheme = 'system'; break
      default: sazTheme = 'dark'
    }
    changeTheme()
    ensureDarkOverrides()
    updateThemeSwitcher()
    saveTheme()
    body.fadeIn(200)
  })
}

function saveTheme () {
  setTimeout(() => localStorage.setItem('prantlf/sazview-theme', sazTheme))
}

export { initializeThemeSwitcher }
