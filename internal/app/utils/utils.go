package utils

import (
	"time"

	"github.com/google/uuid"
)

type RequestPaginationFilterData struct {
	Take int
	Skip int
}

type RequestDateTimeFilterData struct {
	Take      int
	Skip      int
	StartDate time.Time
	EndDate   time.Time
	Timezone  int
}

func GetTimezoneOffset(timezone int) time.Duration {
	return time.Duration((timezone / 60) * int(time.Hour))
}

func TimeFormatAsDate(datetime time.Time, timezone int) string {
	format := "2006-01-02"

	if timezone == 0 {
		return datetime.Format(format)
	}

	return datetime.Add(GetTimezoneOffset(timezone)).Format(format)
}

func ValidateUUIDFromString(str string) (*string, error) {
	_, err := uuid.Parse(str)
	if err != nil {
		return nil, err
	}

	return &str, nil
}

func ParsePaginationFilter(take int, skip int) (*RequestPaginationFilterData, error) {
	if take > 100 || take <= 0 {
		take = 100
	}

	if skip < 0 {
		skip = 0
	}

	return &RequestPaginationFilterData{
		Take: take,
		Skip: skip,
	}, nil
}

func ParseRequestDateTimeFilter(take int, skip int, startDate string, endDate string, timezone int) (*RequestDateTimeFilterData, error) {
	parsePagination, err := ParsePaginationFilter(take, skip)
	if err != nil {
		return nil, err
	}
	take = parsePagination.Take
	skip = parsePagination.Skip

	offset := GetTimezoneOffset(timezone)
	format := "2006-01-02"
	currentDate := time.Now().UTC().Add(offset).Format(format)

	if startDate == "" {
		startDate = currentDate
	}

	parseStartDate, err := time.Parse(format, startDate)

	if err != nil {
		return nil, err
	}

	if endDate == "" {
		endDate = parseStartDate.Add(24 * time.Hour).Format(format)
	}

	parseEndDate, err := time.Parse(format, endDate)

	if err != nil {
		parseEndDate = parseStartDate.Add(24 * time.Hour)
	}

	if parseEndDate.Before(parseStartDate) {
		parseEndDate = parseStartDate.Add(24 * time.Hour)
	}

	if parseEndDate.Equal(parseStartDate) {
		parseEndDate = parseStartDate.Add(24 * time.Hour)
	}

	// parseStartDate = parseStartDate.Add(-offset)
	// parseEndDate = parseEndDate.Add(-offset)

	return &RequestDateTimeFilterData{
		Take:      take,
		Skip:      skip,
		StartDate: parseStartDate,
		EndDate:   parseEndDate,
		Timezone:  timezone,
	}, err
}
