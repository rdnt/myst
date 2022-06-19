import type {Invitation} from "@/api";
import api, {ApiError} from "@/api";
import {readable} from "svelte/store";

const getInvitations = (setFunc) => {
  api.getInvitations().then((invs: Invitation[]) => {
    invs = invs.sort((a, b) => {
      return a.id < b.id ? 1 : -1;
    })

    setFunc(invs)
  }).catch((e: ApiError) => {
    console.error(e)
    setFunc([])
  })
}

export const invitations = readable<Invitation[]>([], (set) => {
  getInvitations(set)

  let interval = window.setInterval(() => {getInvitations(set)}, 10000)

  return () => {
    window.clearInterval(interval)
    set([]);
  }
});
