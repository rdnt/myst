import api, {ApiError} from "@/api";
import type {Invitation, User} from "@/api";
import {readable} from "svelte/store";

const getCurrentUser = (setFunc) => {
  api.currentUser().then((u: User) => {
    setFunc(u)
  }).catch((e: ApiError) => {
    console.error(e)
    setFunc(null)
  })
}

export const currentUser = readable<User | null>(null, (set) => {
  getCurrentUser(set)

  let interval = window.setInterval(() => {getCurrentUser(set)}, 10000)

  return () => {
    window.clearInterval(interval)
    set(null);
  }
});
