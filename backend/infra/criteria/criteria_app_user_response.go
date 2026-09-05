package criteria

import "bitbucket.org/lyndus/backend/global/aggregate"

type AppUserResponse struct {
	CResponse
	Items []aggregate.AppUserAggreate `json:"items"`
}


