/* global localStorage */

const configuration = {}

function loadConfiguration () {
  Object.assign(configuration, JSON.parse(localStorage.getItem('prantlf/sazview') || '{}'))
  initializeFixedSettings()
  ensureDefaultConfiguration()
}

function initializeFixedSettings () {
  configuration.hiddenColumns = [
    'Scheme', 'Host', 'Port', 'HostAndPort', 'Path', 'Query', 'PathAndQuery',
    'BeginTime', 'EndTime', 'SendingTime', 'RespondingTime', 'ReceivingTime',
    'ContentType'
  ]
  configuration.timerNames = [
    'ClientConnected', 'ClientBegin', 'ClientBeginRequest', 'GotRequestHeaders',
    'ClientDoneRequest', 'GatewayTime', 'DNSTime', 'TCPConnectTime',
    'HTTPSHandshakeTime', 'RequestResponseTime', 'RequestSendTime',
    'ServerProcessTime', 'ResponseReceiveTime', 'ServerConnected',
    'FiddlerBeginRequest', 'ServerGotRequest', 'ServerBeginResponse',
    'GotResponseHeaders', 'ServerDoneResponse', 'ClientBeginResponse',
    'ClientDoneResponse'
  ]
}

function ensureDefaultConfiguration () {
  let { columns, order, search, help, colorfulSessions } = configuration
  if (columns === undefined) {
    configuration.columns = columns = {}
  }
  for (const column of configuration.hiddenColumns) {
    if (columns[column] === undefined) {
      columns[column] = { visible: false }
    }
  }
  if (order === undefined) {
    if (columns.Number && columns.Number.visible !== false) {
      configuration.order = { column: 'Number' }
    } else {
      configuration.order = {}
    }
  }
  if (search === undefined) {
    configuration.search = ''
  }
  if (help === undefined) {
    configuration.help = {}
  }
  if (colorfulSessions === undefined) {
    configuration.colorfulSessions = true
  }
}

function saveConfiguration () {
  setTimeout(() =>
    localStorage.setItem('prantlf/sazview', JSON.stringify(configuration)))
}

export { configuration, loadConfiguration, saveConfiguration }
