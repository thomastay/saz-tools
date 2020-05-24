function parseTime (time) {
  return new Date(time)
}

function parseDuration (duration) {
  const [rest, microseconds] = duration.split('.')
  const [hours, minutes, seconds] = rest.split(':')
  return {
    hours: +hours, minutes: +minutes, seconds: +seconds, microseconds: +microseconds
  }
}

function convertMillisecondsToDuration (duration) {
  const hours = Math.floor(duration / 1000 / 60 / 60)
  duration -= hours * 60 * 60 * 1000
  const minutes = Math.floor(duration / 1000 / 60)
  duration -= minutes * 60 * 1000
  const seconds = Math.floor(duration / 1000)
  const microseconds = duration - seconds * 1000
  return { hours, minutes, seconds, microseconds }
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

export { parseTime, parseDuration, convertMillisecondsToDuration, parseError }
