import socket

fixed_size = 1024


def server(generator):
    print(f" Generator: {generator}")
    try:
        server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        print("Socket successfully created!")
        port = 8181
        server_socket.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)

        # server_socket.bind(("localhost", port))
        server_socket.bind(("0.0.0.0", port))
        print("socket binded to %s" % (port))

        # put the socket into listening mode
        server_socket.listen(1)  # 1 Connection at a time can be handled
        print("Socket is listening...")

        conn, addr = server_socket.accept()
        print("Waiting for connetion...")
        print('Got connection from ', conn, " ", addr)
        while True:
            conn.send(b"Thanks for connecting!")
            for msg in generator:
                print(f"Sending: {msg}")
                conn.send(msg.encode())
                # conn.send(msg)
            break
    except KeyboardInterrupt:
        print("Socket closing...")
        conn.close()
        server_socket.close()


def pad_the_message(msg):
    if len(msg) >= fixed_size:
        return msg[:fixed_size]  # cutting the message
    pad = fixed_size - len(msg)
    padded = msg + b"\x00" * pad
    return padded
