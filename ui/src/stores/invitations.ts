import type {Invitation} from "@/api";
import api, {ApiError} from "@/api";
import {currentUser} from "@/stores/user";
import {derived, readable} from "svelte/store";

let setInvitations;
export const invitations = readable<Invitation[]>([], (set) => {
  api.getInvitations().then((invs: Invitation[]) => {
    invs = invs.sort((a, b) => {
      return a.id < b.id ? 1 : -1;
    })

    set(invs)
  }).catch((e: ApiError) => {
    console.error(e)
    set([])
  })

  return () => {
    set([]);
  }
});

// export const invitations = derived([invitations, currentUser], ([invitations, currentUser]) => {
//   return invitations.filter((inv) => currentUser ? inv.inviteeId === currentUser.id : false)
// }, []);

let setInvitationFunc;
export const invitation = readable<Invitation>(undefined, (set) => {
  setInvitationFunc = set

  return () => {
    // set(undefined);
  }
});

export const getInvitation = (id: string) => {
  return api.getInvitation({invitationId: id}).then((inv: Invitation) => {
    if (setInvitationFunc) {
      setInvitationFunc(inv);
    }

    return Promise.resolve(inv)
  }).catch((error: Response) => {
    if (setInvitationFunc) {
      setInvitationFunc(undefined);
    }

    return Promise.reject(error)
  });
}
