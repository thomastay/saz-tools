import $ from './jquery.js'
import 'chardin.js'
import { configuration, saveConfiguration } from './configuration.js'

let phase

function initializeHelp () {
  $('#help').on('click', showHelp)
  prepareInitialHelp()
  setTimeout(() => {
    if (configuration.help.initial !== false) {
      showHelp()
    }
  }, 500)
}

function prepareInitialHelp () {
  $(document.body).chardinJs({ url: 'json/help-initial.json' })
  phase = 'initial'
}

function prepareSazHelp () {
  $(document.body).chardinJs({ url: 'json/help-saz.json' })
  phase = 'saz'
}

function showHelp () {
  $(document.body).chardinJs('start')
  if (configuration.help[phase] !== false) {
    configuration.help[phase] = false
    saveConfiguration()
  }
}

export { initializeHelp, prepareInitialHelp, prepareSazHelp, showHelp }
