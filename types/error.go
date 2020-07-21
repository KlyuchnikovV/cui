package types

// TODO: exrtend with additional info for future errors processing
type ErrorChannel chan error

func NewErrorChannel(capacity int) ErrorChannel {
	if capacity < 0 {
		capacity = 0
	}

	return make(chan error, capacity)
}

func (e ErrorChannel) SendError(err error) {
	e <- err
}

func (e ErrorChannel) Close() {
	close(e)
}

func (e ErrorChannel) Error() error {
	return <-e
}
