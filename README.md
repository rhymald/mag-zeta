# mag-zeta
Another one, 6th trial to design P2P game server.

Delete all images: 
```bash
sudo docker rm $(sudo docker ps -a -f status=exited -q) && sudo docker rmi $(sudo docker images -a -q)
```

Run: 
```bash
sudo docker compose up --build
```