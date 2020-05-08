import $ from 'jquery'
window.jQuery = window.$ = $
import JSZip from 'jszip'
window.JSZip = JSZip
import 'pdfmake'
import dataTable from 'datatables.net-bs4'
dataTable(window, $)
import buttons from 'datatables.net-buttons-bs4'
buttons(window, $)
import columnVisibility from 'datatables.net-buttons/js/buttons.colVis.js'
columnVisibility(window, $)
import buttonsHtml5 from 'datatables.net-buttons/js/buttons.html5.js'
buttonsHtml5(window, $)
import buttonsPrint from 'datatables.net-buttons/js/buttons.print.js'
buttonsPrint(window, $)
import colReorder from 'datatables.net-colreorder-bs4'
colReorder(window, $)
import fixedColumns from 'datatables.net-fixedcolumns-bs4'
fixedColumns(window, $)
import scroller from 'datatables.net-scroller-bs4'
scroller(window, $)
import mimeTypeIcons from './mime-type-icons.js'

let progressWrapper, tableWrapper, infoAlert, errorAlert, table, previousSaz, currentSaz
let columns, hiddenColumns, timerNames, configuration, loadedSaz
const storedSazs = {}

setTimeout(initialize)

function initialize () {
  progressWrapper = $('#progress-wrapper')
  tableWrapper = $('#table-wrapper')
  infoAlert = $('.alert-info')
  errorAlert = $('.alert-danger')
  columns = [
    { data: 'Number', title: '#' },
    { data: 'Timeline', title: 'Timeline' },
    { data: 'Method', title: 'Method' },
    { data: 'StatusCode', title: 'Status' },
    { data: 'URL', title: 'URL' },
    { data: 'Scheme', title: 'Scheme' },
    { data: 'Host', title: 'Host' },
    { data: 'Port', title: 'Port', className: 'dt-right' },
    { data: 'HostAndPort', title: 'Host+Port' },
    { data: 'Path', title: 'Path' },
    { data: 'Query', title: 'Query' },
    { data: 'PathAndQuery', title: 'Path+Query' },
    { data: 'BeginTime', title: 'Begin' },
    { data: 'EndTime', title: 'End' },
    { data: 'Duration', title: 'Duration' },
    { data: 'SendingTime', title: 'Sending' },
    { data: 'RespondingTime', title: 'Responding' },
    { data: 'ReceivingTime', title: 'Receiving' },
    {
      data: 'Size',
      title: 'Size',
      className: 'dt-right',
      render: function ( data, type, row) {
        if (type === 'display' || type === 'filter') {
          return formatSize(data)
        }
        return data
      }
    },
    { data: 'Encoding', title: 'Encoding' },
    { data: 'Caching', title: 'Caching', orderable: false },
    { data: 'Process', title: 'Process' }
  ]
  hiddenColumns = [
    'Scheme', 'Host', 'Port', 'HostAndPort', 'Path', 'Query', 'PathAndQuery',
    'BeginTime', 'EndTime', 'SendingTime', 'RespondingTime', 'ReceivingTime'
  ]
  timerNames = [
    'ClientConnected', 'ClientBeginRequest', 'GotRequestHeaders',
    'ClientDoneRequest', 'GatewayTime', 'DNSTime', 'TCPConnectTime',
    'HTTPSHandshakeTime', 'RequestResponseTime', 'RequestSendTime',
    'ServerProcessTime', 'ResponseReceiveTime', 'ServerConnected',
    'FiddlerBeginRequest', 'ServerGotRequest', 'ServerBeginResponse',
    'GotResponseHeaders', 'ServerDoneResponse', 'ClientBeginResponse',
    'ClientDoneResponse'
  ]
  $('#theme-switcher').on('click', switchTheme)
  currentSaz = $('#saz-file').on('change', selectSaz)
  previousSaz = $('#previous-saz').on('change', sazChanged)
  $(document.documentElement)
    .on('dragenter', preventDefaults)
    .on('dragover', showDropEffect)
    .on('dragleave', preventDefaults)
    .on('drop', dropSaz)
  loadConfiguration()
  progressWrapper.hide()
}

