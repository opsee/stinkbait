package service

import (
	"github.com/opsee/basic/schema"
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/stinkbait/checker"
	"golang.org/x/net/context"
)

func (s *service) TestCheck(ctx context.Context, req *opsee.TestCheckRequest) (*opsee.TestCheckResponse, error) {
	var errstr string

	resp, err := checker.NewRequest(req.Check.Target.Address, req.Check.GetHttpCheck()).Do()
	if err != nil {
		errstr = err.Error()
	}

	testCheckResp := &opsee.TestCheckResponse{
		Responses: []*schema.CheckResponse{
			{
				Target: req.Check.Target,
				Error:  errstr,
				Reply: &schema.CheckResponse_HttpResponse{
					HttpResponse: resp,
				},
			},
		},
		Error: errstr,
	}

	return testCheckResp, nil
}
