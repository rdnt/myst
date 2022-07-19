import type {Invitation} from "@/api";
import api, {ApiError} from "@/api";
import {currentUser} from "@/stores/user";
import {get, readable} from "svelte/store";

let setFunc
export const getInvitations = () => {
  if (!get(currentUser)) {
    return []
  }

  return api.getInvitations().then((invs: Invitation[]) => {
    invs = invs.sort((a, b) => {
      return a.updatedAt < b.updatedAt ? 1 : -1;
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

  let interval = window.setInterval(getInvitations, 5000)

  return () => {
    window.clearInterval(interval)
    setFunc([]);
    setFunc = undefined
  }
});
