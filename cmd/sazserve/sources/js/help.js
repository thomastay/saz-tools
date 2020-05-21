import $ from './jquery.js'
import 'chardin.js'
import { configuration, saveConfiguration } from './configuration.js'

let phase

function initializeHelp () {
  $(document.body)
    .on('chardinJs:start', enableMessage)
    .on('chardinJs:stop', disableMessage)
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

function enableMessage () {
  $('#help-message').show()
}

function disableMessage () {
  $('#help-message').hide()
}

export { initializeHelp, preparePageHelp, prepareTableHelp, showHelp }
