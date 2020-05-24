import $ from './jquery.js'
import 'popper.js'
import 'bootstrap'
import 'pdfmake'
import JSZip from 'jszip'
import dataTables from 'datatables.net-bs4'
import buttons from 'datatables.net-buttons-bs4'
import columnVisibility from 'datatables.net-buttons/js/buttons.colVis.js'
import buttonsHtml5 from 'datatables.net-buttons/js/buttons.html5.js'
import buttonsPrint from 'datatables.net-buttons/js/buttons.print.js'
import colReorder from 'datatables.net-colreorder-bs4'
import fixedHeader from 'datatables.net-fixedheader-bs4'
import jsonExport from './json-export.js'
import colorfulSessions from './colorful-sessions.js'
import { displayDetails, openSession } from './session-details.js'
import { formatColumnStats } from './footer-formatters.js'
import {
  formatURL, formatScheme, formatPort, formatPathOrQuery, reformatDuration,
  reformatTime, formatSize
} from './data-formatters.js'
import { sazStore } from './saz-store.js'
import { configuration, saveConfiguration } from './configuration.js'
import { prepareTableHelp, showHelp } from './help.js'

let tableWrapper, dataTable, columns

function destroySessionTable () {
  if (dataTable) {
    const oldTable = dataTable
    dataTable = undefined
    oldTable.destroy()
    tableWrapper.html('')
  }
}

function initializeSessionTable () {
  tableWrapper = $('#table-wrapper')
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
    { data: 'ContentType', title: 'Type' },
    {
      data: 'ContentLength',
      title: 'Size',
      className: 'dt-right',
      render: function (data, type, row) {
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
  window.JSZip = JSZip
  dataTables(window, $)
  buttons(window, $)
  columnVisibility(window, $)
  buttonsHtml5(window, $)
  buttonsPrint(window, $)
  colReorder(window, $)
  fixedHeader(window, $)
  jsonExport(window, $)
  colorfulSessions(window, $)
}

function displaySessionTable (sessions) {
  const { columns: columnSettings, order: orderSettings, search } = configuration
  configureColumns(columnSettings)
  const order = convertOrder(orderSettings)
  const data = prepareData(sessions)
  const detailRows = []
  let classes = 'table table-sm table-striped table-hover nowrap compact display'
  if (configuration.colorfulSessions) {
    classes += ' colorful'
  }
  const table = $(`<table class="${classes}">`)
  table.append('<thead></thead>')
  table.append('<tbody></tbody>')
  table.append('<tfoot><tr></tr></tfoot>')
  dataTable = table
    .on('column-visibility.dt', columnVisibilityChanged)
    .on('search.dt', filterChanged)
    .on('order.dt', orderChanged)
    .appendTo(tableWrapper)
    .DataTable({
      columns,
      data,
      order,
      search: { search },
      dom: '<"top"ifB>rt',
      fixedHeader: true,
      colReorder: true,
      paging: false,
      buttons: [
        {
          extend: 'colvis',
          text: '\uea71',
          align: 'button-right',
          className: 'toggles'
        },
        {
          extend: 'collection',
          text: '\uea7d',
          align: 'button-right',
          buttons: ['copy', 'print',
            {
              extend: 'json',
              className: 'divide-at-top'
            }, 'csv', 'excel', 'pdf']
        },
        {
          extend: 'collection',
          text: '\ue994',
          align: 'button-right',
          className: 'toggles',
          buttons: [
            'colorful',
            {
              text: 'Redraw table',
              className: 'divide-at-top',
              action: function (event, dataTable, button, definition) {
                destroySessionTable()
                displaySessionTable(sazStore.loaded.Sessions)
              }
            }
          ]
        }
      ],
      rowCallback,
      footerCallback,
      executeSearch
    })
    .on('click', 'tbody tr td:not([colspan])', function (event) {
      event.preventDefault()
      displayDetails(dataTable, $(this), detailRows, sazStore.loaded)
    })
    .on('click', '[data-session]', function (event) {
      event.preventDefault()
      const link = $(this)
      openSession(sazStore.loaded,
        sazStore.loaded.Sessions[link.data().session - 1], link.attr('download'))
    })
  $('.dataTables_wrapper')
    .find('.dataTables_info')
    .addClass('form-control-sm')
    .attr('data-intro', '#info')
    .end()
    .find('.dataTables_filter')
    .attr('data-intro', '#filter')
    .end()
    .find('.dt-buttons')
    .attr('data-intro', '#buttons')
    .end()
    .find('table thead th[aria-label!="#"]')
    .first()
    .attr('data-intro', '#column')
  prepareTableHelp()
  if (configuration.help.saz !== false) {
    setTimeout(showHelp, 500)
  }
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
    Timeline: reformatDuration(session.Timeline, durationPrecision),
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
    BeginTime: reformatTime(session.Timers.ClientBegin),
    EndTime: reformatTime(session.Timers.ClientDoneResponse),
    Duration: reformatDuration(session.Timers.RequestResponseTime, durationPrecision),
    SendingTime: reformatDuration(session.Timers.RequestSendTime, durationPrecision),
    RespondingTime: reformatDuration(session.Timers.ServerProcessTime, durationPrecision),
    ReceivingTime: reformatDuration(session.Timers.ResponseReceiveTime, durationPrecision),
    ContentType: session.Response.ContentType,
    ContentLength: session.Response.ContentLength,
    Encoding: session.Response.Encoding,
    Caching: session.Response.Caching,
    Process: session.Request.Process
  }))
}

