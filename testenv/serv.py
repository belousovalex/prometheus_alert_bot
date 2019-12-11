from http import server
import time

start_time = time.time()

class Handler(server.SimpleHTTPRequestHandler):

    def do_GET(self):
        self.send_response(200)
        self.send_header('Content-type', 'text/plain')
        self.end_headers()
        self.wfile.write((('# HELP my_metric_count test metrics\n'
                           '# TYPE my_metric_count counter\n'
                           'my_metric_total{which="my"} %s\n'
                           'my_metric_total{which="other"} 130'
                         ) % ((time.time() - start_time) % 500)).encode("utf-8"))

def run(server_class=server.HTTPServer, handler_class=Handler):
    server_address = ('', 8000)
    httpd = server_class(server_address, handler_class)
    httpd.serve_forever()

if __name__ == "__main__":
    run()
