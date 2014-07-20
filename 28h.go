package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
	"time"
)

const (
	// TODO(mpl): allow for awake and asleep as parameters, as long as awake + asleep == 28
	h24            = 24 * time.Hour
	h28            = 28 * time.Hour
	awake          = 19 * time.Hour
	asleep         = 9 * time.Hour
	firstDayOfWeek = monday
	saneKitchen    = "15:04"
)

var (
	dayFrom = flag.String("day", "monday", "day to build the week from")
	wake    = flag.String("wake", "08:00", "time to wake up on the day to build the week from")
)

const (
	monday int = iota
	tuesday
	wednesday
	thursday
	friday
	saturday
	sunday
)

var toString = map[int]string{
	monday:    "monday",
	tuesday:   "tuesday",
	wednesday: "wednesday",
	thursday:  "thursday",
	friday:    "friday",
	saturday:  "saturday",
	sunday:    "sunday",
}

func main() {
	flag.Parse()
	wakeTime, err := time.Parse(saneKitchen, *wake)
	if err != nil {
		log.Fatal(err)
	}
	// TODO(mpl): ask about origin != time.Time{}
	origin, err := time.Parse(saneKitchen, "00:00")
	if err != nil {
		log.Fatal(err)
	}
	initDay, err := initialDay(wakeTime, origin)
	if err != nil {
		log.Fatal(err)
	}
	wk := week(initDay, origin)
	wk = shift(wk, firstDayOfWeek)

	for _, v := range wk {
		fmt.Printf("%v	%v\n", toString[v.name][:3], v.String())
	}

}

type weekDay struct {
	name   int
	wakeup time.Time
	bed    time.Time
}

func (w weekDay) String() string {
	if w.wakeup.IsZero() {
		return fmt.Sprintf("up %v bed", w.bed.Format(saneKitchen))
	}
	if w.bed.IsZero() {
		return fmt.Sprintf("bed %v up", w.wakeup.Format(saneKitchen))
	}
	if w.wakeup.Before(w.bed) {
		return fmt.Sprintf("bed %v up %v bed", w.wakeup.Format(saneKitchen), w.bed.Format(saneKitchen))
	}
	return fmt.Sprintf("up %v bed %v up", w.bed.Format(saneKitchen), w.wakeup.Format(saneKitchen))
}

func initialDay(wakeTime, origin time.Time) (weekDay, error) {
	var initial weekDay
	fromString := func(daystring string) int {
		for k, v := range toString {
			if strings.EqualFold(v, daystring) {
				return k
			}
		}
		return -1
	}
	name := fromString(*dayFrom)
	if name == -1 {
		return initial, fmt.Errorf("Invalid day name: %v", *dayFrom)
	}
	yesterday := origin
	if wakeTime.Before(yesterday) {
		return initial, fmt.Errorf("wakeTime is in the past")
	}
	if wakeTime.Equal(yesterday) {
		bedTime := wakeTime.Add(awake)
		wakeTime = bedTime.Add(asleep)
		return weekDay{
			name:   name,
			wakeup: wakeTime,
			bed:    bedTime,
		}, nil
	}
	bedTime := wakeTime.Add(-asleep)
	if bedTime.Before(yesterday) || bedTime.Equal(yesterday) {
		bedTime = wakeTime.Add(awake)
	}
	initial = weekDay{
		name:   name,
		wakeup: wakeTime,
		bed:    bedTime,
	}
	return initial, nil
}

func week(initialDay weekDay, origin time.Time) []weekDay {
	yesterday := origin
	tomorrow := yesterday.Add(h24)
	wk := make([]weekDay, 7)
	wakeup := initialDay.wakeup
	bed := initialDay.bed
	name := initialDay.name
	wakes := []time.Time{wakeup}
	beds := []time.Time{bed}
	for i := 0; i < 7; i++ {
		wakeup = wakeup.Add(h28)
		bed = bed.Add(h28)
		wakes = append(wakes, wakeup)
		beds = append(beds, bed)
	}
	iw, ib := 0, 0
	for k, _ := range wk {
		newDay := weekDay{
			name: name,
		}
		if wakes[iw].After(yesterday) && (wakes[iw].Before(tomorrow) || wakes[iw].Equal(tomorrow)) {
			newDay.wakeup = wakes[iw]
			iw++
		}
		if beds[ib].After(yesterday) && (beds[ib].Before(tomorrow) || beds[ib].Equal(tomorrow)) {
			newDay.bed = beds[ib]
			ib++
		}
		wk[k] = newDay
		name++
		if name == 7 {
			name = 0
		}
		yesterday = yesterday.Add(h24)
		tomorrow = tomorrow.Add(h24)
	}
	return wk
}

func shift(week []weekDay, asFirst int) []weekDay {
	firstPos := 7 - week[0].name + asFirst
	return append(week[firstPos:], week[:firstPos]...)
}
