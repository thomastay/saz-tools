import { parseTime, parseDuration, convertMillisecondsToDuration } from './data-parsers.js'
import { formatParsedTime, formatParsedDuration, formatSize } from './data-formatters.js'
import { sazStore } from './saz-store.js'

function formatTimeStats (data, visibleRows, timer) {
  let minimum = parseTime(data[visibleRows[0]].Timers[timer])
  let maximum = parseTime(data[visibleRows[visibleRows.length - 1]].Timers[timer])
  if (minimum > maximum) {
    const temp = minimum
    minimum = maximum
    maximum = temp
  }
  const difference = convertMillisecondsToDuration(maximum - minimum)
  const precision = difference.hours > 0 ? 0 : difference.minutes > 0 ? 3 : 6
  return `<span>Min: ${formatParsedTime(minimum)}</span><br>
<span>Max: ${formatParsedTime(maximum)}</span><br>
<span>Dif: ${formatParsedDuration(difference, precision)}</span><br>
<br>
<br>`
}

function formatDurationStats (data, visibleRows, timer) {
  const count = visibleRows.length
  let durations = new Array(count)
  for (let i = 0; i < count; ++i) {
    durations[i] = parseDuration(data[visibleRows[i]].Timers[timer])
  }
  durations = durations.sort(compareDurations)
  const minimum = durations[0]
  const maximum = durations[count - 1]
  const total = { hours: 0, minutes: 0, seconds: 0, microseconds: 0 }
  for (const duration of durations) {
    addDuration(total, duration)
  }
  const mean = divideDuration(total, count)
  const medianIndex = Math.floor(count / 2)
  let median
  if (count % 2) {
    median = durations[medianIndex]
  } else {
    const medianTotal = { hours: 0, minutes: 0, seconds: 0, microseconds: 0 }
    addDuration(medianTotal, durations[medianIndex - 1])
    addDuration(medianTotal, durations[medianIndex])
    median = divideDuration(medianTotal, 2)
  }
  const precision = total.hours > 0 ? 0 : total.minutes > 0 ? 3 : 6
  return `<span>Min: ${formatParsedDuration(minimum, precision)}</span><br>
<span>Max: ${formatParsedDuration(maximum, precision)}</span><br>
<span>Avg: ${formatParsedDuration(mean, precision)}</span><br>
<span>Med: ${formatParsedDuration(median, precision)}</span><br>
<span>Tot: ${formatParsedDuration(total, precision)}</span><br>`
}

