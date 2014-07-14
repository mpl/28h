package main

import (
	"flag"
	"fmt"
	"log"
	"time"
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

const (
	h24            = 24 * time.Hour
	h28            = 28 * time.Hour
	awake          = 19 * time.Hour
	asleep         = 9 * time.Hour
	firstDayOfWeek = monday
	saneKitchen    = "15:04"
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

func weekFrom(first weekDay) []weekDay {
	// TODO(mpl): report that bug
	//	yesterday := time.Time{}
	yesterday, err := time.Parse(time.Kitchen, "00:00AM")
	if err != nil {
		log.Fatal(err)
	}
	tomorrow := yesterday.Add(h24)
	wk := make([]weekDay, 7)
	wakeup := first.wakeup
	bed := first.bed
	name := first.name
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
		println()
		println(wakes[iw].Format(time.RFC3339))
		println(yesterday.Format(time.RFC3339))
		if wakes[iw].After(yesterday) && (wakes[iw].Before(tomorrow) || wakes[iw].Equal(tomorrow)) {
			newDay.wakeup = wakes[iw]
			iw++
		}
		println(beds[ib].Format(time.RFC3339))
		println(tomorrow.Format(time.RFC3339))
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

func shift(week []weekDay, first int) []weekDay {
	firstPos := 7 - week[0].name + first
	return append(week[firstPos:], week[:firstPos]...)
}

func week(from weekDay) {
	wk := weekFrom(from)
	wk = shift(wk, firstDayOfWeek)
}

var (
	dayFrom = flag.String("day", "monday", "day to build the week from")
	// Kitchen     = "3:04PM
	// TODO(mpl): use a time
	wake = flag.String("wake", "17:00", "time to wake up on the day to build the week from")
)

func main() {
	flag.Parse()
	wakeTime, err := time.Parse(saneKitchen, *wake)
	if err != nil {
		log.Fatal(err)
	}
	fromString := func(daystring string) int {
		for k, v := range toString {
			if v == daystring {
				return k
			}
		}
		return -1
	}
	from := weekDay{
		// TODO(mpl): safe func. case insensitive.
		name:   fromString(*dayFrom),
		wakeup: wakeTime,
	}
	yesterday, err := time.Parse(time.Kitchen, "00:00AM")
	if err != nil {
		log.Fatal(err)
	}
	bedTime := wakeTime.Add(-asleep)
	if bedTime.Before(yesterday) {
		bedTime = wakeTime.Add(awake)
	}
	from.bed = bedTime

	wk := weekFrom(from)

	for _, v := range wk {
		fmt.Printf("%v	", v.name)
		println(v.String())
		/*
			if v.wakeup.IsZero() {
				fmt.Printf("nope	")
			} else {
				fmt.Printf("%v	", v.wakeup.Format(saneKitchen))
			}
			if v.bed.IsZero() {
				fmt.Printf("nope	")
			} else {
				fmt.Printf("%v	", v.bed.Format(saneKitchen))
			}
			println()
		*/
	}
}
