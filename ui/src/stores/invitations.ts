import type {Invitation} from "@/api";
import api from "@/api";
import {currentUser} from "@/stores/user";
import {derived, readable} from "svelte/store";

let setInvitations;
export const invitations = readable<Invitation[]>([], (set) => {
  setInvitations = set

  return () => {
    set([]);
  }
});

export const myInvitations = derived([invitations, currentUser], ([invitations, currentUser]) => {
  console.log(invitations, currentUser)
  console.log(invitations.filter((inv) => currentUser ? inv.inviteeId === currentUser.id : false))
  return invitations.filter((inv) => currentUser ? inv.inviteeId === currentUser.id : false)
});

export const getInvitations = () => {
  return api.getInvitations().then((invs: Invitation[]) => {
    invs = invs.sort((a, b) => {
      return a.id < b.id ? 1 : -1;
    })


    if (setInvitations) {
      setInvitations(invs);
    }

    return Promise.resolve(invs)
  }).catch((error: Response) => {
    if (setInvitations) {
      setInvitations([]);
    }

    return Promise.reject(error)
  });
}
