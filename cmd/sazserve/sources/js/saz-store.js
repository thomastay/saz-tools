import $ from './jquery.js'
import { updatePreviousSaz } from './previous-saz.js'

const sazStore = {
  loaded: null,
  stored: {}
}

function uploadSaz (file) {
  return postSaz(file).then(storeSaz)
}

function postSaz (file) {
  return $.ajax({
    method: 'POST',
    url: '/api/saz',
    data: file,
    contentType: 'application/octet-stream',
    processData: false
  }).then(response => ({ File: file, ...response }))
}

function storeSaz (response) {
  sazStore.loaded = sazStore.stored[response.File.name] = response
  updatePreviousSaz()
  return response
}

function downloadSaz (key) {
  return $.ajax({ url: `/api/saz/${key}` })
}

export { sazStore, uploadSaz, downloadSaz }
