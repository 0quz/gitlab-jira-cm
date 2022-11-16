package service

import "github.com/go-playground/validator/v10"

// It checks the whole close and approves If it is valid.
func (c *Close) CloseValidate() error {
	v := validator.New()
	v.RegisterValidation("name", validateStatus)
	return v.Struct(c)
}

// It validates the close->issue->status->name
func validateStatus(fl validator.FieldLevel) bool {
	status := false
	if fl.Field().String() == "Done" || fl.Field().String() == "Reject" {
		status = true
	}
	return status
}

// It checks the whole merge request and approves If it is valid.
func (mr *MergeRequest) MergeRequestValidate() error {
	v := validator.New()
	v.RegisterValidation("target_branch", validateTargetBranch)
	v.RegisterValidation("state", validateState)
	return v.Struct(mr)
}

// It validates the mergeRequest->objectAttributes->sourceBranch
func validateTargetBranch(fl validator.FieldLevel) bool {
	status := false
	if fl.Field().String() == "prod-release" ||
		fl.Field().String() == "prod-release-2" {
		status = true
	}
	return status
}

// It validates the mergeRequest->objectAttributes->state
func validateState(fl validator.FieldLevel) bool {
	status := false
	if fl.Field().String() == "opened" {
		status = true
	}
	return status
}
