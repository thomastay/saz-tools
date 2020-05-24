import $ from './jquery.js'
import { sazStore, uploadSaz, downloadSaz } from './saz-store.js'
import { loadConfiguration } from './configuration.js'
import { initializeThemeSwitcher } from './theme-switcher.js'
import { initializeHelp } from './help.js'
import { initializePageState, displayInformation, resetPage, displayError } from './page-state.js'
import { initializePreviousSaz, updatePreviousSaz } from './previous-saz.js'
import { initializeProgress, startProgress, stopProgress } from './progress.js'
import { initializeSessionTable, displaySessionTable } from './session-table.js'

let currentSaz

setTimeout(initialize)

function initialize () {
  initializeProgress()
  initializeSessionTable()
  initializePageState()
  initializeHelp()
  initializeThemeSwitcher()
  initializePreviousSaz(previousSazChanged)
  currentSaz = $('#saz-file').on('change', currentSazChanged)
  $(document.documentElement)
    .on('dragenter', allowDragDrop)
    .on('dragover', showDropEffect)
    .on('dragleave', allowDragDrop)
    .on('drop', dropSaz)
  loadConfiguration()
  stopProgress()
}

function allowDragDrop (event) {
  event.preventDefault()
  event.stopPropagation()
}

function showDropEffect (event) {
  allowDragDrop(event)
  event.originalEvent.dataTransfer.dropEffect = 'copy'
}

function dropSaz (event) {
  allowDragDrop(event)
  currentSaz.val('')
  processSazs(event.originalEvent.dataTransfer.files)
}

function currentSazChanged (event) {
  processSazs(event.target.files)
}

function processSazs (files) {
  if (!files.length) {
    resetPage()
    return displayInformation()
  }
  startProgress()
  const promises = []
  for (const file of files) {
    promises.push(uploadSaz(file))
  }
  Promise
    .allSettled(promises)
    .then(results => {
      let response
      for (const { value, reason } of results) {
        if (reason) {
          return displayError(reason)
        }
        if (value) {
          response = value
        }
      }
      sazStore.loaded = sazStore.stored[response.File.name]
      resetPage()
      displaySessionTable(sazStore.loaded.Sessions)
      updatePreviousSaz()
    })
    .then(stopProgress)
}

function previousSazChanged (event) {
  const name = Object.keys(sazStore.stored)[$(event.target).prop('selectedIndex')]
  const saz = sazStore.stored[name]
  startProgress()
  downloadSaz(saz)
    .then(sessions => {
      resetPage()
      displaySessionTable(sessions)
    })
    .catch(response => {
      if (response instanceof Error || response.status !== 404) {
        throw response
      }
      return uploadSaz(saz.File)
        .then(({ Sessions: sessions }) => {
          resetPage()
          displaySessionTable(sessions)
        })
    })
    .catch(displayError)
    .then(stopProgress)
}
