import { configuration, saveConfiguration } from './configuration.js'

export default function (window, $) {
  $.fn.dataTable.ext.buttons.colorful = {
    text: 'Colourful sessions',
    init: function (dataTable, button, definition) {
      this.active(configuration.colorfulSessions)
    },
    action: function (event, dataTable, button, definition) {
      configuration.colorfulSessions = !configuration.colorfulSessions
      this.active(configuration.colorfulSessions)
      dataTable.table().body().parentNode.classList.toggle('colorful')
      saveConfiguration()
    }
  }
}
