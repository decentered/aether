// Code generated by protoc-gen-go. DO NOT EDIT.
// source: feobjects/feobjects.proto

/*
Package feobjects is a generated protocol buffer package.

It is generated from these files:
	feobjects/feobjects.proto

It has these top-level messages:
	CompiledBoardEntity
	CompiledThreadEntity
	CompiledPostEntity
	CompiledUserEntity
	CompiledContentSignalsEntity
	ExplainedSignalEntity
	CompiledUserSignalsEntity
	AmbientBoardEntity
*/
package feobjects

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CompiledBoardEntity struct {
	Fingerprint            string                        `protobuf:"bytes,1,opt,name=Fingerprint" json:"Fingerprint,omitempty"`
	SelfCreated            bool                          `protobuf:"varint,2,opt,name=SelfCreated" json:"SelfCreated,omitempty"`
	Name                   string                        `protobuf:"bytes,3,opt,name=Name" json:"Name,omitempty"`
	Description            string                        `protobuf:"bytes,4,opt,name=Description" json:"Description,omitempty"`
	CompiledContentSignals *CompiledContentSignalsEntity `protobuf:"bytes,5,opt,name=CompiledContentSignals" json:"CompiledContentSignals,omitempty"`
	Owner                  *CompiledUserEntity           `protobuf:"bytes,6,opt,name=Owner" json:"Owner,omitempty"`
	BoardOwners            []string                      `protobuf:"bytes,7,rep,name=BoardOwners" json:"BoardOwners,omitempty"`
	Creation               int64                         `protobuf:"varint,8,opt,name=Creation" json:"Creation,omitempty"`
	LastUpdate             int64                         `protobuf:"varint,9,opt,name=LastUpdate" json:"LastUpdate,omitempty"`
	Meta                   string                        `protobuf:"bytes,10,opt,name=Meta" json:"Meta,omitempty"`
	ChildThreads           []*CompiledThreadEntity       `protobuf:"bytes,11,rep,name=ChildThreads" json:"ChildThreads,omitempty"`
	ThreadsCount           int32                         `protobuf:"varint,12,opt,name=ThreadsCount" json:"ThreadsCount,omitempty"`
	UserCount              int32                         `protobuf:"varint,13,opt,name=UserCount" json:"UserCount,omitempty"`
	Subscribed             bool                          `protobuf:"varint,14,opt,name=Subscribed" json:"Subscribed,omitempty"`
	Notify                 bool                          `protobuf:"varint,15,opt,name=Notify" json:"Notify,omitempty"`
	LastSeen               int64                         `protobuf:"varint,16,opt,name=LastSeen" json:"LastSeen,omitempty"`
	Whitelisted            bool                          `protobuf:"varint,17,opt,name=Whitelisted" json:"Whitelisted,omitempty"`
}

func (m *CompiledBoardEntity) Reset()                    { *m = CompiledBoardEntity{} }
func (m *CompiledBoardEntity) String() string            { return proto.CompactTextString(m) }
func (*CompiledBoardEntity) ProtoMessage()               {}
func (*CompiledBoardEntity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *CompiledBoardEntity) GetFingerprint() string {
	if m != nil {
		return m.Fingerprint
	}
	return ""
}

func (m *CompiledBoardEntity) GetSelfCreated() bool {
	if m != nil {
		return m.SelfCreated
	}
	return false
}

func (m *CompiledBoardEntity) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CompiledBoardEntity) GetDescription() string {
	if m != nil {
		return m.Description
	}
	return ""
}

func (m *CompiledBoardEntity) GetCompiledContentSignals() *CompiledContentSignalsEntity {
	if m != nil {
		return m.CompiledContentSignals
	}
	return nil
}

func (m *CompiledBoardEntity) GetOwner() *CompiledUserEntity {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *CompiledBoardEntity) GetBoardOwners() []string {
	if m != nil {
		return m.BoardOwners
	}
	return nil
}

func (m *CompiledBoardEntity) GetCreation() int64 {
	if m != nil {
		return m.Creation
	}
	return 0
}

func (m *CompiledBoardEntity) GetLastUpdate() int64 {
	if m != nil {
		return m.LastUpdate
	}
	return 0
}

func (m *CompiledBoardEntity) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func (m *CompiledBoardEntity) GetChildThreads() []*CompiledThreadEntity {
	if m != nil {
		return m.ChildThreads
	}
	return nil
}

func (m *CompiledBoardEntity) GetThreadsCount() int32 {
	if m != nil {
		return m.ThreadsCount
	}
	return 0
}

func (m *CompiledBoardEntity) GetUserCount() int32 {
	if m != nil {
		return m.UserCount
	}
	return 0
}

func (m *CompiledBoardEntity) GetSubscribed() bool {
	if m != nil {
		return m.Subscribed
	}
	return false
}

func (m *CompiledBoardEntity) GetNotify() bool {
	if m != nil {
		return m.Notify
	}
	return false
}

func (m *CompiledBoardEntity) GetLastSeen() int64 {
	if m != nil {
		return m.LastSeen
	}
	return 0
}

func (m *CompiledBoardEntity) GetWhitelisted() bool {
	if m != nil {
		return m.Whitelisted
	}
	return false
}

