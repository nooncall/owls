package task

const (
	WaitApproval = "wait_approval"
	Pass         = "pass"
	Failed       = "failed"
	Cancel       = "cancel"
	Reject       = "reject"
)

func ExecStatus() []string {
	return []string{WaitApproval}
}

func SubmitStatus() []string {
	return []string{WaitApproval, Reject}
}

func HistoryStatus() []string {
	return []string{Pass, Failed, Cancel}
}
