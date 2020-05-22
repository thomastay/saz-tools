/* global Blob */

export default function (window, $) {
  $.fn.dataTable.ext.buttons.json = {
    text: 'JSON',
    action: function (event, dataTable, button, definition) {
      this.processing(true)
      const { header, body } = dataTable.buttons.exportData(definition.exportOptions)
      const output = JSON.stringify({ header, body }, undefined, 2)
      const info = dataTable.buttons.exportInfo(definition)
      $.fn.dataTable.fileSave(
        new Blob([output], { type: 'application/json' }), info.filename, true
      )
      this.processing(false)
    },
    filename: '*',
    extension: '.json',
    exportOptions: {}
  }
}
