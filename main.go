package main

import "fmt"

type MeetingID string

func main() {
	var title string = "Go Basics Meetup"
	participantCount := 3
	var isPublished bool
	const maxParticipants int = 50
	var meetingID MeetingID = "meetup-001"

	fmt.Printf("title=%q, type=%T\n", title, title)
	fmt.Printf("participantCount=%d, type=%T\n", participantCount, participantCount)
	fmt.Printf("isPublished=%t, type=%T\n", isPublished, isPublished)
	fmt.Printf("maxParticipants=%d, type=%T\n", maxParticipants, maxParticipants)
	fmt.Printf("meetingID=%q, type=%T\n", meetingID, meetingID)

	rawID := string(meetingID)
	fmt.Println(rawID)

	summary, freeSlots := buildMeetingSummary(meetingID, title, participantCount, maxParticipants)
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

	for i := 0; i < participantCount; i++ {
		fmt.Printf("participant #%d\n", i+1)
	}
}

func buildMeetingSummary(meetingID MeetingID, title string, participantCount int, maxParticipants int) (string, int) {
	return string(meetingID) + ": " + title, maxParticipants - participantCount
}
