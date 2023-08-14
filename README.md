# This projects is closed

> MAG-eta started as 7th attempt!

1. Alpha: Not published - primitives
2. Beta: Not published - primitives, fighting mechanics
3. Gamma: [here](https://github.com/rhymald/mag-gamma/tree/MBF-elemental-state-refactoring) - primitives and character
4. Delta: [here](https://github.com/rhymald/mag-delta/tree/N33-player-refactoring) - fighting mechanics, block tree
5. Epsilon: [here](https://github.com/rhymald/mag-epsilon/tree/N3G-character) - interactive CLI, trying transactional
6. Zeta: current repo - successfully transactional, with movements across world grid
7. Eta: [here](https://github.com/rhymald/mag-eta) - transactional, block tree

## N7S: Movement and world grid

Coordinates in spaces: 
- DB cache - to filter and seek targets, with row autodeletion
  - Key, TTI => ID, X, Y, time
- Trace in state - to store latest positions, with half-life deletion
  - Time => Direction, X, Y
- Simplified output
  - Direction, X, Y

### Step:
- Fetch last in trace
- If moved recently - move larger and faster

### Focus:
- Fetch last focus:
- if recently - focus fa-faster and preciser

## N5U: Effects

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

## Scripts

List all funcs and types:
```bash
grep -r "^\(func\)\|^\(type\)" . | grep Dot
```

Delete and cleanups: 
```bash
sudo docker rm $(sudo docker ps -a -f status=exited -q) && sudo docker rmi $(sudo docker images -a -q)
sudo sh -c 'truncate -s 0 /var/lib/docker/containers/*/*-json.log'
set -eu ; LANG=en_US.UTF-8 snap list --all | awk '/disabled/{print $1, $3}' | while read snapname revision; do ; snap remove "$snapname" --revision="$revision" ; done
```

### Loadtest

- __Global Change__
- `Fail point`

|When|App RAM|DB Storage|DB RAM|Results|
|:-|:-:|:-:|:-:|-:|
|__No DB__||||__Traces not cleaned up__|
|N5U|`13 GiB`|_none_|_none_|71000 pl + npc (~30 sec.)|
|N5U|`6 GiB`|_none_|_none_|45000 pl + npc (~30 sec.)|
|__2 writes per move__||__In memory__|
|N7S|`2 GiB`|0.04 / 2 GiB|2 GiB|1400 pl + 3900 npc (etern.)|
|N7S|`4 GiB`|0.06 / 2 GiB|2 GiB|1737 pl + 5362 npc (etern.)|
|N7S|`6 GiB`|0.1 / 2 GiB|_unlim._|1903 pl + 5876 npc (etern.)|
|__Trace write per move__|
|N7S|`6 GiB`|1.27 / 2 GiB|_unlim._|1494 pl + 4542 npc (etern.)|
|__Odd + Even traces, with traces cleanup__|
|N7S|`6 GiB`|1.35 / 2 GiB|_unlim._|1774 pl + 5373 npc (etern.)|
|__Step: 200 to 250, Gape: 60s to 30s__|
|N7S|`6 GiB`|1.2 / 2 GiB|_unlim._|1841 pl + 5596 npc (etern.)|
|__Step: 250 to 400, Gape: 30s to 16s__|
|N7S|`6 GiB`|1.2 / 2 GiB|_unlim._|1842 pl + 5630 npc (etern.)|

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
