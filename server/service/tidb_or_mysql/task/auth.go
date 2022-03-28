package task

type authTools interface {
	GetReviewer(userName string) (reviewerName string, err error)
	IsDba(userName string) (isDba bool, err error)
}

var AuthTool authTools

func SetAuthTools(impl authTools) {
	AuthTool = impl
}

type EditAuth struct {
	SubmitReviewEnable        bool `json:"submit_review_enable"`
	SysReviewEnable           bool `json:"sys_review_enable"`
	WithdrawEnable            bool `json:"withdraw_enable"`
	TurnDownEnable            bool `json:"turn_down_enable"`
	SubmitExecEnable          bool `json:"submit_exec_enable"`
	ExecEnable                bool `json:"exec_enable"`
	TurnDownExecEnable        bool `json:"turn_down_exec_enable"`
	ReviewEnable              bool `json:"review_enable"`
	ReviewPassEnable          bool `json:"review_pass_enable"`
	CreatorSubmitReviewEnable bool `json:"creator_submit_review_enable"`
	CreatorTurnDownEnable     bool `json:"creator_turn_down_enable"`
}

// GetTaskOpAbility根据task状态，获取角色的可操作类型
func GetTaskOperateAuth(detail, isCreator, isReviewer, isDba bool, task *OwlTask) *EditAuth {
	switch task.Status {
	case CheckFailed:
		if isCreator {
			if detail {
				return &EditAuth{
					SysReviewEnable:       true,
					CreatorTurnDownEnable: true,
					WithdrawEnable:        true,
				}
			}
			return &EditAuth{
				CreatorSubmitReviewEnable: true,
				CreatorTurnDownEnable:     true,
				WithdrawEnable:            true,
			}
		}
	case CheckPass:
		if (isReviewer || isDba) && isCreator {
			if detail {
				return &EditAuth{
					ReviewPassEnable: true,
					TurnDownEnable:   true,
					WithdrawEnable:   true,
				}
			}
			return &EditAuth{
				ReviewPassEnable: true,
				TurnDownEnable:   true,
				WithdrawEnable:   true,
			}
		}
		if isCreator {
			if detail {
				return &EditAuth{
					SysReviewEnable:       true,
					SubmitReviewEnable:    true,
					CreatorTurnDownEnable: true,
					WithdrawEnable:        true,
				}
			}
			return &EditAuth{
				CreatorSubmitReviewEnable: true,
				CreatorTurnDownEnable:     true,
				WithdrawEnable:            true,
			}
		}
		//既是创建人，又是审核人；撤销和驳回的意义不大了
		if isCreator && isReviewer {
			return &EditAuth{
				ReviewPassEnable: true,
				WithdrawEnable:   true,
			}
		}
		if isReviewer || isDba {
			if detail {
				return &EditAuth{
					ReviewPassEnable: true,
					TurnDownEnable:   true,
				}
			}
			return &EditAuth{
				ReviewPassEnable: true,
				TurnDownEnable:   true,
			}
		}
	case Reject:
		//既是创建人，又是审核人；撤销和驳回的意义不大了
		if isCreator && isReviewer {
			if detail {
				return &EditAuth{
					SysReviewEnable: true,
				}
			}
			// 非detail模式，展示的还是审核
			return &EditAuth{
				ReviewPassEnable: true,
			}
		}
		if isCreator {
			return &EditAuth{
				WithdrawEnable: true,
			}
		}
		if isReviewer {
			if detail {
				return &EditAuth{
					SysReviewEnable: true,
					TurnDownEnable:  true,
					WithdrawEnable:  true,
				}
			}
			return &EditAuth{
				ReviewPassEnable: true,
				TurnDownEnable:   true,
				WithdrawEnable:   true,
			}
		}

		if isDba {
			if detail {
				return &EditAuth{
					SysReviewEnable: true,
					TurnDownEnable:  true,
				}
			}
			return &EditAuth{
				ReviewPassEnable: true,
				TurnDownEnable:   true,
			}
		}
	case ReviewPass, Executing, ExecWait:
		//如果是dml，则不需要dba审核
		if isDba {
			if detail {
				return &EditAuth{
					ReviewPassEnable: true,
					TurnDownEnable:   true,
					SubmitExecEnable: true,
					ExecEnable:       true,
				}
			}
			return &EditAuth{
				ReviewPassEnable: true,
				TurnDownEnable:   true,
				SubmitExecEnable: true,
			}
		}

		if isReviewer && allIsDmlTask(task) {
			if detail {
				return &EditAuth{
					TurnDownEnable:   true,
					SubmitExecEnable: true,
				}
			}
			return &EditAuth{
				TurnDownEnable:   true,
				SubmitExecEnable: true,
				ExecEnable:       true,
			}
		}

		//既是创建人，又是审核人；撤销和驳回的意义不大了
		if isCreator && isReviewer {
			if detail {
				return &EditAuth{
					ReviewPassEnable: true,
					SysReviewEnable:  true,
				}
			}
			return &EditAuth{
				ReviewPassEnable: true,
			}
		}
		if isReviewer {
			if detail {
				return &EditAuth{
					ReviewPassEnable: true,
					SysReviewEnable:  true,
					TurnDownEnable:   true,
					WithdrawEnable:   true,
				}
			}
			return &EditAuth{
				ReviewPassEnable: true,
				TurnDownEnable:   true,
				WithdrawEnable:   true,
			}
		}
		if isCreator {
			return &EditAuth{
				WithdrawEnable: true,
			}
		}
	case DBAPass:
		if isDba {
			if detail {
				return &EditAuth{
					SysReviewEnable:  true,
					TurnDownEnable:   true,
					SubmitExecEnable: true,
					ExecEnable:       true,
				}
			}
			return &EditAuth{
				ReviewPassEnable: true,
				TurnDownEnable:   true,
			}
		}
		if isCreator {
			return &EditAuth{
				WithdrawEnable: true,
			}
		}

		//如果是dml 给leader执行权限
		if isReviewer && allIsDmlTask(task) {
			if detail {
				return &EditAuth{
					ExecEnable:         true,
					TurnDownExecEnable: true,
				}
			}
			return &EditAuth{
				ExecEnable:         true,
				TurnDownExecEnable: true,
			}
		}

		if isReviewer {
			return &EditAuth{
				WithdrawEnable: true,
			}
		}
	case ExecFailed:
		if isDba {
			if detail {
				return &EditAuth{
					ExecEnable:         false,
					TurnDownExecEnable: true,
				}
			}
			return &EditAuth{
				ExecEnable:         true,
				TurnDownExecEnable: true,
			}
		}
	}

	return &EditAuth{}
}

func allIsDmlTask(task *OwlTask) bool {
	for _, v := range task.SubTasks {
		if v.TaskType != DML {
			return false
		}
	}
	return true
}

func IsDba(userName string) (isDba bool, err error) {
	return AuthTool.IsDba(userName)
}
