package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// setup would normally be an init() function, however, there seems
// to be something awry with the testing framework when we set the
// global Logger from an init()
func setupLog() {

	zerolog.TimeFieldFormat = ""

	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func main() {

	setupLog()
	log.Info().Msg("Calendar starting ...")

	log.Debug().Msg(fmt.Sprintf("Go-syntax: %#v", calendar))

	numberOfMonths, year, weekNumbering :=  0, false, false

	// Access specific arguments
	log.Info().Msg("Reading command line arguments...")
	if len(os.Args) > 1 {
		for i, arg := range os.Args[1:] {
			fmt.Printf("Arg %d: %s\n", i+1, arg)
			if os.Args[i+1] == "-1" { numberOfMonths = 1 }
			if os.Args[i+1] == "-3" { numberOfMonths = 3 }
			if os.Args[i+1] == "--months" { numberOfMonths, _ = strconv.Atoi(os.Args[i+2]) }
			if os.Args[i+1] == "--week-numbering" { weekNumbering = true }
			if os.Args[i+1] == "--year" { year = true }
		}
	} else {
		log.Info().Msg("No arguments provided.")
	}

	// Initialize calendar to show
	log.Debug().Msg("Initializing calendarToShow ...")
	var calendarToShow = make(calendarType)
	calendarToShow[time.Now().Year()] = make(YearAsAMapOfMonthsType, MonthsPerYear)
	calendarToShow[time.Now().Year() + 1] = make(YearAsAMapOfMonthsType, MonthsPerYear)

	log.Debug().Msg(fmt.Sprintf("Go-syntax: %#v", calendarToShow))

	if numberOfMonths != 0 {
		for i := 0; i < numberOfMonths; i++ {
			currentTime := time.Now()
			nextMonthTime := currentTime.AddDate(0, i, 0)
			calendarToShow[nextMonthTime.Year()][nextMonthTime.Month()] = calendar[nextMonthTime.Year()][nextMonthTime.Month()]
		}
	} else if year {
		currentYear := time.Now().Year()
		for month := time.January; month <= time.December; month++ {
			calendarToShow[currentYear][month] = calendar[currentYear][month]
		}		
	} else if weekNumbering {
		// TODO
		fmt.Println(WeekHeader)
	}	
	
	log.Debug().Msg(fmt.Sprintf("Go-syntax: %#v", calendarToShow))
	fmt.Println(calendarToShow.String())
}

const WeekHeader = "Su Mo Tu We Th Fr Sa"
const NumberOfDaysPerWeek = 7
const MonthsPerYear = 13
const MaxNumberOfWeeksPerMonth = 6

// Array of days per week
type WeekAsAMapOfWeekDaysType map[time.Weekday]int

func (week WeekAsAMapOfWeekDaysType) String() string {
	returnString := ""
	for dayIndex := time.Sunday; dayIndex <= time.Saturday; dayIndex++ {
		if week[dayIndex] == 0 {
			returnString += "   "
		} else {
			returnString += fmt.Sprintf("%2d ", week[dayIndex])
		}
	}
	return returnString
}

// Map of weeks per month
type MonthAsAMapOfWeeks map[int]WeekAsAMapOfWeekDaysType

func (month MonthAsAMapOfWeeks) String() string {
	returnString := ""
	for weekIndex := 0; weekIndex < len(month); weekIndex++ {
		returnString += month[weekIndex].String() + "\n"
	}
	return returnString
}

// Array of months in a year
type YearAsAMapOfMonthsType map[time.Month]MonthAsAMapOfWeeks

func (year YearAsAMapOfMonthsType) String() string {
	returnString := ""
	for monthIndex := time.January; monthIndex <= time.December; monthIndex++ {
		if year[monthIndex] == nil {continue}
		returnString += "    " + time.Month(monthIndex).String() + "\n"
		returnString += " " + WeekHeader + "\n"
		returnString += year[monthIndex].String()
	}
	return returnString
}

type calendarType map[int]YearAsAMapOfMonthsType

func (calendar calendarType) String() string {

	// When iterating over a map using a for range loop in Go, the order of iteration is random and not guaranteed.
	// The Go runtime intentionally randomizes the iteration order over maps for security and performance reasons, so you should not rely on any specific order.

	// Extract and sort calendar keys to iterate in order
	calendarKeys := make([]int, 0, len(calendar))
	for key := range calendar {
		calendarKeys = append(calendarKeys, key)
	}
	sort.Ints(calendarKeys)

	returnString := ""

	// Prints month + year header
	for weekIndex := 0; weekIndex < MaxNumberOfWeeksPerMonth; weekIndex++ {
		for _, yearNumber := range calendarKeys {
			yearMapOfMonths := calendar[yearNumber]
			for monthIndex := time.January; monthIndex <= time.December; monthIndex++ {
				if yearMapOfMonths[monthIndex] == nil {continue}
				if yearMapOfMonths[monthIndex][weekIndex] == nil {continue}
				if weekIndex == 0 {
					auxString := " "
					auxString += time.Month(monthIndex).String()
					auxString += " "
					auxString += strconv.Itoa(yearNumber)
					returnString += centerString(auxString, len(WeekHeader))
					returnString += "   "
				}
			}
		}
	}
	returnString += "\n"

	// Prints week headers
	for weekIndex := 0; weekIndex < MaxNumberOfWeeksPerMonth; weekIndex++ {
		for _, yearNumber := range calendarKeys {
			yearMapOfMonths := calendar[yearNumber]
			for monthIndex := time.January; monthIndex <= time.December; monthIndex++ {
				if yearMapOfMonths[monthIndex] == nil {continue}
				if yearMapOfMonths[monthIndex][weekIndex] == nil {continue}
				if weekIndex == 0 {
					returnString += " "
					returnString += WeekHeader
					returnString += "  "
				}
			}
		}
	}
	returnString += "\n"
	
	// Print weeks
	for weekIndex := 0; weekIndex < MaxNumberOfWeeksPerMonth; weekIndex++ {
		for _, yearNumber := range calendarKeys {
			yearMapOfMonths := calendar[yearNumber]
			for monthIndex := time.January; monthIndex <= time.December; monthIndex++ {
				if yearMapOfMonths[monthIndex] == nil {continue}
				if yearMapOfMonths[monthIndex][weekIndex] == nil {continue}
				returnString += " "
				returnString += yearMapOfMonths[monthIndex][weekIndex].String()
				returnString += " "
			}
		}
		returnString += "\n"
	}
	return returnString
}

var calendar calendarType

// Initialize calendar of current year
func init()	{
	log.Debug().Msg("Initializing calendar ...")

	calendar = make(map[int]YearAsAMapOfMonthsType)

	// For two years in calendar
	for yearNumber := time.Now().Year(); yearNumber <= time.Now().Year() + 1; yearNumber++ {
		log.Debug().Msg("Initializing year: " + strconv.Itoa(yearNumber) + " ...")
		calendar[yearNumber] = YearAsAMapOfMonthsType{}

		// For each month in a year
		for monthIndex := time.January; monthIndex <= time.December; monthIndex++ {
			log.Debug().Msg("Initializing month: " + time.Month(monthIndex).String() + "...")
			calendar[yearNumber][monthIndex] = MonthAsAMapOfWeeks{}

			// Gets the week day of the first day of the month
			firstWeekDayInMonth := time.Date(yearNumber, time.Month(monthIndex), 1, 0, 0, 0, 0, time.UTC).Weekday()
			log.Debug().Msg("First week day is: " + firstWeekDayInMonth.String())

			numberOfDaysInMonth := GetDaysInMonth(yearNumber, time.Month(monthIndex))
			log.Debug().Msg("Number of days in month is: " + strconv.Itoa(numberOfDaysInMonth))

			// Fills the days of the month
			weekIndex := 0
			calendar[yearNumber][monthIndex][weekIndex] = WeekAsAMapOfWeekDaysType{}
			weekDayIndex := firstWeekDayInMonth
			for dayNumber := 1; dayNumber <= numberOfDaysInMonth; dayNumber++ {
				log.Debug().Msg("Initializing day ... month number and day number: " + time.Month(monthIndex).String() + " " + strconv.Itoa(dayNumber) + " ...")
				log.Debug().Msg("- Month: " + time.Month(monthIndex).String())
				log.Debug().Msg("- Day number: "  + strconv.Itoa(dayNumber))
				log.Debug().Msg("- Week day  : " + weekDayIndex.String())
				calendar[yearNumber][monthIndex][weekIndex][weekDayIndex] = dayNumber
				if (weekDayIndex != time.Saturday) {
					weekDayIndex++
				} else {
					weekDayIndex = time.Sunday
					weekIndex++
					calendar[yearNumber][monthIndex][weekIndex] = WeekAsAMapOfWeekDaysType{}
				}
			}
		}
	}
}


// GetDaysInMonth calculates the number of days in a given month and year
func GetDaysInMonth(year int, month time.Month) int {
	// Create a date at the start of the given month
	firstDayOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)

	// Add one month, then subtract one day to get the last day of the month
	lastDayOfMonth := firstDayOfMonth.AddDate(0, 1, -1)

	return lastDayOfMonth.Day()
}

func centerString(s string, width int) string {
	if len(s) >= width {
		// If the string is already as wide or wider than the specified width, return it as is
		return s
	}

	// Calculate padding
	totalPadding := width - len(s)
	leftPadding := totalPadding / 2
	rightPadding := totalPadding - leftPadding

	// Create the centered string
	return fmt.Sprintf("%s%s%s", strings.Repeat(" ", leftPadding), s, strings.Repeat(" ", rightPadding))
}