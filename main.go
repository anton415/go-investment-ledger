package main

import "fmt"

type MeetingID string

type Meeting struct {
	ID              MeetingID
	Title           string
	Participants    []string
	MaxParticipants int
	IsPublished     bool
}

func main() {
	meeting := Meeting{
		ID:              "meetup-001",
		Title:           "Go Basics Meetup",
		Participants:    []string{"Alice", "Bob", "Charlie"},
		MaxParticipants: 50,
		IsPublished:     false,
	}
	meeting.Participants = append(meeting.Participants, "David")

	fmt.Printf("title=%q, type=%T\n", meeting.Title, meeting.Title)
	fmt.Printf("participantCount=%d, type=%T\n", len(meeting.Participants), len(meeting.Participants))
	fmt.Printf("isPublished=%t, type=%T\n", meeting.IsPublished, meeting.IsPublished)
	fmt.Printf("maxParticipants=%d, type=%T\n", meeting.MaxParticipants, meeting.MaxParticipants)
	fmt.Printf("meetingID=%q, type=%T\n", meeting.ID, meeting.ID)

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
	case freeSlots <= 5:
		availability = "almost full"
	default:
		availability = "available"
	}
	fmt.Println(availability)

	participantsMap := map[MeetingID][]string{meeting.ID: meeting.Participants}
	meetingParticipants, ok := participantsMap[meeting.ID]

	if ok {
		for index, participant := range meetingParticipants {
			fmt.Printf("participant #%d %s\n", index+1, participant)
		}
		fmt.Printf("participants=%v, found=%t\n", meetingParticipants, ok)
	}

	missingParticipants, found := participantsMap[MeetingID("meetup-999")]
	fmt.Printf("participants=%v, found=%t\n", missingParticipants, found)
}

func (m Meeting) buildMeetingSummary() (string, int) {
	return string(m.ID) + ": " + m.Title, m.MaxParticipants - len(m.Participants)
}
