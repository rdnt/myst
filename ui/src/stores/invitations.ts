import api from "@/api";
import type {Invitation} from "@/api";
import {readable} from "svelte/store";

let setInvitations;
export const invitations = readable<Invitation[]>([], (set) => {
  setInvitations = set

  return () => {
    set([]);
  }
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
