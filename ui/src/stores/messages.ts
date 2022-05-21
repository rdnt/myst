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

let setMessagesFunc;
export const messages = readable<Message[]>([], (set) => {
  setMessagesFunc = set

  return () => {
    setMessagesFunc([]);
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
  setMessagesFunc(msgs);

  window.setTimeout(() => {
    const msgs = get(messages);
    msgs.shift();
    setMessagesFunc(msgs)
  }, 5000);
}
