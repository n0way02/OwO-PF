package ffxiv

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Listings struct {
	Listings []*Listing
}

type Listing struct {
	DataCentre  string
	Id          string
	PfCategory  string
	Duty        string `selector:".left .duty"`
	Tags        string `selector:".left .description span"`
	TagsColor   string `selector:".left .description span" attr:"class"`
	Description string
	MinIL       string `selector:".middle .stat .value"`
	Creator     string `selector:".right .creator .text"`
	World       string `selector:".right .world .text"`
	Expires     string `selector:".right .expires .text"`
	Updated     string `selector:".right .updated .text"`
	Party       []*Slot
}

type Slot struct {
	Roles  Roles
	Job    Job
	Filled bool
}

func NewSlot() *Slot {
	return &Slot{
		Roles: Roles{Roles: []Role{}},
	}
}

func (ls *Listings) ForDutyAndDataCentre(duty string, dataCentre string) *Listings {
	listings := &Listings{Listings: []*Listing{}}

	for _, l := range ls.Listings {
		if l.Duty == duty {
			if l.DataCentre == dataCentre {
				listings.Listings = append(listings.Listings, l)
			}
		}
	}

	return listings
}

func (ls *Listings) MostRecentUpdated() (*Listing, error) {
	var mostRecentUpdated time.Time
	var mostRecent *Listing

	for _, l := range ls.Listings {
		updatedAt, err := l.UpdatedAt()
		if err != nil {
			return nil, fmt.Errorf("Could not find most recent update time: %w", err)
		}
		if updatedAt.After(mostRecentUpdated) {
			mostRecentUpdated = updatedAt
			mostRecent = l
		}
	}

	return mostRecent, nil
}

func (ls *Listings) UpdatedWithinLast(duration time.Duration) (*Listings, error) {
	listings := &Listings{Listings: []*Listing{}}
	now := time.Now()

	for _, l := range ls.Listings {
		updatedAt, err := l.UpdatedAt()
		if err != nil {
			return nil, fmt.Errorf("Could not find most recent update time: %w", err)
		}
		if now.Add(-duration).Before(updatedAt) {
			listings.Listings = append(listings.Listings, l)
		}
	}

	return listings, nil
}

func (ls *Listings) Add(l *Listing) {
	for _, existingListing := range ls.Listings {
		if existingListing.Id == l.Id {
			return
		}
	}

	ls.Listings = append(ls.Listings, l)
}

func (l *Listing) PartyDisplay() string {
	var party strings.Builder

	for _, slot := range l.Party {
		if slot.Filled {
			party.WriteString(slot.Job.Emoji() + " ")
		} else {
			party.WriteString(slot.Roles.Emoji() + " ")
		}
	}

	return party.String()

}

func (l *Listing) GetExpires() string {
	return "<:haruamimir:1358896108222284018> " + l.Expires
}

func (l *Listing) GetUpdated() string {
	return "<:harufaca:1358896115352604964>" + l.Updated
}

func (l *Listing) GetTags() string {
	if len(l.Tags) == 0 {
		return "_ _"
	}
	return l.Tags
}

func (l *Listing) GetDescription() string {
	return "```" + l.Description + "```"
}

var expiresSecondsRegexp = regexp.MustCompile(`em (\d+) seconds`)
var expiresMinutesRegexp = regexp.MustCompile(`em (\d+) minutes`)
var expiresHoursRegexp = regexp.MustCompile(`em (\d+) hours`)

func (l *Listing) ExpiresAt() (time.Time, error) {
	now := time.Now()

	if l.Expires == "" {
		return now, nil
	}

	if l.Expires == "agora" {
		return now, nil
	}

	if l.Expires == "em um segundo" {
		return now.Add(time.Duration(1) * time.Second), nil
	}

	if l.Expires == "em um minuto" {
		return now.Add(time.Duration(1) * time.Minute), nil
	}

	if l.Expires == "em uma hora" {
		return now.Add(time.Duration(1) * time.Hour), nil
	}

	match := expiresSecondsRegexp.FindStringSubmatch(l.Expires)
	if len(match) != 0 {
		seconds, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("Could not parse time %v: %w", l.Expires, err)
		}
		return now.Add(time.Duration(seconds) * time.Second), nil
	}

	match = expiresMinutesRegexp.FindStringSubmatch(l.Expires)
	if len(match) != 0 {
		minutes, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("Could not parse time %v: %w", l.Expires, err)
		}
		return now.Add(time.Duration(minutes) * time.Minute), nil
	}

	match = expiresHoursRegexp.FindStringSubmatch(l.Expires)
	if len(match) != 0 {
		hours, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("Could not parse time %v: %w", l.Expires, err)
		}
		return now.Add(time.Duration(hours) * time.Hour), nil
	}

	return now, fmt.Errorf("Failed to parse time %v", l.Expires)
}

var updatedSecondsRegexp = regexp.MustCompile(`(\d+) seconds ago`)
var updatedMinutesRegexp = regexp.MustCompile(`(\d+) minutes ago`)
var updatedHoursRegexp = regexp.MustCompile(`(\d+) hours ago`)

func (l *Listing) UpdatedAt() (time.Time, error) {
	now := time.Now()

	if l.Updated == "" {
		return now, nil
	}

	if l.Updated == "now" {
		return now, nil
	}

	if l.Updated == "a second ago" {
		return now.Add(time.Duration(-1) * time.Second), nil
	}

	if l.Updated == "a minute ago" {
		return now.Add(time.Duration(-1) * time.Minute), nil
	}

	if l.Updated == "an hour ago" {
		return now.Add(time.Duration(-1) * time.Hour), nil
	}

	match := updatedSecondsRegexp.FindStringSubmatch(l.Updated)
	if len(match) != 0 {
		seconds, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("Could not parse time %v: %w", l.Updated, err)
		}
		return now.Add(time.Duration(-seconds) * time.Second), nil
	}

	match = updatedMinutesRegexp.FindStringSubmatch(l.Updated)
	if len(match) != 0 {
		minutes, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("Could not parse time %v: %w", l.Updated, err)
		}
		return now.Add(time.Duration(-minutes) * time.Minute), nil
	}

	match = updatedHoursRegexp.FindStringSubmatch(l.Updated)
	if len(match) != 0 {
		hours, err := strconv.Atoi(match[1])
		if err != nil {
			return now, fmt.Errorf("Could not parse time %v: %w", l.Updated, err)
		}
		return now.Add(time.Duration(-hours) * time.Hour), nil
	}

	return now, fmt.Errorf("Failed to parse time %v", l.Updated)
}