function preventDefaults (event) {
  event.preventDefault()
  event.stopPropagation()
}

function showDropEffect (event) {
  event.originalEvent.dataTransfer.dropEffect = 'copy'
  preventDefaults(event)
}

function dropSaz (event) {
  preventDefaults(event)
  currentSaz.val('')
  processSazs(event.originalEvent.dataTransfer.files)
}

function selectSaz (event) {
  processSazs(event.target.files)
}

function processSazs (files) {
  if (!files.length) {
    resetPage()
    return infoAlert.show()
  }
  progressWrapper.show()
  const promises = []
  for (const file of files) {
    promises.push(uploadSaz(file))
  }
  Promise
    .allSettled(promises)
    .then(results => {
      let response, error
      for (const { value, reason } of results) {
        if (reason) {
          error = reason
          return displayError(reason)
        }
        if (value) {
          response = value
        }
      }
      loadedSaz = storedSazs[response.File.name]
      displaySaz(loadedSaz.Sessions)
      updatePreviousSaz()
    })
    .then(() => progressWrapper.hide())
}

function uploadSaz (file) {
  return postSaz(file).then(storeSaz)
}

function postSaz (file) {
  const formData = new FormData()
  formData.append('saz', file)
  return $.ajax({
    method: 'POST',
    url: '/api/saz',
    data: formData,
    contentType: false,
    processData: false
  }).then (response => ({ File: file, ...response }))
}

function storeSaz (response) {
  loadedSaz = storedSazs[response.File.name] = response
  updatePreviousSaz()
  return response
}

function resetPage () {
  infoAlert.hide()
  errorAlert.hide()
  if (table) {
    const oldTable = table
    table = undefined
    oldTable.destroy()
    tableWrapper.html('')
  }
}

function updatePreviousSaz () {
  setTimeout(() => {
    previousSaz.html('')
    for (const name in storedSazs) {
      const option = $('<option>').text(name)
      if (name === loadedSaz.File.name) {
        option.attr('selected', 'selected')
      }
      previousSaz.append(option)
    }
  })
}

