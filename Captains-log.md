## A collection of thoughts

### Aug 16, 2025
> The old captains log was lost(so ein dreck es regt mich noch immer auf)\
> I reset my computer and forgot to push so I lost a few hours of progress.\
> But I had this file on my .gitignore so it was never pushed.
> I lost all progress on the file on that day,\
> and stopped writing


### Aug 27, 2025

- planning the OS and inputs

Reading input: MxsxllOS

it isn't a single task but depends on what aplication is using it. <br>
The most common one being the console.

it reads input with a task that loops reading the keyboard buffer. <br>
When the task  is spawned it allocates a space in the heap, <br>
and uses the current cursor postion to find out where the char should be placed in the Videochar table see below


### Videochar table:

**IDEA 1**

[this is for better placement of chars on the screen when the video mode is TILE<br>
it represents the frame buffer, (256 x 256 2bpp so 16KB and ≈65K pixels )<br>
as a 32 x 32 because 256 / 8(char size) table or 1024/1KB.<br>
Each entry is 1 byte, as it contains the char(repr. as 0-256) in the space.<br>
3 byte infront of the table is the rollptr.<br>
It signals to the redrawing task what line is the top most to make scrolling easier.<br>

**Imagine this as the table**

3<br>
00000000<br>
00000000<br>
00000000<br>
00000000<br>
00000000<br>
00000000<br>
00000000<br>
00000000<br>

so if the byte is 3 it will now be rendered with 3 as the top and clear line 5 while putting it to the top. <br>
This does limit the screen to only scrolling down as scrolling up would require far more space<br>
The cursor is one word before the table and is byte 1 pos x byte 2 pos y]

**IDEA 2**

[this is for better placement of chars on the screen when the video mode is TILE<br>
it represents the frame buffer, (256 x 256 2bpp so 16KB and ≈65K pixels ) as a<br>
64 x 64 because 1024 total chars * wordsize table or 2048/2KB.<br>
Each one of the cells holds 1 word/addr to the heap where the chars/strings are stored.<br>
when adding to this the lenght of the previous string is respected and the next word has the ammount of free cells,<br>
so that the chars aren't drawn over each other many would belive that this would be a problem because of the inability<br>
to tell wether smth is a string or just the offset. But because a whole word is used for the cell the possible values range from<br>
0 to ≈65K but the heap's highest addr is 23964 or in binary 0101110110011100 when using the full 16 bits. !notice that the<br>
highest bit isn't used. this now functions as a flag to tell if the next word is an entry or an offset.<br>
the scrolling/cursor logic is the same]<br>


**Pros V1:**<br>
smaller: half the size 1KB<br>
faster: renderer is directly able to draw the char to the screen<br>
simpler: easy implementation all reading is the same<br>

**Cons V1:**<br>
incomplete: when the renderer clears a line or the screen all strings should be freed to avoid the heap filling up<br>

**V1-extra:**<br>
always the same speed no matter how much is displayed on screen<br>


**Pros V2:**<br>
heap-storage: renderer or clean-up task is able to free the heap directly without having to do any work arrounds<br>
gets faster the less is on screen or the bigger the individual strings are.<br>

**Cons V2:**<br>
bigger    : double the size<br>
slower    : renderer has to read from the heap<br>
complexer: extra logic required to skip over offsets<br>