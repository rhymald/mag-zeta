# N7S: Movement and world grid

Coordinates in spaces: 
- DB cache - to filter and seek targets, with row autodeletion
  - Key, TTI => ID, X, Y, time
- Trace in state - to store latest positions, with half-life deletion
  - Time => Direction, X, Y
- Simplified output
  - Direction, X, Y

## Step:
- Fetch last in trace
- If moved recently - move larger and faster

## Focus:
- Fetch last focus:
- if recently - focus fa-faster and preciser

# N5U: Effects

Regeneration routine on each char:
- Each stream emits dot => to an effect (players only)
- Each commit emits hp => to an effect

> Regen routine creates effects that not needed an action to be created:
> - HP
> - Dots  
>
> TBD: Calm - to chill the heat  
> TBD: Antiregen for barrier

Effects consumer on each char:
- Take a portions from queue:
  - Portion depends on `Now() - Effect.Time()` differences.
  - If sum of differencec more than threshold - portion is enough.
- Sort by type:
  - Instant.
  - Conditions (TBD).
  - Delayed.
- Transform:
  - Conditions => instant + leftover condition.
  - Delayed => instant - if time is close or late.
- Transfer not consumed delayed back to queue.
- Consume instant.
- Clean queue:
  - Delete by all keys in portion.
- Sleep and restrat dependin on queue size. 

--- 

# Scripts

List all funcs and types:
```bash
grep -r "^\(func\)\|^\(type\)" . | grep Dot
```

Delete all images: 
```bash
sudo docker rm $(sudo docker ps -a -f status=exited -q) && sudo docker rmi $(sudo docker images -a -q)
```

## Loadtest

|When|App RAM|DB Storage|DB RAM|Results|
|-:|:-:|:-:|:-:|:-|
|N5U: No DB|13 GiB|-|-|71000 pl+npc (~30 sec.)|
|N5U: No DB|6 GiB|-|-|45000 pl+npc (~30 sec.)|
|N7S: 2 write per move|2 GiB|2 GiB|2 GiB|1400pl + 3900npc (etern.)|
|N7S: 2 write per move|4 GiB|2 GiB|2 GiB|1737pl + 5362npc (etern.) |
|N7S: 2 write per move|6 GiB|2 GiB|2 GiB|1903pl + 5876npc (etern.) |

## Build
Another one, 6th trial to design P2P game server.

Container: 

```bash
docker buildx build .
docker tag 33415c rhymald/mag:latest
docker push rhymald/mag:latest
```

## Run

```bash
sudo docker compose up --build
```
