package ansible

// Event produced by Ansible when running a playbook
type Event interface {
	// Type is the name of the event type
	Type() string
}

type namedEvent struct {
	Name string
}

type runnerResult struct {
	// Command is the command that was run
	Command []string `json:"cmd"`
	// Stdout captured when the command was run
	Stdout string
	// Stderr captured when the command was run
	Stderr string
}

type runnerResultEvent struct {
	Host   string
	Result runnerResult
}

type PlaybookStartEvent struct {
	namedEvent
}

func (e *PlaybookStartEvent) Type() string {
	return "Playbook Start"
}

type PlayStartEvent struct {
	namedEvent
}

func (e *PlayStartEvent) Type() string {
	return "Play Start"
}

type TaskStartEvent struct {
	namedEvent
}

func (e *TaskStartEvent) Type() string {
	return "Task Start"
}

type HandlerTaskStartEvent struct {
	namedEvent
}

func (e *HandlerTaskStartEvent) Type() string {
	return "Handler Task Start"
}

type RunnerOKEvent struct {
	runnerResultEvent
}

func (e *RunnerOKEvent) Type() string {
	return "Runner OK"
}

type RunnerFailedEvent struct {
	runnerResultEvent
}

func (e *RunnerFailedEvent) Type() string {
	return "Runner Failed"
}

type RunnerItemOKEvent struct {
	runnerResultEvent
}

func (e *RunnerItemOKEvent) Type() string {
	return "Runner Item OK"
}

type RunnerItemRetryEvent struct {
	runnerResultEvent
}

func (e *RunnerItemRetryEvent) Type() string {
	return "Runner Item Retry"
}

type RunnerSkippedEvent struct {
	runnerResultEvent
}

func (e *RunnerSkippedEvent) Type() string {
	return "Runner Skipped"
}

type RunnerUnreachableEvent struct {
	runnerResultEvent
}

func (e *RunnerUnreachableEvent) Type() string {
	return "Runner Unreachable"
}
