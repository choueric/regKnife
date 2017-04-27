# regKnife

This tools is useful for people who are embedded software engineers, because I am.

For example, there is a register in chip like:

```
Bit 0    : 1 = enable, 0 = disable.
BIt 1-2  : 0 = A channel, 1 = B channel, 2 = C channel, 3 = D channel.
Bit 3-5  : reserved.
Bit 6-15 : byte count.
```

It is a 16-bits register and there are three fields are useful. When we read this
register in driver code and get value is 0x05cd, what does it mean ? It is time 
to check the datasheet now. It is a bording and anoying work, especailly when the
register is 32-bits.

Using this tool makes this work much eaiser:

1. Start this program, `-l 16` set the register length to 16 (default is 32). 
   There is a shell-like UI.
   
   $ regknife -l 16

2. Input the register value by 'value' command, and see the output like:
   
   >>> value 0x05cd
   bin: 0000,0101,1100,1101
   dec: 1485
   hex: 0x5cd
   
3. See the value of bit 0, and, of course, it is 1:
   
   >>> 0
   bin: 1
   dec: 1
   hex: 0x1
   
4. See the byte count filed, it's transferred 23 bytes:

   >>> 6:15
   bin: 00,0001,0111
   dec: 23
   hex: 0x17
   
5. You can set bit 2:
	
   >>> s 2
   bin: 0000,0101,1100,1101
   dec: 1485
   hex: 0x5cd
