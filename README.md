# goinetd

Connection redirection daemon written in Go.

Inspired by rinetd which was written in c using select to handle multiple connections. In the experience of using rinetd, it cannot reach the bandwith (1000Mbps) when redirecting a ssh connection initiated by scp.

With the goroutines and non-blocking I/O provided by golang, it's easy to write a cross platform, high performance internet daemon.

## Usage

Configuration file format:

```
bind_addr   bind_port     upstream_addr   upstream_port
```

For example:

```
192.168.1.9     443     www.baidu.com   443
192.168.1.9     80      ::1             80
0.0.0.0         8000    192.168.1.22    8000
```

forwards tcp connections
- from `192.168.1.9:443` to `www.baidu.com:443`
- from `192.168.1.9:80` to `[::1]:80`
- from `0.0.0.0:8000` to `192.168.1.22:8000`

Run:

```bash
goinetd -c goinetd.conf
```

## Project Goals

- [x] Redirect TCP Connection
- [ ] IPv6 Support (Partially at the moment)
- [ ] Redirect UDP Diagram
- [ ] Rinetd compatibility (Optional)
- [ ] Service/Systemd Support
- [ ] Binary releases

