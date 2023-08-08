import { useState } from "react";
import { atom, useAtom } from "jotai";
import { Button } from "./components/ui/button";
import { PlusCircle } from "lucide-react";

interface Author {
  name: string;
  profilePicture: string;
}

interface LastChat {
  from: string;
  content: string;
}

interface Chat {
  id: number;
  from: Author;
  lastChat: LastChat;
}
const counter = atom(0);

function App() {
  const [count, setCount] = useAtom(counter);
  const [addChatOpened, setAddChatOpened] = useState(false);
  const [chats, setChats] = useState<Chat[] | null>([
    {
      id: 123,
      from: {
        name: "Mercur",
        profilePicture:
          "https://as2.ftcdn.net/v2/jpg/00/64/67/63/1000_F_64676383_LdbmhiNM6Ypzb3FM4PPuFP9rHe7ri8Ju.jpg",
      },
      lastChat: {
        from: "You",
        content: "pepek lorem ipsum sit dolor amet sit sit",
      },
    },
  ]);

  return (
    <>
      {addChatOpened && (
        <>
          <div className="w-screen h-screen fixed z-20 grid place-items-center">
            <Button
              variant={"secondary"}
              onClick={() => {
                setAddChatOpened(false);
              }}
            >
              Add
            </Button>
          </div>
          <div className="w-screen h-screen fixed z-10 bg-black opacity-50"></div>
        </>
      )}
      <main className="grid place-items-center">
        <div className="max-w-screen-2xl w-full h-full flex">
          <aside className="bg-black w-1/4 min-h-screen h-full flex flex-col">
            {chats ? (
              chats.map((c) => (
                <a
                  key={c.id}
                  href={`#${c.from.name}`}
                  className="text-white p-5 hover:bg-purple-900 hover:font-bold transition-all w-full overflow-hidden"
                >
                  <div className="grid">
                    <span>{c.from.name}</span>
                    <span className="text-sm opacity-50">
                      {c.lastChat.from}: {c.lastChat.content}
                    </span>
                  </div>
                </a>
              ))
            ) : (
              <span>No chats</span>
            )}
            <button
              className="text-white p-5 hover:bg-purple-900 hover:font-bold transition-all"
              onClick={() => {
                setAddChatOpened(!addChatOpened);
              }}
            >
              <div className="flex">
                <PlusCircle />
                <span className="pl-4">Add a new chat</span>
              </div>
            </button>
          </aside>
          <div className="w-full bg-zinc-600 flex flex-col items-center justify-center gap-5">
            <Button
              variant="default"
              size="lg"
              onClick={() => {
                setCount((c) => c + 1);
              }}
            >
              Add Count
            </Button>
            <span className="text-3xl">{count}</span>
          </div>
        </div>
      </main>
    </>
  );
}

export default App;
