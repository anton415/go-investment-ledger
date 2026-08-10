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

	err := meeting.addParticipant("Bob")
	if err != nil {
		fmt.Printf("%s %v\n", "Bob", err)
	} else {
		fmt.Printf("%s - %s\n", "Bob", "добавлен")
	}
	err = meeting.addParticipant("Alice")
	if err != nil {
		fmt.Printf("%s %v\n", "Alice", err)
	} else {
		fmt.Printf("%s - %s\n", "Alice", "добавлен")
	}
	err = meeting.addParticipant("Carol")
	if err != nil {
		fmt.Printf("%s %v\n", "Carol", err)
	} else {
		fmt.Printf("%s - %s\n", "Carol", "добавлен")
	}
	err = meeting.addParticipant("David")
	if err != nil {
		fmt.Printf("%s %v\n", "David", err)
	} else {
		fmt.Printf("%s - %s\n", "David", "добавлен")
	}

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

func (m *Meeting) addParticipant(participant string) error {
	if len(m.Participants) >= m.MaxParticipants {
		return fmt.Errorf("свободных мест нет")
	}
	if contains(m.Participants, participant) {
		return fmt.Errorf("участник уже добавлен")
	}
	m.Participants = append(m.Participants, participant)
	return nil
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
