import type {User} from "@/api";
import api, {ApiError} from "@/api";
import {writable} from "svelte/store";
import type {Readable} from "svelte/store";

const {subscribe, update} = writable<User | null | undefined>(undefined);

export const currentUser = {subscribe} as Readable<User | null>;

export const getCurrentUser = async (): Promise<User | null> => {
    try {
        return await api.currentUser().then((u: User) => {
            update(() => u);
            return u;
        }).catch((e: ApiError) => {
            if (e.status === 404) {
                update(() => null);
                return null;
            }

            console.error(e)
            return Promise.reject(e)
        })
    } catch (e: any) {
        console.error(e);
        return Promise.reject(e)
    }
}

export {};
