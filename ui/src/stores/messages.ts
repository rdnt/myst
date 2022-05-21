import {readable} from "svelte/store";

export interface MystError {
  title: string;
  message: string;
}

export enum MessageType {
  Status = "status",
  Info = "info",
  Error = "error",
}

export interface Message {
  title: string;
  message: string;
  type: MessageType;
}

let setMessagesFunc;
export const messages = readable<Message[]>(undefined, (set) => {
  setMessagesFunc = set

  return () => {
    setMessagesFunc([]);
  }
});

export const showMessage( ) => {
  setMessagesFunc([msg]);
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
