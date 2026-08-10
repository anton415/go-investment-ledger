package main

import (
	"errors"
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

var (
	ErrMeetingFull             = errors.New("свободных мест нет")
	ErrParticipantAlreadyAdded = errors.New("участник уже добавлен")
)

func main() {
	const maxParticipants = 3
	meeting := createMeeting("meetup-001", "Go Basics", maxParticipants, []string{"Alice"})

	err := meeting.addParticipant("Bob")
	switch err {
	case nil:
		fmt.Printf("%s - %s\n", "Bob", "добавлен")
	case ErrMeetingFull:
		fmt.Println("свободных мест нет")
	case ErrParticipantAlreadyAdded:
		fmt.Println("участник уже добавлен")
	default:
		fmt.Printf("Bob — неизвестная ошибка: %v\n", err)
	}

	err = meeting.addParticipant("Alice")
	switch err {
	case nil:
		fmt.Printf("%s - %s\n", "Alice", "добавлен")
	case ErrMeetingFull:
		fmt.Println("свободных мест нет")
	case ErrParticipantAlreadyAdded:
		fmt.Printf("Alice — %v\n", err)
	default:
		fmt.Printf("Alice — неизвестная ошибка: %v\n", err)
	}

	err = meeting.addParticipant("Carol")
	switch err {
	case nil:
		fmt.Printf("%s - %s\n", "Carol", "добавлен")
	case ErrMeetingFull:
		fmt.Println("свободных мест нет")
	case ErrParticipantAlreadyAdded:
		fmt.Println("участник уже добавлен")
	default:
		fmt.Printf("Carol — неизвестная ошибка: %v\n", err)
	}

	err = meeting.addParticipant("David")
	switch err {
	case nil:
		fmt.Printf("%s - %s\n", "David", "добавлен")
	case ErrMeetingFull:
		fmt.Println("свободных мест нет")
	case ErrParticipantAlreadyAdded:
		fmt.Println("участник уже добавлен")
	default:
		fmt.Printf("David — неизвестная ошибка: %v\n", err)
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
		return ErrMeetingFull
	}
	if contains(m.Participants, participant) {
		return ErrParticipantAlreadyAdded
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
