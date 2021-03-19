package v1

type OperatorStatus string

const (
	STATUS_ERROR       OperatorStatus = "ERROR"
	STATUS_RECONCILING OperatorStatus = "RECONCILING"
	STATUS_HEALTHY     OperatorStatus = "HEALTHY"
	STATUS_READY       OperatorStatus = "READY"
	STATUS_REMOVING    OperatorStatus = "REMOVING"
)
