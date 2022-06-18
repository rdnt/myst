import api, {ApiError} from "@/api";
import type {Invitation, User} from "@/api";
import {readable} from "svelte/store";

export const currentUser = readable<User>(null, (set) => {
  api.currentUser().then((u: User) => {
    set(u)
  }).catch((e: ApiError) => {
    console.error(e)
    set(null)
  })

  return () => {
    set(null);
  }
})
