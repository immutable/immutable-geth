// Copyright 2024 The Immutable go-ethereum Authors
// This file is part of the Immutable go-ethereum library.
//
// The Immutable go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The Immutable go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the Immutable go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package role

import "fmt"

var (
	Validator     = Role{"validator"}
	Boot          = Role{"boot"}
	RPC           = Role{"rpc"}
	Partner       = Role{"partner"}
	PartnerPublic = Role{"partner-public"}
	AllRoles      = []Role{Validator, Boot, RPC, Partner, PartnerPublic}
)

// Count returns the number of supported roles
func Count() int {
	return len(AllRoles)
}

// Role is the role of a node in the network
type Role struct {
	string
}

// String returns the string representation of the role
func (r Role) String() string {
	return r.string
}

// Labels returns the key/value labels for the role
// which is assigned to k8s resources
func (r Role) Labels() map[string]string {
	return map[string]string{
		"app":   "zkevm-geth",
		"zkevm": fmt.Sprintf("geth-%s", r.String()),
	}
}

// ExternalSecretName returns the name of the k8s external secret
// used for the role
func (r Role) ExternalSecretName() string {
	return fmt.Sprintf("zkevm-geth-%s-chain", r.String())
}

// NewFromString returns the role from a string
// which must be one of the predefined Role* values
func NewFromString(name string) (Role, error) {
	switch name {
	case Validator.String():
		return Validator, nil
	case Boot.String():
		return Boot, nil
	case RPC.String():
		return RPC, nil
	case Partner.String():
		return Partner, nil
	case PartnerPublic.String():
		return PartnerPublic, nil
	default:
		return Role{}, fmt.Errorf("no such role for: %s", name)
	}
}
