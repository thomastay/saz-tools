/* global Blob */

import $ from './jquery.js'
import mimeTypeIcons from './mime-type-icons.js'
import { configuration } from './configuration.js'
import { uploadSaz } from './saz-store.js'
import { startProgress, stopProgress } from './progress.js'
import { parseError } from './data-parsers.js'
import { padThousands } from './data-formatters.js'

function loadDetails (key, loadedSaz, session) {
  const url = `/api/saz/${key}/${session.Number}?scope=extras`
  return $.ajax({ url })
    .catch(response => {
      if (response instanceof Error || response.status !== 404) {
        throw response
      }
      return uploadSaz(loadedSaz.File).then(() => $.ajax({ url }))
    })
    .then(response => {
      Object.assign(session.Request, response.Request)
      Object.assign(session.Response, response.Response)
      Object.assign(session.Timers, response.Timers)
      session.Flags = response.Flags
    })
}

function formatDetails (session, key) {
  const header = session.Response.Header
  const contentType = getHeader(header, 'Content-Type') || 'unknown'
  let mimeType = contentType.replace(/;.*$/, '').replace(/\//g, '-')
  const requestHeader = formatHeader(session.Request.Header)
  const responseHeader = formatHeader(header)
  const sessionSummary = formatSessionSummary(session)
  const sessionTimers = formatSessionTimers(session.Timers)
  const sessionFlags = formatSessionFlags(session.Flags)
  const sessionUrl = `/api/saz/${key}/${session.Number}/`
  let requestLinks = ''
  let responseLinks = ''
  if (session.Request.ContentLength > 0) {
    const requestUrl = `${sessionUrl}request/body`
    const requestName = mimeType === 'application-json' ? 'body.json'
      : mimeType === 'application-xml' || mimeType === 'text-xml' ? 'body.xml'
        : 'body.txt'
    requestLinks = `Request Body:
    <a href=${requestUrl} target=_blank>Open</a>
    <a href=${requestUrl} download=${requestName}>Download</a>
    &nbsp;&nbsp;
`
  }
  if (session.Response.ContentLength > 0) {
    const responseUrl = `${sessionUrl}response/body`
    const responseName = session.Request.URL.Path.replace(/^.+\/([^/]+)$/, '$1')
    responseLinks = `Response Body:
    <a href=${responseUrl} target=_blank>Open</a>
    <a href=${responseUrl} download=${responseName}>Download</a>
    &nbsp;&nbsp;
`
  }
  if (!mimeTypeIcons.has(mimeType)) {
    mimeType = 'unknown'
  }
  const dummyUrl = 'javascript:void 0'
  return `<div class=media>
  <img class=mr-3 src=png/${mimeType}.png alt=${contentType}>
  <div class=media-body>
    <h5 class=mt-0>Type ${contentType}</h5>
    <p>
      Session Details:
      <a href=${dummyUrl} data-session=${session.Number} target=_blank>Open</a>
      <a href=${dummyUrl} data-session=${session.Number} download=details-${session.Number}.json>Download</a>
      &nbsp;&nbsp;
      ${requestLinks}${responseLinks}
    </p>
    <h5>Session Summary</h5>
    ${sessionSummary}
    <h5>Request Header</h5>
    ${requestHeader}
    <h5>Response Header</h5>
    ${responseHeader}
    <h5>Session Timers</h5>
    ${sessionTimers}
    <h5>Session Flags</h5>
    ${sessionFlags}
  </div>
</div>`
}

function getHeader (header, name) {
  const value = header[name]
  return Array.isArray(value) ? value.join(', ') : value || ''
}

function formatHeader (header) {
  let list = `<ul class="list-group list-group-flush">
`
  for (const name in header) {
    const value = getHeader(header, name)
    list += `    <li class="list-group-item py-1"><strong>${name}</strong>: ${value}</li>
`
  }
  return list + '</ul>'
}

function formatSessionSummary (session) {
  let scheme = session.Request.URL.Scheme
  scheme = scheme ? `${scheme}:` : ''
  const values = {
    Method: session.Request.Method,
    Origin: `${scheme}//${session.Request.URL.HostAndPort}`,
    Path: session.Request.URL.Path,
    Query: session.Request.URL.Query,
    Status: session.Response.StatusCode
  }
  let list = `<ul class="list-group list-group-flush">
`
  for (const name in values) {
    const value = values[name]
    list += `    <li class="list-group-item py-1"><strong>${name}</strong>: ${value}</li>
`
  }
  return list + '</ul>'
}

function formatSessionTimers (timers) {
  let list = `<ul class="list-group list-group-flush">
`
  for (const name of configuration.timerNames) {
    let value = timers[name]
    if (!name.endsWith('Time')) {
      const date = new Date(value)
      value = `${date.toLocaleString()}.${padThousands(date.getMilliseconds())}`
    }
    list += `    <li class="list-group-item py-1"><strong>${name}</strong>: ${value}</li>
`
  }
  return list + '</ul>'
}

function formatSessionFlags (flags) {
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

function displayDetails (dataTable, td, detailRows, loadedSaz) {
  const tr = td.closest('tr')
  const row = dataTable.row(tr)
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
      showDetails(() => formatDetails(session, loadedSaz.Key))
    } else {
      startProgress()
      loadDetails(loadedSaz.Key, loadedSaz, session)
        .then(() => showDetails(() => formatDetails(session, loadedSaz.Key)))
        .catch(error => showDetails(() => formatError(error)))
        .then(stopProgress)
    }
  }
  function showDetails (formatDetails) {
    tr.addClass('details')
    row.child(formatDetails()).show()
    if (index < 0) {
      detailRows.push(id)
    }
  }
}

function openSession (loadedSaz, session, fileName) {
  const url = `/api/saz/${loadedSaz.Key}/${session.Number}?scope=extras`
  return $.ajax({ url, method: 'HEAD' })
    .catch(response => {
      if (response instanceof Error || response.status !== 404) {
        throw response
      }
      startProgress()
      return uploadSaz(loadedSaz.File)
    })
    .then(() => {
      stopProgress()
      const output = JSON.stringify(session, undefined, 2)
      const file = new Blob([output], { type: 'application/json' })
      if (fileName) {
        $.fn.dataTable.fileSave(file, fileName, true)
      } else {
        const fileURL = URL.createObjectURL(file)
        window.open(fileURL)
        setTimeout(() => URL.revokeObjectURL(fileURL))
      }
    })
    .catch(response => {
      stopProgress()
      const { title, text } = parseError(response)
      $('#error-alert-title').text(title || 'Error')
      $('#error-alert-body').text(text)
      $('#error-alert').modal()
    })
}

export { displayDetails, openSession }
