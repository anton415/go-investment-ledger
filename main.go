package main

import (
	"fmt"
	"slices"
)

type MeetingID string

type Meeting struct {
	ID              MeetingID
	Title           string
	Participants    []string
	MaxParticipants int
	IsPublished     bool
}

func main() {
	const maxParticipants = 3
	meeting := createMeeting("meetup-001", "Go Basics", maxParticipants, []string{"Alice"})

	fmt.Printf("Attempt to add participant #%s %t\n", "Bob", meeting.addParticipant("Bob"))
	fmt.Printf("Attempt to add participant #%s %t\n", "Alice", meeting.addParticipant("Alice"))
	fmt.Printf("Attempt to add participant #%s %t\n", "Carol", meeting.addParticipant("Carol"))
	fmt.Printf("Attempt to add participant #%s %t\n", "David", meeting.addParticipant("David"))

	fmt.Printf("before publish: %t\n", meeting.IsPublished)
	meeting.publish()
	fmt.Printf("after publish: %t\n", meeting.IsPublished)

	meetingsMap := map[MeetingID]Meeting{meeting.ID: meeting}
	meetingParticipants, ok := meetingsMap[meeting.ID]

	if ok {
		for index, participant := range meetingParticipants.Participants {
			fmt.Printf("participant #%d %s\n", index+1, participant)
		}
		fmt.Printf("participants=%v, found=%t\n", meetingParticipants, ok)
	}

	savedParticipants, found := meetingsMap[MeetingID("meetup-001")]
	fmt.Printf("participants=%v, found=%t\n", savedParticipants, found)

	missingParticipants, found := meetingsMap[MeetingID("meetup-999")]
	fmt.Printf("participants=%v, found=%t\n", missingParticipants, found)

	rawID := string(meeting.ID)
	fmt.Println(rawID)

	summary, freeSlots := meeting.buildMeetingSummary()
	fmt.Printf("summary=%q, freeSlots=%d\n", summary, freeSlots)

	if freeSlots > 0 {
		fmt.Println("registration is open")
	} else {
		fmt.Println("meeting is full")
	}

	var availability string
	switch {
	case freeSlots <= 0:
		availability = "full"
	case freeSlots <= meeting.MaxParticipants/2:
		availability = "almost full"
	default:
		availability = "available"
	}
	fmt.Println(availability)
}

func (m Meeting) buildMeetingSummary() (string, int) {
	return string(m.ID) + ": " + m.Title, m.MaxParticipants - len(m.Participants)
}

func (m *Meeting) publish() {
	m.IsPublished = true
}

func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
}

func (m *Meeting) addParticipant(participant string) bool {
	if len(m.Participants) < m.MaxParticipants && !contains(m.Participants, participant) {
		m.Participants = append(m.Participants, participant)
		return true
	}
	return false
}

func createMeeting(id MeetingID, title string, maxParticipants int, participants []string) Meeting {
	return Meeting{
		ID:              id,
		Title:           title,
		Participants:    participants,
		MaxParticipants: maxParticipants,
		IsPublished:     false,
	}
}