function displaySaz (sessions) {
  const { columns: columnSettings, order: orderSettings, search } = configuration
  configureColumns(columnSettings)
  const order = convertOrder(orderSettings)
  const data = prepareData(sessions)
  resetPage()
  const detailRows = []
  table = $('<table class="table table-sm table-striped table-hover nowrap compact display">')
      .on('column-visibility.dt', columnVisibilityChanged)
      .on('search.dt', filterChanged)
      .on( 'order.dt', orderChanged)
      .appendTo(tableWrapper)
      .DataTable({
        columns,
        data,
        order,
        search: { search },
        dom: '<"top"ifBR>rtS',
        scrollX: true,
        scrollY: '60vh',
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
      .on('click', 'tbody tr td:not([colspan])', function () {
        const tr = $(this).closest('tr')
        const row = table.row(tr)
        const id = tr.attr('id')
        const index = detailRows.indexOf(id)
        if (row.child.isShown()) {
          tr.removeClass('details')
          row.child.hide()
          detailRows.splice(index, 1)
        } else {
          const data = row.data()
          const session = loadedSaz.Sessions[data.Number - 1]
          if (session.Request.Header) {
            showDetails(session, formatDetails)
          } else {
            loadDetails(loadedSaz.Key, session)
              .then(() => showDetails(session, formatDetails))
              .catch(error => showDetails(error, formatError))
          }
        }
        function showDetails (data, format) {
          tr.addClass('details')
          row.child(format(data)).show()
          if (index < 0) {
            detailRows.push(id)
          }
        }
      })
}

function configureColumns (columnSettings) {
  for (const column of columns) {
    const settings = columnSettings[column.data]
    column.visible = settings && settings.visible
  }
}

function convertOrder (orderSettings) {
  const { column: orderColumn, descending } = orderSettings
  return [
    columns.findIndex(({ data }) => data === orderColumn),
    descending ? 'desc' : 'asc'
  ]
}

function prepareData (response) {
  const lastTimeLine = response[response.length - 1].Timeline
  const durationPrecision = lastTimeLine.startsWith('00:00')
    ? 6 : lastTimeLine.startsWith('00') ? 3 : 0
  return response.map(session => ({
    Number: session.Number,
    Timeline: formatDuration(session.Timeline, durationPrecision),
    Method: session.Request.Method,
    StatusCode: session.Response.StatusCode,
    URL: formatURL(session.Request.URL.Full),
    Scheme: formatScheme(session.Request.URL.Scheme),
    Host: session.Request.URL.Host,
    Port: formatPort(session.Request.URL.Port, session.Request.URL.Scheme),
    HostAndPort: session.Request.URL.HostAndPort,
    Path: formatPathOrQuery(session.Request.URL.Path),
    Query: formatPathOrQuery(session.Request.URL.Query),
    PathAndQuery: formatPathOrQuery(session.Request.URL.PathAndQuery),
    BeginTime: formatTime(session.Timers.ClientBeginRequest),
    EndTime: formatTime(session.Timers.ClientDoneResponse),
    Duration: formatDuration(session.Timers.RequestResponseTime, durationPrecision),
    SendingTime: formatDuration(session.Timers.RequestSendTime, durationPrecision),
    RespondingTime: formatDuration(session.Timers.ServerProcessTime, durationPrecision),
    ReceivingTime: formatDuration(session.Timers.ResponseReceiveTime, durationPrecision),
    Size: session.Response.ContentLength,
    Encoding: session.Response.Encoding,
    Caching: session.Response.Caching,
    Process: session.Request.Process
  }))
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

function formatTime (time) {
  return time.substr(14, 12)
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

function loadDetails (key, session) {
  return $.ajax({ url: `/api/saz/${key}/${session.Number}` })
    .then (response => {
      Object.assign(session.Request, response.Request)
      Object.assign(session.Response, response.Response)
    })
}

function formatDetails (session) {
  const header = session.Response.Header
  let contentType = header['Content-Type']
  contentType = contentType ? contentType[0] : 'unknown'
  let mimeType = contentType.replace(/;.*$/, '').replace(/\//g, '-')
  if (!mimeTypeIcons.has(mimeType)) {
    mimeType = 'unknown'
  }
  const requestHeader = formatHeader(session.Request.Header)
  const responseHeader = formatHeader(header)
  const timers = formatTimers(session.Timers)
  const flags = formatFlags(session.Flags)
  return `<div class=media>
  <img class=mr-3 src=png/${mimeType}.png alt=${contentType}>
  <div class=media-body>
    <h5 class=mt-0>${contentType}</h5>
    <h6>Request</h6>
    ${requestHeader}
    <h6>Response</h6>
    ${responseHeader}
    <h6>Timers</h6>
    ${timers}
    <h6>Flags</h6>
    ${flags}
  </div>
</div>`
}

function formatHeader (header) {
  let list = `<ul class="list-group list-group-flush">
`
  for (const name in header) {
    const value = header[name].join(', ')
    list += `    <li class="list-group-item py-1"><strong>${name}</strong>: ${value}</li>
`
  }
  return list + '</ul>'
}

function formatTimers(timers) {
  let list = `<ul class="list-group list-group-flush">
`
  for (const name of timerNames) {
    let value = timers[name]
    if (!name.endsWith('Time')) {
      value = new Date(value).toLocaleString()
    }
    list += `    <li class="list-group-item py-1"><strong>${name}</strong>: ${value}</li>
`
  }
  return list + '</ul>'
}

function formatFlags (flags) {
  let list = `<ul class="list-group list-group-flush">
`
  for (const name in flags) {
    list += `    <li class="list-group-item py-1"><strong>${name}</strong>: ${flags[name]}</li>
`
  }
  return list + '</ul>'
}

function formatError (response) {
  const { title, text } = parseError(response)
  return `<div class="alert alert-warning" role=alert>
  <h4 class=alert-heading>${title}</h4>
  <p>${text}</p>
</div>`
}

function displayError (response) {
  const { title, text } = parseError(response)
  resetPage()
  previousSaz.prop('selectedIndex', -1)
  if (title) {
    errorAlert.find('h4').show().text(title)
  } else {
    errorAlert.find('h4').hide()
  }
  errorAlert.show().find('p').text(text)
}

function parseError (response) {
  let title, text
  if (response instanceof Error) {
    text = response.message
  } else {
    title = response.status && `${response.status} (${response.statusText})`
    text = response.responseText || 'Connection failed.'
  }
  return { title, text }
}

function loadConfiguration () {
  configuration = JSON.parse(localStorage.getItem('prantlf/sazview') || '{}')
  ensureDefaultConfiguration()
}

function ensureDefaultConfiguration () {
  let { columns, order, search } = configuration
  if (columns === undefined) {
    configuration.columns = columns = {}
  }
  for (const column of hiddenColumns) {
    if (columns[column] === undefined) {
      columns[column] = { visible: false }
    }
  }
  if (order === undefined) {
    configuration.order = { column: 'Number' }
  }
  if (search === undefined) {
    configuration.search = ''
  }
}

function saveConfiguration () {
  setTimeout(() =>
    localStorage.setItem('prantlf/sazview', JSON.stringify(configuration)))
}

function columnVisibilityChanged (event, settings, column, state) {
  if (table) {
    configuration.columns = settings.aoColumns.reduce((columns, column) => {
      columns[column.mData] = { visible: column.bVisible }
      return columns
    }, {})
    saveConfiguration()
  }
}

function filterChanged (event, settings) {
  if (table) {
    var search = settings.oPreviousSearch.sSearch
    if (search !== configuration.search) {
      configuration.search = search
      saveConfiguration()
    }
  }
}

function orderChanged (event, settings, state) {
  if (table) {
    const order = settings.aaSorting[0];
    const newOrder = { column: columns[order[0]].data, descending: order[1] === 'desc' }
    oldOrder = configuration.order
    if (newOrder.column !== oldOrder.column || newOrder.descending !== oldOrder.descending) {
      configuration.order = newOrder
      saveConfiguration()
    }
  }
}

function sazChanged () {
  const name = Object.keys(storedSazs)[previousSaz.prop('selectedIndex')]
  const saz = storedSazs[name]
  progressWrapper.show()
  downloadSaz(saz.Key)
    .then(displaySaz)
    .catch(response => {
      if (!(response instanceof Error || response.status !== 404)) {
        throw response
      }
      return uploadSaz(saz.File)
        .then(({ Sessions: sessions }) => displaySaz(sessions))
    })
    .catch(displayError)
    .then(() => progressWrapper.hide())
}

function downloadSaz (key) {
  return $.ajax({ url: `/api/saz/${key}` })
}

function switchTheme () {
  const body = $(document.body)
  body.fadeOut(200, function () {
    $('#theme,#dark-overrides').remove()
    switch (sazTheme) {
      case 'dark': sazTheme = 'light'; break
      case 'light': sazTheme = 'system'; break
      default: sazTheme = 'dark'
    }
    changeTheme()
    ensureDarkOverrides()
    updateThemeSwitcher()
    saveTheme()
    body.fadeIn(200)
  })
}

function saveTheme () {
  setTimeout(() => localStorage.setItem('prantlf/sazview-theme', sazTheme))
}
