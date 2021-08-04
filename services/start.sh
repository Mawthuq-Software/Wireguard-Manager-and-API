#!/bin/bash
sudo iptables -t nat -A POSTROUTING -o eth0 -j MASQUERADE
sudo iptables -A FORWARD -i wg0 -j ACCEPT;
sudo systemctl start coredns
sudo systemctl enable coredns
ls -a ../
sudo ../main
