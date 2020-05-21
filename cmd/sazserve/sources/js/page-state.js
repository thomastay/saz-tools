import $ from './jquery.js'
import { preparePageHelp } from './help.js'
import { resetPreviousSaz } from './previous-saz.js'
import { destroySessionTable } from './session-table.js'
import { parseError } from './data-parsers.js'

let infoAlert, errorAlert

function initializePageState () {
  infoAlert = $('.alert-info')
  errorAlert = $('.alert-danger')
}

function displayInformation () {
  infoAlert.show()
}

function resetPage () {
  infoAlert.hide()
  errorAlert.hide()
  destroySessionTable()
}

function displayError (response) {
  const { title, text } = parseError(response)
  resetPage()
  resetPreviousSaz()
  if (title) {
    errorAlert.find('h4').show().text(title)
  } else {
    errorAlert.find('h4').hide()
  }
  errorAlert.show().find('p').text(text)
  preparePageHelp()
}

export { initializePageState, displayInformation, resetPage, displayError }
