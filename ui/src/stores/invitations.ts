import type {Invitation} from "@/api";
import api, {ApiError} from "@/api";
import {readable} from "svelte/store";

let setFunc
export const getInvitations = () => {
  api.getInvitations().then((invs: Invitation[]) => {
    invs = invs.sort((a, b) => {
      return a.id < b.id ? 1 : -1;
    })

    if (setFunc) {
      setFunc(invs)
    }

    return Promise.resolve(invs)
  }).catch((e: ApiError) => {
    console.error(e)

    if (setFunc) {
      setFunc([])
    }

    return Promise.reject(e)
  })
}

export const invitations = readable<Invitation[]>([], (set) => {
  setFunc = set
  getInvitations()

  let interval = window.setInterval(getInvitations, 1000)

  return () => {
    window.clearInterval(interval)
    setFunc([]);
    setFunc = undefined
  }
});
