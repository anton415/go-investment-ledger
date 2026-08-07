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
}
