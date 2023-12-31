import { useEffect, useRef, useState } from "react";
import { useAtom } from "jotai";
import { Button } from "./components/ui/button";
import { Network, PlusCircle, Send } from "lucide-react";
import { type ChatEntry, chatEntriesAtom, countAtom } from "./store";
import useWebSocket, { ReadyState } from "react-use-websocket";

type SocketError = {
  field: string;
  tag: string;
  value: string;
};

type SocketMessageResponse = {
  type: "message";
  from: string;
  body: string;
};

type SocketErrorResponse = {
  type: "error";
  errors: SocketError[];
};

type SocketResponse = SocketMessageResponse | SocketErrorResponse;

function App() {
  const [count, setCount] = useAtom(countAtom);
  const [chatEntries, setChatEntries] = useAtom(chatEntriesAtom);

  const [socketUrl, setSocketUrl] = useState("ws://127.0.0.1:3000/ws");
  const [chatInput, setChatInput] = useState<string>("");
  const [addChatOpened, setAddChatOpened] = useState(false);
  const [connected, setConnected] = useState<boolean>(false);
  const [chatHistory, setChatHistory] = useState<SocketMessageResponse[]>([]);

  const chatContainerRef = useRef<HTMLDivElement>(null);

  const { sendJsonMessage, readyState } = useWebSocket(
    `${socketUrl}/1`,
    {
      onMessage(evt) {
        const jsonData = JSON.parse(evt.data) as SocketResponse;
        console.log(jsonData);
        if (jsonData.type === "message") {
          setChatHistory((prev) => [...prev, jsonData]);
        } else {
          console.error(jsonData.errors);
        }
      },
      onOpen(evt) {
        console.log("Opened connection to socket");
      },
      onClose(evt) {
        setConnected(false);
        console.log("Closed connection to socket");
      },
    },
    connected
  );

  useEffect(() => {
    if (chatContainerRef.current)
      chatContainerRef.current.scroll({ top: chatContainerRef.current.scrollHeight });
  }, [chatHistory]);

  const handleDisconnect = () => {
    setConnected(false);
  };

  const handleConnect = () => {
    setConnected(true);
  };

  const handleSendChat = () => {
    if (chatInput === "") {
      alert("Message must have at least 1 character");
      return;
    }
    sendJsonMessage({ body: chatInput });
    setChatInput("");
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
                <span className="pl-4 text-left">Add a new chat</span>
              </div>
            </button>
            {connected && readyState === ReadyState.OPEN ? (
              <button
                className="text-white p-5 hover:bg-purple-900 hover:font-bold transition-all"
                onClick={handleDisconnect}
              >
                <div className="flex items-center">
                  <Network />
                  <span className="pl-4 text-left">Disconnect from socket</span>
                </div>
              </button>
            ) : (
              <button
                className="text-white p-5 hover:bg-purple-900 hover:font-bold transition-all"
                onClick={handleConnect}
              >
                <div className="flex items-center">
                  <Network />
                  <span className="pl-4 text-left">Connect to socket</span>
                </div>
              </button>
            )}
          </aside>
          <div className="w-full bg-zinc-600 flex flex-col items-center justify-center gap-5">
            {connected && readyState === ReadyState.OPEN ? (
              <div className="flex flex-col w-full h-screen bg-purple-950">
                <div ref={chatContainerRef} className="overflow-y-scroll h-full">
                  {chatHistory.length > 0 &&
                    chatHistory.map((c, i) => (
                      <p key={i} className="text-white text-lg">
                        <span className="opacity-60 text-sm">{c.from}</span>: {c.body}
                      </p>
                    ))}
                </div>
                <div className="w-full flex">
                  <input
                    value={chatInput}
                    onChange={(e) => setChatInput(e.target.value)}
                    type="text"
                    name="chat-input"
                    className="w-full h-11 px-2 outline-none"
                  />
                  <button
                    className="w-14 h-full grid place-items-center bg-black"
                    onClick={handleSendChat}
                  >
                    <Send className="text-white" />
                  </button>
                </div>
              </div>
            ) : (
              <>
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
                <span className="text-3xl">Disonnected from socket</span>
              </>
            )}
          </div>
        </div>
      </main>
    </>
  );
}

export default App;
