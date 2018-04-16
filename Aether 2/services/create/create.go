// Create
// This package deals with the creation of entities. This is a higher level service that is composed of lower level services in the services directory.

package create

import (
	"aether-core/io/api"
	"aether-core/services/globals"
	// "aether-core/services/verify"
	"errors"
	"fmt"
	"time"
)

// Bake is the function that handles the core signature / pow / fingerprint trio.
func Bake(entity api.Provable) error {
	// 1) Signature
	// 2) PoW
	// 3) Fingerprint
	err := entity.CreateSignature(globals.FrontendConfig.GetUserKeyPair())
	if err != nil {
		return errors.New(fmt.Sprintf(
			"Entity creation failed. Error: %s, Entity: %#v\n", err, entity))
	}
	err2 := *new(error)
	switch ent := entity.(type) {
	case *api.Board:
		err2 = ent.CreatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().Board)
	case *api.Thread:
		err2 = ent.CreatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().Thread)
	case *api.Post:
		err2 = ent.CreatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().Post)
	case *api.Vote:
		err2 = ent.CreatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().Vote)
	case *api.Key:
		err2 = ent.CreatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().Key)
	case *api.Truststate:
		err2 = ent.CreatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().Truststate)
	}
	if err2 != nil {
		return errors.New(fmt.Sprintf(
			"Entity creation failed. Error: %s, Entity: %#v\n", err2, entity))
	}
	entity.CreateFingerprint()
	return nil
}

// Rebake saves the updates to the entity and updates the signature and pow accordingly based on given fields.

func Rebake(entity api.Updateable) error {
	err := entity.CreateUpdateSignature(globals.FrontendConfig.GetUserKeyPair())
	if err != nil {
		return errors.New(fmt.Sprintf(
			"Update signature creation failed. Error: %s, Entity: %#v\n", err, entity))
	}
	err2 := *new(error)
	switch ent := entity.(type) {
	case *api.Board:
		err2 = ent.CreateUpdatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().BoardUpdate)
	case *api.Thread:
		err2 = ent.CreateUpdatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().ThreadUpdate)
	case *api.Post:
		err2 = ent.CreateUpdatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().PostUpdate)
	case *api.Vote:
		err2 = ent.CreateUpdatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().VoteUpdate)
	case *api.Key:
		err2 = ent.CreateUpdatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().KeyUpdate)
	case *api.Truststate:
		err2 = ent.CreateUpdatePoW(globals.FrontendConfig.GetUserKeyPair(), globals.BackendConfig.GetMinimumPoWStrengths().TruststateUpdate)
	}
	if err2 != nil {
		return errors.New(fmt.Sprintf(
			"Entity creation failed. Error: %s, Entity: %#v\n", err2, entity))
	}
	return nil
}

// Create sub-entities

func CreateBoardOwner(
	keyFingerprint api.Fingerprint,
	expiry api.Timestamp,
	level uint8,
) (api.BoardOwner, error) {

	var bo api.BoardOwner
	bo.KeyFingerprint = keyFingerprint
	bo.Expiry = expiry
	bo.Level = level
	return bo, nil
}

// Create main entities

func CreateBoard(
	boardName string,
	ownerFp api.Fingerprint,
	boardOwners []api.BoardOwner,
	description string,
	meta string,
	realmId api.Fingerprint,
) (api.Board, error) {

	var entity api.Board
	entity.Creation = api.Timestamp(time.Now().Unix())
	entity.Name = boardName
	entity.Owner = ownerFp
	entity.BoardOwners = boardOwners
	entity.Description = description
	entity.EntityVersion = globals.BackendTransientConfig.EntityVersions.Board
	entity.Meta = meta
	entity.RealmId = realmId
	err := Bake(&entity)
	if err != nil {
		var blankEntity api.Board
		return blankEntity, err
	}
	return entity, nil
}

