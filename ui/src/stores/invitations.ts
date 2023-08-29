import type {Invitation} from "@/api/generated/models/Invitation";
import api, {ApiError} from "@/api";
import {writable} from "svelte/store";
import type {Readable} from "svelte/store";

const {subscribe, update} = writable<Invitation[]>(undefined);

export const invitations = {subscribe} as Readable<Invitation[]>;

export const getInvitations = async () => {
  try {
    await api.getInvitations().then((invitations: Invitation[]) => {
      // sort by updatedAtDesc: TODO: backend should return defined order
      invitations = invitations.sort((a, b) => {
        return a.updatedAt < b.updatedAt ? 1 : -1;
      })

      update(() => invitations);
      return Promise.resolve(invitations)
    }).catch((e: ApiError) => {
      console.error(e)
      return Promise.reject(e)
    })
  } catch (e: any) {
    console.error(e);
    return Promise.reject(e)
  }
}
