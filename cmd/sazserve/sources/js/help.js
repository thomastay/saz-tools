import $ from './jquery.js'
import 'chardin.js'
import { configuration, saveConfiguration } from './configuration.js'

let phase, shown

function initializeHelp () {
  $(document.body)
    .on('chardinJs:start', enableMessage)
    .on('chardinJs:stop', disableMessage)
    .on('keydown', hideHelp)
  $('#help').on('click', showHelp)
  preparePageHelp()
  setTimeout(() => {
    if (configuration.help.initial !== false) {
      showHelp()
    }
  }, 500)
}

function preparePageHelp () {
  $(document.body).chardinJs({ url: 'json/help-page.json' })
  phase = 'initial'
}

function prepareTableHelp () {
  $(document.body).chardinJs({ url: 'json/help-table.json' })
  phase = 'saz'
}

function showHelp () {
  $(document.body).chardinJs('start')
  if (configuration.help[phase] !== false) {
    configuration.help[phase] = false
    saveConfiguration()
  }
}

function hideHelp () {
  if (shown) {
    $(document.body).chardinJs('stop')
  }
}

function enableMessage () {
  $('#help-message').show()
  shown = true
}

function disableMessage () {
  $('#help-message').hide()
  shown = false
}

export { initializeHelp, preparePageHelp, prepareTableHelp, showHelp }