func CreateThread(
	boardFp api.Fingerprint,
	name string,
	body string,
	link string,
	ownerFp api.Fingerprint,
	meta string,
	realmId api.Fingerprint,
) (api.Thread, error) {

	var entity api.Thread
	entity.Creation = api.Timestamp(time.Now().Unix())
	entity.Board = boardFp
	entity.Name = name
	entity.Body = body
	entity.Link = link
	entity.Owner = ownerFp
	entity.EntityVersion = globals.BackendTransientConfig.EntityVersions.Thread
	entity.Meta = meta
	entity.RealmId = realmId
	err := Bake(&entity)
	if err != nil {
		var blankEntity api.Thread
		return blankEntity, err
	}
	return entity, nil
}

func CreatePost(
	boardFp api.Fingerprint,
	threadFp api.Fingerprint,
	parentFp api.Fingerprint,
	body string,
	ownerFp api.Fingerprint,
	meta string,
	realmId api.Fingerprint,
) (api.Post, error) {

	var entity api.Post
	entity.Creation = api.Timestamp(time.Now().Unix())
	entity.Board = boardFp
	entity.Thread = threadFp
	entity.Parent = parentFp
	entity.Body = body
	entity.Owner = ownerFp
	entity.EntityVersion = globals.BackendTransientConfig.EntityVersions.Post
	entity.Meta = meta
	entity.RealmId = realmId
	err := Bake(&entity)
	if err != nil {
		var blankEntity api.Post
		return blankEntity, err
	}
	return entity, nil
}

func CreateVote(
	boardFp api.Fingerprint,
	threadFp api.Fingerprint,
	targetFp api.Fingerprint,
	ownerFp api.Fingerprint,
	voteType int,
	meta string,
	realmId api.Fingerprint,
) (api.Vote, error) {

	var entity api.Vote
	entity.Creation = api.Timestamp(time.Now().Unix())
	entity.Board = boardFp
	entity.Thread = threadFp
	entity.Target = targetFp
	entity.Owner = ownerFp
	entity.Type = voteType
	entity.EntityVersion = globals.BackendTransientConfig.EntityVersions.Vote
	entity.Meta = meta
	entity.RealmId = realmId
	err := Bake(&entity)
	if err != nil {
		var blankEntity api.Vote
		return blankEntity, err
	}
	return entity, nil
}

func CreateAddress(
	loc api.Location,
	subloc api.Location,
	locType uint8,
	port uint16,
	addrType uint8,
	lastOnline api.Timestamp,
	protVMajor uint8,
	protVMinor uint16,
	subprotocols []api.Subprotocol,
	clientVMajor uint8,
	clientVMinor uint16,
	clientVPatch uint16,
	clientName string,
	realmId api.Fingerprint,
) (api.Address, error) {

	var entity api.Address
	entity.Location = loc
	entity.Sublocation = subloc
	entity.LocationType = locType
	entity.Port = port
	entity.Type = addrType
	entity.LastOnline = lastOnline
	var prot api.Protocol
	prot.VersionMajor = protVMajor
	prot.VersionMinor = protVMinor
	prot.Subprotocols = subprotocols
	var client api.Client
	client.VersionMajor = clientVMajor
	client.VersionMinor = clientVMinor
	client.VersionPatch = clientVPatch
	client.ClientName = clientName
	entity.Protocol = prot
	entity.Client = client
	entity.EntityVersion = globals.BackendTransientConfig.EntityVersions.Thread
	entity.RealmId = realmId
	return entity, nil
}

func CreateKey(
	keyType string,
	key string,
	name string,
	info string,
	meta string,
	realmId api.Fingerprint,
) (api.Key, error) {

	var entity api.Key
	entity.Creation = api.Timestamp(time.Now().Unix())
	entity.Type = keyType
	entity.Key = key
	entity.Name = name
	entity.Info = info
	entity.EntityVersion = globals.BackendTransientConfig.EntityVersions.Thread
	entity.Meta = meta
	entity.RealmId = realmId
	err := Bake(&entity)
	if err != nil {
		var blankEntity api.Key
		return blankEntity, err
	}
	return entity, nil
}

