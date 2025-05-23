package reqrespupdate_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/temporalio/samples-go/reqrespupdate"
	"go.temporal.io/sdk/testsuite"
)

func TestUppercaseWorkflow(t *testing.T) {
	// Create env
	var suite testsuite.WorkflowTestSuite
	env := suite.NewTestWorkflowEnvironment()
	env.RegisterWorkflow(reqrespupdate.UppercaseWorkflow)
	env.RegisterActivity(reqrespupdate.UppercaseActivity)

	// Use delayed callbacks to send enoguh updates to cause a continue as new
	for i := 0; i < 550; i++ {
		i := i
		env.RegisterDelayedCallback(func() {
			env.UpdateWorkflow(reqrespupdate.UpdateHandler, fmt.Sprintf("test id %d", i), &testsuite.TestUpdateCallback{
				OnAccept: func() {
					if i >= 500 {
						require.Fail(t, "update should fail since it should be trying to continue-as-new")
					}
				},
				OnReject: func(err error) {
					if i < 500 {
						require.Fail(t, "this update should not fail")
					}
					require.Error(t, err)
				},
				OnComplete: func(response interface{}, err error) {
					if i >= 500 {
						require.Fail(t, "update should fail since it should be trying to continue-as-new")
					}
					require.NoError(t, err)
					require.Equal(t, reqrespupdate.Response{Output: fmt.Sprintf("FOO %d", i)}, response)
				},
			}, &reqrespupdate.Request{Input: fmt.Sprintf("foo %d", i)})
		}, 1)
	}

	// Run workflow
	env.ExecuteWorkflow(reqrespupdate.UppercaseWorkflow, true)
}
