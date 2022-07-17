import api, {ApiError} from "@/api";
import type {Keystore} from "@/api";
import {readable} from "svelte/store";

let setFunc
export const getKeystores = () => {
  return api.keystores().then((keystores: Keystore[]) => {
    keystores = keystores.sort((a, b) => {
      return a.id < b.id ? 1 : -1;
    })

    keystores.forEach((keystore, i) => {
      keystore.entries = keystore.entries.sort((a, b) => {
        return a.id < b.id ? 1 : -1;
      })
    })

    if (setFunc) {
      setFunc(keystores)
    }

    return Promise.resolve(keystores)
  }).catch((e: ApiError) => {
    console.error(e)

    if (setFunc) {
      setFunc([])
    }

    return Promise.reject(e)
  })
}

export const keystores = readable<Keystore[]>([], (set) => {
  setFunc = set
  getKeystores()

  let interval = window.setInterval(getKeystores, 10000)

  return () => {
    window.clearInterval(interval)
    setFunc([]);
    setFunc = undefined
  }
});
