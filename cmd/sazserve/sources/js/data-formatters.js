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

function reformatDuration (duration, precision) {
  duration = duration.substr(precision)
  let [seconds, milliseconds] = duration.split('.')
  milliseconds = padThousands(Math.round(+milliseconds / 1000))
  return `${seconds}.${milliseconds}`
}

function formatParsedDuration (duration, precision) {
  const seconds = `${padHundreds(duration.hours)}:${padHundreds(duration.minutes)}:${padHundreds(duration.seconds)}`
  const milliseconds = padThousands(Math.round(duration.microseconds / 1000))
  return `${seconds.substr(precision)}.${milliseconds}`
}

function reformatTime (time) {
  return time.substr(14, 9)
}

function formatParsedTime (time) {
  const timeOnly = time.toTimeString()
  return `${timeOnly.substr(0, 8)}.${padThousands(time.getMilliseconds())}`
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

function padHundreds (number) {
  return number > 9 ? number : '0' + number
}

export {
  formatURL, formatScheme, formatPort, formatPathOrQuery, reformatDuration,
  formatParsedDuration, reformatTime, formatParsedTime, formatSize,
  padThousands, padHundreds
}
