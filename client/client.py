import json
import _thread

import websocket
import rel


def on_message(ws, message):
    print(message)
    print("hello")

# def listen(ws):
#     while True:
#         msg = ws.recv()
#         print(msg)
#         sleep(.3)

if __name__ == "__main__":
    websocket.enableTrace(True)
    ws = websocket.WebSocketApp("ws://localhost:9001/socket", on_message=on_message)
    ws.run_forever(dispatcher=rel)

    # t = threading.Thread(target=listen, args=(ws,))
    # t.start()

    while True:
        text = input("Input message: ")
        ws.send(json.dumps({"message": text}))
        if text == "/exit":
            break
    ws.close()
