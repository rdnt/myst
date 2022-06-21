import api, {ApiError} from "@/api";
import type {User} from "@/api";
import {readable} from "svelte/store";

let setFunc
export const getCurrentUser = () => {
  return api.currentUser().then((u: User) => {
    if (setFunc) {
      setFunc(u)
    }

    return Promise.resolve(u)
  }).catch((e: ApiError) => {
    console.error(e)

    if (setFunc) {
      setFunc(null)
    }

    return Promise.reject(e)
  })
}

export const currentUser = readable<User>(null, (set) => {
  setFunc = set
  getCurrentUser()

  let interval = window.setInterval(getCurrentUser, 60000)

  return () => {
    window.clearInterval(interval)
    setFunc = undefined
  }
});
