const tableContainer = $('#table-container')
const infoContainer = $('.alert-info')
const errorContainer = $('.alert-danger')

function viewSaz (saz) {
  resetPage()
  const file = saz.files[0]
  if (file) {
    loadSaz(file)
      .then(displaySaz)
      .catch(displayError)
  } else {
    infoContainer.show()
  }

  function resetPage () {
    infoContainer.hide()
    errorContainer.hide()
    if (window.sazTable) {
      window.sazTable.destroy()
      tableContainer.html('')
      window.sazTable = undefined
    }
  }

  function loadSaz (file) {
    const formData = new FormData()
    formData.append('saz', file)
    return $.ajax({
      method: 'POST',
      url: '/saz',
      data: formData,
      contentType: false,
      processData: false
    })
  }

  function displaySaz (response) {			
    var data = response.map(session => [
      session.Number,
      formatDuration(session.Timeline),
      session.Request.Method,
      session.Response.StatusCode,
      formatURL(session.Request.URL),
      formatHost(session.Request.URL),
      formatPath(session.Request.URL),
      formatTime(session.Timers.ClientBeginRequest),
      formatTime(session.Timers.ClientDoneResponse),
      formatDuration(session.Duration),
      session.Response.ContentLength,
      session.Encoding, session.Caching, session.Flags.Process
    ])
    var columns = [
      { title: '#' },
      { title: 'Timeline' },
      { title: 'Method' },
      { title: 'Status' },
      { title: 'URL' },
      { title: 'Host', visible: false },
      { title: 'Path', visible: false },
      { title: 'Begin' },
      { title: 'End' },
      { title: 'Duration' },
      {
        title: 'Size',
        className: 'dt-right',
        render: function ( data, type, row) {
          if (type === 'display' || type === 'filter') {
            return formatSize(data)
          }
          return data
        }
      },
      { title: 'Encoding' },
      { title: 'Caching', orderable: false },
      { title: 'Process' }
    ]
    window.sazTable = $('<table class="table table-sm table-striped table-hover nowrap compact display">')
        .appendTo(tableContainer)
        .DataTable({
          columns,
          data,
          dom: '<"top"ifBR>rtS',
          scrollX: true,
          scrollY: '65vh',
          scrollCollapse: true,
          deferRender: true,
          scroller: true,
          colReorder: true,
          fixedColumns: { leftColumns: 1 },
          buttons: [
            {
              extend: 'colvis',
              text: 'Columns'
            },
            'copy', 'print',
            {
              extend: 'collection',
              text: 'Export',
              buttons: [ 'csv', 'excel', 'pdf' ]
            }
          ]
        })
  }

  function displayError (response) {
    const title = response.status && `${response.status} (${response.statusText})`
    const text = response.responseText || 'Connection failed.'
    if (title) {
      errorContainer.find('h4').show().text(title)
    } else {
      errorContainer.find('h4').hide()
    }
    errorContainer.show().find('p').text(text)
  }

  function formatURL (url) {
    return url.replace(/^(.{160}).*$/, '$1...')
  }

  function formatHost (url) {
    return url.replace(/^(?:\w+:)?\/\/([^:/]+).*$/, '$1')
  }

  function formatPath (url) {
    return url.replace(/^(?:\w+:)?\/\/[^/]+(.*)$/, '$1').replace(/^(.{80}).*$/, '$1')
  }

  function formatDuration (duration) {
    return duration.substr(6)
  }

  function formatTime (duration) {
    return duration.substr(14, 12)
  }

  function formatSize (size) {
    if (size < 0) {
      return 'N/A'
    }
    return formatNumber(size)
  }

  function formatNumber (number) {
    let output = number.toString()
    for (let outputIndex = output.length; outputIndex > 3;) {
      outputIndex -= 3
      output = output.substr(0, outputIndex) + ',' + output.substr(outputIndex)
    }
    return output
  }
}
