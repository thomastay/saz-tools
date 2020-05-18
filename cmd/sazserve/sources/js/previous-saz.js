import $ from './jquery.js'
import { sazStore } from './saz-store.js'

let previousSaz

function initializePreviousSaz (previousSazChanged) {
  previousSaz = $('#previous-saz').on('change', previousSazChanged)
}

function resetPreviousSaz () {
  previousSaz.prop('selectedIndex', -1)
}

function updatePreviousSaz () {
  setTimeout(() => {
    previousSaz.html('')
    for (const name in sazStore.stored) {
      const option = $('<option>').text(name)
      if (name === sazStore.loaded.File.name) {
        option.attr('selected', 'selected')
      }
      previousSaz.append(option)
    }
  })
}

export { initializePreviousSaz, resetPreviousSaz, updatePreviousSaz }