type CompiledThreadEntity struct {
	Fingerprint            string                        `protobuf:"bytes,1,opt,name=Fingerprint" json:"Fingerprint,omitempty"`
	Board                  string                        `protobuf:"bytes,2,opt,name=Board" json:"Board,omitempty"`
	SelfCreated            bool                          `protobuf:"varint,3,opt,name=SelfCreated" json:"SelfCreated,omitempty"`
	Name                   string                        `protobuf:"bytes,4,opt,name=Name" json:"Name,omitempty"`
	Body                   string                        `protobuf:"bytes,5,opt,name=Body" json:"Body,omitempty"`
	Link                   string                        `protobuf:"bytes,6,opt,name=Link" json:"Link,omitempty"`
	CompiledContentSignals *CompiledContentSignalsEntity `protobuf:"bytes,7,opt,name=CompiledContentSignals" json:"CompiledContentSignals,omitempty"`
	Owner                  *CompiledUserEntity           `protobuf:"bytes,8,opt,name=Owner" json:"Owner,omitempty"`
	Creation               int64                         `protobuf:"varint,9,opt,name=Creation" json:"Creation,omitempty"`
	LastUpdate             int64                         `protobuf:"varint,10,opt,name=LastUpdate" json:"LastUpdate,omitempty"`
	Meta                   string                        `protobuf:"bytes,11,opt,name=Meta" json:"Meta,omitempty"`
	ChildPosts             []*CompiledPostEntity         `protobuf:"bytes,12,rep,name=ChildPosts" json:"ChildPosts,omitempty"`
	PostsCount             int32                         `protobuf:"varint,13,opt,name=PostsCount" json:"PostsCount,omitempty"`
}

func (m *CompiledThreadEntity) Reset()                    { *m = CompiledThreadEntity{} }
func (m *CompiledThreadEntity) String() string            { return proto.CompactTextString(m) }
func (*CompiledThreadEntity) ProtoMessage()               {}
func (*CompiledThreadEntity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *CompiledThreadEntity) GetFingerprint() string {
	if m != nil {
		return m.Fingerprint
	}
	return ""
}

func (m *CompiledThreadEntity) GetBoard() string {
	if m != nil {
		return m.Board
	}
	return ""
}

func (m *CompiledThreadEntity) GetSelfCreated() bool {
	if m != nil {
		return m.SelfCreated
	}
	return false
}

func (m *CompiledThreadEntity) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CompiledThreadEntity) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *CompiledThreadEntity) GetLink() string {
	if m != nil {
		return m.Link
	}
	return ""
}

func (m *CompiledThreadEntity) GetCompiledContentSignals() *CompiledContentSignalsEntity {
	if m != nil {
		return m.CompiledContentSignals
	}
	return nil
}

func (m *CompiledThreadEntity) GetOwner() *CompiledUserEntity {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *CompiledThreadEntity) GetCreation() int64 {
	if m != nil {
		return m.Creation
	}
	return 0
}

func (m *CompiledThreadEntity) GetLastUpdate() int64 {
	if m != nil {
		return m.LastUpdate
	}
	return 0
}

func (m *CompiledThreadEntity) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func (m *CompiledThreadEntity) GetChildPosts() []*CompiledPostEntity {
	if m != nil {
		return m.ChildPosts
	}
	return nil
}

func (m *CompiledThreadEntity) GetPostsCount() int32 {
	if m != nil {
		return m.PostsCount
	}
	return 0
}

type CompiledPostEntity struct {
	Fingerprint            string                        `protobuf:"bytes,1,opt,name=Fingerprint" json:"Fingerprint,omitempty"`
	Board                  string                        `protobuf:"bytes,2,opt,name=Board" json:"Board,omitempty"`
	Thread                 string                        `protobuf:"bytes,3,opt,name=Thread" json:"Thread,omitempty"`
	Parent                 string                        `protobuf:"bytes,4,opt,name=Parent" json:"Parent,omitempty"`
	SelfCreated            bool                          `protobuf:"varint,5,opt,name=SelfCreated" json:"SelfCreated,omitempty"`
	Body                   string                        `protobuf:"bytes,6,opt,name=Body" json:"Body,omitempty"`
	CompiledContentSignals *CompiledContentSignalsEntity `protobuf:"bytes,7,opt,name=CompiledContentSignals" json:"CompiledContentSignals,omitempty"`
	Owner                  *CompiledUserEntity           `protobuf:"bytes,8,opt,name=Owner" json:"Owner,omitempty"`
	Creation               int64                         `protobuf:"varint,9,opt,name=Creation" json:"Creation,omitempty"`
	LastUpdate             int64                         `protobuf:"varint,10,opt,name=LastUpdate" json:"LastUpdate,omitempty"`
	Meta                   string                        `protobuf:"bytes,11,opt,name=Meta" json:"Meta,omitempty"`
	Children               []*CompiledPostEntity         `protobuf:"bytes,12,rep,name=Children" json:"Children,omitempty"`
}

func (m *CompiledPostEntity) Reset()                    { *m = CompiledPostEntity{} }
func (m *CompiledPostEntity) String() string            { return proto.CompactTextString(m) }
func (*CompiledPostEntity) ProtoMessage()               {}
func (*CompiledPostEntity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *CompiledPostEntity) GetFingerprint() string {
	if m != nil {
		return m.Fingerprint
	}
	return ""
}

func (m *CompiledPostEntity) GetBoard() string {
	if m != nil {
		return m.Board
	}
	return ""
}

func (m *CompiledPostEntity) GetThread() string {
	if m != nil {
		return m.Thread
	}
	return ""
}

func (m *CompiledPostEntity) GetParent() string {
	if m != nil {
		return m.Parent
	}
	return ""
}

func (m *CompiledPostEntity) GetSelfCreated() bool {
	if m != nil {
		return m.SelfCreated
	}
	return false
}

func (m *CompiledPostEntity) GetBody() string {
	if m != nil {
		return m.Body
	}
	return ""
}

