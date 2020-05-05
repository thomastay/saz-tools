const tableContainer = $('#table-container')
const infoContainer = $('.alert-info')
const errorContainer = $('.alert-danger')

function viewSaz (saz) {
  const file = saz.files[0]
  if (file) {
    loadSaz(file)
      .then(displaySaz)
      .catch(displayError)
  } else {
    resetPage()
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
    const columns = [
      { title: '#' },
      { title: 'Timeline' },
      { title: 'Method' },
      { title: 'Status' },
      { title: 'URL' },
      { title: 'Scheme', visible: false },
      { title: 'Host', visible: false },
      { title: 'Port', className: 'dt-right', visible: false },
      { title: 'Host+Port', visible: false },
      { title: 'Path', visible: false },
      { title: 'Query', visible: false },
      { title: 'Path+Query', visible: false },
      { title: 'Begin', visible: false },
      { title: 'End', visible: false },
      { title: 'Duration' },
      { title: 'Sending', visible: false },
      { title: 'Processing', visible: false },
      { title: 'Receiving', visible: false },
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
    const lastTimeLine = response[response.length - 1].Timeline
    const durationPrecision = lastTimeLine.startsWith('00:00')
      ? 6 : lastTimeLine.startsWith('00') ? 3 : 0
    const data = response.map(session => [
      session.Number,
      formatDuration(session.Timeline, durationPrecision),
      session.Request.Method,
      session.Response.StatusCode,
      formatURL(session.Request.URL.Full),
      formatScheme(session.Request.URL.Scheme),
      session.Request.URL.Host,
      formatPort(session.Request.URL.Port, session.Request.URL.Scheme),
      session.Request.URL.HostAndPort,
      formatPathOrQuery(session.Request.URL.Path),
      formatPathOrQuery(session.Request.URL.Query),
      formatPathOrQuery(session.Request.URL.PathAndQuery),
      formatTime(session.Timers.ClientBeginRequest),
      formatTime(session.Timers.ClientDoneResponse),
      formatDuration(session.Timers.RequestResponseTime, durationPrecision),
      formatDuration(session.Timers.RequestSendTime, durationPrecision),
      formatDuration(session.Timers.ServerProcessTime, durationPrecision),
      formatDuration(session.Timers.ResponseReceiveTime, durationPrecision),
      session.Response.ContentLength,
      session.Flags.Encoding, session.Flags.Caching, session.Flags.Process
    ])
    resetPage()
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
    resetPage()
    if (title) {
      errorContainer.find('h4').show().text(title)
    } else {
      errorContainer.find('h4').hide()
    }
    errorContainer.show().find('p').text(text)
  }

  function formatURL (url) {
    return shortenString(url, 160)
  }

  function formatScheme (scheme) {
    return scheme || 'N/A'
  }

  function formatPort (port, scheme) {
    return port ? +port : scheme === 'https' ? 443 : 80
  }

  function formatPathOrQuery (path) {
    return shortenString(path, 80)
  }

  function formatDuration (duration, precision) {
    duration = duration.substr(precision)
    let [seconds, milliseconds] = duration.split('.')
    milliseconds = padThousands(Math.round(+milliseconds / 1000))
    return `${seconds}.${milliseconds}`
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

  function shortenString (string, length) {
    const regexp = new RegExp(`^(.{${length}}).*$`)
    return string.replace(regexp, '$1...')
  }

  function padThousands (number) {
    return number > 99 ? number : number > 9 ? '0' + number : '00' + number
  }
}
