PPU的作用：用来绘图的芯片
PPU的复杂度远远大于CPU

https://zhuanlan.zhihu.com/p/44236701

PPU的概念：
Background: 背景
Sprite: 精灵

PatternTable: 
存储的是图像数据。
总线映射：0x0000~0x1FFF, 8k

应该全部都映射到卡带上
4k背景用
4k Sprite用
16byte 一个Tile, 
4k能表示256个Tile

Tile是8x8的小像素块

也即16*16个Tile

PatternTable是静态的图案库。(当然有些游戏可以动态的改动)


VRAM的4k指的控制BG的4K, 也即NameTable的4k

NameTable: 
控制的是背景图像。
总线映射：0x2000~0x2FFF, 4k
NameTable为的PatternTable的索引
NameTable + Attribute 就是显存里的数据了
Name Table是就是PPU的显存。
NameTable中的每个字节就确定了特定的Tile（Tile的索引)


背景每个图块是8x8像素. 而FC总共是256x240即有32x30个背景图块.
每个图块用1字节表示所以一个背景差不多就需要1kb(32x30=960b).
为了实现1像素滚动, 两个背景连在一起然后用一个偏移量表示就行了:


AttributeTable：(64byte)
控制当前图像高两位palette的偏移量。


控制背景图像的1k显存：
960byte是NameTable, 64byte是Attribute Table
Attribute Table:
控制背景图像的颜色


PPU的内存映射

PPU的寄存器

图像布局：
像素:256x240 
一个tile: 8x8
所以整个图像的Tile分布: 32x30
0x2000~0x23FF: 第1kb,控制第一个图像窗格
0x2400~0x27FF: 第2kb,控制第二个图像窗格
0x2800~0x2bFF: 第3kb,
0x2C00~0x2FFF: 第4kb,

外加一个PPU滚动寄存器，就可以实现背景运动的效果了。

PPU只支持64种颜色，所以一种颜色只需要1byte即可。
作为对比，表示一种RGB颜色需要24byte,所以在颜色上压缩了3倍。

Palette(调色板)
背景调色版和精灵调色板分别16字节
0x3F00~0x3F0F: 背景调色板
0x3F10~0x3F1F: 精灵调色板

这些调色版每一个字节表示一种颜色(颜色的索引).

那么一个像素点只要用4bit索引调色板的颜色。而且这4bit还可以通过共用进一步压缩。

相比RGB的24byte, 压缩了48倍



---
1. palette完成
2. 从NES文件中还原出PatternTable
3. 从NES文件中还原出Palettes(Palettes可以有可能不是静态的)

画图流程：
从NameTable开始解析
NameTabel的一个字节代表一个Tile
使用该字节+PPUCTRL的指示，去PatternTable取出颜色Tile数据(16byte)
该16byte决定了Tile中每个像素点颜色数据的后两位bit (8*8*2 bit = 16byte)

从AttributeTable中获取每个像素点的前两个bit.

然后根据Tile的颜色索引，去拿到实际的颜色。
至此，一个那么至此一个Tile的图像数据还原完毕。

NameTable组织形式：一个字节代表一个Tile，该字节为PatternTable中的索引
一个32*30个Tile。

PatternTable组织形式: 16byte为一个单位，每16byte表示一个Tile。包含Tile中
每个像素点的后两bit。

AttributeTable: 64字节。
按Tile把屏幕分成4*4的方块。
屏幕中共有64个方块(最后一行的方块是不完整的)
AttributeTable的表示64个区域的颜色数据。
每个区域(16个tile), 1byte,
每个区域再分为4份，每一份2*2个tile
这1byte的每2bit平分给这四份。控制着这些Tile像素点颜色的前两位。





