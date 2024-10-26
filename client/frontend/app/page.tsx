"use client";
import Image from "next/image";
import { useState, useEffect } from "react";
export default function Home() {
  // connect to websocket server first
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [pc, setPs] = useState<RTCPeerConnection | null>(null);

  useEffect(() => {
    const ws = new WebSocket(process.env.NEXT_PUBLIC_WEBSOCKET_SERVER_API!);
    setSocket(ws);

    ws.onopen = () => {
      console.log("WebSocket connection established");
    };
    ws.onclose = () => {
      console.log("WebSocket connection closed");
    };

    ws.onmessage = (event) => {
      console.log(event.data);
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <div className="grid grid-rows-[20px_1fr_20px] items-center justify-items-center min-h-screen p-8 pb-20 gap-16 sm:p-20 font-[family-name:var(--font-geist-sans)]">
      <main className="flex flex-col gap-8 row-start-2 items-center sm:items-start">
        hi
      </main>
    </div>
  );
}
