import {get, readable} from "svelte/store";

export enum MessageType {
  Info = "info",
  Error = "error",
}

export interface Message {
  type: MessageType;
  title: string;
  message: any;
}

let setFunc;
export const messages = readable<Message[]>([], (set) => {
  setFunc = set

  return () => {
    setFunc = undefined
  }
});

export const showMessage = (title: string, message?: any) => {
  showMsg({type: MessageType.Info, title, message});
}

export const showError = (title: string, message?: any) => {
  showMsg({type: MessageType.Error, title, message});
}

const showMsg = (msg: Message) => {
  const msgs = get(messages);
  msgs.push(msg);
  if (setFunc) {
    setFunc(msgs);
  }

  window.setTimeout(() => {
    const msgs = get(messages);
    msgs.shift();

    if (setFunc) {
      setFunc(msgs)
    }
  }, 5000);
}
