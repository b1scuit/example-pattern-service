package core

// Dont have to do this this way, just saves a long func call
type Task1Input struct {
	To      string
	From    string
	Number  string
	Subject string
	Body    string
}

// Provides clear, simple to read calls to make decisions from / define behaviour
func (ti *Task1Input) IsNumberSet() bool {
	return ti.Number != ""
}
