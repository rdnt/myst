import type {Keystore} from "@/api";
import api, {ApiError} from "@/api";
import {writable} from "svelte/store";
import type {Readable} from "svelte/store";

const {subscribe, update} = writable<Keystore[]>(undefined);

export const keystores = {subscribe} as Readable<Keystore[]>;

export const getKeystores = async () => {
    try {
        await api.keystores().then((keystores: Keystore[]) => {
            // sort by id: TODO: backend should return defined order
            keystores = keystores.sort((a, b) => {
                return a.id < b.id ? 1 : -1;
            })

            // sort entries by id: TODO: backend should return defined order
            keystores.forEach((keystore, i) => {
                keystore.entries = keystore.entries.sort((a, b) => {
                    return a.id < b.id ? 1 : -1;
                })
            })

            update(() => keystores);
            return Promise.resolve(keystores)
        }).catch((e: ApiError) => {
            console.error(e)
            return Promise.reject(e)
        })
    } catch (e: any) {
        console.error(e);
        return Promise.reject(e)
    }
}

export {};
