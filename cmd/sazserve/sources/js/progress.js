import $ from './jquery.js'

let progressWrapper

function initializeProgress () {
  progressWrapper = $('#progress-wrapper')
}

function startProgress (event) {
  progressWrapper.show()
}

function stopProgress (event) {
  progressWrapper.hide()
}

export { initializeProgress, startProgress, stopProgress }
