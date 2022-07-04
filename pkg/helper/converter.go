package helper

import (
	"encoding/json"
	"strconv"
	"time"
)

const (
	DateTimeFormatDefault = "2006-01-02 15:04:05"
	DateFormatDefault = "2006-01-02"
)


func FloatToString(inputNum float64) string {
	// to convert a float number to a string
	if inputNum != 0 {
		return strconv.FormatFloat(inputNum, 'f', 0, 64)
	} else {
		return "0"
	}
}

func StringNullableToFloat(value *string) float64 {
	if value != nil {
		res, _ := strconv.ParseFloat(*value, 64)
		return res
	}
	return 0
}
func StringToFloat(value string) float64 {
	if value != "" {
		res, _ := strconv.ParseFloat(value, 64)
		return res
	}
	return 0
}
func FloatNullableToString(inputNum *float64) string {
	// to convert a float number to a string
	if inputNum != nil {
		return strconv.FormatFloat(*inputNum, 'f', 0, 64)
	} else {
		return ""
	}
}
func FloatNullableToFloat(value *float64) float64 {
	if value != nil {
		return *value
	}
	return 0
}
func FloatToFloatNullable(value float64) *float64 {
	return &value
}
func DateTimeToDateTimeNullable(value time.Time) *time.Time {
	return &value
}
func DateTimeNullableToDateTime(value *time.Time) time.Time {
	if value == nil {
		return time.Time{}
	}
	return *value
}
func IntToIntNullable(value int) *int {
	return &value
}
func IntNullableToInt64(value *int64) int64 {
	if value == nil {
		return 0
	}
	return *value
}
func IntNullableToInt(value *int) int {
	if value == nil {
		return 0
	}
	return *value
}
func StringToStringNullable(value string) *string {
	return &value
}
func ObjectToString(value interface{}) string {
	result, _ := json.Marshal(value)
	return string(result)
}
func StringNullableToString(value *string) string {
	if value != nil {
		return *value
	}
	return ""
}
func IntNullableToStringNullable(value *int) *string {

	if value != nil {
		result := strconv.Itoa(*value)
		return &result
	}
	return nil
}
func IntNullableToString(value *int) string {

	if value != nil {
		result := strconv.Itoa(*value)
		return result
	}
	return "0"
}

func IntToString(value int) string {

	if value != 0 {
		result := strconv.Itoa(value)
		return result
	}
	return "0"
}
func StringToIntNullable(value string) *int {

	if value != "" {
		result, _ := strconv.Atoi(value)
		return &result
	}
	return nil
}
func Int64NullableToInt(value *int64) int {

	if value != nil {
		result := int(*value)
		return result
	}
	return 0
}
func StringToInt(value string) int {

	if value != "" {
		result, _ := strconv.Atoi(value)
		return result
	}
	return 0
}
func StringNullableToInt(value *string) int {

	if value != nil {
		result, _ := strconv.Atoi(*value)
		return result
	}
	return 0
}
func StringNullableToDateTimeNullable(value *string) *time.Time {
	if value != nil {
		var layoutFormat string
		var date time.Time

		layoutFormat = "2006-01-02 15:04:05"
		date, _ = time.Parse(layoutFormat, *value)
		return &date
	}

	return nil
}

func DateTimeNullableToStringNullable(value *time.Time) *string {
	if value != nil {
		layoutFormat := "2006-01-02 15:04:05"
		date := value.Format(layoutFormat)
		return &date
	}

	return nil
}

func DateTimeToStringNullable(value time.Time) *string {
	layoutFormat := "2006-01-02 15:04:05"
	date := value.Format(layoutFormat)
	return &date
}
func DateTimeToStringWithFormat(value time.Time, format string) string {
	if !value.IsZero() {
		layoutFormat := format
		date := value.Format(layoutFormat)
		return date
	}

	return ""
}
func DateTimeNullableToStringNullableWithFormat(value *time.Time, format string) *string {
	if value != nil {
		layoutFormat := format
		date := value.Format(layoutFormat)
		return &date
	}

	return nil
}

func StringNullableToStringDefaultFormatDate(value *string) *string {
	if value != nil {
		var layoutFormat string
		var date time.Time

		layoutFormat = "2006-01-02T15:04:05Z"
		date, _ = time.Parse(layoutFormat, *value)
		dateString := date.Format(DateTimeFormatDefault)
		return &dateString
	}

	return nil
}
func StringNullableToDateTime(value *string) time.Time {
	if value != nil {
		var layoutFormat string
		var date time.Time
		layoutFormat = "2006-01-02T15:04:05Z"
		date, err := time.Parse(layoutFormat, *value)
		if err != nil {
			return time.Time{}
		}
		return date
	}

	return time.Time{}
}
func StringToDateTimeNullable(value string) time.Time {
	if value != "" {
		var layoutFormat string
		var date time.Time
		layoutFormat = "2006-01-02T15:04:05.999999999Z07:00"
		date, err := time.Parse(layoutFormat, value)
		if err != nil {
			return time.Time{}
		}
		return date
	}

	return time.Time{}
}
func StringToDateWithFormat(value string,format string) time.Time {
	if value != "" {
		var layoutFormat string
		var date time.Time

		layoutFormat = format
		date, _ = time.Parse(layoutFormat, value)
		return date
	}

	return time.Time{}
}
func StringToDate(value string) time.Time {
	if value != "" {
		var layoutFormat string
		var date time.Time

		layoutFormat = DateFormatDefault
		date, _ = time.Parse(layoutFormat, value)
		return date
	}

	return time.Time{}
}
func StringNullableToDateNullable(value *string) *string {
	if value != nil {
		var layoutFormat string
		var date time.Time

		layoutFormat = "20060102"
		date, _ = time.Parse(layoutFormat, *value)
		dateString := date.Format("20060102")
		return &dateString
	}

	return nil
}
func ConvertIntBool(value *int) bool {
	if value != nil {
		if *value == 1 {
			return true
		}
	}
	return false
}



