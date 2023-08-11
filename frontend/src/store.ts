import { atom } from "jotai";

interface Sender {
  name: string;
  isGroup: boolean;
  profilePicture: string;
}

interface LastChat {
  from: string;
  content: string;
}

export interface ChatEntry {
  id: number;
  from: Sender;
  lastChat: LastChat;
}

export const countAtom = atom<number>(0);

const initialChatEntries: ChatEntry[] = [
  {
    id: 1,
    from: { name: "Mercur", profilePicture: "default.png", isGroup: false },
    lastChat: { content: "Pler", from: "Mercur" },
  },
];

export const chatEntriesAtom = atom<ChatEntry[]>(initialChatEntries);
export const newChatEntryAtom = atom(
  () => "",
  (get, set) => {
    get(chatEntriesAtom);
    set(chatEntriesAtom, []);
  }
);
