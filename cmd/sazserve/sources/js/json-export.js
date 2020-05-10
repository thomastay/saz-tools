export default function (window, $) {
  $.fn.dataTable.ext.buttons.json = {
    available: function () {
      return window.FileReader !== undefined && window.Blob
    },
    text: 'JSON',
    action: function (event, dataTable, button, config) {
      this.processing(true)
      const { header, body } = dataTable.buttons.exportData(config.exportOptions)
      const output = JSON.stringify({ header, body }, undefined, 2)
      const info = dataTable.buttons.exportInfo(config)
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
