import zmq
import json
context = zmq.Context.instance()

sock = context.socket(zmq.PUSH)
sock.bind('tcp://*:5000')

sock2 = context.socket(zmq.PUSH)
sock2.bind('tcp://*:6000')

while True:
    data = json.dumps({
                '@source': "file://logstash.vwdl.bcinfra.net/some/file/name.log",
                '@type': "apache",
                '@tags': ["some","tags"],
                '@fields': {"some": "field"},
                '@timestamp': "12345",
                '@source_host': "logstash.vwdl.bcinfra.net",
                '@source_path': "/some/file/name.log",
                '@message': "Four score and seventy years ago",
    })
    sock.send(data)
    sock2.send(data)
    #response = sock.recv()