func (m *CompiledPostEntity) GetCompiledContentSignals() *CompiledContentSignalsEntity {
	if m != nil {
		return m.CompiledContentSignals
	}
	return nil
}

func (m *CompiledPostEntity) GetOwner() *CompiledUserEntity {
	if m != nil {
		return m.Owner
	}
	return nil
}

func (m *CompiledPostEntity) GetCreation() int64 {
	if m != nil {
		return m.Creation
	}
	return 0
}

func (m *CompiledPostEntity) GetLastUpdate() int64 {
	if m != nil {
		return m.LastUpdate
	}
	return 0
}

func (m *CompiledPostEntity) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

func (m *CompiledPostEntity) GetChildren() []*CompiledPostEntity {
	if m != nil {
		return m.Children
	}
	return nil
}

type CompiledUserEntity struct {
	Fingerprint         string                     `protobuf:"bytes,1,opt,name=Fingerprint" json:"Fingerprint,omitempty"`
	NonCanonicalName    string                     `protobuf:"bytes,2,opt,name=NonCanonicalName" json:"NonCanonicalName,omitempty"`
	Creation            int64                      `protobuf:"varint,3,opt,name=Creation" json:"Creation,omitempty"`
	LastUpdate          int64                      `protobuf:"varint,4,opt,name=LastUpdate" json:"LastUpdate,omitempty"`
	LastRefreshed       int64                      `protobuf:"varint,5,opt,name=LastRefreshed" json:"LastRefreshed,omitempty"`
	CompiledUserSignals *CompiledUserSignalsEntity `protobuf:"bytes,6,opt,name=CompiledUserSignals" json:"CompiledUserSignals,omitempty"`
	Expiry              int64                      `protobuf:"varint,7,opt,name=Expiry" json:"Expiry,omitempty"`
	Info                string                     `protobuf:"bytes,8,opt,name=Info" json:"Info,omitempty"`
	Meta                string                     `protobuf:"bytes,9,opt,name=Meta" json:"Meta,omitempty"`
}

