import {readable} from "svelte/store";
import type {Keystore} from "../api/generated";
import {DefaultService} from "../api/generated";
// import api from "../api";

let setKeystores;
export const keystores = readable<Keystore[]>([], (set) => {
  setKeystores = set

  getKeystores().then((keystores) => {
    set(keystores);
  }).catch(() => {
    set([]);
  });

  return () => {
    set([]);
  }
});

export const getKeystores = () => {
  return DefaultService.keystores().then((keystores: Keystore[]) => {
    keystores = keystores.sort((a, b) => {
      return a.id < b.id ? 1 : -1;
    })

    keystores.forEach((keystore, i) => {
      keystore.entries = keystore.entries.sort((a, b) => {
        return a.id < b.id ? 1 : -1;
      })
    })

    if (setKeystores) {
      setKeystores(keystores);
    }

    return Promise.resolve(keystores)
  }).catch((error: Response) => {
    console.log(error)

    if (setKeystores) {
      setKeystores([]);
    }

    return Promise.reject(error)
  });
}