function rowCallback (row, data, displayNum, displayIndex, dataIndex) {
  row = $(row)
  const statusCode = data.StatusCode
  if (statusCode >= 400) {
    return row.addClass('text-danger')
  }
  if (statusCode >= 300) {
    return row.addClass('text-warning')
  }
  if (data.Method === 'CONNECT') {
    return row.addClass('text-muted')
  }
  const contentType = data.ContentType
  if (contentType.startsWith('application/json') || contentType.startsWith('application/xml') ||
      contentType.startsWith('text/xml') || contentType.startsWith('text/plain')) {
    return row.addClass('text-info')
  }
  if (contentType.startsWith('text/css') || contentType.startsWith('text/javascript') ||
      contentType.startsWith('application/javascript') || contentType.startsWith('image/') ||
      contentType.startsWith('font/')) {
    return row.addClass('text-success')
  }
}

function footerCallback (row, data, start, end, display) {
  const table = this.api().table()
  const headers = table.header().querySelectorAll('tr th')
  const foot = $(table.footer()).html('')
  const footerRow = $('<tr>').appendTo(foot)
  if (display.length) {
    for (let i = 0, length = headers.length; i < length; ++i) {
      const header = headers[i]
      const columnIndex = +header.getAttribute('data-column-index')
      const column = columns[columnIndex]
      const footer = header.cloneNode(true)
      footer.innerHTML = `${formatColumnStats(display, column)}
<span>${column.title}</span>`
      footerRow.append(footer)
    }
  }
}

function executeSearch (settings, data, dataString, input) {
  input = input.trim()
  if (!input) return true
  const filters = input
    .toLowerCase()
    .split(/\s+/)
    .map(part => {
      let [column, term] = part.split(':')
      if (!term) {
        term = column
        column = null
      }
      return term.startsWith('-') ? { column, exclude: true, term: term.substr(1) }
        : term.startsWith('+') ? { column, term: term.substr(1) } : { column, term }
    })
  dataString = dataString.trim().toLowerCase()
  for (const filter of filters) {
    const column = filter.column
    const dataContent = column ? getColumnValue(column) : dataString
    const term = filter.term
    if (filter.exclude) {
      if (term.length) {
        if (dataContent.includes(term)) return false
      } else {
        if (!dataContent.length) return false
      }
    } else {
      if (term.length) {
        if (!dataContent.includes(term)) return false
      } else {
        if (dataContent.length) return false
      }
    }
  }
  return true
  function getColumnValue (column) {
    for (const name in data) {
      if (name.toLowerCase() === column) {
        return data[name].toString().toLowerCase()
      }
    }
    return ''
  }
}

function columnVisibilityChanged (event, settings, column, state) {
  if (dataTable) {
    configuration.columns = settings.aoColumns.reduce((columns, column) => {
      columns[column.mData] = { visible: column.bVisible }
      return columns
    }, {})
    saveConfiguration()
  }
}

function filterChanged (event, settings) {
  if (dataTable) {
    var search = settings.oPreviousSearch.sSearch
    if (search !== configuration.search) {
      configuration.search = search
      saveConfiguration()
    }
  }
}

function orderChanged (event, settings, state) {
  if (dataTable) {
    const order = settings.aaSorting[0]
    const newOrder = { column: columns[order[0]].data, descending: order[1] === 'desc' }
    const oldOrder = configuration.order
    if (newOrder.column !== oldOrder.column || newOrder.descending !== oldOrder.descending) {
      configuration.order = newOrder
      saveConfiguration()
    }
  }
}

export { initializeSessionTable, displaySessionTable, destroySessionTable }
