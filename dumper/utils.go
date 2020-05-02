package sazdumper

import (
	"strconv"
	"time"
)

func parseTime(dateTime string) (time.Time, error) {
	return time.Parse(time.RFC3339, dateTime)
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
