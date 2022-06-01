import api from "@/api";
import type {Invitation, User} from "@/api";
import {readable} from "svelte/store";

let setCurrentUser;
export const currentUser = readable<User>(undefined, (set) => {
  setCurrentUser = set

  return () => {
    set(undefined);
  }
});

export const getCurrentUser = () => {
  return api.currentUser().then((u: User) => {
    if (setCurrentUser) {
      setCurrentUser(u);
    }

    return Promise.resolve(u)
  }).catch((error: Response) => {
    if (setCurrentUser) {
      setCurrentUser(undefined);
    }

    return Promise.reject(error)
  });
}
