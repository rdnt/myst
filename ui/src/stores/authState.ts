import type {User} from "@/api";
import api, {ApiError} from "@/api";
import {writable} from "svelte/store";
import type {Readable} from "svelte/store";

export enum AuthState {
    Onboarding = "ONBOARDING",
    SignedOut = "SIGNED_OUT",
    SignedIn = "SIGNED_IN",
}

const {subscribe, update} = writable<AuthState | undefined>(undefined);

export const authState = {subscribe} as Readable<AuthState | undefined>;

export const setAuthState = update;

export const getAuthState = async (): Promise<AuthState | undefined> => {
    try {
        return await api.enclave().then(() => {
            update(() => AuthState.SignedIn)
            return AuthState.SignedIn
        }).catch(err => {
            if (err.status == 404) {
                update(() => AuthState.Onboarding)
                return AuthState.Onboarding
            } else if (err.status == 401) {
                // if (state === "signedIn") {
                //     // was logged out, notify
                //     showError("Signed out", "You have been signed out due to inactivity.");
                // }
                update(() => AuthState.SignedOut)
                return AuthState.SignedOut
            }
        })
    } catch (e: any) {
        console.error(e);
        return Promise.reject(e)
    }
}
