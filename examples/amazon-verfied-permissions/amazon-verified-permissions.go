package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/verifiedpermissions"
	"github.com/aws/aws-sdk-go-v2/service/verifiedpermissions/types"
	"github.com/strata-io/service-extension/orchestrator"
)

// The following values should be updated for your environment.
const (
	// policyStoreID: The ID of your Amazon Verified Permissions policy store.
	policyStoreID = "your-verified-permissions-store-id"

	// The session value from your IdP used as the principal ID in the call to
	// Amazon Verified Permissions.
	principalID = "Azure-OIDC.mail"
	// The action type to be used in the isAuthorized request.
	actionType = "Action"
	// The principal type to be used in the isAuthorized request.
	principalType = "User"
	// The resource type to be used in the isAuthorized request.
	resourceType = "Endpoint"
	//The action value to be used in the isAuthorized request.
	actionId = "ReadEndpoint"
)

// IsAuthorized is called after you log in with your IdP.  This function calls the
// Amazon Verified Permissions API with the associated principalID and endpoint to
// determine the authorization decision.
func IsAuthorized(api orchestrator.Orchestrator, _ http.ResponseWriter, req *http.Request) bool {
	session, _ := api.Session()
	logger := api.Logger()
	email, _ := session.GetString(principalID)
	logger.Info("requesting isAuthorized decision for " + email + " at " + req.URL.Path)
	avpReq, err := createVerifiedPermissionsRequest("aadkins@sonarsystems.co", req.URL.Path, api)
	//avpReq, err := createVerifiedPermissionsRequest(email, req.URL.Path, api)
	if err != nil {
		logger.Info("error creating request: " + err.Error())
		return false
	}
	if len(avpReq.Errors) > 0 {
		for _, e := range avpReq.Errors {
			logger.Error("se", "Error during IsAuthorized: "+*e.ErrorDescription)
		}
	}
	logger.Info("se", "The following policy id's contributed to the decision:")
	for _, dp := range avpReq.DeterminingPolicies {
		logger.Info("se", "Determing policy id for the decision: "+*dp.PolicyId)
	}
	for _, d := range avpReq.Decision.Values() {
		logger.Info("se", "isAuthorized decision from Amazon verified permissions: "+d)
		return d == "ALLOW"
	}
	return false
}

// createVerifiedPermissionsRequest builds a new verified permissions API request with the supplied
// principal and path.
func createVerifiedPermissionsRequest(principal, path string, api orchestrator.Orchestrator) (*verifiedpermissions.IsAuthorizedOutput, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	pType := principalType
	rType := resourceType
	aType := actionType
	psId := policyStoreID
	aId := actionId

	vpClient := verifiedpermissions.NewFromConfig(cfg)

	output, err := vpClient.IsAuthorized(context.TODO(), &verifiedpermissions.IsAuthorizedInput{

		Principal:     &types.EntityIdentifier{EntityId: &principal, EntityType: &pType},
		Resource:      &types.EntityIdentifier{EntityId: &path, EntityType: &rType},
		Action:        &types.ActionIdentifier{ActionId: &aId, ActionType: &aType},
		PolicyStoreId: &psId})
	if err != nil {
		log.Fatal(err)
	}
	return output, err
}

type Request struct {
	PolicyStoreID string    `json:"policyStoreId"`
	Action        Action    `json:"action"`
	Principal     Principal `json:"principal"`
	Resource      Resource  `json:"resource"`
}

type Action struct {
	ActionId   string `json:"actionId"`
	ActionType string `json:"actionType"`
}
type Principal struct {
	EntityId   string `json:"entityId"`
	EntityType string `json:"entityType"`
}

type Resource struct {
	EntityId   string `json:"entityId"`
	EntityType string `json:"entityType"`
}

type Response struct {
	Decision            string        `json:"decision"`
	DeterminingPolicies []interface{} `json:"determiningPolicies"`
	Errors              []interface{} `json:"errors"`
}
