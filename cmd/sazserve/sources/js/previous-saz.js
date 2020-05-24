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
    const names = Object
      .keys(sazStore.stored)
      .map((name, index) => ({ name, index }))
      .sort((left, right) => left.name.localeCompare(right.name))
    for (const { name, index } of names) {
      const option = $('<option>')
        .text(name)
        .attr('value', index.toString())
      if (name === sazStore.loaded.File.name) {
        option.attr('selected', 'selected')
      }
      previousSaz.append(option)
    }
  })
}

export { initializePreviousSaz, resetPreviousSaz, updatePreviousSaz }
