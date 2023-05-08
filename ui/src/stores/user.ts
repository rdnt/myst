import api, {ApiError} from "@/api";
import type {User} from "@/api";
import {readable} from "svelte/store";
import {getInvitations} from "@/stores/invitations";

let setFunc
export const getCurrentUser = () => {
  return api.currentUser().then((u: User) => {
    if (setFunc) {
      setFunc(u)
    }

    return Promise.resolve(u)
  }).catch((e: ApiError) => {
    if (e.status === 404) {
      if (setFunc) {
        setFunc(null)
      }
    } else if (e.status === 401) {
      if (setFunc) {
        setFunc(undefined)
      }
    } else {
      console.log(e)
      return Promise.reject(e)
    }
  })
}

export const currentUser = readable<User>(undefined, (set) => {
  setFunc = set
  getCurrentUser().then(() => {
    getInvitations()
  })

  let interval = window.setInterval(getCurrentUser, 5000)

  return () => {
    window.clearInterval(interval)
    setFunc = undefined
  }
});