func (m *CompiledUserEntity) Reset()                    { *m = CompiledUserEntity{} }
func (m *CompiledUserEntity) String() string            { return proto.CompactTextString(m) }
func (*CompiledUserEntity) ProtoMessage()               {}
func (*CompiledUserEntity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

func (m *CompiledUserEntity) GetFingerprint() string {
	if m != nil {
		return m.Fingerprint
	}
	return ""
}

func (m *CompiledUserEntity) GetNonCanonicalName() string {
	if m != nil {
		return m.NonCanonicalName
	}
	return ""
}

func (m *CompiledUserEntity) GetCreation() int64 {
	if m != nil {
		return m.Creation
	}
	return 0
}

func (m *CompiledUserEntity) GetLastUpdate() int64 {
	if m != nil {
		return m.LastUpdate
	}
	return 0
}

func (m *CompiledUserEntity) GetLastRefreshed() int64 {
	if m != nil {
		return m.LastRefreshed
	}
	return 0
}

func (m *CompiledUserEntity) GetCompiledUserSignals() *CompiledUserSignalsEntity {
	if m != nil {
		return m.CompiledUserSignals
	}
	return nil
}

func (m *CompiledUserEntity) GetExpiry() int64 {
	if m != nil {
		return m.Expiry
	}
	return 0
}

func (m *CompiledUserEntity) GetInfo() string {
	if m != nil {
		return m.Info
	}
	return ""
}

func (m *CompiledUserEntity) GetMeta() string {
	if m != nil {
		return m.Meta
	}
	return ""
}

type CompiledContentSignalsEntity struct {
	TargetFingerprint string `protobuf:"bytes,1,opt,name=TargetFingerprint" json:"TargetFingerprint,omitempty"`
	// ATD
	Upvotes            int32                    `protobuf:"varint,2,opt,name=Upvotes" json:"Upvotes,omitempty"`
	Downvotes          int32                    `protobuf:"varint,3,opt,name=Downvotes" json:"Downvotes,omitempty"`
	SelfUpvoted        bool                     `protobuf:"varint,4,opt,name=SelfUpvoted" json:"SelfUpvoted,omitempty"`
	SelfDownvoted      bool                     `protobuf:"varint,5,opt,name=SelfDownvoted" json:"SelfDownvoted,omitempty"`
	SelfATDCreation    int64                    `protobuf:"varint,6,opt,name=SelfATDCreation" json:"SelfATDCreation,omitempty"`
	SelfATDLastUpdate  int64                    `protobuf:"varint,7,opt,name=SelfATDLastUpdate" json:"SelfATDLastUpdate,omitempty"`
	SelfATDFingerprint string                   `protobuf:"bytes,8,opt,name=SelfATDFingerprint" json:"SelfATDFingerprint,omitempty"`
	Reports            []*ExplainedSignalEntity `protobuf:"bytes,9,rep,name=Reports" json:"Reports,omitempty"`
	// bool SelfReported = 10;
	ModBlocks        []*ExplainedSignalEntity `protobuf:"bytes,11,rep,name=ModBlocks" json:"ModBlocks,omitempty"`
	ModApprovals     []*ExplainedSignalEntity `protobuf:"bytes,12,rep,name=ModApprovals" json:"ModApprovals,omitempty"`
	ByMod            bool                     `protobuf:"varint,13,opt,name=ByMod" json:"ByMod,omitempty"`
	ByFollowedPerson bool                     `protobuf:"varint,14,opt,name=ByFollowedPerson" json:"ByFollowedPerson,omitempty"`
	ByBlockedPerson  bool                     `protobuf:"varint,15,opt,name=ByBlockedPerson" json:"ByBlockedPerson,omitempty"`
	ByOP             bool                     `protobuf:"varint,16,opt,name=ByOP" json:"ByOP,omitempty"`
	ModBlocked       bool                     `protobuf:"varint,17,opt,name=ModBlocked" json:"ModBlocked,omitempty"`
	ModApproved      bool                     `protobuf:"varint,18,opt,name=ModApproved" json:"ModApproved,omitempty"`
	LastRefreshed    int64                    `protobuf:"varint,19,opt,name=LastRefreshed" json:"LastRefreshed,omitempty"`
}

func (m *CompiledContentSignalsEntity) Reset()                    { *m = CompiledContentSignalsEntity{} }
func (m *CompiledContentSignalsEntity) String() string            { return proto.CompactTextString(m) }
func (*CompiledContentSignalsEntity) ProtoMessage()               {}
func (*CompiledContentSignalsEntity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

func (m *CompiledContentSignalsEntity) GetTargetFingerprint() string {
	if m != nil {
		return m.TargetFingerprint
	}
	return ""
}

func (m *CompiledContentSignalsEntity) GetUpvotes() int32 {
	if m != nil {
		return m.Upvotes
	}
	return 0
}

func (m *CompiledContentSignalsEntity) GetDownvotes() int32 {
	if m != nil {
		return m.Downvotes
	}
	return 0
}

func (m *CompiledContentSignalsEntity) GetSelfUpvoted() bool {
	if m != nil {
		return m.SelfUpvoted
	}
	return false
}

func (m *CompiledContentSignalsEntity) GetSelfDownvoted() bool {
	if m != nil {
		return m.SelfDownvoted
	}
	return false
}

func (m *CompiledContentSignalsEntity) GetSelfATDCreation() int64 {
	if m != nil {
		return m.SelfATDCreation
	}
	return 0
}

func (m *CompiledContentSignalsEntity) GetSelfATDLastUpdate() int64 {
	if m != nil {
		return m.SelfATDLastUpdate
	}
	return 0
}

func (m *CompiledContentSignalsEntity) GetSelfATDFingerprint() string {
	if m != nil {
		return m.SelfATDFingerprint
	}
	return ""
}

func (m *CompiledContentSignalsEntity) GetReports() []*ExplainedSignalEntity {
	if m != nil {
		return m.Reports
	}
	return nil
}

func (m *CompiledContentSignalsEntity) GetModBlocks() []*ExplainedSignalEntity {
	if m != nil {
		return m.ModBlocks
	}
	return nil
}

func (m *CompiledContentSignalsEntity) GetModApprovals() []*ExplainedSignalEntity {
	if m != nil {
		return m.ModApprovals
	}
	return nil
}

func (m *CompiledContentSignalsEntity) GetByMod() bool {
	if m != nil {
		return m.ByMod
	}
	return false
}

func (m *CompiledContentSignalsEntity) GetByFollowedPerson() bool {
	if m != nil {
		return m.ByFollowedPerson
	}
	return false
}

func (m *CompiledContentSignalsEntity) GetByBlockedPerson() bool {
	if m != nil {
		return m.ByBlockedPerson
	}
	return false
}

func (m *CompiledContentSignalsEntity) GetByOP() bool {
	if m != nil {
		return m.ByOP
	}
	return false
}

func (m *CompiledContentSignalsEntity) GetModBlocked() bool {
	if m != nil {
		return m.ModBlocked
	}
	return false
}

func (m *CompiledContentSignalsEntity) GetModApproved() bool {
	if m != nil {
		return m.ModApproved
	}
	return false
}

func (m *CompiledContentSignalsEntity) GetLastRefreshed() int64 {
	if m != nil {
		return m.LastRefreshed
	}
	return 0
}

type ExplainedSignalEntity struct {
	SourceFp   string `protobuf:"bytes,1,opt,name=SourceFp" json:"SourceFp,omitempty"`
	Reason     string `protobuf:"bytes,2,opt,name=Reason" json:"Reason,omitempty"`
	Creation   int64  `protobuf:"varint,3,opt,name=Creation" json:"Creation,omitempty"`
	LastUpdate int64  `protobuf:"varint,4,opt,name=LastUpdate" json:"LastUpdate,omitempty"`
}

func (m *ExplainedSignalEntity) Reset()                    { *m = ExplainedSignalEntity{} }
func (m *ExplainedSignalEntity) String() string            { return proto.CompactTextString(m) }
func (*ExplainedSignalEntity) ProtoMessage()               {}
func (*ExplainedSignalEntity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ExplainedSignalEntity) GetSourceFp() string {
	if m != nil {
		return m.SourceFp
	}
	return ""
}

func (m *ExplainedSignalEntity) GetReason() string {
	if m != nil {
		return m.Reason
	}
	return ""
}

func (m *ExplainedSignalEntity) GetCreation() int64 {
	if m != nil {
		return m.Creation
	}
	return 0
}

func (m *ExplainedSignalEntity) GetLastUpdate() int64 {
	if m != nil {
		return m.LastUpdate
	}
	return 0
}

type CompiledUserSignalsEntity struct {
	TargetFingerprint      string `protobuf:"bytes,1,opt,name=TargetFingerprint" json:"TargetFingerprint,omitempty"`
	Domain                 string `protobuf:"bytes,2,opt,name=Domain" json:"Domain,omitempty"`
	FollowedBySelf         bool   `protobuf:"varint,3,opt,name=FollowedBySelf" json:"FollowedBySelf,omitempty"`
	BlockedBySelf          bool   `protobuf:"varint,4,opt,name=BlockedBySelf" json:"BlockedBySelf,omitempty"`
	FollowerCount          int32  `protobuf:"varint,5,opt,name=FollowerCount" json:"FollowerCount,omitempty"`
	CanonicalName          string `protobuf:"bytes,6,opt,name=CanonicalName" json:"CanonicalName,omitempty"`
	CNameSourceFingerprint string `protobuf:"bytes,7,opt,name=CNameSourceFingerprint" json:"CNameSourceFingerprint,omitempty"`
	SelfPEFingerprint      string `protobuf:"bytes,8,opt,name=SelfPEFingerprint" json:"SelfPEFingerprint,omitempty"`
	SelfPECreation         int64  `protobuf:"varint,9,opt,name=SelfPECreation" json:"SelfPECreation,omitempty"`
	SelfPELastUpdate       int64  `protobuf:"varint,10,opt,name=SelfPELastUpdate" json:"SelfPELastUpdate,omitempty"`
	MadeModBySelf          bool   `protobuf:"varint,11,opt,name=MadeModBySelf" json:"MadeModBySelf,omitempty"`
	MadeNonModBySelf       bool   `protobuf:"varint,12,opt,name=MadeNonModBySelf" json:"MadeNonModBySelf,omitempty"`
	MadeModByDefault       bool   `protobuf:"varint,13,opt,name=MadeModByDefault" json:"MadeModByDefault,omitempty"`
	MadeModByNetwork       bool   `protobuf:"varint,14,opt,name=MadeModByNetwork" json:"MadeModByNetwork,omitempty"`
	MadeNonModByNetwork    bool   `protobuf:"varint,15,opt,name=MadeNonModByNetwork" json:"MadeNonModByNetwork,omitempty"`
}

func (m *CompiledUserSignalsEntity) Reset()                    { *m = CompiledUserSignalsEntity{} }
func (m *CompiledUserSignalsEntity) String() string            { return proto.CompactTextString(m) }
func (*CompiledUserSignalsEntity) ProtoMessage()               {}
func (*CompiledUserSignalsEntity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CompiledUserSignalsEntity) GetTargetFingerprint() string {
	if m != nil {
		return m.TargetFingerprint
	}
	return ""
}

func (m *CompiledUserSignalsEntity) GetDomain() string {
	if m != nil {
		return m.Domain
	}
	return ""
}

func (m *CompiledUserSignalsEntity) GetFollowedBySelf() bool {
	if m != nil {
		return m.FollowedBySelf
	}
	return false
}

func (m *CompiledUserSignalsEntity) GetBlockedBySelf() bool {
	if m != nil {
		return m.BlockedBySelf
	}
	return false
}

func (m *CompiledUserSignalsEntity) GetFollowerCount() int32 {
	if m != nil {
		return m.FollowerCount
	}
	return 0
}

func (m *CompiledUserSignalsEntity) GetCanonicalName() string {
	if m != nil {
		return m.CanonicalName
	}
	return ""
}

func (m *CompiledUserSignalsEntity) GetCNameSourceFingerprint() string {
	if m != nil {
		return m.CNameSourceFingerprint
	}
	return ""
}

func (m *CompiledUserSignalsEntity) GetSelfPEFingerprint() string {
	if m != nil {
		return m.SelfPEFingerprint
	}
	return ""
}

func (m *CompiledUserSignalsEntity) GetSelfPECreation() int64 {
	if m != nil {
		return m.SelfPECreation
	}
	return 0
}

func (m *CompiledUserSignalsEntity) GetSelfPELastUpdate() int64 {
	if m != nil {
		return m.SelfPELastUpdate
	}
	return 0
}

func (m *CompiledUserSignalsEntity) GetMadeModBySelf() bool {
	if m != nil {
		return m.MadeModBySelf
	}
	return false
}

func (m *CompiledUserSignalsEntity) GetMadeNonModBySelf() bool {
	if m != nil {
		return m.MadeNonModBySelf
	}
	return false
}

func (m *CompiledUserSignalsEntity) GetMadeModByDefault() bool {
	if m != nil {
		return m.MadeModByDefault
	}
	return false
}

func (m *CompiledUserSignalsEntity) GetMadeModByNetwork() bool {
	if m != nil {
		return m.MadeModByNetwork
	}
	return false
}

func (m *CompiledUserSignalsEntity) GetMadeNonModByNetwork() bool {
	if m != nil {
		return m.MadeNonModByNetwork
	}
	return false
}

type AmbientBoardEntity struct {
	Fingerprint string `protobuf:"bytes,1,opt,name=Fingerprint" json:"Fingerprint,omitempty"`
	Name        string `protobuf:"bytes,2,opt,name=Name" json:"Name,omitempty"`
	LastUpdate  int64  `protobuf:"varint,3,opt,name=LastUpdate" json:"LastUpdate,omitempty"`
	LastSeen    int64  `protobuf:"varint,4,opt,name=LastSeen" json:"LastSeen,omitempty"`
}

func (m *AmbientBoardEntity) Reset()                    { *m = AmbientBoardEntity{} }
func (m *AmbientBoardEntity) String() string            { return proto.CompactTextString(m) }
func (*AmbientBoardEntity) ProtoMessage()               {}
func (*AmbientBoardEntity) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *AmbientBoardEntity) GetFingerprint() string {
	if m != nil {
		return m.Fingerprint
	}
	return ""
}

func (m *AmbientBoardEntity) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *AmbientBoardEntity) GetLastUpdate() int64 {
	if m != nil {
		return m.LastUpdate
	}
	return 0
}

