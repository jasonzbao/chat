import json

import asyncio
import websockets

async def hello():
    uri = "ws://localhost:9001/socket"
    async with websockets.connect(uri) as websocket:
        async def consumer_handler(websocket):
            async for message in websocket:
                print(message)

        while True:
            text = input("Input message: ")
            await websocket.send(json.dumps({"message": text}))
            if text == "/exit":
                break

        greeting = await websocket.recv()
        print(f"<<< {greeting}")

if __name__ == "__main__":
    asyncio.run(hello())