func CreateTruststate(
	targetFp api.Fingerprint,
	ownerFp api.Fingerprint,
	tsType int,
	domains []api.Fingerprint,
	expiry api.Timestamp,
	meta string,
	realmId api.Fingerprint,
) (api.Truststate, error) {

	var entity api.Truststate
	entity.Creation = api.Timestamp(time.Now().Unix())
	entity.Target = targetFp
	entity.Owner = ownerFp
	entity.Type = tsType
	entity.Domains = domains
	entity.Expiry = expiry
	entity.EntityVersion = globals.BackendTransientConfig.EntityVersions.Thread
	entity.Meta = meta
	entity.RealmId = realmId
	err := Bake(&entity)
	if err != nil {
		var blankEntity api.Truststate
		return blankEntity, err
	}
	return entity, nil
}

// The functions below cannot be methods on the api types because they are defined in the api package, not here. If I try to extend that here, I get an error. If I try to import the create from api, it won't compile because of circular imports.

type BoardUpdateRequest struct {
	Entity             *api.Board
	BoardOwnersUpdated bool
	NewBoardOwners     []api.BoardOwner
	DescriptionUpdated bool
	NewDescription     string
}

func UpdateBoard(request BoardUpdateRequest) error {
	if request.BoardOwnersUpdated {
		request.Entity.BoardOwners = request.NewBoardOwners
	}
	if request.DescriptionUpdated {
		request.Entity.Description = request.NewDescription
	}
	request.Entity.LastUpdate = api.Timestamp(time.Now().Unix())
	err := Rebake(request.Entity)
	if err != nil {
		return err
	}
	return nil
}

type ThreadUpdateRequest struct {
	Entity      *api.Thread
	BodyUpdated bool
	NewBody     string
}

func UpdateThread(request ThreadUpdateRequest) error {
	if request.BodyUpdated {
		request.Entity.Body = request.NewBody
	}
	request.Entity.LastUpdate = api.Timestamp(time.Now().Unix())
	err := Rebake(request.Entity)
	if err != nil {
		return err
	}
	return nil
}

type PostUpdateRequest struct {
	Entity      *api.Post
	BodyUpdated bool
	NewBody     string
}

func UpdatePost(request PostUpdateRequest) error {
	if request.BodyUpdated {
		request.Entity.Body = request.NewBody
	}
	request.Entity.LastUpdate = api.Timestamp(time.Now().Unix())
	err := Rebake(request.Entity)
	if err != nil {
		return err
	}
	return nil
}

type VoteUpdateRequest struct {
	Entity      *api.Vote
	TypeUpdated bool
	NewType     int
}

func UpdateVote(request VoteUpdateRequest) error {
	if request.TypeUpdated {
		request.Entity.Type = request.NewType
	}
	request.Entity.LastUpdate = api.Timestamp(time.Now().Unix())
	err := Rebake(request.Entity)
	if err != nil {
		return err
	}
	return nil
}

type KeyUpdateRequest struct {
	Entity      *api.Key
	InfoUpdated bool
	NewInfo     string
}

func UpdateKey(request KeyUpdateRequest) error {
	if request.InfoUpdated {
		request.Entity.Info = request.NewInfo
	}
	request.Entity.LastUpdate = api.Timestamp(time.Now().Unix())
	err := Rebake(request.Entity)
	if err != nil {
		return err
	}
	return nil
}

type TruststateUpdateRequest struct {
	Entity         *api.Truststate
	TypeUpdated    bool
	NewType        int
	DomainsUpdated bool
	NewDomains     []api.Fingerprint
	ExpiryUpdated  bool
	NewExpiry      api.Timestamp
}

func UpdateTruststate(request TruststateUpdateRequest) error {
	if request.TypeUpdated {
		request.Entity.Type = request.NewType
	}
	if request.DomainsUpdated {
		request.Entity.Domains = request.NewDomains
	}
	if request.ExpiryUpdated {
		request.Entity.Expiry = request.NewExpiry
	}
	request.Entity.LastUpdate = api.Timestamp(time.Now().Unix())
	err := Rebake(request.Entity)
	if err != nil {
		return err
	}
	return nil
}
