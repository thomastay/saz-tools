package sazdumper

import (
	"strconv"
	"strings"
	"time"
)

func parseTime(dateTime string) (time.Time, error) {
	return time.Parse(time.RFC3339Nano, dateTime)
}

func parseDuration(duration string) (time.Duration, error) {
	duration = strings.Replace(duration, ":", "h", 1)
	duration = strings.Replace(duration, ":", "m", 1)
	duration = strings.Replace(duration, ".", "s", 1)
	return time.ParseDuration(duration + "us")
}

func formatTime(dateTime time.Time) string {
	return dateTime.Format("15:04:05.000")
}

func formatDuration(duration time.Duration) string {
	var wholeHour time.Time
	timeInHour := wholeHour.Add(duration)
	return timeInHour.Format("04:05.000")
}

func formatSize(size int) string {
	if size < 0 {
		return "N/A"
	}
	input := strconv.Itoa(size)
	inputLength := len(input)
	numberOfDigits := inputLength
	numberOfCommas := (numberOfDigits - 1) / 3
	if numberOfCommas == 0 {
		return input
	}
	outputLength := inputLength + numberOfCommas
	output := make([]byte, outputLength)
	for inputIndex, outputIndex, indexInGroup := inputLength-1, outputLength-1, 0; ; {
		output[outputIndex] = input[inputIndex]
		if inputIndex == 0 {
			return string(output)
		}
		if indexInGroup++; indexInGroup == 3 {
			outputIndex--
			indexInGroup = 0
			output[outputIndex] = ','
		}
		inputIndex--
		outputIndex--
	}
}