const columnStatsFormatters = {
  formatNumberStats (data, visibleRows) {
    const count = visibleRows.length
    let minimum = data[visibleRows[0]].Number
    let maximum = data[visibleRows[count - 1]].Number
    if (minimum > maximum) {
      const temp = minimum
      minimum = maximum
      maximum = temp
    }
    return `<span>Min: ${minimum}</span><br>
<span>Max: ${maximum}</span><br>
<span>Cnt: ${count}</span><br>
<br>
<br>`
  },

  formatTimelineStats (data, visibleRows) {
    let minimum = parseDuration(data[visibleRows[0]].Timeline)
    let maximum = parseDuration(data[visibleRows[visibleRows.length - 1]].Timeline)
    if (compareDurations(minimum, maximum) > 0) {
      const temp = minimum
      minimum = maximum
      maximum = temp
    }
    const difference = subtractDuration(maximum, minimum)
    const precision = maximum.hours > 0 ? 0 : maximum.minutes > 0 ? 3 : 6
    return `<span>Min: ${formatParsedDuration(minimum, precision)}</span><br>
<span>Max: ${formatParsedDuration(maximum, precision)}</span><br>
<span>Dif: ${formatParsedDuration(difference, precision)}</span><br>
<br>
<br>`
  },

  formatBeginTimeStats: (data, visibleRows) => formatTimeStats(data, visibleRows, 'ClientBegin'),

  formatEndTimeStats: (data, visibleRows) => formatTimeStats(data, visibleRows, 'ClientDoneResponse'),

  formatDurationStats: (data, visibleRows) => formatDurationStats(data, visibleRows, 'RequestResponseTime'),

  formatSendingTimeStats: (data, visibleRows) => formatDurationStats(data, visibleRows, 'RequestSendTime'),

  formatRespondingTimeStats: (data, visibleRows) => formatDurationStats(data, visibleRows, 'ServerProcessTime'),

  formatReceivingTimeStats: (data, visibleRows) => formatDurationStats(data, visibleRows, 'ResponseReceiveTime'),

  formatContentLengthStats (data, visibleRows) {
    const count = visibleRows.length
    let sizes = new Array(count)
    for (let i = 0; i < count; ++i) {
      sizes[i] = data[visibleRows[i]].Response.ContentLength || 0
    }
    sizes = sizes.sort((left, right) => left - right)
    const minimum = sizes[0]
    const maximum = sizes[count - 1]
    let total = 0
    for (const size of sizes) {
      total += size
    }
    const mean = Math.round(total / count)
    const medianIndex = Math.floor(count / 2)
    let median
    if (count % 2) {
      median = sizes[medianIndex]
    } else {
      median = Math.round((sizes[medianIndex - 1] + sizes[medianIndex]) / 2)
    }
    return `<span>Min: ${formatSize(minimum)}</span><br>
<span>Max: ${formatSize(maximum)}</span><br>
<span>Avg: ${formatSize(mean)}</span><br>
<span>Med: ${formatSize(median)}</span><br>
<span>Tot: ${formatSize(total)}</span><br>`
  }
}

function formatColumnStats (visibleRows, column) {
  const data = sazStore.loaded.Sessions
  const formatter = columnStatsFormatters[`format${column.data}Stats`]
  return (formatter && formatter(data, visibleRows, column)) || `<br>
<br>
<br>
<br>
<br>`
}

function compareDurations (left, right) {
  return left.hours < right.hours ? -1 : left.hours > right.hours ? 1
    : left.minutes < right.minutes ? -1 : left.minutes > right.minutes ? 1
      : left.seconds < right.seconds ? -1 : left.seconds > right.seconds ? 1
        : left.microseconds < right.microseconds ? -1 : left.microseconds > right.microseconds ? 1 : 0
}

function divideDuration ({ hours, minutes, seconds, microseconds }, divider) {
  let duration = Math.floor((hours * 60 * 60 * 1000000 + minutes * 60 * 1000000 +
    seconds * 1000000 + microseconds) / divider)
  hours = Math.floor(duration / 1000000 / 60 / 60)
  duration -= hours * 60 * 60 * 1000000
  minutes = Math.floor(duration / 1000000 / 60)
  duration -= minutes * 60 * 1000000
  seconds = Math.floor(duration / 1000000)
  microseconds = duration - seconds * 1000000
  return { hours, minutes, seconds, microseconds }
}

function addDuration (total, part) {
  total.hours += part.hours
  total.minutes += part.minutes
  total.seconds += part.seconds
  total.microseconds += part.microseconds
  if (total.microseconds > 999999) {
    total.microseconds -= 1000000
    ++total.seconds
  }
  if (total.seconds > 59) {
    total.seconds -= 60
    ++total.minutes
  }
  if (total.minutes > 59) {
    total.minutes -= 60
    ++total.hours
  }
}

function subtractDuration (end, begin) {
  let { hours, minutes, seconds, microseconds } = {
    hours: end.hours - begin.hours,
    minutes: end.minutes - begin.minutes,
    seconds: end.seconds - begin.seconds,
    microseconds: end.microseconds - begin.microseconds
  }
  if (microseconds < 0) {
    microseconds = 999999 + microseconds
    --seconds
  }
  if (seconds < 0) {
    seconds = 60 + seconds
    --minutes
  }
  if (minutes < 0) {
    minutes = 60 + minutes
    --hours
  }
  return { hours, minutes, seconds, microseconds }
}

export { formatColumnStats }
