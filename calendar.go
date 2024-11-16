package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// setup would normally be an init() function, however, there seems
// to be something awry with the testing framework when we set the
// global Logger from an init()
func setup() {

	zerolog.TimeFieldFormat = ""

	zerolog.TimestampFunc = func() time.Time {
		return time.Date(2008, 1, 8, 17, 5, 05, 0, time.UTC)
	}
	log.Logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
}

func main() {

	setup()
	log.Info().Msg("hello world")

	showOneMonth, showThreeMonths, months, weekNumbering := false, false, 0, false

	// Access specific arguments
	if len(os.Args) > 1 {
		for i, arg := range os.Args[1:] {
			fmt.Printf("Arg %d: %s\n", i+1, arg)
			if os.Args[i+1] == "-1" { showOneMonth = true }
			if os.Args[i+1] == "-3" { showThreeMonths = true }
			if os.Args[i+1] == "--months" { months, _ = strconv.Atoi(os.Args[i+2]) }
			if os.Args[i+1] == "--week-numbering" { weekNumbering = true }
		}
	} else {
		log.Info().Msg("No arguments provided.")
	}

	calendarToShow := make(YearAsAnSlideOfMonthsType, MonthsPerYear)

	if showOneMonth {
		calendarToShow[time.Now().Month()] = currentYearCalendar[time.Now().Month()]
	} else if showThreeMonths {
		calendarToShow[time.Now().Month()] = currentYearCalendar[time.Now().Month()]
		calendarToShow[time.Now().Month()+1] = currentYearCalendar[time.Now().Month()+1]
		// calendarToShow[time.Now().Month()+2] = currentYearCalendar[time.Now().Month()+2]
	} else if months != 0 {
		for i := 0; i < months; i++ {
			// TODO: this is pending to fix
			calendarToShow[time.Now().Month()] = currentYearCalendar[time.Now().Month()]
		}
	} else if weekNumbering {
		fmt.Println(WeekHeader)
		fmt.Println(currentYearCalendar[time.Now().Month()].String())
	}
	fmt.Println(calendarToShow.String())
}

const WeekHeader = "Su Mo Tu We Th Fr Sa"
const NumberOfDaysPerWeek = 7
const MonthsPerYear = 13

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
type YearAsAnSlideOfMonthsType []MonthAsAMapOfWeeks

func (year YearAsAnSlideOfMonthsType) String() string {
	returnString := ""
	for monthIndex := 0; monthIndex < len(year); monthIndex++ {
		if year[monthIndex] == nil {continue}
		returnString += "    " + time.Month(monthIndex).String() + " " + strconv.Itoa(currentYearNumber) + "\n"
		returnString += " " + WeekHeader + "\n"
		returnString += year[monthIndex].String()
	}
	return returnString
}

var currentYearNumber int
var currentYearCalendar YearAsAnSlideOfMonthsType

// Initialize calendar of current year
func init()	{
	fmt.Println("Initializing calendar...")

	currentYearNumber = time.Now().Year()
	currentYearCalendar = make(YearAsAnSlideOfMonthsType, MonthsPerYear)

	// For each month in a year
	for month := time.January; month <= time.December; month++ {
		fmt.Println("Initializing month: ", month, "...")
		currentYearCalendar[month] = MonthAsAMapOfWeeks{}

		// Gets the week day of the first day of the month
		firstWeekDayInMonth := time.Date(currentYearNumber, month, 1, 0, 0, 0, 0, time.UTC).Weekday()
		log.Debug().Msg("First week day is: " + firstWeekDayInMonth.String())

		numberOfDaysInMonth := GetDaysInMonth(currentYearNumber, month)
		log.Debug().Msg("Number of days in month is: " + strconv.Itoa(numberOfDaysInMonth))

		// Fills the days of the month
		weekIndex := 0
		currentYearCalendar[month][weekIndex] = WeekAsAMapOfWeekDaysType{}
		weekDayIndex := firstWeekDayInMonth
		for dayNumber := 1; dayNumber <= numberOfDaysInMonth; dayNumber++ {
			log.Debug().Msg("Initializing month number and day number: " + month.String() + strconv.Itoa(dayNumber) + "...")
			currentYearCalendar[month][weekIndex][weekDayIndex] = dayNumber
			if weekDayIndex == 6 {
				weekDayIndex = 0
				weekIndex++
				currentYearCalendar[month][weekIndex] = WeekAsAMapOfWeekDaysType{}
			} else {
				weekDayIndex++
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