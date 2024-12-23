package notification

import "fmt"

type FakeNotifier struct {
}

func NewFakeNotifier() (Notifier, error) {
	return &FakeNotifier{}, nil
}

func (n FakeNotifier) Notify(data interface{}, to string) error {
	fmt.Println(data)
	fmt.Println("Sent a message to: " + to)
	return nil
}
