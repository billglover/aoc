package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

type event struct {
	ts  time.Time
	msg string
}

func main() {
	f, err := os.Open("input.txt")
	if err != nil {
		fmt.Println("unable to open file:", err)
		os.Exit(1)
	}
	defer f.Close()

	events, err := parseEvents(f)
	if err != nil {
		fmt.Println("unable to parse events:", err)
		os.Exit(1)
	}

	// Events are listed in the order found, and not in time order. If we are
	// to associate events with guards, we need to sort these chronologically.
	sort.Slice(events, func(i, j int) bool { return events[i].ts.Before(events[j].ts) })

	// Assign each guard an event log. Each day in the event log contains a
	// string representation of the shift pattern. A '#' indicates the guard
	// fell asleep, a '.' indicates the guard woke up.
	//
	// Note: we assume there are no errors in the logs and that guards can only
	// wake up if they are asleep and vice versa.
	guards := make(map[int]map[string][60]byte)
	var re = regexp.MustCompile(`(?m)^.+#(\d+).+$`)

	var g int
	for _, e := range events {

		day := e.ts.Format("01-02")

		// The first byte of the log line is enough to differentiate events.
		switch e.msg[0] {

		// A new guard starts their shift
		case 'G':
			matches := re.FindStringSubmatch(e.msg)
			g, err = strconv.Atoi(matches[1])
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			if _, ok := guards[g]; ok == false {
				guards[g] = make(map[string][60]byte)
			}

		// The guard falls asleep
		case 'f':
			if _, ok := guards[g][day]; ok == false {
				log := [60]byte{}
				for i := range log {
					log[i] = ' '
				}
				guards[g][day] = log
			}
			log := guards[g][day]
			log[e.ts.Minute()] = '#'
			guards[g][day] = log

		// The guard wakes up
		case 'w':
			if _, ok := guards[g][day]; ok == false {
				log := [60]byte{}
				for i := range log {
					log[i] = ' '
				}
				guards[g][day] = log
			}
			log := guards[g][day]
			log[e.ts.Minute()] = '.'
			guards[g][day] = log

		default:
			fmt.Println("unknown event")
		}
	}

	// We fill in the gaps in the daily log to make it easier to determine
	// if a guard was asleep or awake in a given minute.
	for g := range guards {
		for d := range guards[g] {
			next := byte('.')
			for m := range guards[g][d] {
				switch guards[g][d][m] {
				case ' ':
					log := guards[g][d]
					log[m] = next
					guards[g][d] = log
				default:
					next = guards[g][d][m]
				}
			}
		}
	}

	// Count minutes asleep and find the guard who slept most.
	asleep := map[int]int{}
	sleepyHead, sleepyHeadDur := 0, 0
	for g := range guards {
		asleep[g] = 0
		for d := range guards[g] {
			for m := range guards[g][d] {
				if guards[g][d][m] == '#' {
					asleep[g]++
				}
			}
		}

		if asleep[g] > sleepyHeadDur {
			sleepyHead, sleepyHeadDur = g, asleep[g]
		}
	}

	// For our sleepiest guard, find the minute they are most likely to be asleep.
	mins := map[int]int{}
	for d := range guards[sleepyHead] {
		for m := range guards[sleepyHead][d] {
			if guards[sleepyHead][d][m] == '#' {
				mins[m]++
			}
		}
	}

	likelyMinute, max := 0, 0
	for k, v := range mins {
		if v > max {
			likelyMinute, max = k, v
		}
	}
	fmt.Println()
	fmt.Println("Part One:")
	fmt.Printf("Guard: %d, Minute: %d, Answer: %d\n", sleepyHead, likelyMinute, sleepyHead*likelyMinute)
	fmt.Println()

	// For all guards, find the guard who is most commonly asleep in any given minute.
	maxGuard, guardMinute, guardCount := 0, 0, 0
	for g := range guards {
		minuteCount := map[int]int{}
		for d := range guards[g] {
			for m := range guards[g][d] {
				if guards[g][d][m] == '#' {
					minuteCount[m]++
				}
			}
		}
		maxMinute, maxCount := 0, 0
		for m, count := range minuteCount {
			if count > maxCount {
				maxMinute, maxCount = m, count
			}
		}

		if maxCount > guardCount {
			maxGuard, guardMinute, guardCount = g, maxMinute, maxCount
		}
	}

	fmt.Println("Part Two:")
	fmt.Printf("Guard: %d, Minute: %d, Answer: %d\n", maxGuard, guardMinute, maxGuard*guardMinute)
	fmt.Println()
}

func parseEvents(r io.Reader) ([]event, error) {
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(data), "\n")

	events := make([]event, len(lines))

	var re = regexp.MustCompile(`(?m)^\[(.+)\]\s(.+)$`)
	for i, l := range lines {
		matches := re.FindStringSubmatch(l)
		if matches == nil {
			return nil, fmt.Errorf("unable to parse line: %s", l)
		}

		ts, err := time.Parse("2006-01-02 15:04", matches[1])
		if err != nil {
			return nil, err
		}

		events[i] = event{ts: ts, msg: matches[2]}
	}

	return events, nil
}