func (m *AmbientBoardEntity) GetLastSeen() int64 {
	if m != nil {
		return m.LastSeen
	}
	return 0
}

func init() {
	proto.RegisterType((*CompiledBoardEntity)(nil), "feobjects.CompiledBoardEntity")
	proto.RegisterType((*CompiledThreadEntity)(nil), "feobjects.CompiledThreadEntity")
	proto.RegisterType((*CompiledPostEntity)(nil), "feobjects.CompiledPostEntity")
	proto.RegisterType((*CompiledUserEntity)(nil), "feobjects.CompiledUserEntity")
	proto.RegisterType((*CompiledContentSignalsEntity)(nil), "feobjects.CompiledContentSignalsEntity")
	proto.RegisterType((*ExplainedSignalEntity)(nil), "feobjects.ExplainedSignalEntity")
	proto.RegisterType((*CompiledUserSignalsEntity)(nil), "feobjects.CompiledUserSignalsEntity")
	proto.RegisterType((*AmbientBoardEntity)(nil), "feobjects.AmbientBoardEntity")
}

func init() { proto.RegisterFile("feobjects/feobjects.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 1136 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe4, 0x57, 0x5f, 0x6f, 0xe3, 0x44,
	0x10, 0x57, 0xea, 0x38, 0x8d, 0x37, 0xb9, 0x7f, 0xdb, 0xa3, 0xf2, 0xa1, 0x72, 0x44, 0xd1, 0x09,
	0x22, 0x04, 0x3d, 0x74, 0x27, 0x21, 0x81, 0x04, 0x52, 0x93, 0xb4, 0x12, 0xd2, 0xb5, 0x57, 0xb9,
	0x2d, 0x48, 0xbc, 0x20, 0x27, 0xde, 0x34, 0xa6, 0xee, 0xae, 0xb5, 0xde, 0x5e, 0x2f, 0x5f, 0x00,
	0x9e, 0x91, 0xf8, 0x66, 0xbc, 0xf2, 0x25, 0xf8, 0x02, 0x08, 0xcd, 0xec, 0xda, 0x59, 0xc7, 0x4e,
	0x69, 0x81, 0xb7, 0x7b, 0xdb, 0xf9, 0xed, 0xec, 0x66, 0x66, 0x7e, 0xbf, 0xd9, 0x71, 0xc8, 0x93,
	0x19, 0x13, 0x93, 0x9f, 0xd8, 0x54, 0x65, 0xcf, 0x8b, 0xd5, 0x6e, 0x2a, 0x85, 0x12, 0xd4, 0x2b,
	0x80, 0xfe, 0x6f, 0x2e, 0xd9, 0x1a, 0x89, 0xcb, 0x34, 0x4e, 0x58, 0x34, 0x14, 0xa1, 0x8c, 0xf6,
	0xb9, 0x8a, 0xd5, 0x82, 0xf6, 0x48, 0xe7, 0x20, 0xe6, 0xe7, 0x4c, 0xa6, 0x32, 0xe6, 0xca, 0x6f,
	0xf4, 0x1a, 0x03, 0x2f, 0xb0, 0x21, 0xf0, 0x38, 0x61, 0xc9, 0x6c, 0x24, 0x59, 0xa8, 0x58, 0xe4,
	0x6f, 0xf4, 0x1a, 0x83, 0x76, 0x60, 0x43, 0x94, 0x92, 0xe6, 0x51, 0x78, 0xc9, 0x7c, 0x07, 0x0f,
	0xe3, 0x1a, 0x4e, 0x8d, 0x59, 0x36, 0x95, 0x71, 0xaa, 0x62, 0xc1, 0xfd, 0xa6, 0xbe, 0xd7, 0x82,
	0xe8, 0x8f, 0x64, 0x3b, 0x0f, 0x68, 0x24, 0xb8, 0x62, 0x5c, 0x9d, 0xc4, 0xe7, 0x3c, 0x4c, 0x32,
	0xdf, 0xed, 0x35, 0x06, 0x9d, 0x17, 0x1f, 0xef, 0x2e, 0xd3, 0xa9, 0x77, 0xd4, 0x29, 0x04, 0x6b,
	0xae, 0xa1, 0x2f, 0x89, 0xfb, 0xfa, 0x9a, 0x33, 0xe9, 0xb7, 0xf0, 0xbe, 0x0f, 0x6a, 0xee, 0x3b,
	0xcb, 0x98, 0x34, 0xb7, 0x68, 0x5f, 0x88, 0x1b, 0xcb, 0x83, 0x56, 0xe6, 0x6f, 0xf6, 0x1c, 0x88,
	0xdb, 0x82, 0xe8, 0xfb, 0xa4, 0x8d, 0x89, 0x43, 0x5a, 0xed, 0x5e, 0x63, 0xe0, 0x04, 0x85, 0x4d,
	0x9f, 0x12, 0xf2, 0x2a, 0xcc, 0xd4, 0x59, 0x1a, 0x85, 0x8a, 0xf9, 0x1e, 0xee, 0x5a, 0x08, 0x54,
	0xea, 0x90, 0xa9, 0xd0, 0x27, 0xba, 0x52, 0xb0, 0xa6, 0x23, 0xd2, 0x1d, 0xcd, 0xe3, 0x24, 0x3a,
	0x9d, 0x4b, 0x16, 0x46, 0x99, 0xdf, 0xe9, 0x39, 0x83, 0xce, 0x8b, 0x0f, 0x6b, 0xa2, 0xd5, 0x1e,
	0x26, 0xde, 0xd2, 0x21, 0xda, 0x27, 0x5d, 0xb3, 0x1c, 0x89, 0x2b, 0xae, 0xfc, 0x6e, 0xaf, 0x31,
	0x70, 0x83, 0x12, 0x46, 0x77, 0x88, 0x07, 0xf9, 0x6a, 0x87, 0x7b, 0xe8, 0xb0, 0x04, 0x20, 0xf4,
	0x93, 0xab, 0x09, 0xd0, 0x33, 0x61, 0x91, 0x7f, 0x1f, 0x59, 0xb6, 0x10, 0xba, 0x4d, 0x5a, 0x47,
	0x42, 0xc5, 0xb3, 0x85, 0xff, 0x00, 0xf7, 0x8c, 0x05, 0xe5, 0x80, 0x04, 0x4f, 0x18, 0xe3, 0xfe,
	0x43, 0x5d, 0x8e, 0xdc, 0x86, 0x62, 0x7e, 0x3f, 0x8f, 0x15, 0x4b, 0xe2, 0x0c, 0xa4, 0xf3, 0x48,
	0x4b, 0xc7, 0x82, 0xfa, 0x7f, 0x3a, 0xe4, 0x71, 0x5d, 0x7a, 0xb7, 0xd0, 0xe5, 0x63, 0xe2, 0x22,
	0x2d, 0xa8, 0x48, 0x2f, 0xd0, 0xc6, 0xaa, 0x5a, 0x9d, 0xf5, 0x6a, 0x6d, 0x5a, 0x6a, 0xa5, 0xa4,
	0x39, 0x14, 0xd1, 0x02, 0x95, 0xe7, 0x05, 0xb8, 0x06, 0xec, 0x55, 0xcc, 0x2f, 0x50, 0x3d, 0x5e,
	0x80, 0xeb, 0x1b, 0x34, 0xbb, 0xf9, 0x3f, 0x6b, 0xb6, 0x7d, 0x07, 0xcd, 0xda, 0x8a, 0xf4, 0x6e,
	0x54, 0x24, 0x59, 0xab, 0xc8, 0x8e, 0xa5, 0xc8, 0xaf, 0x09, 0x41, 0x71, 0x1d, 0x8b, 0x4c, 0x65,
	0x7e, 0x17, 0xf5, 0x58, 0x17, 0x09, 0xec, 0x9b, 0x48, 0xac, 0x03, 0xf0, 0x93, 0xb8, 0xb0, 0x85,
	0x66, 0x21, 0xfd, 0xdf, 0x1d, 0x42, 0xab, 0x57, 0xfc, 0x6b, 0xc6, 0xb7, 0x49, 0x4b, 0x2b, 0xc7,
	0xbc, 0x3f, 0xc6, 0x02, 0xfc, 0x38, 0x94, 0x8c, 0x2b, 0xc3, 0xb4, 0xb1, 0x56, 0x15, 0xe2, 0xd6,
	0x2a, 0x04, 0xd5, 0xd0, 0xb2, 0xd4, 0xf0, 0x6e, 0x30, 0xff, 0x25, 0x69, 0x23, 0x91, 0x92, 0xf1,
	0xdb, 0xf1, 0x5e, 0xb8, 0xf7, 0xff, 0xd8, 0x58, 0xb2, 0xba, 0x0c, 0xf4, 0x16, 0xac, 0x7e, 0x42,
	0x1e, 0x1e, 0x09, 0x3e, 0x0a, 0xb9, 0xe0, 0xf1, 0x34, 0x4c, 0xb0, 0x37, 0x35, 0xc1, 0x15, 0xbc,
	0x94, 0xaf, 0x73, 0x63, 0xbe, 0xcd, 0x4a, 0xbe, 0xcf, 0xc8, 0x3d, 0xb0, 0x02, 0x36, 0x93, 0x2c,
	0x9b, 0x1b, 0xe6, 0x9d, 0xa0, 0x0c, 0xd2, 0xef, 0x96, 0x63, 0x12, 0xb2, 0xc8, 0x49, 0xd6, 0x23,
	0xe4, 0xd9, 0x1a, 0x52, 0xca, 0x0c, 0xd7, 0x5d, 0x00, 0x6a, 0xdc, 0x7f, 0x9b, 0xc6, 0x72, 0x81,
	0x7a, 0x71, 0x02, 0x63, 0x01, 0x0b, 0xdf, 0xf2, 0x99, 0x40, 0xd6, 0xbd, 0x00, 0xd7, 0x05, 0x33,
	0xde, 0x92, 0x99, 0xfe, 0x5f, 0x2e, 0xd9, 0xb9, 0x49, 0x57, 0xf4, 0x53, 0xf2, 0xe8, 0x34, 0x94,
	0xe7, 0x4c, 0x55, 0xcb, 0x5d, 0xdd, 0xa0, 0x3e, 0xd9, 0x3c, 0x4b, 0xdf, 0x08, 0xc5, 0x32, 0xac,
	0xb5, 0x1b, 0xe4, 0x26, 0x4c, 0x89, 0xb1, 0xb8, 0xe6, 0x7a, 0xcf, 0xd1, 0x53, 0xa2, 0x00, 0xf2,
	0xe6, 0xd1, 0xce, 0x11, 0x56, 0xd9, 0x34, 0x8f, 0x81, 0xa0, 0xcc, 0x60, 0xe6, 0x47, 0xf2, 0x06,
	0x2b, 0x83, 0x74, 0x40, 0x1e, 0x00, 0xb0, 0x77, 0x3a, 0x2e, 0xf8, 0x6c, 0x61, 0x5d, 0x56, 0x61,
	0xc8, 0xcb, 0x40, 0x16, 0xbb, 0xba, 0x86, 0xd5, 0x0d, 0xba, 0x4b, 0xa8, 0x01, 0xed, 0x32, 0xe8,
	0xe2, 0xd6, 0xec, 0xd0, 0xaf, 0xc8, 0x66, 0xc0, 0x52, 0x21, 0x55, 0xe6, 0x7b, 0xa8, 0xf7, 0x9e,
	0x45, 0xf1, 0xfe, 0xdb, 0x34, 0x09, 0x63, 0xce, 0x22, 0x5d, 0x69, 0x43, 0x6f, 0x7e, 0x80, 0x7e,
	0x43, 0xbc, 0x43, 0x11, 0x0d, 0x13, 0x31, 0xbd, 0xc8, 0xa7, 0xf6, 0x3f, 0x9f, 0x5e, 0x1e, 0xa1,
	0x63, 0xd2, 0x3d, 0x14, 0xd1, 0x5e, 0x9a, 0x4a, 0xf1, 0x06, 0x34, 0xd6, 0xbd, 0xe5, 0x15, 0xa5,
	0x53, 0xf8, 0x28, 0x2e, 0x0e, 0x45, 0x84, 0x0f, 0x6d, 0x3b, 0xd0, 0x06, 0x34, 0xd5, 0x70, 0x71,
	0x20, 0x92, 0x44, 0x5c, 0xb3, 0xe8, 0x98, 0xc9, 0x4c, 0x70, 0x33, 0xd3, 0x2b, 0x38, 0x70, 0x31,
	0x5c, 0x60, 0x4c, 0x85, 0xab, 0x1e, 0xf1, 0xab, 0x30, 0x3e, 0x8c, 0x8b, 0xd7, 0xc7, 0x38, 0xe7,
	0xdb, 0x01, 0xae, 0xa1, 0xed, 0xf2, 0x94, 0x8a, 0x11, 0x6f, 0x21, 0xa0, 0x98, 0x22, 0x5e, 0x16,
	0xf9, 0x54, 0x2b, 0xc6, 0x82, 0xaa, 0x8d, 0xb9, 0x55, 0xd3, 0x98, 0xfd, 0x5f, 0x1a, 0xe4, 0xbd,
	0xda, 0x7a, 0xc0, 0xa3, 0x70, 0x22, 0xae, 0xe4, 0x94, 0x1d, 0xa4, 0x46, 0xf0, 0x85, 0x0d, 0x6d,
	0x17, 0xb0, 0x10, 0x52, 0xd2, 0x4f, 0x8a, 0xb1, 0xfe, 0xcb, 0x43, 0xd2, 0xff, 0xd5, 0x25, 0x4f,
	0xd6, 0x76, 0xff, 0x1d, 0xfb, 0x70, 0x9b, 0xb4, 0xc6, 0xe2, 0x32, 0x8c, 0x8b, 0xf8, 0xb4, 0x45,
	0x3f, 0x22, 0xf7, 0x73, 0x96, 0x86, 0x0b, 0xd0, 0xad, 0xf9, 0x92, 0x59, 0x41, 0xa1, 0x76, 0xa6,
	0xd0, 0xc6, 0x4d, 0x77, 0x64, 0x19, 0x04, 0x2f, 0x73, 0xce, 0x7c, 0xfd, 0xb9, 0xd8, 0xd7, 0x65,
	0x10, 0xbc, 0xca, 0xaf, 0xb0, 0x9e, 0x7f, 0x65, 0x90, 0x7e, 0x41, 0xb6, 0x47, 0xb0, 0x30, 0x25,
	0xb6, 0x92, 0xdc, 0x44, 0xf7, 0x35, 0xbb, 0x79, 0x1f, 0x1f, 0xef, 0x57, 0x1b, 0xb3, 0xba, 0x01,
	0xf9, 0x6b, 0x70, 0x65, 0xbc, 0xad, 0xa0, 0xa0, 0x73, 0x8d, 0x54, 0x46, 0x5d, 0x05, 0x87, 0xfc,
	0x0e, 0xc3, 0x88, 0x81, 0x36, 0x75, 0xad, 0x3a, 0xba, 0x56, 0x25, 0x10, 0x6e, 0x04, 0xe0, 0x48,
	0xf0, 0xa5, 0x63, 0x57, 0x77, 0xce, 0x2a, 0x9e, 0xfb, 0x22, 0x30, 0x66, 0xb3, 0xf0, 0x2a, 0x51,
	0xa6, 0x0d, 0x2b, 0x78, 0xc9, 0xf7, 0x88, 0xa9, 0x6b, 0x21, 0x2f, 0xf2, 0x8e, 0x5c, 0xc5, 0xe9,
	0xe7, 0x64, 0xcb, 0xfe, 0xad, 0xdc, 0x5d, 0x77, 0x65, 0xdd, 0x56, 0xff, 0xe7, 0x06, 0xa1, 0x7b,
	0x97, 0x93, 0x98, 0x71, 0x75, 0xb7, 0x7f, 0x77, 0xf9, 0xd7, 0xf0, 0x86, 0xf5, 0x35, 0x5c, 0x6e,
	0x00, 0xa7, 0x32, 0x49, 0xed, 0x4f, 0xfe, 0x66, 0xf9, 0x93, 0x7f, 0xf8, 0xf4, 0x87, 0x9d, 0x90,
	0xa9, 0x39, 0x93, 0x9f, 0x4d, 0x85, 0x64, 0xcf, 0xf1, 0x7f, 0xa8, 0xf5, 0xc7, 0x74, 0xd2, 0x42,
	0xe4, 0xe5, 0xdf, 0x01, 0x00, 0x00, 0xff, 0xff, 0xc3, 0xa1, 0x2b, 0x9c, 0xb6, 0x0e, 0x00, 0x00,
}
