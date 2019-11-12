// Copyright (c) 2019, AT&T Intellectual Property. All rights reserved.
//
// Copyright (c) 2014 by Brocade Communications Systems, Inc.
// All rights reserved.
//
// SPDX-License-Identifier: MPL-2.0

package group

//Lookup looks up a group by its name. If the group cannot be found,
//the error is of type UnknownGroupError.
func Lookup(groupname string) (*Group, error) {
	return lookup(groupname)
}

//LookupId looks up a group by its gid. If the group cannot be found,
//the error is of type UnknownGroupIdError.
func LookupId(gid string) (*Group, error) {
	return lookupId(gid)
}

//List returns a list of all group entries on the system.
func List() ([]*Group, error) {
	return list()
}

//LookupUser returns a list of groups a given username is associated with.
//If the no groups can be found, the error is of type user.UnknownUserError.
func LookupUser(username string) ([]*Group, error) {
	return lookupUser(username)
}

//LookupUId returns a list of groups a given userid is associated with.
//If the no groups can be found, the error is of type user.UnknownUserIdError.
func LookupUid(uid string) ([]*Group, error) {
	return lookupUid(uid)
}
