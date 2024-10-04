sudo systemctl stop apparmor
sudo iptables -t nat -N REDSOCKS
sudo iptables -t nat -A PREROUTING --in-interface wlan0 -p tcp -j REDSOCKS
sudo iptables -t nat -A OUTPUT -o wlan0 -p tcp -j REDSOCKS
sudo iptables -t nat -A REDSOCKS -p tcp -j REDIRECT --to-ports 12345
sudo redsocks -c /etc/redsocks.conf
