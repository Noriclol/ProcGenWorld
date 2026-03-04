package calendar

import "fmt"

const (
	DaysPerMonth     = 30
	MonthsPerYear    = 12
	DayDurationTicks = 600
	WeekDurationDays = 7
)

type Season int

const (
	Spring Season = iota
	Summer
	Autumn
	Winter
)

func (s Season) String() string {
	return [...]string{"Spring", "Summer", "Autumn", "Winter"}[s]
}

type Date struct {
	Year  int
	Month int
	Day   int
}

func (d Date) Season() Season {
	return Season((d.Month - 1) / 3)
}

func (d Date) String() string {
	return fmt.Sprintf("Year %d  Month %d  Day %d  (%s)", d.Year, d.Month, d.Day, d.Season())
}

type Calendar struct {
	TotalTicks int64
	TotalDays  int
	DayTick    int
}

func (c *Calendar) Tick() (newDay, newWeek bool) {
	c.TotalTicks++
	c.DayTick++
	if c.DayTick >= DayDurationTicks {
		c.DayTick = 0
		c.TotalDays++
		newDay = true
		newWeek = c.TotalDays%WeekDurationDays == 0
	}
	return
}

func (c *Calendar) CurrentDate() Date {
	dayOfYear := c.TotalDays % (DaysPerMonth * MonthsPerYear)
	month := dayOfYear/DaysPerMonth + 1
	day := dayOfYear%DaysPerMonth + 1
	year := c.TotalDays/(DaysPerMonth*MonthsPerYear) + 1
	return Date{Year: year, Month: month, Day: day}
}
