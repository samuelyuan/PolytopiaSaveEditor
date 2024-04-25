# PolytopiaSaveEditor

## How to Use

### Decompress the file

The initial save file is compressed using the LZ4 algorithm. Make sure to provide the file with ".state" extension.

```
./PolytopiaSaveEditor -input=[original save file name] -mode=[decompress]
```

The new file will be created and the new filename will be the same as the original file name with ".decomp" appended at the end.

### Modify the file

There are various commands to modify the file. All changes must be be done on the decompressed file because the compressed file removes a lot of bytes for redundancy. If you modify one byte in the compressed file, you might end up modifying multiple bytes in the decompressed file.

Read-only commands:

1. list-cities: List all cities owned by each tribe. The city name and location will be shown.
2. list-units: List all units owned by each tribe. The unit type and location will be shown.

Commands to modify indiviudal tiles

1. modify-unit-tribe: Modify the unit to switch to a different tribe, but keep the same type.

```
./PolytopiaSaveEditor -input=[decompressed file name] -mode=modify-unit-tribe -x=[tile's x value] -y=[tile's y value] -value=[new tribe]
```

2. modify-unit-type: Modify the unit type, but keep the same tribe.

```
./PolytopiaSaveEditor -input=[decompressed file name] -mode=modify-unit-type -x=[tile's x value] -y=[tile's y value] -value=[new tribe]
```

3. reveal-tile: Reveal the tile for a player.

```
./PolytopiaSaveEditor -input=[decompressed file name] -mode=reveal-tile -x=[tile's x value] -y=[tile's y value] -value=[player's tribe]
```

Commands to modify multiple units:

1. convert-tribe: Change all units under one tribe to be under a different tribe. If you use this command, you can convert all units from another player to become your units.

```
./PolytopiaSaveEditor -input=[decompressed file name] -mode=convert-tribe -oldvalue=[old tribe] -value=[new tribe]
```

2. convert-all-units: This will change all units in the game to be under one tribe. This is more efficient to run if you want to convert all tribes in one command.

```
./PolytopiaSaveEditor -input=[decompressed file name] -mode=convert-all-units -value=[new tribe]
```

3. reveal-all-tiles: This will reveal all tiles on the map for one player. Using this command will break the Explorer task because the game won't trigger the logic that updates and checks that you saw all the lighthouses on the map.


```
./PolytopiaSaveEditor -input=[decompressed file name] -mode=reveal-all-tiles -value=[player's tribe]
```

The commands to modify the file will make changes within the decompressed file.

```
./PolytopiaSaveEditor -input=[decompressed file name] -mode=[command]
```

### Compress the file

The file will be compressed with the LZ4 algorithm. Make sure to provide the file with ".decomp" extension, which is the decompressed file.

```
./PolytopiaSaveEditor -input=[decompresssed file name] -mode=[compress]
```

A new file will be created and the new filename will be the same as the decompressed file name with ".comp" appended to the end.

### Overwrite existing save

Make sure you quit your current Polytopia game and go to the main menu before overwriting the save file. If you overwrite the file while the game is still in progress, the game will overwrite the file when you leave and none of your new changes will apply.

1. Change file extension to just be ".state"
2. Copy new file to save directory and overwrite existing save file in Singleplayer/ folder.
3. Go the main menu and click "Resume Game"
4. You should see all your changes take effect in the game. 