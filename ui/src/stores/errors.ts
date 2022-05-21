import {readable} from "svelte/store";

export interface MystError {
  title: string;
  message: string;
}

let setErrorFunc;
export const error = readable<MystError | undefined>(undefined, (set) => {
  setErrorFunc = set

  return () => {
    set(undefined);
  }
});

export const showError = (title: string, message: string) => {
  setErrorFunc({ title, message })
  window.setTimeout(() => {
    setErrorFunc(undefined)
  }, 10000);
}
