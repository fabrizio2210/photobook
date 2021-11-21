import socket

def get_my_ip():
  ips = []
  for ip in socket.gethostbyname_ex(socket.gethostname())[2]:
    if not ip.startswith("127."):
      ips.append(ip)
  # If the previous didn't work
  if not ips:
    s = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
    s.connect(("8.8.8.8", 53))
    ips.append(s.getsockname()[0])
    s.close()
  if ips:
    return ips[0]

