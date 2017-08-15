# regKnife

It is used to check and manipulate the values of bits in a register. This tools is useful for people who are embedded software engineers, because I am.

# usage

This is the built-in usage information:

```
Usage:
  [h]elp          : print this message.
  [p]rint         : show current value.
  [v]alue <v>     : change value to <v>.
  [s]et <f>       : set <f to 1.
  [c]lear <f>     : clear <f> to 0.
  [w]rite <f> <v> : write val <v> into field <f>.
  [l]ist [0]      : list all offsets of '1's or '0's.
  <f>             : read the value of field <f>.
  exit, quit      : exit this program.
  
Two format to represent field:
  single bit  : like 1, 3, 0
  field range : like 0:3, 3:1
```

# example

For example, there is a register in chip like:

```
Bit 0    : 1 = enable, 0 = disable.
BIt 1-2  : 0 = A channel, 1 = B channel, 2 = C channel, 3 = D channel.
Bit 3-5  : reserved.
Bit 6-15 : byte count.
```

It is a 16-bits register and there are three fields are useful. When we read this
register in driver code and get value is `0x05cd`, what does it mean ? It is time 
to check the datasheet. This is usually a bording and even anoying work, especailly
when the length of register is 32-bits.

Using this tool makes this work much easier:

1. Start this program, `-l 16` set the register length to 16 (default is 32). 
   There is a shell-like UI.
   ```
   $ regknife -l 16
   ```

2. Input the register value by 'value' command, and see the output like:
   ```
   >>> value 0x05cd
   bin: 0000,0101,1100,1101
   dec: 1485
   hex: 0x5cd
   ```
   
3. See the value of `enable` field, i.e. [0]:
   ```
   >>> 0
   bin: 1
   dec: 1
   hex: 0x1
   ```
   
4. See the `byte_count` field, i.e. [6:15]:
   ```
   >>> 6:15
   bin: 00,0001,0111
   dec: 23
   hex: 0x17
   ```
   
5. You can clear the `enable` field:
   ```	
   >>> c 0
   bin: 0000,0101,1100,1100
   dec: 1484
   hex: 0x5cc
   
   ```
   
6. You can write the `byte_count` field to 77:
   ```
   >>> w 6:15 77
   0001001101
   bin: 0001,0011,0100,1100
   dec: 4940
   hex: 0x134c
   
   >>> 6:15
   bin: 00,0100,1101
   dec: 77
   hex: 0x4d
   ```
7. After modification, see the current value by:
   ```
   >>> p
   bin: 00,0100,1101
   dec: 77
   hex: 0x4d
   ```
8. You can show all offsets of '1's in this register:
   ```
   >>> l
   6,3,2,0
   ```
   
   Or show all offsets of '0's:
   ```
   >>> l 0
   15,14,13,12,11,10,9,8,7,5,4,1
   ```
