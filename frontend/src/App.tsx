import { useEffect, useState } from "react";
import { useAtom } from "jotai";
import { Button } from "./components/ui/button";
import { Network, PlusCircle } from "lucide-react";
import { type ChatEntry, chatEntriesAtom, countAtom } from "./store";

interface TestEvent {
  id: number;
  data: string;
}

function App() {
  const [count, setCount] = useAtom(countAtom);
  const [chatEntries, setChatEntries] = useAtom(chatEntriesAtom);

  const [addChatOpened, setAddChatOpened] = useState(false);
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [testEvents, setTestEvents] = useState<TestEvent[]>([]);

  const handleConnect = () => {
    setSocket(new WebSocket("ws://127.0.0.1:3000/ws"));
  };

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
            {chatEntries.length > 0 ? (
              chatEntries.map((c) => (
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
              <span className="text-white p-5 hover:bg-purple-900 hover:font-bold transition-all w-full overflow-hidden">
                No chats
              </span>
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
            <button
              className="text-white p-5 hover:bg-purple-900 hover:font-bold transition-all"
              onClick={handleConnect}
            >
              <div className="flex">
                <Network />
                <span className="pl-4">Connect to socket</span>
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
            {socket ? (
              <>
                <span className="text-3xl">Connected to socket</span>
                <Button
                  onClick={() => {
                    socket.send("string");
                  }}
                >
                  Send to socket
                </Button>
              </>
            ) : (
              <span className="text-3xl">Disonnected from socket</span>
            )}
          </div>
        </div>
      </main>
    </>
  );
}

export default App;